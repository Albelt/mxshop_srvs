package main

import (
	pb "albelt.top/mxshop_protos/albelt/stock_srv/go"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop_srvs/stock_srv/config"
	"sync"
)

const (
	configFilePath = "./config-dev.yaml"
)

func InitConfig() *config.Config {
	v := viper.New()
	v.SetConfigFile(configFilePath)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		panic(err)
	}

	return &cfg
}

const (
	Concurrency = 100
)

func main() {
	cfg := InitConfig()

	//连接grpc服务器
	serverAddr := fmt.Sprintf("127.0.0.1:%d", cfg.Server.Port)
	conn, err := grpc.Dial(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	cli := pb.NewStockServiceClient(conn)

	wg := sync.WaitGroup{}
	for i := 0; i < Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := cli.IncreaseStocks(context.Background(), &pb.IncreaseStocksRequest{
				OrderSn: "123",
				Stocks: []*pb.Stock{
					{
						GoodsId: 1,
						Number:  1,
					},
				},
			})
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("ok.")
			}
		}()
	}

	wg.Wait()
	fmt.Println("done.")
}
