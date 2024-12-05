package test

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"shop_srvs/user_srv/proto"
)

var userClient proto.UserClient // 定义 gRPC 客户端
var conn *grpc.ClientConn       // gRPC 连接

func Init() {
	var err error
	// 初始化 gRPC 客户端连接
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn) // 创建用户服务的 gRPC 客户端
}

func TestGetUserList() {
	// 测试获取用户列表
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 5,
	})
	if err != nil {
		panic(err)
	}
	// 遍历用户列表响应并输出用户信息
	for _, user := range rsp.Data {
		fmt.Println(user.Mobile, user.NickName, user.PassWord)
		// 测试校验密码
		checkRsp, err := userClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          "admin123",
			EncryptedPassword: user.PassWord,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(checkRsp.Success)
	}
}

func TestCreateUser() {
	// 测试创建用户
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			NickName: fmt.Sprintf("bobby%d", i),
			Mobile:   fmt.Sprintf("1878222222%d", i),
			PassWord: "admin123",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}
}

func main() {
	Init() // 初始化 gRPC 客户端连接

	// 测试功能
	//TestGetUserList() // 测试获取用户列表
	TestCreateUser() // 测试创建用户

	conn.Close() // 关闭 gRPC 连接
}
