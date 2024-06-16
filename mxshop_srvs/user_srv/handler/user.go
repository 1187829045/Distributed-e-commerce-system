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
	"sale_master/mxshop_srvs/user_srv/global"
	"sale_master/mxshop_srvs/user_srv/model"
	"sale_master/mxshop_srvs/user_srv/proto"
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
		Id:       user.ID,          // 用户 ID
		PassWord: user.Password,    // 用户密码
		NickName: user.NickName,    // 用户昵称
		Gender:   user.Gender,      // 用户性别
		Role:     int32(user.Role), // 用户角色
		Mobile:   user.Mobile,      // 用户手机号
	}
	if user.Birthday != nil {
		userInfoRsp.BirthDay = uint64(user.Birthday.Unix()) // 用户生日，转换为时间戳
	}
	return userInfoRsp
}

// Paginate 返回一个用于设置分页参数的 GORM 函数
// Paginate 函数的作用是为了将分页逻辑封装在一个可复用的函数中，以便在数据库查询中轻松设置和应用分页参数，提高代码的复用性、灵活性和可维护性。
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1 // 如果页码为 0，则设置为 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100 // 如果每页大小超过 100，则设置为 100
		case pageSize <= 0:
			pageSize = 10 // 如果每页大小小于等于 0，则设置为 10
		}
		// 偏移量的作用是告诉数据库从查询结果集的哪个位置开始返回数据，配合 Limit（限制条数），可以有效地实现分页查询。
		offset := (page - 1) * pageSize          // 计算偏移量
		return db.Offset(offset).Limit(pageSize) // 设置分页参数
	}
}

// GetUserList 获取用户列表
func (s *UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	result := global.DB.Find(&users) // 查询所有用户
	if result.Error != nil {
		return nil, result.Error // 如果查询出错，则返回错误
	}
	fmt.Println("用户列表")
	rsp := &proto.UserListResponse{}
	rsp.Total = int32(result.RowsAffected) // 设置总用户数
	//分页查询在应用中通常用于优化用户界面的数据展示和用户体验
	global.DB.Scopes(Paginate(int(req.Pn), int(req.PSize))).Find(&users) // 分页查询用户

	for _, user := range users {
		userInfoRsp := ModelToRsponse(user)
		rsp.Data = append(rsp.Data, &userInfoRsp) // 将用户信息添加到响应数据中
	}
	return rsp, nil
}

// GetUserByMobile 通过手机号码查询用户
func (s *UserServer) GetUserByMobile(ctx context.Context, req *proto.MobileRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user) // 根据手机号查询用户
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在") // 如果用户不存在，则返回错误
	}
	if result.Error != nil {
		return nil, result.Error // 如果查询出错，则返回错误
	}

	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

// GetUserById 通过 ID 查询用户
func (s *UserServer) GetUserById(ctx context.Context, req *proto.IdRequest) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.First(&user, req.Id) // 根据用户 ID 查询用户
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "用户不存在") // 如果用户不存在，则返回错误
	}
	if result.Error != nil {
		return nil, result.Error // 如果查询出错，则返回错误
	}

	userInfoRsp := ModelToRsponse(user)
	return &userInfoRsp, nil
}

// CreateUser 新建用户
func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserInfo) (*proto.UserInfoResponse, error) {
	var user model.User
	result := global.DB.Where(&model.User{Mobile: req.Mobile}).First(&user) // 根据手机号查询用户是否存在
	if result.RowsAffected == 1 {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在") // 如果用户已存在，则返回错误
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
	options := &password.Options{16, 100, 32, sha512.New}                             // 设置密码校验选项
	passwordInfo := strings.Split(req.EncryptedPassword, "$")                         // 拆分加密后的密码
	check := password.Verify(req.Password, passwordInfo[2], passwordInfo[3], options) // 校验密码
	return &proto.CheckResponse{Success: check}, nil                                  // 返回校验结果
}
