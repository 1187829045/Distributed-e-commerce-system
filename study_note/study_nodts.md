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
sudo docker run -d -p 8500:8500 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0
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
docker run --name nacos-standalone -e MODE=standalone -e JVM_XMS=512m -e JVM_XMX=512m -e JVM_XMN=256m -p 8848:8848 -d nacos/nacos-server:latest
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
## 第十三天，第十四天,第十五天,十六天
## 订单和购物车服务 service层，web层，用户操作

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

此时分布式事务还没有解决后续还进行解决。
具体到你提到的 github.com/mbobakov/grpc-consul-resolver 这个包，它是一个 gRPC 的 Consul 服务发现解析器。在 gRPC 应用中，
如果你需要使用 Consul 作为服务发现机制，并且想要在 gRPC 客户端中使用这个解析器，就需要在代码中引入这个包。

DeleteCartItem函数实现有点问题，进行修改
原函数 if result := global.DB. Delete(&model.ShoppingCart{},req.Id); result.RowsAffected == 0
容易出现问题，如果web端不做权限验证，那么这个Id如果是其他商品的Id，就有可能删错，就意味着web端必须做权限验证，为了严谨起见修改成下面的版本
if result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{}); result.RowsAffected == 0
通过web端传递userId和商品Id，前端拿商品ID更容易所以就不用主键ID。


claims 通常是指从 JWT（JSON Web Token）中提取的声明信息。JWT 是一种用于认证和信息交换的紧凑且自包含的方式，在许多 Web 应用中用于用户认证和会话管理。
JWT 中的 Claims,JWT 的结构由三部分组成：Header、Payload 和 Signature。claims 是 JWT 的 Payload 部分，包含了一组声明（claims），
这些声明用于存储关于用户和其他元数据的信息。常见的声明包括用户 ID、用户名、角色、权限等。
获取 Claims：在 Gin 中，你可以通过中间件解析 JWT，并将解析后的 claims 存储在 context 中，方便后续的处理函数使用。


### 公钥和私钥
涉及到钱的情况下往往会使用加密，公钥和私钥是一对的，首先我们的电商生成两个公钥并生成两个与之匹配的私钥，其中一对给支付宝，本身用私钥对支付宝的url进行加密，支付宝可以通过公钥解密，
并知道是哪个url发过来的。然后支付宝又通过私钥加密，发送回电商系统，比如告诉我客户的支付信息，然后电商通过公钥解密。第一对公钥和私钥是自己生成的，而后一对是支付宝生成的。

### 收藏的功能，收货地址，留言服务
这些就是认为是用户的操作的微服务。增删改查的操作。


## 第十七天，第十八天

### 商品的搜索功能，elasticsearch

#### 通过mysql进行搜索会出现以下问题：

性能底下，没有相关性排名，无法全文搜索，搜索不准确没有分词。

### 什么是全文搜索？
全文搜索（Full-Text Search）是一种用于在文档或数据库中查找特定词语或短语的技术。与传统的基于字段或标签的搜索不同，全文搜索会扫描文档的全部内容，
以找到匹配的词语或短语。这种搜索方式通常用于处理大量文本数据，如文章、书籍、电子邮件、网页等。

#### 全文搜索的关键特性包括：
关键词搜索：用户可以输入一个或多个关键词，系统会返回包含这些关键词的文档或记录。
布尔搜索：允许使用布尔运算符（如 AND、OR、NOT）来组合多个搜索条件，提高搜索的精确性。
短语搜索：可以搜索包含特定短语的文档，而不仅仅是单个词语。
模糊搜索：可以处理拼写错误或词形变化，返回与关键词近似的结果。
权重和排序：根据关键词出现的频率和位置，对搜索结果进行排序，确保最相关的结果排在前面。
分词和词干提取：将文档中的文本分解成独立的词语，并提取词干（如将“running”简化为“run”），以提高搜索的覆盖面。

#### 应用场景

全文搜索广泛应用于各种场景，包括但不限于：
搜索引擎：如Google、Bing等，通过全文搜索技术来查找和索引网页内容。
内容管理系统（CMS）：如WordPress、Drupal等，用于在网站内容中进行快速搜索。
电子邮件系统：如Gmail，通过全文搜索来查找邮件中的特定信息。
文档管理系统：如SharePoint，用于在大量文档中查找相关内容。

#### 技术实现

全文搜索可以通过多种技术和工具实现，包括：

数据库中的全文索引：许多关系型数据库（如MySQL、PostgreSQL）都提供内置的全文搜索功能，通过创建全文索引来加速搜索。
专用搜索引擎：如Elasticsearch、Apache Solr，它们专门设计用于处理大规模的全文搜索任务，具有高效的索引和查询能力。
库和框架：如Lucene，它是一个高性能、可扩展的全文搜索库，许多搜索引擎（如Elasticsearch和Solr）都基于Lucene构建。
总的来说，全文搜索是一种强大且灵活的搜索技术，能够在海量文本数据中快速找到所需的信息。


