package handler

import (
	pb "albelt.top/mxshop_protos/albelt/stock_srv/go"
	"mxshop_srvs/stock_srv/model"
)

type StockService struct {
	pb.UnimplementedStockServiceServer
}

func NewStockFromPb(s *pb.Stock) model.Stock {
	if s == nil {
		return model.Stock{}
	}
	return model.Stock{
		GoodsId: s.GoodsId,
		Number:  s.Number,
	}
}

func NewPbStock(s *model.Stock) *pb.Stock {
	if s == nil {
		return nil
	}
	return &pb.Stock{
		GoodsId: s.GoodsId,
		Number:  s.Number,
	}
}
