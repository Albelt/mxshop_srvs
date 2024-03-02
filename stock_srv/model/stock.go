package model

import (
	pb "albelt.top/mxshop_protos/albelt/stock_srv/go"
	"database/sql/driver"
	"encoding/json"
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

// 库存表
type Stock struct {
	BaseModel
	GoodsId int32 `gorm:"column:goods_id;type:int;not null;index;comment:关联goods的id"`
	Number  int32 `gorm:"column:number;type:int;not null;default:0;comment:库存数量"`
	Version int32 `gorm:"column:version;type:int;default:0;comment:乐观锁版本号"`
}

// 库存扣减历史(用于扣减/归还库存时做幂等性检查)
type StockHistory struct {
	OrderSn string                `gorm:"column:order_sn;type:string;comment:order_sn;index:unique"`
	Status  pb.StockHistoryStatus `gorm:"column:status;type:int;comment:1已扣减,2已归还"`
	Details GoodsDetailList       `gorm:"column:details;type:string;comment:json格式,各个商品的扣减数量"`
}

type GoodsDetail struct {
	GoodsId int32 `json:"goods_id,omitempty"`
	Num     int32 `json:"num,omitempty"`
}

type GoodsDetailList []GoodsDetail

func (l *GoodsDetailList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &l)
}

func (l GoodsDetailList) Value() (driver.Value, error) {
	return json.Marshal(l)
}
