## 第一天

1.安装环境
2.rpc :
rpc 叫做远程过程调用，一个节点请求另一个节点提供的服务
对应rpc的是本地过程调用，函数调用是最常见的本地过程调用
将本地过程调用变成远程过程调用会面临各种问题

rpc 技术在架构设计上由四部分 客户端，客户端存根，服务端，服务端存根

## 第二天

第四周第三章,第四章
安装学习了grpc

gRPC（gRPC Remote Procedure Call）是一个高性能、通用的开源远程过程调用（RPC）框架，由Google开发。
它基于HTTP/2协议和Protocol Buffers（protobuf）数据序列化格式，旨在简化跨网络的函数调用，
使得不同系统之间的通信更高效、更可靠。gRPC适用于构建分布式系统和微服务架构中的服务间通信。
服务端的数据流模式：这种模式是客户端发起一次请求，服务端返回一段连续的数据流。
客户端数据来模式：与上面相反，客户端源源不断向服务端发送数据流
双向数据流模式：上面综合，比如聊天软件
第五周第一章
学习protobuf grpc

# 第三天

## 第五周

grpc的metadata
在 gRPC 中，拦截器（interceptor）是一种中间件机制，可以在 gRPC 方法调用之前或之后执行一些通用的逻辑。

interceptor := func(ctx context.Context, req interface{},
info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
fmt.Println("接受到一新的请求")
return handler(ctx, req)
}

ctx 是一个上下文对象，它携带了请求的范围和生命周期信息
req 是请求参数，表示传递给 gRPC 方法的具体请求消息
info 提供了关于当前调用的一些信息，例如被调用的方法的全名等。
handler 是实际处理请求的函数，当拦截器完成它的工作之后，可以调用这个函数来继续处理请求。
拦截器在 gRPC 中扮演着类似于中间件的角色，允许在 gRPC 方法调用之前或之后插入自定义逻辑。这对于实现许多通用功能非常有用

grpc 通过metadata传输密码，放在拦截器中，可以形成验证的中间件方式

### 验证器生成命令

protoc -I. --go_out=. helloworld.proto
protoc -I. --go-grpc_out=. helloworld.proto
protoc -I. --validate_out="lang=go:." helloworld.proto

grpc客户端与服务端流程

### 服务端

// 创建一个新的 gRPC 服务器实例
g := grpc.NewServer()
// 注册 Greeter 服务到 gRPC 服务器
// 将 Server 的实例传递给生成的 RegisterGreeterServer 方法
pb.RegisterGreeterServer(g, &Server{})
// 监听所有网络接口上的 8080 端口
lis, err := net.Listen("tcp", "0.0.0.0:8080")
// 启动 gRPC 服务器以监听传入的连接
// 如果服务器运行过程中出现错误，记录错误并退出程序
if err := g.Serve(lis); err != nil {
log.Fatalf("failed to serve: %v", err)
}

### 客户端

// 创建一个与 gRPC 服务器的连接
conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure()) // 改为 8080 端口
// 延迟关闭连接，确保程序结束前连接会被关闭
defer conn.Close()
// 创建一个新的 Greeter 客户端实例
c := proto.NewGreeterClient(conn)
// 调用 服务端定义的方法
r, err := c.SayHello(context.Background(), &proto.HelloRequest{Name: "bobby"})

## 第六周
### 单体应用

单体应用是一种传统的软件架构模式，它是指将整个应用程序作为一个单一的、 完整的单元来开发、部署和维护。
在单体应用中， 所有的功能模块、业务逻辑以及数据访问都被打包在一起，并部署到同一个运行环境中。
可扩展性差：随着应用规模的增长，单体应用的可扩展性会变差。由于所有的组件都耦合在一起，很难对某个特定组件进行水平扩展。
难以维护：随着时间的推移，单体应用往往会变得越来越庞大、复杂，导致代码难以理解和维护。
技术栈固定：由于整个应用都使用同一种技术栈，因此很难引入新的技术或更新现有技术。

### 微服务应用
微服务应用是一种分布式系统架构模式，它将应用程序拆分成多个小型服务，每个服务都有自己独立的业务功能，
并通过轻量级的通信机制来相互协作。 每个微服务都可以独立部署、扩展和维护，从而提高了系统的灵活性、可扩展性和可维护性。
前后端分离系统开发接口管理痛点

