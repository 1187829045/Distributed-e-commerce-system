package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"sale_master/study_note/grpc_token_auth_test/proto"
)

type customCredentials struct{}

func (c customCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "101010",
		"appkey": "I am key",
	}, nil
}
func (c customCredentials) RequireTransportSecurity() bool {
	return false
}
func main() {
	grpc.WithPerRPCCredentials(customCredentials{})
	// 定义 gRPC 客户端连接选项
	var opts []grpc.DialOption
	// 设置连接为不安全模式（不使用 TLS）
	opts = append(opts, grpc.WithInsecure())
	// 将拦截器添加到连接选项中
	opts = append(opts, grpc.WithPerRPCCredentials(customCredentials{}))

	// 创建一个与 gRPC 服务器的连接
	conn, err := grpc.Dial("127.0.0.1:8080", opts...) // 改为 8080 端口
	if err != nil {
		// 如果连接失败，则 panic 终止程序
		panic(err)
	}
	// 延迟关闭连接，确保程序结束前连接会被关闭
	defer conn.Close()

	// 创建一个新的 Greeter 客户端实例
	c := proto.NewGreeterClient(conn)

	// 调用 SayHello 方法，向服务器发送一个 HelloRequest 消息
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "bobby"})
	if err != nil {
		// 如果调用失败，则 panic 终止程序
		panic(err)
	}

	// 打印服务器返回的消息
	fmt.Println(r.Message)
}
