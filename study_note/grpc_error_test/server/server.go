package main

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"

	"google.golang.org/grpc"

	"sale_master/study_note/grpc_error_test/proto"
)

// 定义一个空的结构体 Server，用于实现 gRPC 服务接口

type Server struct{}

// 实现 SayHello 方法，这是 gRPC 服务中定义的一个 RPC 方法

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	// 返回一个错误状态，表示没有找到记录，并将请求中的 Name 信息包含在错误消息中
	return nil, status.Errorf(codes.NotFound, "记录未找到：%s", request.Name)

	// 如果上面的返回语句被注释掉，下面这段代码将会返回一个正常的应答
	// return &proto.HelloReply{
	// 	Message: "hello " + request.Name,
	// }, nil
}

func main() {
	// 创建一个新的 gRPC 服务器实例
	g := grpc.NewServer()

	// 将实现了 gRPC 服务接口的 Server 注册到 gRPC 服务器中
	proto.RegisterGreeterServer(g, &Server{})

	// 监听 TCP 连接，地址为 0.0.0.0:50051
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		// 如果监听失败，打印错误信息并终止程序
		panic("failed to listen:" + err.Error())
	}

	// 启动 gRPC 服务器，开始监听客户端的请求
	err = g.Serve(lis)
	if err != nil {
		// 如果服务器启动失败，打印错误信息并终止程序
		panic("failed to start grpc:" + err.Error())
	}
}
