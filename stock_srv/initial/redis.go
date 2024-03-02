package initial

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"mxshop_srvs/stock_srv/global"
	"time"
)

func InitRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:         global.Config.Redis.Addr,
		Password:     global.Config.Redis.Password,
		DB:           int(global.Config.Redis.Db),
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}

	global.RdsCli = rdb
	zap.S().Info("InitRedis ok.")
}
