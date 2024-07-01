package model

import "time"

//购物车的表，将商品加入购物车

type ShoppingCart struct {
	BaseModel
	User    int32 `gorm:"type:int;index"` //在购物车列表中我们需要查询当前用户的购物车记录
	Goods   int32 `gorm:"type:int;index"` //索引不是越多越好，加索引：我们需要查询时候， 会带来负面问题，1. 会影响插入性能 2. 会占用磁盘
	Nums    int32 `gorm:"type:int"`       //商品数量
	Checked bool  //是否选中
}

func (ShoppingCart) TableName() string {
	return "shoppingcart"
}

type OrderInfo struct {
	BaseModel

	User    int32  `gorm:"type:int;index"`                                          //根据用户查询用户订单
	OrderSn string `gorm:"type:varchar(30);index"`                                  //订单号，我们平台自己生成的订单号
	PayType string `gorm:"type:varchar(20) comment 'alipay(支付宝)， wechat(微信)'"` //便于后期查账

	//status可以考虑使用iota来做
	Status     string `gorm:"type:varchar(20)  comment 'PAYING(待支付), TRADE_SUCCESS(成功)， TRADE_CLOSED(超时关闭), WAIT_BUYER_PAY(交易创建), TRADE_FINISHED(交易结束)'"`
	TradeNo    string `gorm:"type:varchar(100) comment '交易号'"` //交易号就是支付宝的订单号 查账
	OrderMount float32
	PayTime    *time.Time `gorm:"type:datetime"`

	Address      string `gorm:"type:varchar(100)"`
	SignerName   string `gorm:"type:varchar(20)"`
	SingerMobile string `gorm:"type:varchar(11)"`
	Post         string `gorm:"type:varchar(20)"`
}

func (OrderInfo) TableName() string {
	return "orderinfo"
}

//订单商品 一个订单多个商品
// 订单表，如果有根据商品查询订单，放入上面表中，效率就很低，所以单独拿一张表

type OrderGoods struct {
	BaseModel

	Order int32 `gorm:"type:int;index"`
	Goods int32 `gorm:"type:int;index"`

	//把商品的信息保存下来了 ， 字段冗余，但是高并发系统中我们一般都不会遵循三范式 ，如果这里不保存，如果后期我要展示订单的商品，就意味着我要垮服务去
	//访问商品信息，会增加流量，如果后期根据商品的名称去查询订单，就会很麻烦先跨服务找ID，反而性能降低，为了降低跨服务所以地段冗余，
	//还有做镜像的时候，用户购买商品知道名称图片价格，所以要保存下来，到时候客户可以根据图片价格名称查询 ，做下记录

	GoodsName  string `gorm:"type:varchar(100);index"`
	GoodsImage string `gorm:"type:varchar(200)"`
	GoodsPrice float32
	Nums       int32 `gorm:"type:int"`
}

func (OrderGoods) TableName() string {
	return "ordergoods"
}
