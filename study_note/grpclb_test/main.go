package main

import (
	"context"
	"fmt"
	"log"
	"sale_master/study_note/grpclb_test/proto"

	_ "github.com/mbobakov/grpc-consul-resolver" // Import the Consul resolver package

	"google.golang.org/grpc"
)

func main() {
	// 连接到 Consul 注册中心的 gRPC 服务
	conn, err := grpc.Dial(
		"consul://192.168.128.128:8500/user-srv?wait=14s&tag=srv", // Consul 地址及服务标签
		grpc.WithInsecure(), // 使用不安全连接，仅用于示例目的，生产环境应使用安全连接或证书
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 设置默认的负载均衡策略为轮询
	)
	if err != nil {
		log.Fatal(err) // 连接失败时输出错误并退出程序
	}
	defer conn.Close() // 确保在函数返回前关闭 gRPC 连接

	// 创建 gRPC 客户端
	for i := 0; i < 10; i++ {
		userSrvClient := proto.NewUserClient(conn) // 创建用户服务的 gRPC 客户端实例
		rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
			Pn:    1,
			PSize: 2,
		})
		if err != nil {
			panic(err) // 处理调用服务方法时的错误
		}
		// 输出服务端返回的数据
		for index, data := range rsp.Data {
			fmt.Println(index, data)
		}
	}
}