### elasticsearch 是什么
Elasticsearch 是一个开源的搜索和分析引擎，基于Apache Lucene构建。它提供了分布式、多租户功能的全文搜索引擎，并且具备RESTful web接口，
使其能够进行实时的搜索和数据分析。Elasticsearch在处理大规模数据集方面表现优异，因其高效的索引和查询能力而广受欢迎。
#### 主要特点

分布式架构：Elasticsearch可以水平扩展，通过将数据分布在多个节点上来处理和存储海量数据，从而提高系统的可扩展性和可靠性。
实时搜索和分析：Elasticsearch能够实时地索引和搜索数据，使其特别适用于需要即时搜索和分析的大数据应用。
RESTful API：Elasticsearch提供了简单的RESTful接口，方便与各种编程语言和平台集成。
全文搜索：基于Apache Lucene，Elasticsearch提供强大的全文搜索功能，包括关键词搜索、布尔搜索、短语搜索、模糊搜索等。
聚合分析：Elasticsearch提供丰富的数据聚合功能，可以进行复杂的数据分析和统计。
高度可配置：用户可以根据具体需求自定义索引和查询的行为，以优化性能和结果的相关性。
插件支持：Elasticsearch拥有广泛的插件生态系统，可以扩展其功能，如安全性、监控、分析等。

#### 典型应用场景

日志和事件数据分析：通过收集和分析日志数据来监控系统和应用的健康状况和性能。
全文搜索引擎：用于网站、应用程序或企业内部搜索，实现快速和准确的全文搜索。
数据分析和可视化：与Kibana结合使用，Elasticsearch能够实现强大的数据分析和可视化功能。
推荐系统：通过分析用户行为数据，提供个性化的推荐服务。

#### 技术实现

索引和分片：数据在Elasticsearch中存储为索引，索引可以进一步分为多个分片（Shards），每个分片可以独立存储和查询。
倒排索引：Elasticsearch使用倒排索引来加速搜索，存储了词语到其所在文档的映射。
集群管理：Elasticsearch集群由多个节点组成，每个节点可以充当主节点、数据节点或协调节点等不同角色，以分担不同的任务。
查询DSL：提供了强大的查询DSL（Domain Specific Language），用户可以通过JSON格式的请求来构建复杂的查询。

#### 启动并禁用

sudo systemctl stop firewalld.service
sudo systemctl disable firewalld.service//防止重启后又启动
sudo systemctl status firewalld.service

#### 通过docker安装elasticsearch

#新建es的config配置⽂件夹
sudo mkdir -p /data/elasticsearch/config
#新建es的data⽬录
sudo mkdir -p /data/elasticsearch/data
#新建es的plugins⽬录
sudo mkdir -p /data/elasticsearch/plugins
#给⽬录设置权限
sudo chmod 777 -R /data/elasticsearch
#写⼊配置到elasticsearch.yml中， 下⾯的 > 表示覆盖的⽅式写⼊， >>表示追加的⽅式写⼊，但是要确
sudo echo "http.host: 0.0.0.0" >> /data/elasticsearch/config/elasticsearch.yml
#安装es
sudo docker run --name elasticsearch -p 9200:9200 -p 9300:9300 \
-e "discovery.type=single-node" \
-e ES_JAVA_OPTS="-Xms128m -Xmx256m" \
-v /data/elasticsearch/config/elasticsearch.yml:/usr/share/elasticsearch/config/elasticsearch.yml \
-v /data/elasticsearch/data:/usr/share/elasticsearch/data \
-v /data/elasticsearch/plugins:/usr/share/elasticsearch/plugins \
-d elasticsearch:7.10.1

#### 通过docker安装kibana
docker run -d --name kibana -e ELASTICSEARCH_HOSTS="http://192.168.128.128:9200" -p 5601:5601 kibana:7.10.1

批量操作：单个操作每次都要简历http连接，会减慢，所以有了批量操作

### 倒排索引：
Elasticsearch 中的倒排索引（Inverted Index）是其核心特性之一，它是用于快速搜索和定位文档的关键组成部分。下面我会详细解释倒排索引的概念及其在 Elasticsearch 中的应用。

#### 1. 什么是倒排索引？

倒排索引是一种数据结构，它将文档中的内容进行统计和索引，以便快速查找文档。通常，文档被分解为词项（terms），并且记录每个词项在哪些文档中出现。
它的名称"倒排"来源于其不同于传统的直接索引，传统的索引是记录每个文档中的词项，而倒排索引是记录每个词项出现在哪些文档中。

#### 2. 倒排索引的组成部分

在 Elasticsearch 中，每个索引都包含一个或多个分片（shard），每个分片都包含一个完整的倒排索引。以下是倒排索引的主要组成部分：

