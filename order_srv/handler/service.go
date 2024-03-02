package handler

import (
	pb "albelt.top/mxshop_protos/albelt/order_srv/go"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"mxshop_srvs/order_srv/model"
	"time"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
}

func NewPbCartItem(i *model.ShoppingCartItem) *pb.ShoppingCartItem {
	if i == nil {
		return nil
	}

	return &pb.ShoppingCartItem{
		Id:      i.ID,
		UserId:  i.UserId,
		GoodsId: i.GoodsId,
		Num:     i.Num,
		Checked: i.Checked,
	}
}

func NewPbOrder(order *model.Order, orderGoods []*model.OrderGoods) *pb.Order {
	if order == nil {
		return nil
	}

	ret := &pb.Order{
		Id:          order.ID,
		No:          order.No,
		PayType:     order.PayType,
		Status:      order.Status,
		TradeNo:     order.TradeNo,
		Amount:      order.Amount,
		UserId:      order.UserId,
		UserName:    order.UserName,
		UserMobile:  order.UserMobile,
		UserAddress: order.UserAddress,
	}

	if order.PayTime != nil {
		ret.PayTime = timestamppb.New(*order.PayTime)
	}

	if orderGoods != nil {
		for _, g := range orderGoods {
			ret.Goodss = append(ret.Goodss, &pb.OrderGoods{
				Id:         g.ID,
				OrderId:    g.OrderId,
				GoodsId:    g.GoodsId,
				GoodsNum:   g.GoodsNum,
				GoodsName:  g.GoodsName,
				GoodsPrice: g.GoodsPrice,
			})
		}
	}

	return ret
}

func genOrderNo() string {
	now := time.Now()
	return fmt.Sprintf("no:%.4d%.2d%.2d%.2d%.2d%.2d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Unix()%1000000)
}
