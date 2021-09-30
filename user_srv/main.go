package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mbSrvs/user_srv/global"
	"mbSrvs/user_srv/handler"
	"mbSrvs/user_srv/initialize"
	"mbSrvs/user_srv/proto"
	"net"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")

	flag.Parse()

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()

	zap.S().Info("ip:%s", *IP)
	zap.S().Info("port:%s", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))

	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()

	consulInfo := global.ServerConfig.ConsulInfo

	cfg.Address = fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", consulInfo.Host, consulInfo.Port),
		Timeout:                        "50s",
		Interval:                       "50s",
		DeregisterCriticalServiceAfter: "100s",
	}

	//生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = global.ServerConfig.Name
	registration.Port = *Port
	registration.Tags = []string{"cowboy", "user", "srv"}
	registration.Address = consulInfo.Host
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)

	if err != nil {
		zap.S().Panicf(err.Error())
	}

	err = server.Serve(listener)

	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}

}
