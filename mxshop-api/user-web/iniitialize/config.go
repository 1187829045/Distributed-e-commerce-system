package iniitialize

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper" // 导入用于处理配置文件的 viper 包

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"go.uber.org/zap"          // 导入 zap 日志库
	"shop-api/user-web/global" // 导入全局变量包
)

// GetEnvInfo 从环境变量中获取配置信息
func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

// InitConfig 初始化配置文件
func InitConfig() {
	debug := GetEnvInfo("shop_DEBUG")                                       // 检查是否处于调试模式
	configFilePrefix := "config"                                            // 配置文件前缀
	configFileName := fmt.Sprintf("user-web/%s-pro.yaml", configFilePrefix) // 生产环境配置文件路径
	if debug {
		configFileName = fmt.Sprintf("user-web/%s-debug.yaml", configFilePrefix) // 调试环境配置文件路径
	}

	v := viper.New()
	v.SetConfigFile(configFileName)          // 设置配置文件路径
	if err := v.ReadInConfig(); err != nil { // 读取配置文件
		panic(err) // 如果读取配置文件失败，则抛出错误
	}
	// 这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(global.NacosConfig); err != nil { // 将配置文件解析到全局变量 NacosConfig
		panic(err)
	}
	zap.S().Infof("配置信息: &v", global.NacosConfig) // 打印配置信息
	// 从 nacos 中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host, // Nacos 服务器地址
			Port:   global.NacosConfig.Port, // Nacos 服务器端口
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // 如果需要支持多namespace，我们可以创建多个 client，它们有不同的 NamespaceId
		TimeoutMs:           5000,                         // 客户端超时时间
		NotLoadCacheAtStart: true,                         // 启动时不加载缓存
		LogDir:              "tmp/nacos/log",              // 日志目录
		CacheDir:            "tmp/nacos/cache",            // 缓存目录
		LogLevel:            "debug",                      // 日志级别
	}

	// 创建配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc, // 服务器配置
		"clientConfig":  cc, // 客户端配置
	})
	if err != nil {
		panic(err) // 如果创建配置客户端失败，则抛出错误
	}

	// 从 nacos 获取配置信息
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId, // 数据 ID
		Group:  global.NacosConfig.Group,  // 配置组
	})
	if err != nil {
		panic(err) // 如果获取配置信息失败，则抛出错误
	}

	// 想要将一个 JSON 字符串转换成 struct，需要设置这个 struct 的 tag
	err = json.Unmarshal([]byte(content), &global.ServerConfig) // 解析 JSON 配置信息到 ServerConfig 结构体
	if err != nil {
		zap.S().Fatalf("读取 nacos 配置失败： %s", err.Error()) // 如果解析失败，记录错误日志并终止程序
	}
	fmt.Println(&global.ServerConfig) // 打印配置信息
}
