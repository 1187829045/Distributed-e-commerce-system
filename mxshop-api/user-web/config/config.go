package config

// 一些配置信息结构体

// UserSrvConfig 用于存储用户服务的配置信息
type UserSrvConfig struct {
	Host string `mapstructure:"host"` // 用户服务的主机地址
	Port int    `mapstructure:"port"` // 用户服务的端口号
	Name string `mapstructure:"name" json:"name"`
}

// JWTConfig 用于存储 JWT 的配置信息
type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"` // JWT 的签名密钥
}

// AliSmsConfig 用于存储阿里短信服务的配置信息
type AliSmsConfig struct {
	ApiKey     string `mapstructure:"key" json:"key"`         // 阿里短信服务的 API 密钥
	ApiSecrect string `mapstructure:"secrect" json:"secrect"` // 阿里短信服务的 API 秘钥
}

// ConsulConfig 用于存储 Consul 的配置信息
type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"` // Consul 的主机地址
	Port int    `mapstructure:"port" json:"port"` // Consul 的端口号
}

// RedisConfig 用于存储 Redis 的配置信息
type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`     // Redis 的主机地址
	Port   int    `mapstructure:"port" json:"port"`     // Redis 的端口号
	Expire int    `mapstructure:"expire" json:"expire"` // Redis 的过期时间
}

// ServerConfig 用于存储服务器的配置信息
type ServerConfig struct {
	Name        string        `mapstructure:"name"`               // 服务器名称
	Host        string        `mapstructure:"host" json:"host"`   // 服务器主机地址
	Tags        []string      `mapstructure:"tags" json:"tags"`   // 服务器标签
	Port        int           `mapstructure:"port" json:"port"`   // 服务器端口号
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv"`           // 用户服务配置信息
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt"`     // JWT 配置信息
	AliSmsInfo  AliSmsConfig  `mapstructure:"sms" json:"sms"`     // 阿里短信服务配置信息
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"` // Redis 配置信息
	// ConsulInfo 被注释掉的 Consul 配置信息
	// ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
}

// NacosConfig 用于存储 Nacos 的配置信息
type NacosConfig struct {
	Host      string `mapstructure:"host"`      // Nacos 的主机地址
	Port      uint64 `mapstructure:"port"`      // Nacos 的端口号
	Namespace string `mapstructure:"namespace"` // Nacos 的命名空间
	User      string `mapstructure:"user"`      // Nacos 的用户名
	Password  string `mapstructure:"password"`  // Nacos 的密码
	DataId    string `mapstructure:"dataid"`    // Nacos 的数据 ID
	Group     string `mapstructure:"group"`     // Nacos 的分组
}
