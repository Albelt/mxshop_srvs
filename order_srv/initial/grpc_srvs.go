package initial

import (
	pb_goods "albelt.top/mxshop_protos/albelt/good_srv/go"
	pb_stock "albelt.top/mxshop_protos/albelt/stock_srv/go"
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop_srvs/order_srv/global"
)

const (
	// %s:consul地址  %s:服务名称
	//grpcConsulUrlFormat = "consul://ubuntu-learn:8500/user-srv?wait=14s"
	grpcConsulUrlFormat = "consul://%s/%s?wait=14s"
)

func InitGoodsSrvCli() {
	grpcConsulUrl := fmt.Sprintf(grpcConsulUrlFormat, global.Config.Consul.Addr, global.Config.GrpcSrvs.GoodsSrvName)

	conn, err := grpc.Dial(
		grpcConsulUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic(err)
	}

	global.GoodsSrvCli = pb_goods.NewGoodsServiceClient(conn)
	zap.S().Info("InitGoodsSrvCli ok.")
}

func InitStockSrvCli() {
	grpcConsulUrl := fmt.Sprintf(grpcConsulUrlFormat, global.Config.Consul.Addr, global.Config.GrpcSrvs.StockSrvName)

	conn, err := grpc.Dial(
		grpcConsulUrl,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic(err)
	}

	global.StockSrvCli = pb_stock.NewStockServiceClient(conn)
	zap.S().Info("InitStockSrvCli ok.")
}
