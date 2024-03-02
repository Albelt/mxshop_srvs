package initial

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"go.uber.org/zap"
	"mxshop_srvs/stock_srv/global"
	"mxshop_srvs/stock_srv/handler"
)

func InitRocketMqConsumer() {
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{global.Config.RocketMq.NameServer}),
		consumer.WithGroupName(global.Config.RocketMq.ConsumerGroupName),
	)
	if err != nil {
		panic(err)
	}

	err = c.Subscribe(global.Config.RocketMq.Topic, consumer.MessageSelector{}, handler.RestoreStocks)
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}

	zap.S().Info("InitRocketMqConsumer ok.")
}