# 第四天
## 第六周
### yapi
Yapi 是一个可视化的接口管理工具，主要用于接口的管理、文档的生成和测试的执行.
它提供了一个用户友好的界面，让开发团队能够更轻松地管理接口，生成接口文档，并进行接口测试。
Yapi 的一些主要功能和作用：
接口管理： Yapi 提供了一个可视化的界面，让用户能够方便地创建、编辑和删除接口，管理接口的请求和响应参数，设置接口的状态和分类等。
接口文档生成： Yapi 能够自动根据接口的定义生成接口文档，包括接口的描述、请求参数、响应参数等信息，使团队成员能够更清晰地了解接口的使用方法和规范。
接口测试： Yapi 允许用户在界面中直接进行接口测试，通过输入参数和发送请求来测试接口的功能和性能，并查看测试结果和日志。
权限管理： Yapi 支持多用户协作，管理员可以根据需要设置不同用户的权限，控制其对接口的访问、编辑和管理权限。
团队协作： Yapi 提供了团队协作的功能，允许多个团队成员共同使用和管理接口，提高团队的协作效率和开发质量。

### orm
ORM（Object-Relational Mapping）是一种编程技术，用于将面向对象的程序中的对象与关系型数据库中的表之间建立映射关系，
从而实现对象和数据库之间的数据转换和交互。
在传统的关系型数据库中，数据是以表的形式存储的，每个表包含多个列，每行代表一个记录。而在面向对象编程中，数据以对象的形式表示，每个对象包含多个属性。
ORM 技术的目标是解决关系型数据库和面向对象编程之间的差异，让开发者可以使用面向对象的方式来操作数据库，而不必直接编写 SQL 语句。

#### ORM 主要有以下几个核心概念：

对象（Object）： 在面向对象编程中表示一个实体或数据模型的抽象，通常由类或结构体表示
关系型数据库表（Relational Database Table）： 数据库中的数据以表的形式组织，每个表包含多个行和列，每行代表一个记录，每列代表一个属性。
映射（Mapping）： 将对象的属性映射到数据库表的列，使得对象和数据库表之间可以相互转换和交互。
持久化（Persistence）： 将对象的状态持久化到数据库中，使得对象的数据能够在程序重启后得以保留。

#### ORM 的优势包括：

提高开发效率： ORM 让开发者不必直接编写 SQL 语句，减少了与数据库交互的复杂性，提高了开发效率。
提高可维护性： ORM 使得代码更加清晰和易于理解，降低了程序的耦合度，提高了代码的可维护性。
跨平台兼容性： ORM 提供了对多种数据库的支持，使得应用程序能够在不同的数据库管理系统之间迁移和切换。
### gorm
ORM 是一个 Go 语言的 ORM（Object-Relational Mapping）库，用于简化 Go 语言程序与数据库的交互。
ORM 是一种编程技术，它允许开发者使用面向对象的方式来操作数据库，而不必直接编写 SQL 语句。
GORM 提供了一组方法和结构体，用于定义和操作数据库模型（Model），并提供了对常见数据库操作的支持，包括增删改查（CRUD）、数据迁移、事务处理等功能。
它支持多种主流的关系型数据库，如 MySQL、PostgreSQL、SQLite、SQL Server 等。
#### GORM 主要特点包括：
简洁易用： GORM 提供了简洁的 API 和丰富的文档，使得开发者能够轻松上手，并且提高了代码的可读性和可维护性。
自动迁移： GORM 支持自动迁移功能，能够根据模型定义自动创建、更新数据库表结构，方便开发者进行数据库的版本管理和升级。
链式操作： GORM 的 API 设计采用链式操作的风格，使得开发者可以通过链式调用方法来构建复杂的数据库查询和操作。
事务支持： GORM 提供了事务处理功能，允许开发者对数据库操作进行事务管理，保证数据的一致性和完整性。
Callback： GORM 提供了丰富的回调函数，允许开发者在模型生命周期的各个阶段插入自定义逻辑，实现更灵活的业务处理。
// Delete - 删除 product,并没有执行delete语句，逻辑删除
db.Delete(&product, 1)

