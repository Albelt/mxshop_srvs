package model

import (
	"context"
	"gorm.io/gorm"
	"mxshop_srvs/good_srv/global"
	"strconv"
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primaryKey"`
	CreatedAt time.Time `gorm:"column:create_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

// 商品分类
type Category struct {
	BaseModel
	Name             string      `gorm:"type:varchar(32);not null" json:"name"`
	ParentCategoryID *int32      `gorm:"" json:"parent_category_id"`
	ParentCategory   *Category   `json:"-"`
	SubCategory      []*Category `gorm:"foreignKey:ParentCategoryID;references:ID" json:"sub_category"`
	Level            int32       `gorm:"type:int;not null;default:1" json:"level"`
	IsTab            bool        `gorm:"default:false;not null" json:"is_tab"`
}

// 品牌
type Brand struct {
	BaseModel
	Name string `gorm:"type:varchar(32);not null"`
	Logo string `gorm:"type:varchar(256);default:'';not null"`
}

// 分类-品牌关联
type CategoryBrandRelation struct {
	BaseModel
	CategoryID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Category   Category

	BrandsID int32 `gorm:"type:int;index:idx_category_brand,unique"`
	Brands   Brand
}

// 首页轮播图
type Banner struct {
	BaseModel
	Image string `gorm:"type:varchar(256);not null"`
	Url   string `gorm:"type:varchar(256);not null"`
	Index int32  `gorm:"type:int;default:1;not null"`
}

// 商品
type Goods struct {
	BaseModel

	CategoryID int32 `gorm:"type:int;not null"`
	Category   Category
	BrandID    int32 `gorm:"type:int;not null"`
	Brand      Brand

	OnSale   bool `gorm:"default:false;not null"`
	ShipFree bool `gorm:"default:false;not null"`
	IsNew    bool `gorm:"default:false;not null"`
	IsHot    bool `gorm:"default:false;not null"`

	Name            string  `gorm:"type:varchar(50);not null"`
	GoodsSn         string  `gorm:"type:varchar(50);not null"`
	ClickNum        int32   `gorm:"type:int;default:0;not null"`
	SoldNum         int32   `gorm:"type:int;default:0;not null"`
	FavNum          int32   `gorm:"type:int;default:0;not null"`
	MarketPrice     float32 `gorm:"type:float;not null"`
	ShopPrice       float32 `gorm:"type:float;not null"`
	GoodsBrief      string  `gorm:"type:varchar(256);not null"`
	GoodsFrontImage string  `gorm:"type:varchar(128);not null"`
}

func (g *Goods) AfterCreate(tx *gorm.DB) error {
	// 同步goods数据到es
	esG := EsGoods{
		ID:          g.ID,
		CategoryID:  g.CategoryID,
		BrandID:     g.BrandID,
		OnSale:      g.OnSale,
		ShipFree:    g.ShipFree,
		IsNew:       g.IsNew,
		IsHot:       g.IsHot,
		Name:        g.Name,
		ClickNum:    g.ClickNum,
		SoldNum:     g.SoldNum,
		FavNum:      g.FavNum,
		MarketPrice: g.MarketPrice,
		ShopPrice:   g.ShopPrice,
		GoodsBrief:  g.GoodsBrief,
	}

	// 使用DB的主键作为es的id，避免数据重复
	_, err := global.EsCli.Index().Index(esG.IndexName()).
		BodyJson(esG).Id(strconv.FormatInt(int64(g.ID), 10)).Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (g *Goods) AfterUpdate(tx *gorm.DB) error {
	return g.AfterCreate(tx)
}

func (g *Goods) AfterDelete(tx *gorm.DB) error {
	_, err := global.EsCli.Delete().
		Index(EsGoods{}.IndexName()).
		Id(strconv.FormatInt(int64(g.ID), 10)).
		Do(context.Background())

	return err
}
