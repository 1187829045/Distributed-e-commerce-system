package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"shop_srvs/userop_srv/handler"
	"shop_srvs/userop_srv/utils/register/consul"
	"syscall"

	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"shop_srvs/userop_srv/global"
	"shop_srvs/userop_srv/initialize"
	"shop_srvs/userop_srv/proto"
	"shop_srvs/userop_srv/utils"
)

func main() {
	// 通过命令行标志解析ip地址，默认值为 "0.0.0.0"
	IP := flag.String("ip", "0.0.0.0", "ip地址")
	// 通过命令行标志解析端口号，默认值为 0
	Port := flag.Int("port", 0, "端口号")

	// 初始化日志记录器
	initialize.InitLogger()
	// 初始化配置文件加载
	initialize.InitConfig()
	// 初始化数据库连接
	initialize.InitDB()
	// 输出加载的全局配置
	zap.S().Info(global.ServerConfig)

	// 解析命令行标志，将命令行参数赋值给变量 IP 和 Port
	flag.Parse()
	// 输出解析后的 IP 地址
	zap.S().Info("ip: ", *IP)
	// 如果端口号为 0（即未指定端口号），则获取一个可用的随机端口
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	// 输出最终使用的端口号
	zap.S().Info("port: ", *Port)

	// 创建一个新的 gRPC 服务器实例
	server := grpc.NewServer()
	// 注册 gRPC 服务，具体实现由 handler.UserOpServer 提供
	proto.RegisterAddressServer(server, &handler.UserOpServer{})
	proto.RegisterMessageServer(server, &handler.UserOpServer{})
	proto.RegisterUserFavServer(server, &handler.UserOpServer{})

	// 启动网络监听，监听指定的 IP 和端口
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	if err != nil {
		// 如果监听失败，则抛出错误
		panic("failed to listen:" + err.Error())
	}

	// 注册 gRPC 服务的健康检查服务
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 启动 gRPC 服务器
	go func() {
		// 如果服务器启动失败，则抛出错误
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	// 创建一个 Consul 客户端用于服务注册，使用全局配置中的 Consul 地址和端口
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	// 生成一个唯一的服务 ID
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	// 将服务注册到 Consul
	err = register_client.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		// 如果注册失败，记录错误并抛出异常
		zap.S().Panic("服务注册失败:", err.Error())
	}
	// 服务成功启动，输出调试信息
	zap.S().Debugf("启动服务器, 端口： %d", *Port)

	// 创建一个用于接收操作系统信号的通道
	quit := make(chan os.Signal)
	// 将系统中断信号（SIGINT）和终止信号（SIGTERM）通知给该通道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞主线程，直到收到终止信号
	<-quit
	// 注销服务
	if err = register_client.DeRegister(serviceId); err != nil {
		// 如果注销失败，记录错误信息
		zap.S().Info("注销失败:", err.Error())
	} else {
		// 注销成功，记录成功信息
		zap.S().Info("注销成功:")
	}
}
