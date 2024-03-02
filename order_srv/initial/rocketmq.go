package initial

import (
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/handler"
)

func InitRocketMqProducer() {
	// 创建producer并启动
	p, err := rocketmq.NewProducer(
		producer.WithNameServer([]string{global.Config.RocketMq.NameServer}),
		producer.WithRetry(2),
	)
	if err != nil {
		panic(err)
	}

	if err = p.Start(); err != nil {
		panic(err)
	}

	global.RocketMqProducer = p

	zap.S().Infof("InitRocketMqProducer ok.")
}

func InitRocketMqConsumer() {
	// 创建consumer并启动
	c, err := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{global.Config.RocketMq.NameServer}),
		consumer.WithGroupName(global.Config.RocketMq.ConsumerGroupName),
	)
	if err != nil {
		panic(err)
	}

	err = c.Subscribe(global.Config.RocketMq.DelayTopic, consumer.MessageSelector{}, handler.HandleOrderTimeout)
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}

	zap.S().Info("InitRocketMqConsumer ok.")
}
