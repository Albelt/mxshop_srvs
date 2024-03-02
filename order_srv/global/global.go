package global

import (
	pb_goods "albelt.top/mxshop_protos/albelt/good_srv/go"
	pb_stock "albelt.top/mxshop_protos/albelt/stock_srv/go"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"
	"mxshop_srvs/order_srv/config"
)

var (
	Config           *config.Config
	Mysql            *gorm.DB
	RdsCli           *redis.Client
	GoodsSrvCli      pb_goods.GoodsServiceClient
	StockSrvCli      pb_stock.StockServiceClient
	RocketMqProducer rocketmq.Producer
	Tracer           opentracing.Tracer
)
