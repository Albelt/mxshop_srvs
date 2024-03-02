package initial

import (
	"fmt"
	consul_api "github.com/hashicorp/consul/api"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mxshop_srvs/order_srv/global"
)

var (
	serviceId string
	consulCli *consul_api.Client
)

// 将grpc服务注册到注册中心
func RegisterConsul(grpcSrv *grpc.Server, externalIp string, port int) {
	var (
		err error
	)

	// 注册grpc的健康检查服务
	grpc_health_v1.RegisterHealthServer(grpcSrv, health.NewServer())

	// 创建consul客户端
	consulCfg := consul_api.DefaultConfig()
	consulCfg.Address = global.Config.Consul.Addr
	consulCli, err = consul_api.NewClient(consulCfg)
	if err != nil {
		panic(err)
	}

	// 在consul中注册当前的grpc服务
	serviceId = uuid.NewV4().String()
	srv := &consul_api.AgentServiceRegistration{
		ID:      serviceId,
		Name:    global.Config.Server.Name,
		Tags:    global.Config.Server.Tags,
		Port:    port,
		Address: externalIp,
		Check: &consul_api.AgentServiceCheck{
			Interval:                       global.Config.Consul.HealthCheck.Interval,
			Timeout:                        global.Config.Consul.HealthCheck.Timeout,
			GRPC:                           fmt.Sprintf("%s:%d", externalIp, port),
			DeregisterCriticalServiceAfter: global.Config.Consul.HealthCheck.Deregister,
		},
	}

	if err = consulCli.Agent().ServiceRegister(srv); err != nil {
		panic(err)
	}

	zap.S().Info("RegisterConsul ok.")
}

// 将grpc服务从consul注销
func DeregisterConsul() {
	if err := consulCli.Agent().ServiceDeregister(serviceId); err != nil {
		zap.S().Errorf("DeregisterConsul failed:%s", err.Error())
	} else {
		zap.S().Info("DeregisterConsul ok.")
	}
}
