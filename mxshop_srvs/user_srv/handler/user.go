package handler

import (
	"context"                                   // 导入用于处理上下文的包
	"crypto/sha512"                             // 导入用于密码加密的哈希算法包
	"fmt"                                       // 导入格式化 I/O 包
	"github.com/anaskhan96/go-password-encoder" // 导入用于密码加密的包
	"github.com/golang/protobuf/ptypes/empty"   // 导入 gRPC 的空消息类型
	"google.golang.org/grpc/codes"              // 导入 gRPC 状态码包
	"google.golang.org/grpc/status"             // 导入 gRPC 错误状态包
	"gorm.io/gorm"                              // 导入 GORM ORM 库
	"sale_master/mxshop_srvs/user_srv/global"   // 导入全局变量
	"sale_master/mxshop_srvs/user_srv/model"    // 导入数据模型
	"sale_master/mxshop_srvs/user_srv/proto"    // 导入 gRPC 协议文件
	"strings"                                   // 导入字符串处理包
	"time"                                      // 导入时间处理包
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
// 分页参数是用于数据库查询中控制返回结果数量和位置的设置，通常包括两个主要的参数：
// 页码（Page Number）：
// 页码表示要获取的数据位于结果集中的哪一页。通常从 1 开始计数，即第一页。
// 每页大小（Page Size）：
// 每页大小指定每次查询返回的数据条数。例如，每页大小为 10 表示每次查询返回最多 10 条数据。
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
	// Scopes 将指定的作用域函数（func）应用于当前数据库操作对象，并返回一个新的数据库操作对象（tx）。
	// 这些作用域函数允许在查询或事务中应用一系列预定义的条件或操作。
	//func (db *DB) Scopes(funcs ...func(*DB) *DB) (tx *DB) {
	//	 获取当前数据库操作对象的克隆实例
	//	tx = db.getInstance()
	//	 将传入的作用域函数（funcs）追加到语句（Statement）的作用域列表（scopes）中
	//	tx.Statement.scopes = append(tx.Statement.scopes, funcs...)
	//	 返回应用作用域函数后的新数据库操作对象（tx）
	//	return tx
	//}
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

	// 密码加密
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
