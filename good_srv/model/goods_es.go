package model

const (
	goodsIndexName = "goods"
	goodsMapping   = `
{
    "mappings": {
        "properties": {
            "brand_id": {
                "type": "integer"
            },
            "category_id": {
                "type": "integer"
            },
            "click_num": {
                "type": "integer"
            },
            "fav_num": {
                "type": "integer"
            },
            "goods_brief": {
                "type": "text",
                "analyzer": "ik_smart"
            },
            "id": {
                "type": "integer"
            },
            "is_hot": {
                "type": "boolean"
            },
            "is_new": {
                "type": "boolean"
            },
            "market_price": {
                "type": "float"
            },
            "name": {
                "type": "text",
                "analyzer": "ik_max_word"
            },
            "on_sale": {
                "type": "boolean"
            },
            "ship_free": {
                "type": "boolean"
            },
            "shop_price": {
                "type": "float"
            },
            "sold_num": {
                "type": "integer"
            }
        }
    },
    "settings": {
        "index": {
            "number_of_shards": "1",
            "number_of_replicas": "1"
        }
    }
}`
)

type EsGoods struct {
	ID          int32   `json:"id"`
	CategoryID  int32   `json:"category_id"`
	BrandID     int32   `json:"brand_id"`
	OnSale      bool    `json:"on_sale"`
	ShipFree    bool    `json:"ship_free"`
	IsNew       bool    `json:"is_new"`
	IsHot       bool    `json:"is_hot"`
	Name        string  `json:"name"`
	ClickNum    int32   `json:"click_num"`
	SoldNum     int32   `json:"sold_num"`
	FavNum      int32   `json:"fav_num"`
	MarketPrice float32 `json:"market_price"`
	ShopPrice   float32 `json:"shop_price"`
	GoodsBrief  string  `json:"goods_brief"`
}

func (g EsGoods) IndexName() string {
	return goodsIndexName
}

func (g EsGoods) Mapping() string {
	return goodsMapping

}
