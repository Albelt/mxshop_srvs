package handler

import (
	pb "albelt.top/mxshop_protos/albelt/stock_srv/go"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"mxshop_srvs/stock_srv/data/rds_lock"
	"mxshop_srvs/stock_srv/global"
	"mxshop_srvs/stock_srv/model"
)

/*
	使用redis分布式锁-解决并发一致性问题
*/

const (
	increOrDecreStockMutexFmt    = "incre_or_decre_stock_mutex_%d"
	increOrDecreStockMutexExpire = 10
)

// 扣减库存
func (s *StockService) DecreseStocksV3(ctx context.Context, req *pb.DecreseStocksRequest) (*pb.DecreseStocksResponse, error) {
	err := global.Mysql.Transaction(func(tx *gorm.DB) error {
		for _, pbStock := range req.Stocks {
			if err2 := decreaseOneStock(tx, pbStock.GoodsId, pbStock.Number); err2 != nil {
				return err2
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.DecreseStocksResponse{}, nil
}

func decreaseOneStock(tx *gorm.DB, goodsId, decreseNum int32) error {
	var err error

	// 针对goods_id加锁,同一商品只允许一个客户端操作
	mutex := rds_lock.NewRedisLock(global.RdsCli, fmt.Sprintf(increOrDecreStockMutexFmt, goodsId), increOrDecreStockMutexExpire)
	if err = mutex.Lock(); err != nil {
		return err
	}
	defer func(mutex *rds_lock.RedisLock) {
		err := mutex.Unlock()
		if err != nil {
			zap.S().Warnf("Unlock失败：%s", err.Error())
		}
	}(mutex)

	// 查询stock
	var stock model.Stock
	res := tx.Where(&model.Stock{GoodsId: goodsId}).
		First(&stock)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return status.Errorf(codes.Internal, "get stock failed:%s", res.Error.Error())
	}
	if res.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "库存(goods_id=%d)不存在", goodsId)
	}

	// 扣减库存
	if stock.Number < decreseNum {
		return status.Errorf(codes.NotFound,
			"库存(goods_id=%d)不足,需要扣减%d,剩余%d",
			goodsId, decreseNum, stock.Number)
	}

	stock.Number -= decreseNum
	if err = tx.Model(&stock).Update("number", stock.Number).Error; err != nil {
		return status.Errorf(codes.Internal, "update stock failed:%s", res.Error.Error())
	}

	return nil
}

// 归还库存
func (s *StockService) IncreaseStocks(ctx context.Context, req *pb.IncreaseStocksRequest) (*pb.IncreaseStocksResponse, error) {
	err := global.Mysql.Transaction(func(tx *gorm.DB) error {
		for _, pbStock := range req.Stocks {
			if err2 := increaseOneStock(tx, pbStock.GoodsId, pbStock.Number); err2 != nil {
				return err2
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &pb.IncreaseStocksResponse{}, nil
}

func increaseOneStock(tx *gorm.DB, goodsId, increseNum int32) error {
	var err error

	// 针对goods_id加锁,同一商品只允许一个客户端操作
	mutex := rds_lock.NewRedisLock(global.RdsCli, fmt.Sprintf(increOrDecreStockMutexFmt, goodsId), increOrDecreStockMutexExpire)
	if err = mutex.Lock(); err != nil {
		return err
	}
	defer func(mutex *rds_lock.RedisLock) {
		err := mutex.Unlock()
		if err != nil {
			zap.S().Warnf("Unlock失败:%s", err.Error())
		}
	}(mutex)

	// 查询stock
	var stock model.Stock
	res := tx.Where(&model.Stock{GoodsId: goodsId}).First(&stock)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return status.Errorf(codes.Internal, "get stock failed:%s", res.Error.Error())
	}
	if res.RowsAffected == 0 {
		return status.Errorf(codes.NotFound, "库存(goods_id=%d)不存在", goodsId)
	}

	// 增加库存
	stock.Number += increseNum
	if err = tx.Model(&stock).Update("number", stock.Number).Error; err != nil {
		return status.Errorf(codes.Internal, "update stock failed:%s", res.Error.Error())
	}

	return nil
}
