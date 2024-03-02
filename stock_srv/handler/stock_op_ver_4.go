package handler

import (
	pb "albelt.top/mxshop_protos/albelt/stock_srv/go"
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mxshop_srvs/stock_srv/global"
	"mxshop_srvs/stock_srv/model"
)

/*
	基于rocketmq的事务消息实现库存扣减的一致性:
	1. 扣减库存时记录历史
	2. 订单创建失败时会收到补偿库存的消息，若库存已经扣减，需要根据历史记录将扣减的库存加回来
*/

// 扣减库存,并且记录扣减库存历史
func (s *StockService) DecreseStocks(ctx context.Context, req *pb.DecreseStocksRequest) (*pb.DecreseStocksResponse, error) {
	err := global.Mysql.Transaction(func(tx *gorm.DB) error {
		// TODO: 幂等性:先查询扣减记录是否存在

		// 库存扣减历史,记录订单的sn和每个商品的扣减数量
		history := model.StockHistory{
			OrderSn: req.OrderSn,
			Status:  pb.StockHistoryStatus_STOCK_HISTORY_STATUS_REDUCED,
		}

		for _, pbStock := range req.Stocks {
			history.Details = append(history.Details, model.GoodsDetail{GoodsId: pbStock.GoodsId, Num: pbStock.Number})

			if err2 := decreaseOneStock(tx, pbStock.GoodsId, pbStock.Number); err2 != nil {
				return err2
			}
		}

		if err2 := tx.Create(&history).Error; err2 != nil {
			return err2
		}

		zap.S().Infof("订单(sn=%s)库存已扣减:%v", history.OrderSn, history.Details)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.DecreseStocksResponse{}, nil
}

// 补偿库存(分布式事务下，若创建订单回滚，则需要将扣减的库存补回来)
func RestoreStocks(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		zap.S().Infof("get msg:%s", msg.Body)

		// 解析消息中的order_sn
		var orderDetail pb.OrderDetails
		if err := json.Unmarshal(msg.Body, &orderDetail); err != nil {
			zap.S().Warnf("decode msg failed:%s", err.Error())
			continue
		}

		err := global.Mysql.Transaction(func(tx *gorm.DB) error {
			// 查询该订单的库存是否已经扣减
			var history model.StockHistory
			res := global.Mysql.
				Where("order_sn = ? AND status = ?", orderDetail.OrderSn, pb.StockHistoryStatus_STOCK_HISTORY_STATUS_REDUCED).
				First(&history)
			if res.RowsAffected == 0 {
				zap.S().Infof("订单(sn=%s)未扣减库存", orderDetail.OrderSn)
				return nil
			}

			// 已扣减则需要补偿库存
			for _, g := range history.Details {
				if err2 := increaseOneStock(tx, g.GoodsId, g.Num); err2 != nil {
					return err2
				}
			}

			// 最后将history的状态改为已归还
			history.Status = pb.StockHistoryStatus_STOCK_HISTORY_STATUS_RESTORED
			err2 := tx.Model(&model.StockHistory{}).
				Where("order_sn = ?", history.OrderSn).
				Update("status", history.Status).Error
			if err2 != nil {
				return err2
			}

			zap.S().Infof("订单(sn=%s)库存已补偿:%v", history.OrderSn, history.Details)

			return nil
		})

		if err != nil {
			zap.S().Warnf("订单(sn=%s)补偿库存失败:%s", orderDetail.OrderSn, err.Error())
			continue
		}
	}

	return consumer.ConsumeSuccess, nil
}