词项列表（Term Dictionary）： 记录了所有不重复的词项及其在倒排索引中的位置信息。

倒排列表（Inverted List）： 每个词项都有一个对应的倒排列表，记录了包含该词项的所有文档及其位置信息（如文档ID、位置等）。

词项频率（Term Frequency）和文档频率（Document Frequency）：

词项频率（TF）：指定词项在文档中的出现次数。
文档频率（DF）：指定词项出现在多少个文档中。

#### 3. 倒排索引的创建过程

在索引文档到 Elasticsearch 时，它会执行以下步骤来创建倒排索引：

分析文档内容：首先，文档内容会经过分析器（Analyzer）处理，将文本分割成适合索引的词项。

构建倒排索引：根据分析结果，构建每个词项的倒排列表，并更新词项列表和其他元数据。

#### 4. 倒排索引的搜索和检索

一旦索引创建完成，可以使用查询（Query）来搜索文档。Elasticsearch 使用倒排索引来加速查询过程，例如：

词项匹配查询（Term Queries）：根据特定词项搜索文档。

短语查询（Phrase Queries）：查找包含特定短语的文档。

全文搜索（Full-text Search）：对文档中的所有词项进行搜索。

#### 5. 倒排索引的优点

快速搜索和检索：倒排索引提供了快速定位文档的能力，特别适用于大规模的文本搜索和分析任务。

空间效率：与直接索引相比，倒排索引通常占用更少的存储空间。

#### 结论

倒排索引是 Elasticsearch 强大搜索引擎的核心技术之一，它提供了高效的文本搜索和检索功能，适用于各种文本数据分析和搜索应用场景。
理解倒排索引的原理和运作方式，有助于更好地利用 Elasticsearch 提供的强大功能来处理和分析数据。

### term查询

term : 这种查询和match在有些时候是等价的，比如我们查询单个的词hello，那么会和match查询结果一
样，但是如果查询"hello world"，结果就相差很大，因为这个输入不会进行分词，就是说查询的时候，是
查询字段分词结果中是否有"hello world"的字样，而不是查询字段中包含"hello world"的字样，
elasticsearch会对字段内容进行分词，“hello world"会被分成hello和world，不存在"hello world”，因此这
里的查询结果会为空。这也是term查询和match的区别。


## 什么是 Mapping？
 Mapping 类似于数据库中的表结构定义 schema，它有以下几个作用：
定义索引中的字段的名称定义字段的数据类型，比如字符串、数字、布尔字段，倒排索引的相关配置，比
如设置某个字段为不被索引、记录 position 等在 ES 早期版本，一个索引下是可以有多个 Type ，从 7.0
开始，一个索引只有一个 Type，也可以说一个 Type 有一个 Mapping 定义。type使keyword不会进行分词。

### 为什么match查询keyword可以查到，不是要先分词吗？

#### 分析器

Analyzer 由三部分组成：Character Filters、Tokenizer、Token Filters

在 Elasticsearch 中，分析器（Analyzer）是处理文本的核心组件。分析器将文本数据分解成独立的词项（tokens），然后这些词项可以被索引和搜索。
一个分析器由以下三部分组成：

##### 字符过滤器（Character Filters）：

功能：字符过滤器在文本被分词之前对文本进行预处理。它们可以删除、替换或添加字符。通常用于处理 HTML 标签、标点符号等。
常见的字符过滤器：
html_strip：移除 HTML 元素。
mapping：通过映射来替换字符。

##### 分词器（Tokenizer）：
功能：分词器将输入的文本流分解为独立的词项（tokens）。这是分析器中最关键的一步，决定了文本被如何分割。
###### 常见的分词器：
standard：标准分词器，按词边界分词。
whitespace：按空白字符分词。
keyword：不进行分词，整个输入作为一个词项。
pattern：基于正则表达式进行分词。

##### 词项过滤器（Token Filters）：
功能：词项过滤器在分词之后对词项进行处理。它们可以修改、删除或添加词项，进行词干提取、同义词替换等。
常见的词项过滤器：
lowercase：将所有词项转换为小写。
stop：移除停用词（如 "the", "and"）。
stemmer：进行词干提取，如将 "running" 转换为 "run"。
synonym：同义词过滤器，将同义词词项标准化。
了解分析器的工作原理和如何配置分析器，可以帮助你更好地控制文本的索引和搜索行为，从而提升 Elasticsearch 的搜索性能和准确性。

执行搜索的时候，查询顺序，指定的分析器->存储的分析器->设置里面的分析器->mapping 的分析器。

IK分词器，中文分词器。集成方便
GET _analyze
{
"text":"中国科学技术大学",
"analyzer": "ik_max_word"
}

### 添加自己的词汇

