package initialize

import (
	"encoding/json"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"shop_srvs/user_srv/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
func InitConfig() {
	// 从环境变量中获取调试模式信息，决定使用哪个配置文件
	debug := GetEnvInfo("SHOP_DEBUG")
	configFilePrefix := "config"
	// 根据调试模式选择配置文件名
	configFileName := fmt.Sprintf("shop_srvs/user_srv/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("shop_srvs/user_srv/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	// 设置配置文件路径和名称
	v.SetConfigFile(configFileName)
	// 读取配置文件，如果失败则触发 panic
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// 将配置文件中的内容反序列化到 global.NacosConfig 结构体中
	if err := v.Unmarshal(&global.NacosConfig); err != nil {
		panic(err)
	}
	// 使用 zap 记录配置信息
	zap.S().Infof("配置信息: %v", global.NacosConfig)

	// 从 Nacos 中读取配置信息
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host, // Nacos 服务器的 IP 地址
			Port:   global.NacosConfig.Port, // Nacos 服务器的端口号
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.NacosConfig.Namespace, // Nacos 命名空间 ID
		TimeoutMs:           5000,                         // 请求超时时间
		NotLoadCacheAtStart: true,                         // 启动时不加载本地缓存
		LogDir:              "tmp/nacos/log",              // 日志目录
		CacheDir:            "tmp/nacos/cache",            // 缓存目录
		LogLevel:            "debug",                      // 日志级别
	}

	// 创建 Nacos 配置客户端
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil {
		panic(err)
	}

	// 从 Nacos 获取配置信息
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId, // Nacos 配置的 DataId
		Group:  global.NacosConfig.Group,  // Nacos 配置的 Group
	})

	if err != nil {
		panic(err)
	}
	// 将 Nacos 中获取的配置内容反序列化到 global.ServerConfig 结构体中
	err = json.Unmarshal([]byte(content), &global.ServerConfig)
	if err != nil {
		zap.S().Fatalf("读取 Nacos 配置失败： %s", err.Error())
	}
	// 打印 global.ServerConfig 的内容
	fmt.Println(&global.ServerConfig)
}
