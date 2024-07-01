package model

const (
	LEAVING_MESSAGES = iota + 1
	COMPLAINT
	INQUIRY
	POST_SALE
	WANT_TO_BUY
)

//留言

type LeavingMessages struct {
	BaseModel

	User int32 `gorm:"type:int;index"`
	//留言的类型
	MessageType int32  `gorm:"type:int comment '留言类型: 1(留言),2(投诉),3(询问),4(售后),5(求购)'"`
	Subject     string `gorm:"type:varchar(100)"` //主题

	Message string //留言的内容
	File    string `gorm:"type:varchar(200)"` //文件
}

func (LeavingMessages) TableName() string {
	return "leavingmessages"
}

type Address struct {
	BaseModel

	User         int32  `gorm:"type:int;index"`
	Province     string `gorm:"type:varchar(10)"` //省
	City         string `gorm:"type:varchar(10)"` //市
	District     string `gorm:"type:varchar(20)"` //区域
	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"` //收货人名称
	SignerMobile string `gorm:"type:varchar(11)"` //收货人的电话号码
}

//收藏

type UserFav struct {
	BaseModel

	User  int32 `gorm:"type:int;index:idx_user_goods,unique"` //一个用户收藏一件商品一次加unique
	Goods int32 `gorm:"type:int;index:idx_user_goods,unique"`
}

func (UserFav) TableName() string {
	return "userfav"
}