llb@llb-virtual-machine:/data/elasticsearch/plugins/ik/config$ mkdir custom
llb@llb-virtual-machine:/data/elasticsearch/plugins/ik/config$ cd custom/
llb@llb-virtual-machine:/data/elasticsearch/plugins/ik/config/custom$ vim mydic.dic
llb@llb-virtual-machine:/data/elasticsearch/plugins/ik/config/custom$ vim extra_stopword.dic
llb@llb-virtual-machine:/data/elasticsearch/plugins/ik/config/custom$ cd ..
llb@llb-virtual-machine:/data/elasticsearch/plugins/ik/config$ vim IKAnalyzer.cfg.xml

### 保存之后的钩子逻辑（after-save hook）

是指在某些编程框架或系统中，当某个数据对象被保存之后，触发的一段自定义代码逻辑。它通常用于在数据保存完成后执行一些额外的操作，比如更新关联数据、发送通知、记录日志等。
具体来说，可以在以下场景中使用保存之后的钩子逻辑：
数据库操作：在关系型数据库中，当一条记录被保存后，可以触发钩子逻辑来更新其他相关表的数据或执行一些业务逻辑。
Web框架：在使用Web框架（如Django、Ruby on Rails等）时，可以定义模型的after-save钩子，在模型对象保存后执行特定操作。
事件驱动系统：在事件驱动的架构中，当某个保存操作完成时，触发相关事件处理器来处理后续的逻辑。


### 数据的一致性

现在像mysql和es写两份数据。如何解决数据一致性。什么叫一致性，ES入库失败，但是Mysql写入成功，ES成功，mysql不存在。
MySQL写入成功，但ES写入失败：这意味着MySQL中存在数据，但Elasticsearch中没有对应的数据。
ES写入成功，但MySQL写入失败：这意味着Elasticsearch中存在数据，但MySQL中没有对应的数据。
事务性保证
如果你的业务允许，可以使用分布式事务来确保MySQL和Elasticsearch的数据一致性。使用两阶段提交（2PC）或类似的分布式事务协议来保证数据的一致性。但这种方法实现复杂，性能开销大，一般在强一致性需求下才会使用。

异步补偿机制
   可以采用异步补偿机制来解决数据一致性问题。当MySQL写入成功但Elasticsearch写入失败时，可以记录失败的操作，并在后台重试写入Elasticsearch。这种方法实现简单，适用于最终一致性的场景。

事件驱动架构
   通过事件驱动架构，可以使用消息队列（如Kafka、RabbitMQ等）来确保数据一致性。具体步骤如下：

发布事件：在MySQL写入成功后，发布一个数据变更事件到消息队列。
消费事件：设置一个消费者来监听消息队列，并将数据写入Elasticsearch。
失败重试：如果Elasticsearch写入失败，可以记录失败事件并重试。
一致性检查
   定期进行一致性检查，比较MySQL和Elasticsearch中的数据，发现不一致时进行修复。例如，可以定期扫描MySQL中的数据，并检查是否在Elasticsearch中存在对应的数据，如果不存在则进行重新写入。
分离读写
   在某些场景下，可以考虑分离读写操作，即MySQL只负责写入数据，Elasticsearch只负责读取数据。这种方式可以避免写入时的数据一致性问题，但需要业务上能够接受这种读写分离的模式。

使用外部工具
   一些外部工具可以帮助管理和同步数据，例如Debezium，它可以捕获MySQL的变更事件，并将其同步到Elasticsearch中。

`//解决数据的一致性
tx := global.DB.Begin()
result := tx.Save(&goods) //自动调用AfterCreate
if result.Error != nil {
tx.Rollback()
return nil, result.Error
}
tx.Commit() //手动提交
return &proto.GoodsInfoResponse{
Id: goods.ID,
}, nil`


### 超时机制

#### 为什么订单会有超时机制

订单超时机制是指在电子商务平台、在线支付系统等场景中，设置一个时间限制，如果用户在规定时间内未完成支付、确认或其他操作，订单将自动取消或变为无效状态。这种机制有以下几个重要的原因和作用：

#####  防止资源占用
订单超时机制可以防止用户长时间占用库存资源，影响其他用户的购买体验。
库存管理：当用户创建订单但未支付时，系统通常会预留相应的库存。如果没有超时机制，这些库存将长时间无法售出，影响其他潜在买家的购买。
服务资源：在某些服务类型的交易中（如预定票务、预约服务等），未支付的订单会占用时间和服务资源。超时机制可以释放这些资源，让其他用户能够使用。

##### 提高用户决策效率
设定一个时间限制可以促使用户在一定时间内完成决策，减少犹豫和拖延。
购买决策：用户在下单后，如果知道有时间限制，会更快地做出支付决定，从而加快交易流程。
避免遗忘：有时用户可能会下单后离开而忘记支付，超时机制可以提醒用户完成操作或重新下单。

