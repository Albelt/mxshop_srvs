package handler

import (
	pb "albelt.top/mxshop_protos/albelt/order_srv/go"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/model"
)

// 获取用户的购物车
func (s *OrderService) GetUserCart(ctx context.Context, req *pb.GetUserCartRequest) (*pb.GetUserCartResponse, error) {
	var (
		err       error
		cartItems []*model.ShoppingCartItem
		resp      pb.GetUserCartResponse
	)

	err = global.Mysql.Where("user_id = ?", req.UserId).Find(&cartItems).Error
	if err != nil {
		return nil, err
	}

	resp.TotalItem = int32(len(cartItems))
	for _, i := range cartItems {
		resp.Items = append(resp.Items, NewPbCartItem(i))
	}

	return &resp, nil
}

// 添加商品到购物车
func (s *OrderService) AddItemToCart(ctx context.Context, req *pb.AddItemToCartRequest) (*pb.AddItemToCartResponse, error) {
	var (
		err  error
		item model.ShoppingCartItem
	)

	res := global.Mysql.Where("user_id = ? AND goods_id = ?", req.Item.UserId, req.Item.GoodsId).First(&item)

	if res.RowsAffected == 0 {
		// 不存在则创建
		item.UserId = req.Item.UserId
		item.GoodsId = req.Item.GoodsId
		item.Num = req.Item.Num
		item.Checked = req.Item.Checked
		err = global.Mysql.Create(&item).Error
	} else {
		// 存在则合并
		item.Num += req.Item.Num
		err = global.Mysql.Model(&item).Update("num", item.Num).Error
	}

	if err != nil {
		return nil, err
	}

	return &pb.AddItemToCartResponse{}, nil
}

// 从购物车删除商品
func (s *OrderService) RemoveItemFromCart(ctx context.Context, req *pb.RemoveItemFromCartRequest) (*pb.RemoveItemFromCartResponse, error) {
	userId := req.UserId
	goodsId := req.GoodsId

	res := global.Mysql.Where("user_id = ? AND goods_id = ?", userId, goodsId).
		Delete(&model.ShoppingCartItem{})
	if res.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录(user_id=%d,goods_id=%d)不存在", userId, goodsId)
	}

	return &pb.RemoveItemFromCartResponse{}, nil
}
