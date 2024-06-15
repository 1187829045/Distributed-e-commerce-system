package main

import (
	"fmt"
	"github.com/spf13/viper" // 导入 viper 包，用于处理配置文件
)

// ServerConfig 定义了服务端配置结构体，使用 mapstructure 标签指定字段与配置文件中的键名对应关系
type ServerConfig struct {
	ServiceName string `mapstructure:"name"` // 对应配置文件中的 "name" 键
	Port        int    `mapstructure:"port"` // 对应配置文件中的 "port" 键
}

func main() {
	v := viper.New() // 创建一个新的 viper 实例
	// 设置配置文件的路径
	v.SetConfigFile("config.yaml")
	// 读取并解析配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(err) // 如果读取配置文件失败，直接抛出异常
	}
	serverConfig := ServerConfig{} // 创建一个 ServerConfig 结构体实例用于存储配置信息
	// 将 viper 实例中的配置信息解析到 ServerConfig 结构体中
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err) // 如果解析配置信息失败，直接抛出异常
	}
	// 打印输出解析后的 ServerConfig 结构体
	fmt.Println(serverConfig)
	// 通过 viper 获取特定键的值并打印输出
	fmt.Printf("%v", v.Get("name"))
}