#####  防止恶意订单
超时机制可以减少恶意用户创建大量订单占用资源的问题。
减少刷单：某些恶意用户可能会通过大量创建未支付订单来占用库存或其他资源，超时机制可以有效减少这种行为。

##### 系统稳定性

通过设置订单超时机制，可以减轻系统的负担，保持系统的稳定性和性能。
资源管理：系统需要管理的有效订单数量减少，有助于优化系统资源和性能。
简化逻辑：订单状态管理更加简化和明确，减少因长时间未支付订单导致
的复杂逻辑处理。

## 事务和分布式事务

### 事务的概念和作用
**事务（Transaction）**是在数据库或其他系统中的一系列操作，这些操作要么全部成功，要么全部失败，确保数据的完整性和一致性。
#### 事务的四大特性通常被称为ACID特性：
原子性（Atomicity）：事务中的所有操作要么全部完成，要么全部不完成。换句话说，事务是一个不可分割的单位。
一致性（Consistency）：事务执行前后，数据库的状态必须保持一致。也就是说，事务执行之后，所有的数据规则都必须得到满足。
隔离性（Isolation）：一个事务的执行不能被其他事务干扰，事务之间的操作是隔离的，互不干扰。
持久性（Durability）：一旦事务提交，事务对数据库的改变是永久的，即使系统崩溃也不会丢失。
#### 事务的作用
数据一致性：确保数据库在事务执行前后保持一致状态。
数据完整性：保护数据的完整性，即使在发生系统故障时。
并发控制：允许多个用户同时访问数据库，而不影响事务的完整性和一致性。
### 分布式事务的概念和作用
分布式事务（Distributed Transaction） 是指在分布式系统中，事务的操作涉及多个独立的事务资源或服务，这些资源可能位于不同的数据库、不同的网络节点甚至不同的系统中。
#### 作用
分布式事务的主要作用是确保在分布式系统中，跨多个系统、服务、数据库的操作能够保持一致性和完整性。分布式事务使得多个参与节点能够协同工作，确保数据的一致性和完整性。
#### 分布式事务的特征
跨越多个节点：分布式事务涉及多个数据库或服务节点。
复杂性：由于涉及多个独立的系统，分布式事务的管理和协调比单机事务要复杂得多。
CAP理论：分布式系统中存在一致性（Consistency）、可用性（Availability）、分区容忍性（Partition Tolerance）三者不可兼得的挑战。
性能影响：分布式事务需要在多个节点间进行协调，会带来额外的网络通信和同步开销，影响系统性能。

## 分布式系统中会出现哪些故障导致数据不一致
在分布式系统中，由于系统的复杂性和多节点协作，数据不一致的问题时常会出现。以下是一些常见的故障和原因：

### 网络故障

原因：
网络分区：网络分区会导致系统中的一部分节点无法与另一部分节点通信，形成孤岛效应。
网络延迟：高延迟可能导致消息超时或顺序错乱。
消息丢失：网络中传输的消息可能会因为各种原因丢失。
影响：
数据可能会不同步或部分节点的数据更新被延迟或丢失。
在网络分区的情况下，不同分区中的数据可能会出现分歧。

### 节点故障
原因：
节点崩溃：服务器宕机、硬件故障、操作系统崩溃等。
服务重启：由于软件升级、维护等原因导致服务重启。
影响：
数据更新操作在中途失败，可能导致部分写入的数据丢失。
某些节点上的数据可能与其他节点不同步，出现数据不一致。

### 数据复制故障

原因：
复制延迟：数据在主节点和从节点之间复制时的延迟。
复制失败：复制过程中出现错误导致数据没有正确复制到从节点。
影响：
从节点的数据可能与主节点不一致，导致读取数据时出现旧数据或不完整数据。

### 并发控制问题

原因：
锁竞争：多个事务竞争同一资源，可能导致死锁或数据更新的顺序问题。
事务隔离级别：不同的事务隔离级别可能导致脏读、不可重复读和幻读问题。
影响：
多个并发事务操作同一数据，可能导致数据不一致。

### 分布式事务问题

原因：
两阶段提交失败：两阶段提交协议中的任一阶段失败，可能导致数据不一致。
超时：事务在规定时间内未能完成，导致回滚或部分提交。
影响：
分布式事务的一部分提交成功，而另一部分提交失败，导致数据不一致。

### 时钟同步问题

原因：
时钟漂移：不同节点的时钟不一致，导致时间戳误差。
时间同步失败：NTP服务故障或网络原因导致时间同步失败。
影响：
时间戳不一致可能导致数据写入和读取的顺序错乱，影响数据的一致性。

### 数据修复和回滚操作失败

原因：
数据修复错误：手动或自动的数据修复操作未能正确执行。
回滚操作失败：事务回滚操作未能成功执行。
影响：
数据修复或回滚失败可能导致数据状态不正确，影响整体数据一致性。

## 应对措施
为了解决这些问题，可以采取以下措施：
 网络分区容忍：
