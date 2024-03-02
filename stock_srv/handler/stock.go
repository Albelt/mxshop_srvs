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

// 设置库存
func (s *StockService) SetStock(ctx context.Context, req *pb.SetStockRequest) (*pb.SetStockResponse, error) {
	var (
		err   error
		stock model.Stock
	)

	// 根据goods_id查询库存
	res := global.Mysql.Where(&model.Stock{GoodsId: req.Stock.GoodsId}).First(&stock)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.Internal, "get stock failed:%s", res.Error.Error())
	}

	if res.RowsAffected == 0 {
		// 不存在则创建
		stock = NewStockFromPb(req.Stock)
		if err = global.Mysql.Create(&stock).Error; err != nil {
			return nil, status.Errorf(codes.Internal, "create stock failed:%s", res.Error.Error())
		}
	} else {
		// 存在则更新
		stock.Number = req.Stock.Number
		if err = global.Mysql.Save(&stock).Error; err != nil {
			return nil, status.Errorf(codes.Internal, "update stock failed:%s", res.Error.Error())
		}
	}

	return &pb.SetStockResponse{Stock: NewPbStock(&stock)}, nil
}

// 获取货物的库存信息
func (s *StockService) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {
	var (
		stock model.Stock
	)

	res := global.Mysql.Where(&model.Stock{GoodsId: req.GoodsId}).First(&stock)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.Internal, "get stock failed:%s", res.Error.Error())
	}

	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "库存(goods_id=%d)不存在", req.GoodsId)
	}

	return &pb.GetStockResponse{Stock: NewPbStock(&stock)}, nil
}
