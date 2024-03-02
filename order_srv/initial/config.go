package initial

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop_srvs/order_srv/config"
	"mxshop_srvs/order_srv/global"
)

const (
	configFilePath = "config/config-dev.yaml"
)

func InitConfig() {
	v := viper.New()
	v.SetConfigFile(configFilePath)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	global.Config = &cfg
	zap.S().Infof("InitConfig from %s, content:\n%s", configFilePath, cfg.String())
}