使用CAP理论中的P（Partition tolerance），设计系统时考虑网络分区情况。
使用一致性算法如Paxos或Raft，确保在网络分区的情况下尽可能维持数据一致性。
数据复制优化：
使用异步复制或多主复制，减小复制延迟。
实施数据校验和修复机制，定期检查和修复数据不一致的问题。
分布式事务管理：
使用分布式事务协议如2PC（Two-Phase Commit）或3PC（Three-Phase Commit）。
使用Saga模式或TCC（Try-Confirm-Cancel）模式，确保分布式事务的一致性。
时钟同步：
部署可靠的NTP服务器，确保系统时钟的同步。
使用逻辑时钟或矢量时钟，减少对物理时钟的依赖。
并发控制：
使用乐观锁或悲观锁机制，控制并发访问。
根据业务需求选择合适的事务隔离级别。
故障检测和恢复：
部署高可用性和故障检测机制，快速识别和恢复故障节点。
实施数据备份和恢复策略，确保在节点故障时数据能够恢复。
通过以上措施，可以在分布式系统中尽量减少数据不一致的风险，确保系统的可靠性和数据的一致性。

## CAP和BASE理论

在分布式系统的设计和实现中，CAP和BASE理论是两种重要的理论，用于理解和权衡一致性、可用性和性能之间的关系。以下是对这两种理论的详细解释：

### CAP理论

CAP理论（也称为Brewer定理）由Eric Brewer在2000年提出，理论指出在分布式数据存储中，不可能同时满足以下三个特性：

##### 一致性（Consistency）：

每次读取都能返回最近写入的结果。即所有节点在同一时刻看到的数据都是相同的。

#### 可用性（Availability）：

每个请求都能收到非错误响应——但不保证它是最新的数据。系统始终能响应请求，即使部分节点出现故障。

#### 分区容忍性（Partition Tolerance）：

系统能够继续运行，即使网络分区使得部分节点之间无法通信。分区容忍性要求系统在网络分区情况下仍能处理请求。
CAP理论指出，在任何分布式系统中，最多只能同时满足这三个特性中的两个。例如：

 CA（一致性 + 可用性）：
系统能够在没有网络分区的情况下保证一致性和可用性，但在出现网络分区时无法保证。
示例：传统的关系型数据库在单节点环境下能满足CA。
CP（一致性 + 分区容忍性）：
系统能够在网络分区情况下保证一致性，但可能会牺牲部分可用性（例如，一些请求可能会被拒绝或阻塞）。
示例：多数一致性协议（如Paxos、Raft）在网络分区时优先保证一致性。
AP（可用性 + 分区容忍性）：
系统能够在网络分区情况下继续提供服务，但无法保证所有读取都是一致的。
示例：一些NoSQL数据库（如Cassandra、DynamoDB）在网络分区时优先保证可用性。

### BASE理论
BASE理论是对CAP理论的一种折中方案，由Amazon提出，用于大规模分布式系统。BASE是以下术语的缩写：

#### 基本可用（Basically Available）：
系统在出现故障时仍能基本可用，即使有部分功能不可用。系统可以在降级模式下运行，提供部分功能而不是完全宕机。
软状态（Soft State）：
系统中的状态可以在一定时间内不同步。节点间的数据同步是渐进的，而不是即时的。

#### 最终一致性（Eventual Consistency）：

系统保证在没有更多更新操作后，最终所有副本的数据会达到一致。最终一致性是一种弱一致性模型，允许数据在短时间内不一致，但保证最终会收敛到一致状态。
BASE理论强调可用性和容错性，适用于对实时一致性要求不高的场景，通常用于大规模互联网应用和NoSQL数据库。

#### 对比与应用
CAP理论：
强调系统在一致性、可用性和分区容忍性之间的权衡，适用于需要严格一致性和容错性的场景。
主要用于设计和选择分布式系统的架构，如分布式数据库、分布式缓存等。
BASE理论：
强调系统的可用性和最终一致性，适用于对实时一致性要求不高但需要高可用性和容错性的场景。
主要用于大规模互联网应用、分布式缓存系统和NoSQL数据库。
实践应用
CA系统：一般适用于强一致性要求的系统，但在分布式环境中很少完全实现。
例子：单节点的关系数据库。
CP系统：适用于需要严格一致性且能够容忍短暂不可用的系统。
例子：Zookeeper、HBase。
AP系统：适用于对可用性要求高且可以容忍短暂不一致的系统。
例子：Cassandra、DynamoDB、CouchDB。
总结
CAP和BASE理论提供了理解和设计分布式系统的重要框架，帮助我们在一致性、可用性和性能之间做出合理的权衡和选择。 CAP理论强调在分布式环境中不可能
同时满足一致性、可用性和分区容忍性，而BASE理论则提供了一种实际可行的折中方案，适用于需要高可用性和容错性的应用场景。

