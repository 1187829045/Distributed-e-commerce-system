package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"

	"sale_master/study_note/grpc_interpretor/proto"
)

func main() {
	// 定义一个拦截器函数，用于拦截 gRPC 请求，并在请求前后执行额外的逻辑
	interceptor := func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		// 记录请求开始时间
		start := time.Now()

		// 调用实际的 gRPC 请求
		err := invoker(ctx, method, req, reply, cc, opts...)

		// 打印请求耗时
		fmt.Printf("耗时：%s\n", time.Since(start))

		// 返回请求结果和可能的错误
		return err
	}

	// 定义一个 gRPC 拨号选项的切片，用于配置客户端
	var opts []grpc.DialOption

	// 配置客户端使用不安全连接（不启用 SSL/TLS）
	opts = append(opts, grpc.WithInsecure())

	// 定义重试策略的选项
	retryOpts := []grpc_retry.CallOption{
		// 最大重试次数为 3 次
		grpc_retry.WithMax(3),
		// 每次重试的超时时间为 1 秒
		grpc_retry.WithPerRetryTimeout(1 * time.Second),
		// 当遇到以下状态码时进行重试：Unknown、DeadlineExceeded、Unavailable
		grpc_retry.WithCodes(codes.Unknown, codes.DeadlineExceeded, codes.Unavailable),
	}

	// 将自定义的拦截器添加到拨号选项中，拦截器用于监控每个请求的耗时
	opts = append(opts, grpc.WithUnaryInterceptor(interceptor))

	// 将重试拦截器添加到拨号选项中，配置客户端在特定条件下进行重试
	opts = append(opts, grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)))

	// 建立与 gRPC 服务器的连接，服务器地址为 127.0.0.1:50051
	conn, err := grpc.Dial("127.0.0.1:50051", opts...)
	if err != nil {
		// 如果连接失败，程序将输出错误信息并终止
		panic(err)
	}
	// 在 main 函数退出前关闭连接
	defer conn.Close()

	// 使用连接创建一个 Greeter 客户端，Greeter 是通过 proto 文件生成的 gRPC 客户端接口
	c := proto.NewGreeterClient(conn)

	// 调用 SayHello 方法，向 gRPC 服务器发送请求，并传递包含 "bobby" 的 HelloRequest 消息
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "bobby"})
	if err != nil {
		// 如果调用失败，程序将输出错误信息并终止
		panic(err)
	}

	// 打印服务器返回的消息
	fmt.Println(r.Message)
}
