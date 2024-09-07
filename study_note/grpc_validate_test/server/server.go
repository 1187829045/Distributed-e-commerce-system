package main

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"

	"google.golang.org/grpc"

	"sale_master/study_note/grpc_validate_test/proto"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, request *proto.Person) (*proto.Person,
	error) {
	return &proto.Person{
		Id: 32,
	}, nil
}

type Validator interface {
	Validate() error
}

func main() {
	// 定义一个gRPC的单向拦截器，用于在处理请求前执行额外的逻辑
	var interceptor grpc.UnaryServerInterceptor

	// 拦截器函数，接收请求并在请求处理前执行验证逻辑
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		// 检查请求类型是否实现了 Validator 接口，用于验证请求数据的有效性
		if r, ok := req.(Validator); ok {
			// 如果实现了 Validator 接口，调用其 Validate 方法进行验证
			if err := r.Validate(); err != nil {
				// 如果验证失败，返回 InvalidArgument 错误和错误信息
				return nil, status.Error(codes.InvalidArgument, err.Error())
			}
		}

		// 如果请求通过验证，调用实际的请求处理函数 handler 并返回结果
		return handler(ctx, req)
	}

	// 创建一个包含拦截器选项的gRPC服务器选项切片
	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	// 使用定义的选项创建一个新的gRPC服务器
	g := grpc.NewServer(opts...)

	// 注册 Greeter 服务到服务器实例 g 中，将请求路由到服务实现
	proto.RegisterGreeterServer(g, &Server{})

	// 监听TCP连接，指定监听的IP地址和端口号
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		// 如果监听失败，使用 panic 终止程序并打印错误信息
		panic("failed to listen:" + err.Error())
	}

	// 启动gRPC服务器，开始接收并处理客户端请求
	err = g.Serve(lis)
	if err != nil {
		// 如果服务器启动失败，使用 panic 终止程序并打印错误信息
		panic("failed to start grpc:" + err.Error())
	}
}
