package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"sale_master/study_note/grpc_interpretor/proto"
)

// 定义一个空的结构体 Server，用于实现 gRPC 服务接口
type Server struct{}

// 实现 SayHello 方法，这是 gRPC 服务中定义的一个 RPC 方法
func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	// 模拟业务处理延迟，休眠 2 秒
	time.Sleep(2 * time.Second)

	// 返回一个 HelloReply 响应，Message 字段包含 "hello " 和请求中的 Name
	return &proto.HelloReply{
		Message: "hello " + request.Name,
	}, nil
}

func main() {
	// 定义一个拦截器函数，用于在处理每个 gRPC 请求时执行额外的逻辑
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 打印日志，表示接收到一个新的请求
		fmt.Println("接收到了一个新的请求")

		// 调用实际的处理函数 handler，并传递上下文和请求
		res, err := handler(ctx, req)

		// 打印日志，表示请求已经处理完成
		fmt.Println("请求已经完成")

		// 返回处理结果和可能的错误
		return res, err
	}

	// 创建一个 gRPC 服务器选项，将拦截器添加到服务器中
	opt := grpc.UnaryInterceptor(interceptor)

	// 创建一个新的 gRPC 服务器实例，并应用拦截器选项
	g := grpc.NewServer(opt)

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
