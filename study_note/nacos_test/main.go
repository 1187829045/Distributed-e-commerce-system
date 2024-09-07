package main

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"sale_master/study_note/nacos_test/config"
	"time"
)

func main() {
	// 创建一个ServerConfig切片，用于存储Nacos服务器的配置
	sc := []constant.ServerConfig{
		{
			// 指定Nacos服务器的IP地址
			IpAddr: "192.168.128.128",
			// 指定Nacos服务器的端口号
			Port: 8848,
		},
	}

	// 创建ClientConfig结构体，用于配置客户端信息
	cc := constant.ClientConfig{
		// 指定NamespaceId，用于区分不同的命名空间
		NamespaceId: "eaa72e66-8bf8-4ed3-a9cb-08067dd75e77", // 如果需要支持多namespace，可以创建多个客户端实例，它们有不同的NamespaceId
		// 指定请求超时时间，单位为毫秒
		TimeoutMs: 5000,
		// 指定是否在启动时加载缓存
		NotLoadCacheAtStart: true,
		// 指定日志的存储目录
		LogDir: "tmp/nacos/log",
		// 指定缓存的存储目录
		CacheDir: "tmp/nacos/cache",
		// 指定日志级别
		LogLevel: "debug",
	}

	// 使用配置创建一个Nacos的配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		// 传入服务器配置
		"serverConfigs": sc,
		// 传入客户端配置
		"clientConfig": cc,
	})
	// 如果创建客户端失败，终止程序并输出错误信息
	if err != nil {
		panic(err)
	}

	// 获取指定配置文件的内容
	content, err := configClient.GetConfig(vo.ConfigParam{
		// 指定配置文件的DataId
		DataId: "user-web.yaml",
		// 指定配置文件所在的组
		Group: "dev",
	})
	// 如果获取配置失败，终止程序并输出错误信息
	if err != nil {
		panic(err)
	}

	// 解析获取到的内容并将其转换为指定的结构体
	serverConfig := config.ServerConfig{}
	// 将JSON格式的字符串转换成Go结构体
	json.Unmarshal([]byte(content), &serverConfig)
	// 打印解析后的结构体
	fmt.Println(serverConfig)

	// 开始监听指定的配置文件，当配置文件发生变化时触发回调函数
	err = configClient.ListenConfig(vo.ConfigParam{
		// 指定监听的配置文件的DataId
		DataId: "user-web.json",
		// 指定监听的配置文件所在的组
		Group: "dev",
		// 当配置文件变化时，执行此回调函数
		OnChange: func(namespace, group, dataId, data string) {
			// 打印提示信息，表明配置文件发生了变化
			fmt.Println("配置文件变化")
			// 打印配置文件的group、dataId和最新的数据内容
			fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
		},
	})
	// 如果监听失败，终止程序并输出错误信息
	if err != nil {
		panic(err)
	}

	// 让程序休眠3000秒，以保持监听状态
	time.Sleep(3000 * time.Second)
}