#### 逻辑删除：
定义： 逻辑删除是指通过修改数据的状态或添加额外的标记来表示数据已被删除，而实际上并未从数据库中移除。
通常会在表中增加一个表示删除状态的字段，比如一个名为 deleted_at 的时间戳字段，用于标记记录的删除时间。
优点： 逻辑删除保留了数据的历史记录，可以方便地进行数据恢复或者审计，同时避免了数据丢失。
缺点： 需要额外的字段来表示删除状态，可能会增加数据库存储空间的消耗，并且需要在查询数据时考虑过滤已删除的记录。
#### 物理删除：
定义： 物理删除是指直接从数据库中删除数据记录，彻底清除数据，使其不再存在于数据库中。
优点： 释放了数据库存储空间，减少了数据库的存储压力。
缺点： 数据一旦被删除就无法恢复，可能会导致数据丢失，不利于数据审计和追溯。


#### 更新0值

Code       sql.NullString // 商品编码
db.Create(&Product{Code: sql.NullString{"D42", true}, Price: 100})
db.Model(&product).Updates(Product{Price: 200, Code: sql.NullString{"", true}})

#### 模型mode
GORM 通过将 Go 结构体（Go structs） 映射到数据库表来简化数据库交互。 了解如何在GORM中定义模型，是充分利用GORM全部功能的基础。
type Product struct {
gorm.Model                // GORM 提供的公共模型，包含 ID、CreatedAt、UpdatedAt、DeletedAt 字段
Code       sql.NullString // 商品编码
Price      uint           // 商品价格
}

##### // gorm.Model 的定义

type Model struct {
ID        uint           `gorm:"primaryKey"`
CreatedAt time.Time
UpdatedAt time.Time
DeletedAt gorm.DeletedAt `gorm:"index"`
}

使用 GORM Migrator 创建表时，不会创建被忽略的字段
update语句可以更新0值，updates语句不可以

#### 解决进更新非零值字段的方法有两种

1.将string设置为*string
2.使用sql的NULxxx来解决

# 第五天

## 第七周

### gin学习

type Person struct {
ID   int    `uri:"id" binding:"required"`    // ID字段从URI中的id参数获取值，并且是必需的
Name string `uri:"name" binding:"required"`  // Name字段从URI中的name参数获取值，并且是必需的
}

if err := c.ShouldBindUri(&person); err != nil {
// 如果绑定失败，返回404状态码
c.Status(404)
return
}

// 处理/post请求的处理函数
func getPost(c *gin.Context) {
// 从查询参数中获取"id"
id := c.Query("id")
// 从查询参数中获取"page"，如果没有提供则默认值为"0"
page := c.DefaultQuery("page", "0")
// 从表单参数中获取"name"
name := c.PostForm("name")
// 从表单参数中获取"message"，如果没有提供则默认值为"信息"
message := c.DefaultPostForm("message", "信息")

ShouldBind 是 gin-gonic/gin 框架中用于绑定请求参数到结构体的一个方法。
它会根据请求的 Content-Type 自动选择合适的绑定方式，将请求参数绑定到结构体中。
ShouldBind 系列方法有多个变种（例如 ShouldBindJSON, ShouldBindQuery, ShouldBindForm 等），每种方法适用于不同类型的请求参数。

validator 是一个 Go 语言的包，用于对结构体中的字段进行验证。它是功能强大的验证库，允许开发者使用标签在结构体字段上定义验证规则。
validator 可以用于验证各种数据类型，并且支持自定义验证规则和国际化翻译。

# 第六天
## 用户相关的微服务
密码加密不用对称加密，因为密码不可以反解，所以采用不对称加密。采用md5算法是一个信息摘要算法
密码如果不可反解，如何找回密码。

### MD5
md5一个信息摘要算法，是一种常用的散列函数，用于生成信息摘要或哈希值。其主要功能是将任意长度的输入（如文件或文本）通过一种不可逆的数学转换，
输出为固定长度（通常是128位，即32个十六进制字符）的哈希值。

#### 用途：

数据完整性验证：在数据传输或存储过程中，通过比对MD5哈希值，可以检测数据是否被篡改。
数字签名：结合公钥加密算法，MD5可以用于数字签名，确保消息的完整性和真实性。
密码存储：尽管不推荐，但有时MD5用于存储用户密码的哈希值（应结合盐值来增加安全性）。

#### 特性

压缩性：任意长度的数据算出的MD5值的长度都是固定的。
容易计算：从源数据计算出MD5值很容易
抗修改性：对源数据进行任何修改，哪怕一个字节，MD5值差异很大。
强碰撞：想找两个不同数据，使他们具有相同的MD5值，非常困难
不可逆性，不可反解

### MD5盐值加密

#### 1.加盐

1.通过生成随机数和MD5生成字符串进行组合
2.数据库同时存储MD5值和salt值，验证正确性使用salt进行MD5即可

##### //举例

func genMd5(code string) string {
// 创建一个新的MD5哈希对象
Md5 := md5.New()
// 将输入字符串写入到MD5哈希对象中进行处理
// io.WriteString 返回两个值（写入的字节数和错误信息），这里用下划线 "_" 忽略它们
_, _ = io.WriteString(Md5, code)
// 计算输入字符串的MD5哈希值，并将其转换为十六进制字符串后返回
return hex.EncodeToString(Md5.Sum(nil))
}

//MD5 彩虹表是一种预计算的哈希值查找表，用于通过反向查找快速破解 MD5 哈希值。为了防止这样的攻击，常见的做法是使用 "盐"（salt）来增强哈希算法的安全性。
fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd) 这一行代码中，$ 符号的作用是作为分隔符，用于将不同部分的字符串连接起来。
具体来说，这个字符串格式化操作的目的是生成一个符合某种规范的加密密码字符串。

