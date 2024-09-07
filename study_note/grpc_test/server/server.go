package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "sale_master/study_note/grpc_test/proto"
)

// Server 是你的 gRPC 服务的实现
// Server struct 实现了生成的 gRPC 服务器接口。
type Server struct {
	pb.UnimplementedGreeterServer // 嵌入未实现的GreeterServer，确保我们遵循接口的所有方法
}

// SayHello 是 GreeterServer 接口中定义的方法
// 实现了 gRPC 定义的 SayHello 方法，接收一个 HelloRequest 并返回一个 HelloReply。
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	// 创建一个新的 gRPC 服务器实例
	g := grpc.NewServer()

	// 注册 Greeter 服务到 gRPC 服务器
	// 将 Server 的实例传递给生成的 RegisterGreeterServer 方法
	pb.RegisterGreeterServer(g, &Server{})

	// 监听所有网络接口上的 8080 端口
	lis, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		// 如果监听失败，记录错误并退出程序
		log.Fatalf("failed to listen: %v", err)
	}

	// 启动 gRPC 服务器以监听传入的连接
	// 如果服务器运行过程中出现错误，记录错误并退出程序
	if err := g.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
