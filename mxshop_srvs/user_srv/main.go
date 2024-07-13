package main

import (
	"flag"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/inner/uuid"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/hashicorp/consul/api"
	"shop_srvs/user_srv/global"
	"shop_srvs/user_srv/handler"
	"shop_srvs/user_srv/initialize"
	"shop_srvs/user_srv/proto"
	"shop_srvs/user_srv/utils"
)

func main() {
	// 定义命令行参数 `ip` 和 `port`
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	Port := flag.Int("port", 0, "端口号")

	// 初始化各种配置和资源
	initialize.InitLogger()           // 初始化日志
	initialize.InitConfig()           // 初始化配置
	initialize.InitDB()               // 初始化数据库
	zap.S().Info(global.ServerConfig) // 打印配置信息

	// 解析命令行参数
	flag.Parse()
	zap.S().Info("ip: ", *IP) // 打印 IP 地址
	if *Port == 0 {
		*Port, _ = utils.GetFreePort() // 获取一个可用端口
	}

	zap.S().Info("port: ", *Port) // 打印端口号

	// 创建 gRPC 服务器实例
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})         // 注册用户服务
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port)) // 开始监听指定的 IP 和端口
	if err != nil {
		panic("failed to listen:" + err.Error()) // 监听失败时触发 panic
	}

	// 注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)

	client, err := api.NewClient(cfg) // 创建 Consul 客户端
	if err != nil {
		panic(err) // 创建客户端失败时触发 panic
	}

	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("192.168.128.128:%d", *Port), // gRPC 检查地址
		Timeout:                        "5s",                                     // 超时时间
		Interval:                       "5s",                                     // 检查间隔
		DeregisterCriticalServiceAfter: "15s",                                    // 服务注销时间
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name // 服务名称
	//UUID v4 是随机生成的，保证了它的唯一性和不可预测性。
	serviceID := fmt.Sprintf("%s", uuid.NewV4) // 服务 ID
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = []string{"llb", "bobby", "user", "srv"} // 服务标签
	registration.Address = "192.168.128.128"
	registration.Check = check

	// 服务注册到 Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err) // 注册失败时触发 panic
	}

	// 启动 gRPC 服务
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error()) // 服务启动失败时触发 panic
		}
	}()

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	// 注销服务
	if err = client.Agent().ServiceDeregister(serviceID); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
