package global

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"mxshop_srvs/user_srv/config"
)

var (
	Config *config.Config
	Mysql  *gorm.DB
	RdsCli *redis.Client
)
