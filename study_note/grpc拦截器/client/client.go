package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"

	"sale_master/study_note/grpc拦截器/proto"
)

// 定义一个Server结构体，该结构体实现了proto.GreeterServer接口

type Server struct{}

// 实现SayHello方法，模拟处理请求时延迟2秒钟

func (s *Server) SayHello(ctx context.Context, request *proto.HelloRequest) (*proto.HelloReply, error) {
	time.Sleep(2 * time.Second) // 模拟处理时间
	return &proto.HelloReply{
		Message: "hello " + request.Name,
	}, nil
}

func main() {
	// 定义一个拦截器，拦截所有的Unary RPC调用
	interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		fmt.Println("接收到了一个新的请求")     // 在接收到请求时打印日志
		res, err := handler(ctx, req) // 调用实际的RPC方法处理请求
		fmt.Println("请求已经完成")         // 在请求处理完成后打印日志
		return res, err               // 返回处理结果
	}

	// 将拦截器作为选项传递给gRPC服务器
	opt := grpc.UnaryInterceptor(interceptor)
	g := grpc.NewServer(opt)

	// 注册Greeter服务到gRPC服务器
	proto.RegisterGreeterServer(g, &Server{})

	// 监听TCP连接
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		panic("failed to listen:" + err.Error()) // 如果监听失败，输出错误并退出
	}

	// 启动gRPC服务器，开始接受和处理连接
	err = g.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error()) // 如果服务器启动失败，输出错误并退出
	}
}