## 分布式事务如何和解决同时成功或失败的
### 两阶段提交(缺点明显一般不会用到)
两阶段提交协议（Two-Phase Commit，2PC）是一种经典的分布式事务协议，用于确保分布式系统中事务的所有参与节点要么全部提交事务，要么全部回滚事务，
以保证数据的一致性。两阶段提交协议通过协调器（Coordinator）和参与者（Participant）的协作，分为两个阶段：准备阶段（Prepare Phase）和提交阶段（Commit Phase）。下面详细介绍2PC的工作流程、优点和缺点。

#### 准备阶段（Prepare Phase）

在准备阶段，协调器向所有参与者发送准备请求，询问他们是否可以准备提交事务。每个参与者执行事务的本地预处理并写入日志，然后做出响应：
协调器向所有参与者发送prepare请求。
每个参与者接收到请求后，执行事务的本地操作，但不提交，只是记录操作日志（写入预备日志）。
如果参与者可以成功执行操作，它们会返回yes（准备好）响应；否则，返回no（失败）响应。 

#### 提交阶段（Commit Phase）

在提交阶段，协调器根据所有参与者的响应决定是提交事务还是回滚事务：
如果所有参与者都返回yes响应，协调器向所有参与者发送commit请求，要求他们正式提交事务。
如果任何参与者返回no响应，协调器向所有参与者发送rollback请求，要求他们回滚事务。

#### 完整工作流程

##### 事务开始：

协调器启动事务，生成唯一的事务ID。

##### 准备阶段：

协调器向所有参与者发送prepare请求。
参与者执行本地预处理并返回yes或no响应。

##### 提交阶段：

如果所有参与者返回yes，协调器向所有参与者发送commit请求。
如果任何参与者返回no，协调器向所有参与者发送rollback请求。

##### 事务结束：

参与者接收到commit请求后正式提交事务，或接收到rollback请求后回滚事务。
协调器在所有参与者完成提交或回滚后，记录事务完成状态。

#### 优点

简单性：2PC协议相对简单，容易实现和理解。
一致性保证：确保所有参与节点在事务提交时的数据一致性。

#### 缺点

同步阻塞：参与者在等待协调器指令时会阻塞，可能导致资源锁定和性能问题。
单点故障：协调器是单点故障，如果协调器崩溃，整个事务将无法完成。
脑裂问题：网络分区或节点故障可能导致参与者处于不确定状态（即不知道该提交还是回滚）。


### TCC分布式事务实现方案

TCC（Try-Confirm-Cancel）是一种分布式事务实现方案，用于确保跨多个服务或资源的事务一致性。TCC模型将一个事务分为三个阶段：尝试（Try）、确认（Confirm）和取消（Cancel）。以下是TCC分布式事务实现方案的详细介绍。
TCC模型概述
TCC模型概述

#### Try阶段：

预留资源或执行事务准备工作，但不提交事务。确保资源可用并锁定资源。

#### Confirm阶段：

在所有操作都成功预留资源后，正式执行事务提交。确保事务的最终一致性。

#### Cancel阶段：

如果任何操作在Try阶段失败或无法确认，则执行事务回滚，释放预留的资源。


### 基于本地消息表的最终一致性，比TCC稍微弱一点的一致性，且更加简单
本地消息表是一种实现分布式事务最终一致性的常用方法，特别适用于需要在多个独立服务之间保证数据一致性的场景。该方法通过将消息存储在本地数据库表中，
以确保消息和业务操作在同一个事务中完成，从而避免分布式事务的一些复杂性和性能开销。
本地消息表的工作原理
本地消息表的最终一致性主要依赖于以下步骤：
业务操作与消息存储在同一事务中：
在执行业务操作的同时，将需要传递的消息存储在本地消息表中，这两个操作在同一个数据库事务中进行，以确保原子性。
定期扫描本地消息表并发送消息：
使用一个定时任务或消息服务扫描本地消息表，发送未发送的消息。
消息消费方确认消息处理结果：
消费方接收到消息后处理业务逻辑，并发送确认反馈。
删除或标记已处理的消息：
发送方在收到消费方的确认反馈后，删除或标记本地消息表中的消息为已处理。

### 最大努力通知(没有上面一个应用广)
最大努力通知（Best Effort Notification）是一种分布式系统中的数据一致性解决方案，它在保证数据一致性的同时，尽量减少对系统性能的影响。
这种方法主要应用于事务的异步处理场景，通过尽最大努力通知的方式，确保尽可能地通知到所有参与方，但不保证一定能够通知成功。
工作原理
最大努力通知的基本思想是，在执行事务后，系统尽可能多次地尝试通知其他参与方执行相应的操作。如果通知失败，系统会记录失败的尝试，并继续重试，直到达到预设的重试次数或超时时间。
主要步骤
事务操作与通知记录：
执行业务操作并记录操作结果。
将需要通知的消息记录在通知表中。
定期扫描通知表并尝试通知：
使用定时任务或消息队列扫描通知表，尝试通知未成功的消息。
处理通知结果：
如果通知成功，更新通知表中的状态。
如果通知失败，记录失败原因并重试，直到达到最大重试次数。

