package handler

import (
	pb_goods "albelt.top/mxshop_protos/albelt/good_srv/go"
	pb "albelt.top/mxshop_protos/albelt/order_srv/go"
	pb_stock "albelt.top/mxshop_protos/albelt/stock_srv/go"
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/model"
)

/*
	从购物车创建订单接口使用分布式事务:
		1.基于rocketMQ事务消息实现库存的一致性
		2.使用rocketMQ的延迟消息来实现订单的超时取消
	具体原理见order_v2.png
*/

type OrderListener struct {
	Order      *model.Order
	OrderGoods []*model.OrderGoods
	Ctx        context.Context
	Err        error
}

func NewOrderListener(order *model.Order, ctx context.Context) *OrderListener {
	return &OrderListener{
		Order: order,
		Ctx:   ctx,
	}
}

//执行本地事务
func (l *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var (
		err        error
		cartItems  []*model.ShoppingCartItem
		orderGoods []*model.OrderGoods
	)

	order := l.Order
	ctx := l.Ctx

	// 链路追踪
	serverSpan := opentracing.SpanFromContext(l.Ctx)

	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		span1 := global.Tracer.StartSpan("get_items_from_cart", opentracing.ChildOf(serverSpan.Context()))

		// 1.从用户的购物车获取已选中的商品
		res := global.Mysql.Where("user_id = ? AND checked = ?", order.UserId, true).Find(&cartItems)
		if res.RowsAffected == 0 {
			return status.Errorf(codes.InvalidArgument, "购物车中没有选中的商品")
		}

		var goodsIds []int32
		goodsId2Num := map[int32]int32{}
		for _, item := range cartItems {
			goodsIds = append(goodsIds, item.GoodsId)
			goodsId2Num[item.GoodsId] = item.Num
		}

		span1.Finish()

		span2 := global.Tracer.StartSpan("cal_total_amount", opentracing.ChildOf(serverSpan.Context()))

		// 2.查询商品金额,计算订单总金额
		batchGetGoodsResp, err := global.GoodsSrvCli.BatchGetGoods(ctx, &pb_goods.BatchGetGoodsRequest{GoodsIds: goodsIds})
		if err != nil {
			return status.Errorf(codes.Internal, "批量查询商品信息失败:%s", err.Error())
		}
		var orderAmount float32
		for _, g := range batchGetGoodsResp.Goods {
			num := goodsId2Num[g.Id]
			orderAmount += g.ShopPrice * float32(num)
		}

		span2.Finish()

		span3 := global.Tracer.StartSpan("decrease_stocks", opentracing.ChildOf(serverSpan.Context()))

		// 3.扣减库存
		decreaseStocksReq := &pb_stock.DecreseStocksRequest{OrderSn: order.No}
		for goodsid, num := range goodsId2Num {
			decreaseStocksReq.Stocks = append(decreaseStocksReq.Stocks, &pb_stock.Stock{
				GoodsId: goodsid,
				Number:  num,
			})
		}
		_, err = global.StockSrvCli.DecreseStocks(ctx, decreaseStocksReq)
		if err != nil {
			return status.Errorf(codes.Internal, "扣减库存失败:%s", err.Error())
		}

		span3.Finish()

		span4 := global.Tracer.StartSpan("create_order_and_clear_cart", opentracing.ChildOf(serverSpan.Context()))

		// 4.创建订单
		order.Amount = orderAmount
		if err = global.Mysql.Create(&order).Error; err != nil {
			return status.Errorf(codes.Internal, "创建订单失败:%s", err.Error())
		}

		for _, g := range batchGetGoodsResp.Goods {
			orderGoods = append(orderGoods, &model.OrderGoods{
				OrderId:    order.ID,
				GoodsId:    g.Id,
				GoodsNum:   goodsId2Num[g.Id],
				GoodsName:  g.Name,
				GoodsPrice: g.ShopPrice,
			})
		}
		if err = global.Mysql.Create(&orderGoods).Error; err != nil {
			return status.Errorf(codes.Internal, "创建订单商品失败:%s", err.Error())
		}

		// 5.从购物车中清除已购买的商品
		if err = global.Mysql.Delete(&cartItems).Error; err != nil {
			return status.Errorf(codes.Internal, "从购物车中清除已购买商品失败:%s", err.Error())
		}

		span4.Finish()

		// 测试失败情况
		//return errors.New("xxx原因失败")

		return nil
	})

	// 本地事务若失败，需要提交补偿库存的消息
	if err != nil {
		zap.S().Info("执行本地事务失败，提交补偿库存的消息")
		l.Err = err
		return primitive.CommitMessageState
	}

	// 本地事务若成功，则回滚补偿库存的消息
	l.OrderGoods = orderGoods
	zap.S().Info("执行本地事务成功，回滚补偿库存的消息")

	// 发送订单超时的延迟消息(发送失败不影响订单创建)
	if err = sendDelayedMsgForOrderTimeout(ctx, order.No); err != nil {
		zap.S().Warnf("发送订单超时的延迟消息失败:%s", err.Error())
	}

	return primitive.RollbackMessageState
}

// 测试回查
//func (l *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
//	fmt.Printf("-- ExecuteLocalTransaction --\n")
//	return primitive.UnknowState
//}

