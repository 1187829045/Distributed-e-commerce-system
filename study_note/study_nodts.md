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
行锁粒度比较细，只会满足符合条件的数据，在没有索引的情况下，行锁升级为表锁。如果没有满足条件的结果，不会锁表（有索引的情况）。

### redsync源码解读

#### setnx的作用
将获取和设置值变成原子操作

#### 如果断电了宕机了怎么办：

造成了死锁，
a.设置过期时间
b.如果设置了过期时间，那么过期时间到了我的业务逻辑没有执行怎么办？
i.在过期之前再刷新一下
ii.需要自己去启动写成完成延时的操作
i.延时的接口肯会带来负面影响-如果其中某一个服务bung住了，2s就能执行完，但是你hung住那么你就会一直去申请延时长锁，导致别人永远无法获取锁。

#### 分布式锁需要解决的问题-lua脚本做

a.互斥性
b.死锁
c.安全性
i.所只能被持有该锁的用户删除，不能被其他人删掉
1.当时设置的value值是多少，只有当时的设置的人才知道,g
2.再删除的时候取出value值对比一下自己拥有的value看看 一不一样

#### 即使这样实现分布式锁还是有问题 redlock
常见情况搭建一个redis集群，提高可用性。
建设有五台redis集群，分别处于不同的服务器上，有一个住redis集群，其他都是从redis集群。如果不同的实例向不同的实例拿锁。能达到同一把锁的目的。
写数据，redis集群会自动同步。当同步的时候master redis宕机了，会选择其中一台作为master,当一台服务器使用这台且数据还没同步，那么其他就认为某个东西没被锁住，
但只是数据没来得及同步，那么就会出问题。

##### 如何解决这个问题嘞？
出现了redlock

##### 核心原理：

Redlock 是一种实现分布式锁的算法，它基于 Redis 提供的高效、可靠的分布式锁实现。Redlock 核心原理主要包括以下几个关键部分：

##### 核心思想

Redlock 的主要目标是确保在分布式系统中锁的独占性，即在同一时间只有一个客户端能够持有锁。它通过在多个独立的 Redis 实例上尝试获取锁来实现这一点，即使其中某些实例不可用，锁依然能保持高可用性。

##### 步骤详解
获取当前时间：
在获取锁的开始，记录当前时间 t1。
依次尝试在 N 个 Redis 实例上获取锁：
使用相同的 key 和具有唯一标识（如 UUID）的值，尝试在每个 Redis 实例上设置锁。设置锁时使用 SET resource_name my_random_value NX PX 30000 命令，其中：
NX 表示只有在 key 不存在时才设置成功（保证锁的独占性）。
PX 后面跟着的数值表示锁的过期时间（以毫秒为单位）。
依次尝试在每个 Redis 实例上获取锁，并记录成功获取锁的实例数目。
计算获取锁的总时间：
记录当前时间 t2，计算获取锁所花费的总时间 t = t2 - t1。
检查获取锁的条件：
如果在至少半数以上的 Redis 实例上成功获取了锁，并且获取锁的总时间 t 小于锁的过期时间，则认为锁获取成功。
如果锁获取失败，则在所有实例上释放锁（以防止持有部分锁的情况）。
锁的续约和释放：
续约：客户端在持有锁的过程中，可以定期续约锁的过期时间，以确保锁不会在处理过程中意外过期。
释放：当客户端完成任务后，应主动释放锁。释放锁时需要验证锁的持有者（即使用 DEL 命令删除锁之前，检查锁的值是否为客户端的唯一标识）。
Redlock 的可靠性
单点故障：Redlock 通过多个 Redis 实例（建议至少 3 个或 5 个）避免了单点故障。如果一个或两个实例不可用，锁依然可以正常运作。
时钟漂移：由于 Redlock 依赖于 Redis 实例的过期时间，系统时钟的漂移可能会影响锁的可靠性，因此需要确保系统时钟同步。
唯一标识：每次获取锁时使用唯一标识（如 UUID），确保即使锁过期或被错误删除，也不会影响其他客户端的锁。
伪代码示例
package main

import (
"fmt"
"time"
"context"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"sync"
)

// 初始化 Redis 客户端和 Redlock 实例
func initRedisClients() []*goredislib.Client {
// 这里应该创建多个 Redis 客户端，连接到不同的 Redis 实例
client1 := goredislib.NewClient(&goredislib.Options{Addr: "192.168.0.101:6379"})
client2 := goredislib.NewClient(&goredislib.Options{Addr: "192.168.0.102:6379"})
client3 := goredislib.NewClient(&goredislib.Options{Addr: "192.168.0.103:6379"})

	return []*goredislib.Client{client1, client2, client3}
}

// 获取分布式锁
func acquireLock(rs *redsync.Redsync, lockName string, ttl time.Duration) (*redsync.Mutex, error) {
mutex := rs.NewMutex(lockName, redsync.WithExpiry(ttl), redsync.WithRetryDelay(time.Millisecond*100))

	if err := mutex.Lock(); err != nil {
		return nil, err
	}

	return mutex, nil
}

// 释放分布式锁
func releaseLock(mutex *redsync.Mutex) error {
ok, err := mutex.Unlock()
if !ok || err != nil {
return fmt.Errorf("unlock failed: %v", err)
}
return nil
}

func main() {
// 初始化 Redis 客户端池
clients := initRedisClients()

	// 创建 Redsync 实例
	pools := []redsync.Pool{}
	for _, client := range clients {
		pools = append(pools, goredis.NewPool(client))
	}
	rs := redsync.New(pools...)

	// 锁的名称和 TTL
	lockName := "resource_lock"
	ttl := 8 * time.Second

	// 启动多个 goroutine 来模拟并发锁请求
	gNum := 2
	var wg sync.WaitGroup
	wg.Add(gNum)
	for i := 0; i < gNum; i++ {
		go func(id int) {
			defer wg.Done()

			fmt.Printf("Goroutine %d: 开始获取锁\n", id)
			mutex, err := acquireLock(rs, lockName, ttl)
			if err != nil {
				fmt.Printf("Goroutine %d: 获取锁失败: %v\n", id, err)
				return
			}
			fmt.Printf("Goroutine %d: 获取锁成功\n", id)

			// 模拟一些工作
			time.Sleep(3 * time.Second)

			fmt.Printf("Goroutine %d: 开始释放锁\n", id)
			if err := releaseLock(mutex); err != nil {
				fmt.Printf("Goroutine %d: 释放锁失败: %v\n", id, err)
			} else {
				fmt.Printf("Goroutine %d: 释放锁成功\n", id)
			}
		}(i)
	}

	wg.Wait()
}



##### 总结

Redlock 算法通过在多个 Redis 实例上尝试获取锁来实现高可用性和可靠性。
它的设计思想确保了即使在部分实例失效的情况下，锁依然能保持独占性，从而在分布式系统中实现可靠的锁机制
下面是 Redlock 算法的伪代码示例
## 第十三天，第十四天,第十五天
## 订单和购物车服务 service层

### 为什么UpdateCartItem只更新Checked和Nums
GORM 的自动更新时间戳管理
GORM 通过钩子函数自动处理时间戳字段：
CreatedAt：记录数据创建的时间。
UpdatedAt：记录数据最后一次更新的时间。
在执行 Save、Create 或 Update 操作时，GORM 会自动更新 UpdatedAt 字段。因此，在代码中不需要手动更新这个字段。
GORM 会在后台处理这一逻辑，确保 UpdatedAt 字段总是记录最新的更新时间。

本地事务
tx := global.DB.Begin()
tx.Save(&inv)
tx.Commit() // 需要自己手动提交操作