// DB 定义了一个包含数据库操作相关信息的结构体
type DB struct {
Config       *Config   // 指向数据库配置信息的指针
Error        error     // 表示数据库操作过程中的错误，如果为 nil 则表示操作成功
RowsAffected int64     // 表示数据库操作影响的行数，如插入、更新、删除操作影响的行数
Statement    *Statement // 指向数据库语句的指针，用于执行特定的 SQL 语句或事务操作
clone        int       // 克隆次数或其他用途的整数字段，具体用途需要根据上下文进一步确认
}
第八周over

## 学习总结:

今天实现了用户的grpc服务，首先将数据库的数据库连接初始化和管理代码放在一个独立的包global里面，

### 为什么要放在一个包里面？
全局可访问性和单例模式：通过将数据库连接对象 DB 定义为包级别的变量，可以确保在整个应用程序中只有一个数据库连接实例。
这种单例模式有助于避免在多个地方重复创建连接，确保数据库资源的有效利用和一致性管理。
解耦和模块化：将数据库连接代码独立为一个包，有助于将数据库操作与业务逻辑分离
统一的初始化和配置：在包的 init() 函数中进行数据库连接的初始化，可以确保在应用程序启动时就完成数据库的连接和配置。
这种方式能够集中处理数据库的连接参数、日志记录设置等，使得管理和调整数据库配置变得更加方便。
跨包访问：定义全局的数据库连接对象允许其他包或模块可以轻松地访问数据库服务。
这对于大型应用程序或者微服务架构尤为重要，不同的服务或模块可以共享同一个数据库连接实例，实现数据的一致性和统一管理。

随后定义了protobuf,定义了用户的常见服务
service User {
// 定义 User 服务
rpc GetUserList(PageInfo) returns (UserListResponse); // 获取用户列表
rpc GetUserByMobile(MobileRequest) returns (UserInfoResponse); // 通过手机号查询用户
rpc GetUserById(IdRequest) returns (UserInfoResponse); // 通过 ID 查询用户
rpc CreateUser(CreateUserInfo) returns (UserInfoResponse); // 添加用户
rpc UpdateUser(UpdateUserInfo) returns (google.protobuf.Empty); // 更新用户
rpc CheckPassWord(PasswordCheckInfo) returns (CheckResponse); // 检查密码
}
学习了 MD5盐值加密

### 随后建立了一个hander包
主要是实现了用户的处理函数
实现了ModelToRsponse 用于将数据库中的用户模型转换为 gRPC 的用户响应消息
数据库中的用户模型转换为 gRPC 的用户响应消息是为了适配不同层次的需求，并确保在不同系统组件之间传递数据时的一致性和安全性
实现了分页函数
Paginate 函数的作用是为了将分页逻辑封装在一个可复用的函数中，以便在数据库查询中轻松设置和应用分页参数，提高代码的复用性、灵活性和可维护性。
实现user服务的函数

### 定义了一个model包
这段代码定义了一个名为 model 的包，主要用于定义应用程序中的数据模型，特别是用户数据模型 User 和基础模型 BaseModel
是一个基模，代表用户的信息