func (l *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	// 解析sn
	var orderDetails pb_stock.OrderDetails
	if err := json.Unmarshal(msg.Body, &orderDetails); err != nil {
		zap.S().Warnf("CheckLocalTransaction -> decode msg err:%s", err.Error())
		return primitive.UnknowState
	}

	// 通过sn查询订单
	var order model.Order
	res := global.Mysql.
		Where("no = ? AND status = ?", orderDetails.OrderSn, pb.OrderStatus_ORDER_STATUS_PAYING).
		First(&order)

	// 未找到待支付的订单，则需要提交补偿库存的消息
	if res.RowsAffected == 0 {
		zap.S().Infof("CheckLocalTransaction CommitMessage")
		return primitive.CommitMessageState
	}

	// 找到待支付的订单，则回滚补偿库存的消息
	zap.S().Infof("CheckLocalTransaction RollbackMessage")
	return primitive.RollbackMessageState
}

// 发送补偿库存的半消息
func sendRestoreStockHalfMsg(ctx context.Context, order *model.Order) (*OrderListener, error) {
	// 创建listener
	l := NewOrderListener(order, ctx)

	// 创建事务型生产者并启动
	p, err := rocketmq.NewTransactionProducer(
		l,
		producer.WithNameServer([]string{global.Config.RocketMq.NameServer}),
		producer.WithGroupName(global.Config.RocketMq.ProducerGroupName),
		producer.WithRetry(int(global.Config.RocketMq.Retry)),
	)

	if err = p.Start(); err != nil {
		return l, err
	}

	// 发送半消息，消息体中包含订单的sn
	orderDetails := pb_stock.OrderDetails{OrderSn: order.No}
	bs, err := json.Marshal(&orderDetails)
	if err != nil {
		return l, err
	}

	// 发送半消息成功后立即执行本地事务
	_, err = p.SendMessageInTransaction(
		ctx,
		primitive.NewMessage(global.Config.RocketMq.CreateTopic, bs),
	)
	// 处理发送消息失败的情况
	if err != nil {
		return l, err
	}

	zap.S().Infof("sendRestoreStockHalfMsg:%s", bs)

	return l, nil
}

// 使用购物车创建订单(分布式事务版本)
func (s *OrderService) CreateOrderFromCart(ctx context.Context, req *pb.CreateOrderFromCartRequest) (*pb.CreateOrderFromCartResponse, error) {
	order := model.Order{
		No:          genOrderNo(),
		PayType:     req.Order.PayType,
		Status:      pb.OrderStatus_ORDER_STATUS_PAYING,
		UserId:      req.UserId,
		UserName:    req.Order.UserName,
		UserMobile:  req.Order.UserMobile,
		UserAddress: req.Order.UserAddress,
	}

	// 发送半消息并执行本地事务
	listener, err := sendRestoreStockHalfMsg(ctx, &order)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "发送半消息失败")
	}
	if listener.Err != nil {
		// 创建订单失败
		return nil, listener.Err
	}

	return &pb.CreateOrderFromCartResponse{Order: NewPbOrder(listener.Order, listener.OrderGoods)}, nil
}

/* -----------------------------------订单超时处理---------------------------------------- */
type OrderTimeoutMsg struct {
	OrderSn string `json:"order_sn"`
}

// 发送订单超时处理的延时消息
func sendDelayedMsgForOrderTimeout(ctx context.Context, orderSn string) error {
	// 发送延时消息
	msgBody := OrderTimeoutMsg{OrderSn: orderSn}
	bs, err := json.Marshal(msgBody)
	if err != nil {
		return err
	}
	msg := primitive.Message{
		Topic: global.Config.RocketMq.DelayTopic,
		Body:  bs,
	}
	msg.WithDelayTimeLevel(int(global.Config.RocketMq.DelayLevel))

	if _, err = global.RocketMqProducer.SendSync(ctx, &msg); err != nil {
		return err
	}

	zap.S().Infof("发送订单超时处理的延时消息:%s", bs)

	return nil
}

// 订单超时处理
func HandleOrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	sendRestoreStockMsg := func(orderSn string) error {
		orderDetail := pb_stock.OrderDetails{OrderSn: orderSn}
		bs, err := json.Marshal(&orderDetail)
		if err != nil {
			return err
		}

		_, err = global.RocketMqProducer.SendSync(
			ctx, primitive.NewMessage(global.Config.RocketMq.CreateTopic, bs),
		)
		if err != nil {
			return err
		}

		zap.S().Infof("订单超时处理，发送补偿库存消息(order_sn=%s)", orderSn)

		return nil
	}

	for _, msg := range msgs {
		// 解析订单的sn
		var msgBody OrderTimeoutMsg
		if err := json.Unmarshal(msg.Body, &msgBody); err != nil {
			zap.S().Warnf("decode msg err:%s", err.Error())
			continue
		}

		err := global.Mysql.Transaction(func(tx *gorm.DB) error {
			// 查询订单
			var order model.Order
			res := tx.Where("no = ? AND status = ?", msgBody.OrderSn, pb.OrderStatus_ORDER_STATUS_PAYING).
				First(&order)
			if res.RowsAffected == 0 {
				return nil
			}

			// 更新订单状态为"超时关闭"
			err2 := tx.Model(&model.Order{}).Where("id = ?", order.ID).
				Update("status", pb.OrderStatus_ORDER_STATUS_CLOSED).Error
			if err2 != nil {
				return err2
			}

			// 发送补偿库存的消息,让库存服务去处理
			if err2 = sendRestoreStockMsg(order.No); err2 != nil {
				return err2
			}

			return nil
		})

		if err != nil {
			zap.S().Warnf("订单(sn=%s)超时处理失败:%s", msgBody.OrderSn, err.Error())
			return consumer.ConsumeRetryLater, nil
		}

	}

	return consumer.ConsumeSuccess, nil
}
