# 第十天
当你查询一个模型的同时需要查询它的关联模型时，如果不使用预加载，就可能会陷入 N+1 查询问题。这种问题会导致额外的查询次数，严重影响性能。
在你的场景中，如果只查询一级类别（例如 Level: 1），然后在迭代每个一级类别时再查询其关联的子类别，就会产生 N+1 查询问题。
每个主查询结果需要额外的查询来获取关联数据，这样的操作会显著增加数据库负载和响应时间。

## Save() 和 Create()的区别

Save() 和 Create() 是用于向数据库插入数据的方法，但它们有一些关键的区别。让我们来详细解释一下它们的区别和适用场景：

1. global.DB.Save(&category)
   Save() 方法用于向数据库中保存（更新或插入）一个对象。它的行为取决于对象的主键（Primary Key）是否已设置：

如果对象的主键已设置（非零值），Save() 方法会执行更新操作。它会将对象的所有字段更新到数据库中已存在的记录中，以保持数据库中的数据与对象的状态一致。

如果对象的主键未设置（零值），Save() 方法会执行插入操作。它会在数据库中创建一条新的记录，并将对象的所有字段插入到该记录中。

适用场景：
当你需要根据对象的主键状态自动执行插入或更新操作时，可以使用 Save() 方法。
适用于处理更新已存在记录和插入新记录的情况，具有更强的灵活性和智能性。
2. global.DB.Create(&category)
   Create() 方法用于向数据库中插入一条新的记录。它总是将对象的所有字段插入到数据库中，并且不会考虑对象的主键状态：

不管对象的主键是否已设置，Create() 方法总是执行插入操作，将对象的数据插入为新的数据库记录。
适用场景：
当你需要明确地执行插入新记录的操作时，可以使用 Create() 方法。
通常用于向数据库中添加新数据，不会更新已存在的记录。
区别总结：
主键处理：Save() 方法根据对象的主键状态决定是更新还是插入操作，而 Create() 方法总是执行插入操作。
灵活性：Save() 方法更加智能，适用于处理不同的对象状态（更新或插入），而 Create() 方法适用于明确的插入操作。
用法场景：根据需要选择合适的方法来满足数据操作的需求，确保数据库操作的正确性和一致性。

# 第十一天

func (g GormList) Value() (driver.Value, error) {
return json.Marshal(g)
}

## 为什么上面代码不是传的指针
type Valuer interface {
Value() (Value, error)
}
这个接口只有一个方法 Value()，返回值是 driver.Value 类型，即 type Value []byte。它的目的是将 Go 类型的值转换为数据库可以存储的原始值。
值接收者：如果你实现的类型（如 GormList）是一个轻量级的值类型，并且不需要修改实例自身的状态，那么使用值接收者是合适的选择。
在这种情况下，Go 会在方法调用时对接收者进行复制，但这不会影响到性能或行为。

指针接收者：如果你的类型是一个复杂的结构体，或者你的方法需要修改接收者的状态，那么使用指针接收者可能更为合适。
指针接收者允许你在方法内部修改接收者的内容，而不是在方法中操作接收者的副本。

## 注册健康检查服务
grpc_health_v1.RegisterHealthServer(server, health.NewServer())

## 配置 Consul 客户端并注册

	// 配置 Consul 客户端
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)

	// 创建 Consul 客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// 配置服务的健康检查信息
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", global.ServerConfig.Host, *Port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	// 配置服务注册信息
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	//作用是生成一个唯一的服务 ID，使用了 github.com/satori/go.uuid 包中的 uuid.NewV4() 函数来生成 UUID（Universally Unique Identifier）。
	serviceID := fmt.Sprintf("%s", uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = global.ServerConfig.Host
	registration.Check = check

	// 注册服务到 Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}

c.Abort() 在 Gin 框架中用于中止当前请求的处理。调用 c.Abort() 后，Gin 将停止执行后续的中间件和处理器，
只执行已经调用过的中间件和处理器。这在处理某些请求时非常有用，特别是在验证和权限检查的场景中。

在 JWTAuth 函数中，我们使用 c.Abort() 来确保在令牌无效的情况下，中止请求的进一步处理，并返回相应的错误响应。
这是一个安全机制，确保未经授权的请求不会继续被处理或访问受保护的资源。

## Unix 时间戳

Unix 时间戳是指从 1970 年 1 月 1 日 00:00:00 UTC（协调世界时）起至现在的秒数。它通常以整数形式表示。
在计算机科学中，Unix 时间戳是一种常用的时间表示方式，用于跨平台和系统之间的时间标准化。
在Go语言中，通过 time.Unix() 函数可以将 Unix 时间戳转换为时间对象，如下所示
timestamp := int64(1624067115)  // 一个Unix时间戳
t := time.Unix(timestamp, 0)   // 将Unix时间戳转换为时间对象

## ISO8601 格式的时间字符串

ISO8601 是国际标准化组织（ISO）制定的日期和时间的表示方法。它提供了一种统一的格式来表示日期和时间，
使得不同国家和地区之间能够更容易地进行时间的交流和解析。ISO8601格式的时间字符串的基本格式如下 ：
YYYY-MM-DDTHH:mm:ssZ
YYYY 表示四位数的年份（例如：2023）；
MM 表示两位数的月份（01到12）；
DD 表示两位数的日期（01到31）；
T 是日期和时间的分隔符；
HH 表示两位数的小时（00到23）；
mm 表示两位数的分钟（00到59）；
ss 表示两位数的秒数（00到59）；
Z 表示时区，通常为UTC时区。

## 内网穿透

内网穿透是一种技术手段，允许将本地网络或私有网络中的服务暴露到公共网络或互联网上，从而使得外部用户可以访问这些服务。它通常用于解决以下问题：

本地开发调试：开发人员在本地开发环境中开发的应用程序需要与外部服务或客户端进行交互，但受限于本地网络的环境无法直接访问互联网，这时可以使用内网穿透将本地服务暴露出去，方便外部访问和测试。

远程办公：员工需要访问公司内部的应用程序或资源，而这些资源受限于公司的内部网络，通过内网穿透技术可以安全地将这些内部资源暴露给远程员工。

设备管理：通过内网穿透技术，可以远程管理和访问设备，例如监控摄像头、物联网设备或者家庭网络中的设备，而不需要物理上连接到同一网络中。
续断 nacos配置要和回调port一致
## 第十二天

