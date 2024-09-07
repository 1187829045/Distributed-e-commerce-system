package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "sale_master/study_note/grpc_token_auth_test/proto"
)

type Server struct {
	pb.UnimplementedGreeterServer // 嵌入未实现的GreeterServer，确保我们遵循接口的所有方法
}

// SayHello 是 GreeterServer 接口中定义的方法
// 实现了 gRPC 定义的 SayHello 方法，接收一个 HelloRequest 并返回一个 HelloReply。
func (s *Server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	// 定义一个拦截器函数，用于在处理 gRPC 请求之前执行一些逻辑
	interceptor := func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		// 打印日志，表示接收到一个新的请求
		fmt.Println("接受到一新的请求")

		// 从上下文中提取元数据（metadata），通常包含请求头中的信息
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			// 如果没有从上下文中提取到元数据，返回未认证错误
			return resp, status.Error(codes.Unauthenticated, "无token认证信息")
		}

		// 定义两个变量来存储从元数据中提取的 appid 和 appkey
		var (
			appid  string
			appkey string
		)

		// 从元数据中提取 "appid" 值，并赋值给变量 appid
		if val, ok := md["appid"]; ok {
			appid = val[0]
		}

		// 从元数据中提取 "appkey" 值，并赋值给变量 appkey
		if val, ok := md["appkey"]; ok {
			appkey = val[0]
		}

		// 验证提取到的 appid 和 appkey 是否匹配预期值
		if appid != "101010" || appkey != "I am key" {
			// 如果 appid 或 appkey 不匹配，返回未认证错误
			return resp, status.Error(codes.Unauthenticated, "无token认证信息")
		}

		// 认证通过后，调用实际的处理函数 handler 进行请求处理
		res, err := handler(ctx, req)

		// 打印日志，表示请求处理完成
		fmt.Println("请求已经完成")

		// 返回处理结果和可能的错误
		return res, err
	}

	// 创建一个拦截器选项，将拦截器应用于 gRPC 服务器
	opt := grpc.UnaryInterceptor(interceptor)

	// 创建一个新的 gRPC 服务器实例，并应用拦截器选项
	g := grpc.NewServer(opt)

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
