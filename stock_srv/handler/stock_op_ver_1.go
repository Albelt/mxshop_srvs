package handler

import (
	pb "albelt.top/mxshop_protos/albelt/stock_srv/go"
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"mxshop_srvs/stock_srv/global"
	"mxshop_srvs/stock_srv/model"
)

/*
	使用MySQL的行锁(悲观锁)解决并发一致性问题
*/

// 扣减库存
func (s *StockService) DecreseStocksV1(ctx context.Context, req *pb.DecreseStocksRequest) (*pb.DecreseStocksResponse, error) {
	var (
		err error
	)

	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		for _, pbStock := range req.Stocks {
			// 查询stock,加行锁(goods_id上有索引,使用FOR UPDATE查询时会加行锁)
			var stock model.Stock
			res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where(&model.Stock{GoodsId: pbStock.GoodsId}).
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

			stock.Number -= pbStock.Number
			if err = tx.Model(&stock).Update("number", stock.Number).Error; err != nil {
				return status.Errorf(codes.Internal, "update stock failed:%s", res.Error.Error())
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
func (s *StockService) IncreaseStocksV1(ctx context.Context, req *pb.IncreaseStocksRequest) (*pb.IncreaseStocksResponse, error) {
	var (
		err error
	)

	err = global.Mysql.Transaction(func(tx *gorm.DB) error {
		for _, pbStock := range req.Stocks {
			// 查询stock,加行锁(goods_id上有索引)
			var stock model.Stock
			res := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Stock{GoodsId: pbStock.GoodsId}).First(&stock)
			if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return status.Errorf(codes.Internal, "get stock failed:%s", res.Error.Error())
			}
			if res.RowsAffected == 0 {
				return status.Errorf(codes.NotFound, "库存(goods_id=%d)不存在", pbStock.GoodsId)
			}

			// 增加库存
			stock.Number += pbStock.Number
			if err = tx.Model(&stock).Update("number", stock.Number).Error; err != nil {
				return status.Errorf(codes.Internal, "update stock failed:%s", res.Error.Error())
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.IncreaseStocksResponse{}, nil
}
