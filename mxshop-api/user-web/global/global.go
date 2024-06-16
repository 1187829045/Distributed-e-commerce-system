package global

import (
	ut "github.com/go-playground/universal-translator"
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/proto"
)

// 全局变量包，用于存储全局配置和对象，方便在整个应用中共享和访问。

var (
	// Trans 是一个全局的翻译器实例，用于处理多语言翻译。
	Trans ut.Translator

	// ServerConfig 是一个全局的服务器配置实例，指向 config 包中的 ServerConfig 结构体。
	// 通过初始化为一个新的 ServerConfig 对象，可以全局共享这个配置。
	ServerConfig *config.ServerConfig = &config.ServerConfig{}

	// NacosConfig 是一个全局的 Nacos 配置实例，指向 config 包中的 NacosConfig 结构体。
	// 通过初始化为一个新的 NacosConfig 对象，可以全局共享这个配置。
	//十周
	NacosConfig *config.NacosConfig = &config.NacosConfig{}

	// UserSrvClient 是一个全局的用户服务客户端实例，用于 gRPC 调用。
	UserSrvClient proto.UserClient
)
