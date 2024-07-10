package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify" // 导入 fsnotify 库，用于监控配置文件变化
	"github.com/spf13/viper"       // 导入 viper 库，用于处理配置文件
	"time"
)

// MysqlConfig 定义 MySQL 配置结构体
type MysqlConfig struct {
	Host string `mapstructure:"host"` // MySQL 主机地址
	Port int    `mapstructure:"port"` // MySQL 端口号
}

// ServerConfig 定义服务器配置结构体
type ServerConfig struct {
	ServiceName string      `mapstructure:"name"`  // 服务名称
	MysqlInfo   MysqlConfig `mapstructure:"mysql"` // MySQL 配置信息
}

// GetEnvInfo 根据环境变量获取配置信息
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()      // 自动加载环境变量
	return viper.GetBool(env) // 获取指定环境变量的布尔值
}

func main() {
	debug := GetEnvInfo("SHOP_DEBUG")                                              // 获取环境变量 shop_DEBUG 的布尔值
	configFilePrefix := "config"                                                   // 配置文件前缀名
	configFileName := fmt.Sprintf("viper_test/ch02/%s-pro.yaml", configFilePrefix) // 生产环境配置文件路径

	if debug {
		configFileName = fmt.Sprintf("viper_test/ch02/%s-debug.yaml", configFilePrefix) // 调试环境配置文件路径
	}

	v := viper.New()                         // 创建一个新的 viper 实例
	v.SetConfigFile(configFileName)          // 设置配置文件路径
	if err := v.ReadInConfig(); err != nil { // 读取配置文件
		panic(err)
	}

	serverConfig := ServerConfig{}                     // 创建一个 ServerConfig 结构体实例，用于存储配置信息
	if err := v.Unmarshal(&serverConfig); err != nil { // 解析配置文件到 ServerConfig 结构体
		panic(err)
	}
	fmt.Println(serverConfig)       // 打印解析后的配置信息
	fmt.Printf("%V", v.Get("name")) // 打印配置文件中的 name 字段值

	// 监听配置文件变化并自动重新加载
	v.WatchConfig()                           // 监听配置文件变化
	v.OnConfigChange(func(e fsnotify.Event) { // 配置文件变化时的回调函数
		fmt.Println("config file changed:", e.Name)
		_ = v.ReadInConfig()           // 重新读取配置文件
		_ = v.Unmarshal(&serverConfig) // 更新 ServerConfig 结构体
		fmt.Println(serverConfig)      // 打印更新后的配置信息
	})

	time.Sleep(time.Second * 300) // 程序睡眠 300 秒，等待配置文件变化
}
