package model

import (
	pb "albelt.top/mxshop_protos/albelt/order_srv/go"
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"column:create_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

// 购物车项目
type ShoppingCartItem struct {
	BaseModel
	UserId  int32 `gorm:"column:user_id;type:int;index;comment:用户id"`
	GoodsId int32 `gorm:"column:goods_id;type:int;index;comment:商品id"`
	Num     int32 `gorm:"column:num;type:int;comment:数量"`
	Checked bool  `gorm:"column:checked;type:bool;comment:是否选中"`
}

// 订单
type Order struct {
	BaseModel
	No      string         `gorm:"column:no;type:string;comment:生成的订单号"`
	PayType pb.PayType     `gorm:"column:pay_type;type:tinyint;comment:支付类型"`
	Status  pb.OrderStatus `gorm:"column:status;type:tinyint;comment:订单状态"`
	UserId  int32          `gorm:"column:user_id;type:int;index;comment:用户id"`
	TradeNo string         `gorm:"column:trade_no;type:string;comment:交易号"`
	Amount  float32        `gorm:"column:amount;type:float;comment:订单金额"`
	PayTime *time.Time     `gorm:"column:pay_time;type:datetime;comment:支付时间"`

	// 用户信息冗余
	UserName    string `gorm:"column:user_name;type:string;comment:用户姓名"`
	UserMobile  string `gorm:"column:user_mobile;type:string;comment:用户手机号"`
	UserAddress string `gorm:"column:user_address;type:string;comment:用户收件地址"`
}

// 订单-商品
type OrderGoods struct {
	BaseModel
	OrderId  int32 `gorm:"column:order_id;type:int;index;comment:订单id"`
	GoodsId  int32 `gorm:"column:goods_id;type:int;index;comment:商品id"`
	GoodsNum int32 `gorm:"column:goods_num;type:int;comment:商品购买数量"`

	// 商品信息冗余
	GoodsName  string  `gorm:"column:goods_name;type:string;comment:商品名称"`
	GoodsPrice float32 `gorm:"column:goods_price;type:float;comment:商品价格"`
}
