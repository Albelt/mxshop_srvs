package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"gorm.io/gorm"
	"mxshop_srvs/good_srv/config"
)

var (
	Config *config.Config
	Mysql  *gorm.DB
	RdsCli *redis.Client
	EsCli  *elastic.Client
)
