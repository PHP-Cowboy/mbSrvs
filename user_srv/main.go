package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"mbSrvs/user_srv/handler"
	"mbSrvs/user_srv/proto"
	"net"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 50051, "端口号")

	flag.Parse()
	fmt.Println("ip:%s", *IP)
	fmt.Println("port:%s", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))

	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	err = server.Serve(listener)

	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}

}
