package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/hashicorp/consul/api"
	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/handler"
	"shop_srvs/goods_srv/initialize"
	"shop_srvs/goods_srv/proto"
	"shop_srvs/goods_srv/utils"
)

func main() {
	// 定义命令行参数
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")
	// 初始化应用程序的各种组件（日志、配置、数据库、Elasticsearch等）
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitEs() // 初始化ES连接
	// 输出全局配置信息到日志
	zap.S().Info(global.ServerConfig)
	// 解析命令行参数
	flag.Parse()
	// 打印命令行参数中指定的IP和Port信息
	zap.S().Info("ip: ", *IP)
	if *Port == 0 {
		// 如果未指定端口号，则调用工具函数获取一个空闲端口
		*Port, _ = utils.GetFreePort()
	}
	zap.S().Info("port: ", *Port)
	// 创建一个 gRPC 服务器实例
	server := grpc.NewServer()
	// 注册商品服务到 gRPC 服务器
	proto.RegisterGoodsServer(server, &handler.GoodsServer{})

	// 监听指定 IP 和端口号的 TCP 连接
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))

	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	// 注册健康检查服务
	//rpc_health_v1.RegisterHealthServer 是一个函数，用于将健康检查服务注册到 gRPC 服务器 (server) 上。
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 配置 Consul 客户端
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	// 创建 Consul 客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 配置服务的健康检查信息
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	// 配置服务注册信息
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	//作用是生成一个唯一的服务 ID，使用了 github.com/satori/go.uuid 包中的 uuid.NewV4() 函数来生成 UUID（Universally Unique Identifier）。
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host
	registration.Check = check

	// 注册服务到 Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

	// 启动 gRPC 服务器
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	// 接收终止信号
	//优雅的退出
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 注销服务
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
