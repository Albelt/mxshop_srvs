package handler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	pb "albelt.top/mxshop_protos/albelt/stock_srv/go"
	"mxshop_srvs/stock_srv/global"
	"mxshop_srvs/stock_srv/model"
)

/*
	使用MySQL的version字段(乐观锁)解决并发一致性问题
*/

// 扣减库存
func (s *StockService) DecreseStocksV2(ctx context.Context, req *pb.DecreseStocksRequest) (*pb.DecreseStocksResponse, error) {
	var (
		err error
	)

	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		for _, pbStock := range req.Stocks {
			// 查询stock(SELECT语句不加锁)
			var stock model.Stock
			res := tx.Where(&model.Stock{GoodsId: pbStock.GoodsId}).
				First(&stock)
			if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.Internal, "get stock failed:%s", res.Error.Error())
			}
			if res.RowsAffected == 0 {
				return status.Errorf(codes.NotFound, "库存(goods_id=%d)不存在", pbStock.GoodsId)
			}

			// 扣减库存
			if stock.Number < pbStock.Number {
				return status.Errorf(codes.NotFound,
					"库存(goods_id=%d)不足,需要扣减%d,剩余%d",
					pbStock.GoodsId, pbStock.Number, stock.Number)
			}

			// UPDATE语句中加入版本号字段的检查
			res = tx.Model(&stock).
				Where("version = ?", stock.Version).
				Updates(map[string]interface{}{
					"number":  stock.Number - pbStock.Number,
					"version": stock.Version + 1,
				})
			if res.Error != nil {
				return status.Errorf(codes.Internal, "update stock failed:%s", err.Error())
			}
			if res.RowsAffected == 0 {
				return status.Errorf(codes.Aborted, "update stock aborted, try again.")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.DecreseStocksResponse{}, nil
}

// 归还库存
func (s *StockService) IncreaseStocksV2(ctx context.Context, req *pb.IncreaseStocksRequest) (*pb.IncreaseStocksResponse, error) {
	var (
		err error
	)

	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		for _, pbStock := range req.Stocks {
			// 查询stock,不加锁
			var stock model.Stock
			res := tx.Where(&model.Stock{GoodsId: pbStock.GoodsId}).First(&stock)
			if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.Internal, "get stock failed:%s", res.Error.Error())
			}
			if res.RowsAffected == 0 {
				return status.Errorf(codes.NotFound, "库存(goods_id=%d)不存在", pbStock.GoodsId)
			}

			// 增加库存,UPDATE语句中检查版本号字段
			res = tx.Model(&stock).
				Where("version = ?", stock.Version).
				Updates(map[string]interface{}{
					"number":  stock.Number + pbStock.Number,
					"version": stock.Version + 1,
				})
			if res.Error != nil {
				return status.Errorf(codes.Internal, "update stock failed:%s", err.Error())
			}
			if res.RowsAffected == 0 {
				return status.Errorf(codes.Aborted, "update stock aborted, try again.")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.IncreaseStocksResponse{}, nil
}
