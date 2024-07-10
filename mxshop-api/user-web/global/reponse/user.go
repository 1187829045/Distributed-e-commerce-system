package reponse

import (
	"fmt"
	"time"
)

// JsonTime 自定义时间类型，用于 JSON 序列化
type JsonTime time.Time

// MarshalJSON 实现了 JsonTime 类型的 JSON 序列化方法
func (j JsonTime) MarshalJSON() ([]byte, error) {
	// 将时间格式化为 "2006-01-02" 的字符串，并加上双引号，然后转换为字节数组返回
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	return []byte(stmp), nil
}

// UserResponse 用户响应结构体，用于返回用户信息的 JSON 序列化
type UserResponse struct {
	Id       int32    `json:"id"`
	NickName string   `json:"name"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}
