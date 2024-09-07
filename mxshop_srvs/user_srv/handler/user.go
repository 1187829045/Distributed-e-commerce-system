package handler

import (
	"context"
	"crypto/sha512" // 导入用于密码加密的哈希算法包
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"shop_srvs/user_srv/global"
	"shop_srvs/user_srv/model"
	"shop_srvs/user_srv/proto"
	"strings"
	"time"
)

// UserServer 结构体实现了 User 服务的所有方法
type UserServer struct {
	proto.UnimplementedUserServer
}

// ModelToRsponse 将数据库中的用户模型转换为 gRPC 的用户响应消息
// 数据库中的用户模型转换为 gRPC 的用户响应消息是为了适配不同层次的需求，并确保在不同系统组件之间传递数据时的一致性和安全性
func ModelToRsponse(user model.User) proto.UserInfoResponse {
	// 在 gRPC 的消息中，字段有默认值，不能随便赋值 nil，容易出错
	// 这里要搞清楚哪些字段有默认值
	userInfoRsp := proto.UserInfoResponse{
		Id:       user.ID,
		PassWord: user.Password,
		NickName: user.NickName,
		Gender:   user.Gender,
		Role:     int32(user.Role),
		Mobile:   user.Mobile,
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix()) // 用户生日，转换为时间戳
	}
	return userInfoRsp
}

// Paginate 函数的作用是为了将分页逻辑封装在一个可复用的函数中
// 从第 (page - 1) * pageSize 条记录开始读取数据。例如，第一页从第 0 条数据开始，第二页从第 pageSize 条数据开始，依此类推。
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		if pageSize < 0 {
			pageSize = 10
		}
		if pageSize > 100 {
			pageSize = 100
		}
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}
}

//获取用户列表

func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	res := &proto.UserListResponse{}
	result := global.DB.Find(&model.User{})
	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}
	fmt.Println("当前是GetUserList")
	res.Total = int32(result.RowsAffected)
	var user []model.User
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&user)
	for _, userInfo := range user {
		userRsp := ModelToRsponse(userInfo)
		res.Data = append(res.Data, &userRsp)
	}
	return res, nil
}

// GetUserByMobile 通过手机号码查询用户
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil

}

// GetUserById 通过 ID 查询用户
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id)
	if result.RowsAffected == 0 {
		fmt.Printf("用户不存在")
	}
	if result.Error != nil {
		fmt.Printf("GetUserById 通过 ID 查询用户 错误")
		return nil, result.Error
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

// CreateUser 新建用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user) // 根据手机号查询用户是否存在
	if result.RowsAffected != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName
	options := &password.Options{16, 100, 32, sha512.New}
	salt, encodedPwd := password.Encode(req.PassWord, options)
	user.Password = fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd) // 格式化加密后的密码

	result = global.DB.Create(&user) // 创建新用户
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error()) // 如果创建用户出错，则返回错误
	}
	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

// UpdateUser 更新用户信息
func (s *UserServer) UpdateUser(ctx context.Context, req *proto.UpdateUserInfo) (*empty.Empty, error) {
	var user model.User
	result := global.DB.First(&user, req.Id) // 根据用户 ID 查询用户
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在") // 如果用户不存在，则返回错误
	}

	birthDay := time.Unix(int64(req.BirthDay), 0) // 将生日时间戳转换为时间
	user.NickName = req.NickName
	user.Birthday = &birthDay
	user.Gender = req.Gender

	result = global.DB.Save(&user) // 保存用户更新信息
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error()) // 如果保存出错，则返回错误
	}
	return &empty.Empty{}, nil
}

// CheckPassWord 校验密码
func (s *UserServer) CheckPassWord(ctx context.Context, req *proto.PasswordCheckInfo) (*proto.CheckResponse, error) {
	// 设置密码校验选项
	//SaltLength（16）：盐的长度，盐是一种随机生成的数据，用来与密码混合后再进行哈希，以增加安全性。
	//Iterations（100）：迭代次数，表示哈希函数在生成加密密码时要重复计算多少次。
	//KeyLength（32）：生成的密钥长度。
	//HashFunction（sha512.New）：使用 SHA-512 哈希函数来加密密码。
	options := &password.Options{16, 100, 32, sha512.New}
	// 拆分加密后的密码
	passwordInfo := strings.Split(req.EncryptedPassword, "$")
	// 校验密码
	//passwordInfo[2] 通常是从加密密码中提取出的盐值。
	//passwordInfo[3] 通常是提取出的实际哈希值。
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options)
	// 返回校验结果
	return &proto.CheckResponse{Success: check}, nil
}
