package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"sale_master/study_note/grpc_test/proto"
)

func main() {
	// 创建一个与 gRPC 服务器的连接
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure()) // 改为 8080 端口
	if err != nil {
		// 如果连接失败，则 panic 终止程序
		panic(err)
	}
	// 延迟关闭连接，确保程序结束前连接会被关闭
	defer conn.Close()

	// 创建一个新的 Greeter 客户端实例
	c := proto.NewGreeterClient(conn)

	// 调用 SayHello 方法，向服务器发送一个 HelloRequest 消息
	r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "llb"})
	if err != nil {
		// 如果调用失败，则 panic 终止程序
		panic(err)
	}

	// 打印服务器返回的消息
	fmt.Println(r.Message)
}
