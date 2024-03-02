package handler

import (
	pb_goods "albelt.top/mxshop_protos/albelt/good_srv/go"
	pb "albelt.top/mxshop_protos/albelt/order_srv/go"
	pb_stock "albelt.top/mxshop_protos/albelt/stock_srv/go"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/model"
)

// 使用购物车创建订单
func (s *OrderService) CreateOrderFromCartV1(ctx context.Context, req *pb.CreateOrderFromCartRequest) (*pb.CreateOrderFromCartResponse, error) {
	var (
		err        error
		cartItems  []*model.ShoppingCartItem
		orderGoods []*model.OrderGoods
		order      model.Order
	)

	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		// 1.从用户的购物车获取已选中的商品
		res := global.Mysql.Where("user_id = ? AND checked = ?", req.UserId, true).Find(&cartItems)
		if res.RowsAffected == 0 {
			return status.Errorf(codes.InvalidArgument, "购物车中没有选中的商品")
		}

		var goodsIds []int32
		goodsId2Num := map[int32]int32{}
		for _, item := range cartItems {
			goodsIds = append(goodsIds, item.GoodsId)
			goodsId2Num[item.GoodsId] = item.Num
		}

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

		// 3.扣减库存
		decreaseStocksReq := &pb_stock.DecreseStocksRequest{}
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

		// 4.创建订单
		order = model.Order{
			No:          genOrderNo(),
			PayType:     req.Order.PayType,
			Status:      pb.OrderStatus_ORDER_STATUS_PAYING,
			UserId:      req.UserId,
			Amount:      orderAmount,
			UserName:    req.Order.UserName,
			UserMobile:  req.Order.UserMobile,
			UserAddress: req.Order.UserAddress,
		}
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

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.CreateOrderFromCartResponse{Order: NewPbOrder(&order, orderGoods)}, nil
}

// 订单详情
func (s *OrderService) GetOrderDetail(ctx context.Context, req *pb.GetOrderDetailRequest) (*pb.GetOrderDetailResponse, error) {
	var (
		order      model.Order
		orderGoods []*model.OrderGoods
	)

	res := global.Mysql.Where("id = ?", req.OrderId).First(&order)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "订单(id=%d)不存在", req.OrderId)
	}

	res = global.Mysql.Where("order_id = ?", req.OrderId).Find(&orderGoods)
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "订单商品(order_id=%d)不存在", req.OrderId)
	}

	return &pb.GetOrderDetailResponse{Order: NewPbOrder(&order, orderGoods)}, nil
}

// 修改订单状态
func (s *OrderService) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponce, error) {
	var query string
	var args []interface{}

	if req.OrderId != 0 {
		query = "id = ?"
		args = append(args, req.OrderId)
	} else {
		query = "no = ?"
		args = append(args, req.OrderNo)
	}

	res := global.Mysql.
		Model(&model.Order{}).Where(query, args...).
		Update("status", int32(req.Status))
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.Internal, "订单(id=%d)不存在", req.OrderId)
	}

	return &pb.UpdateOrderStatusResponce{}, nil
}