## MQ
消息队列（Message Queue，MQ）是一种用于在分布式系统中实现异步通信的中间件，它通过存储和转发消息来解耦生产者和消费者，从而提高系统的可扩展性、可靠性和性能。消息队列在微服务架构中非常常见，可以用于日志处理、事件驱动架构、数据同步等场景。

### 消息队列的基本概念

生产者（Producer）：负责发送消息到消息队列的应用程序或服务。
消费者（Consumer）：负责从消息队列中接收和处理消息的应用程序或服务。
队列（Queue）：存储消息的容器，通常是FIFO（First In First Out）结构。
消息（Message）：生产者发送并由消费者接收和处理的数据单元。
主题（Topic）：用于发布/订阅模型中的消息通道，生产者发布消息到主题，多个消费者可以订阅该主题以接收消息。

### 消息队列的优势

解耦：生产者和消费者之间的松耦合关系，使系统更加灵活和易于维护。
异步处理：允许生产者在发送消息后立即返回，提高系统的响应速度和吞吐量。
负载均衡：通过将消息分发到多个消费者，分散处理负载，提高系统的处理能力。
可靠性：消息队列通常提供消息持久化、重试机制和确认机制，确保消息不丢失。
扩展性：可以方便地添加新的消费者或生产者，提升系统的扩展能力。

### 消息队列的模型

点对点模型（Point-to-Point Model）
在点对点模型中，消息生产者将消息发送到一个队列中，一个消息只能被一个消费者消费。适用于任务队列等场景。
队列：消息存储的容器。
生产者：发送消息到队列。
消费者：从队列中接收并处理消息。
发布/订阅模型（Publish/Subscribe Model）
在发布/订阅模型中，消息生产者将消息发布到一个主题中，所有订阅该主题的消费者都可以接收到消息。适用于广播通知等场景。
主题：消息发布的通道。
发布者：发布消息到主题。
订阅者：订阅主题以接收消息。

安装mq
将压缩包放在/home/llb目录下
解压缩 sudo unzip install.zip
cd install
sudo docker-compose up


### 链路追踪

sudo docker run \
--rm \
--name jaeger \
-p 6831:6831/udp \
-p16686:16686 \
jaegertracing/all-in-one:latest


### kong的安装
docker run -d --name kong-database -p 5432:5432 -e "POSTGRES_USER=kong" -e "POSTGRES_DB=kong" -e "POSTGRES_PASSWORD=kong" -e "POSTGRES_DB=kong" postgres:12

sudo docker run --rm -e "KONG_DATABASE=postgres" -e "KONG_PG_HOST=192.168.128.128" -e "KONG_PG_PASSWORD=kong" -e "POSTGRES_USER=kong" -e "KONG_CASSANDRA_CONTACT_POINTS=kong-database" kong kong migrations bootstrap

sudo curl -Lo kong.2.5.0.amd64.deb "https://download.konghq.com/gateway-2.x-ubuntu-$(lsb_release -cs)/pool/all/k/kong/kong_2.5.0_amd64.deb"
curl -Lo kong-enterprise-edition-2.6.1.0.all.deb "https://packages.konghq.com/public/gateway-legacy/deb/ubuntu/pool/xenial/main/k/ko/kong-enterprise-edition_2.6.1.0/kong-enterprise-edition_2.6.1.0_all.deb"
sudo apt install -y ./kong-enterprise-edition-2.6.1.0.all.deb
sudo systemctl stop firewalld.service
sudo systemctl restart docker
sudo cp /etc/kong/kong.conf.default /etc/kong/kong.conf
sudo vim /etc/kong/kong.conf
sudo kong start -c /etc/kong/kong.conf

### kongga安装

sudo docker run -d -p 1337:1337 --name konga pantsel/konga

### 安装jenkins
https://blog.csdn.net/xhmico/article/details/136535498
cd /soft/jenkins
java -jar jenkins.war --httpPort=8088
3cb81fac55b54e0d83650c8de9700f5c
This may also be found at: /home/llb/.jenkins/secrets/initialAdminPassword

修改插件地址
sudo sed -i 's/https:\/\/updates.jenkins.io\/download/http:\/\/mirrors.tuna.tsinghua.edu.cn\/jenkins/g' /var/lib/jenkins/updates/default.json
sudo sed -i 's/http:\/\/www.google.com/http:\/\/www.baidu.com/g' /var/lib/jenkins/updates/default.json

http://mirrors.tuna.tsinghua.edu.cn/jenkins/updates/update-center.json


