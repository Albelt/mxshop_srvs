package main

import (
	pb "albelt.top/mxshop_protos/albelt/order_srv/go"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"mxshop_srvs/common/middleware"
	"mxshop_srvs/common/middleware/grpc_tracing"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/handler"
	"mxshop_srvs/order_srv/initial"
	"mxshop_srvs/order_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 初始化日志
	initial.InitLogger()

	// 初始化配置
	initial.InitConfig()

	// 初始化mysql, redis
	initial.InitMysql()
	initial.InitRedis()

	// 初始化其他的grpc服务
	initial.InitGoodsSrvCli()
	initial.InitStockSrvCli()

	// 初始化rocketmq的producer和consumer
	initial.InitRocketMqProducer()
	initial.InitRocketMqConsumer()

	// 初始化tracer
	tracer, closer := initial.InitTracer()
	defer closer.Close()
	global.Tracer = tracer

	// 获取端口
	var port int
	var err error
	if global.Config.Server.Debug {
		// debug模式使用固定端口
		port = int(global.Config.Server.Port)
	} else {
		// 获取空闲端口
		port, err = utils.GetFreePort()
		if err != nil {
			panic(err)
		}
	}

	// 创建grpc服务器并初始化
	s := grpc.NewServer(
		// 日志,tracing中间件
		grpc.ChainUnaryInterceptor(
			middleware.ServerLogging(),
			grpc_tracing.OpenTracingServerInterceptor(tracer),
		),
	)
	pb.RegisterOrderServiceServer(s, &handler.OrderService{})
	reflection.Register(s)

	// 注册服务到注册中心
	initial.RegisterConsul(s, global.Config.Server.ExternalIp, port)

	// 启动服务器并配置优雅退出
	addr := fmt.Sprintf("%s:%d", global.Config.Server.Ip, port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	go func() {
		zap.S().Infof("Start grpc server on:%s", addr)
		s.Serve(listen)
	}()

	// 主协程监听退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 从注册中心注销
	initial.DeregisterConsul()
	// grpc服务器优雅退出
	s.GracefulStop()

	zap.S().Info("Stop grpc server")
}
