package main

import (
	"fmt"
	"github.com/spf13/viper"
)

// ServerConfig 定义了服务端配置结构体，使用 mapstructure 标签指定字段与配置文件中的键名对应关系
type ServerConfig struct {
	ServiceName string `mapstructure:"name"`
	Port        int    `mapstructure:"port"`
}

func main() {
	v := viper.New() // 创建一个新的 viper 实例
	// 设置配置文件的路径
	v.SetConfigFile("config.yaml")
	// 读取并解析配置文件
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	serverConfig := ServerConfig{}
	// 将 viper 实例中的配置信息解析到 ServerConfig 结构体中
	if err := v.Unmarshal(&serverConfig); err != nil {
		panic(err)
	}
	// 打印输出解析后的 ServerConfig 结构体
	fmt.Println(serverConfig)
	fmt.Printf("%v", v.Get("name"))

}
