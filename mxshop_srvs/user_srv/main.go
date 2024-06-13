package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"sale_master/mxshop_srvs/user_srv/handler"
	"sale_master/mxshop_srvs/user_srv/proto"
)

func main() {
	IP := flag.String("ip", "0.0.0.0", "ip 地址")
	Port := flag.Int("port", 50051, "端口号")
	flag.Parse()
	fmt.Println("ip:", *IP)
	fmt.Println("port:", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		panic(err)
	}
	server.Serve(lis)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}

}
