package config

// UserSrvConfig 用于存储用户服务的配置信息
type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"` // 用户服务的主机地址，映射配置文件中的 "host" 字段
	Port int    `mapstructure:"port" json:"port"` // 用户服务的端口号，映射配置文件中的 "port" 字段
	Name string `mapstructure:"name" json:"name"` // 用户服务的名称，映射配置文件中的 "name" 字段
}

// JWTConfig 用于存储JWT（JSON Web Token）的配置信息
type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"` // JWT的签名密钥，映射配置文件中的 "key" 字段
}

// AliSmsConfig 用于存储阿里云短信服务的配置信息
type AliSmsConfig struct {
	ApiKey     string `mapstructure:"key" json:"key"`         // 阿里云短信服务的API密钥，映射配置文件中的 "key" 字段
	ApiSecrect string `mapstructure:"secrect" json:"secrect"` // 阿里云短信服务的API密钥，映射配置文件中的 "secrect" 字段
}

// ConsulConfig 用于存储Consul的配置信息
type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"` // Consul服务的主机地址，映射配置文件中的 "host" 字段
	Port int    `mapstructure:"port" json:"port"` // Consul服务的端口号，映射配置文件中的 "port" 字段
}

// RedisConfig 用于存储Redis的配置信息
type RedisConfig struct {
	Host   string `mapstructure:"host" json:"host"`     // Redis服务的主机地址，映射配置文件中的 "host" 字段
	Port   int    `mapstructure:"port" json:"port"`     // Redis服务的端口号，映射配置文件中的 "port" 字段
	Expire int    `mapstructure:"expire" json:"expire"` // Redis缓存的过期时间，映射配置文件中的 "expire" 字段
}

// ServerConfig 用于存储服务器的总体配置信息
type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name"`         // 服务器名称，映射配置文件中的 "name" 字段
	Port        int           `mapstructure:"port" json:"port"`         // 服务器端口号，映射配置文件中的 "port" 字段
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"` // 嵌套的用户服务配置信息，映射配置文件中的 "user_srv" 字段
	JWTInfo     JWTConfig     `mapstructure:"jwt" json:"jwt"`           // 嵌套的JWT配置信息，映射配置文件中的 "jwt" 字段
	AliSmsInfo  AliSmsConfig  `mapstructure:"sms" json:"sms"`           // 嵌套的阿里云短信服务配置信息，映射配置文件中的 "sms" 字段
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`       // 嵌套的Redis配置信息，映射配置文件中的 "redis" 字段
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`     // 嵌套的Consul配置信息，映射配置文件中的 "consul" 字段
}
