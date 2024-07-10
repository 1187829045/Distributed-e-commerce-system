package main

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

// Register 将服务注册到 Consul
func Register(address string, port int, name string, tags []string, id string) error {
	// 创建默认的 Consul 配置
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.128.128:8500" // 指定 Consul 服务器的地址

	// 创建一个新的 Consul 客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		// 如果创建客户端失败，抛出错误
		return err
	}

	// 生成对应的健康检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           "http://192.168.128.128:8021/health", // 健康检查的 HTTP 地址
		Timeout:                        "5s",                                 // 健康检查超时时间
		Interval:                       "5s",                                 // 健康检查间隔时间
		DeregisterCriticalServiceAfter: "10s",                                // 在服务不健康后取消注册的时间
	}

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = name       // 服务名称
	registration.ID = id           // 服务 ID
	registration.Port = port       // 服务端口
	registration.Tags = tags       // 服务标签
	registration.Address = address // 服务地址
	registration.Check = check     // 服务健康检查

	// 将服务注册到 Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		// 如果注册服务失败，抛出错误
		return err
	}
	return nil
}

// AllServices 获取所有注册的服务并打印
func AllServices() {
	// 创建默认的 Consul 配置
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.128.128:8500" // 指定 Consul 服务器的地址

	// 创建一个新的 Consul 客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		// 如果创建客户端失败，抛出错误
		panic(err)
	}

	// 获取所有注册的服务
	data, err := client.Agent().Services()
	if err != nil {
		// 如果获取服务失败，抛出错误
		panic(err)
	}

	// 打印所有服务的键
	for key := range data {
		fmt.Println(key)
	}
}

// FilterService 根据过滤条件获取服务并打印
func FilterService() {
	// 创建默认的 Consul 配置
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.128.128:8500" // 指定 Consul 服务器的地址

	// 创建一个新的 Consul 客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		// 如果创建客户端失败，抛出错误
		panic(err)
	}

	// 根据过滤条件获取服务
	data, err := client.Agent().ServicesWithFilter(`Service == "user-web"`)
	if err != nil {
		// 如果获取服务失败，抛出错误
		panic(err)
	}

	// 打印所有符合条件的服务的键
	for key := range data {
		fmt.Println(key)
	}
}

func main() {
	// 注册服务
	_ = Register("192.168.128.128", 8021, "user-web", []string{"shop", "llb"}, "user-web")

	// 获取并打印所有服务
	AllServices()

	// 根据过滤条件获取并打印服务
	FilterService()

	// 打印格式化的服务名称
	fmt.Println(fmt.Sprintf(`Service == "%s"`, "user-srv"))
}
