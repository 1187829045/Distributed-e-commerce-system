package forms

// 创建订单的Form
type CreateOrderForm struct {
	Address string `json:"address" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Mobile  string `json:"mobile" binding:"required,mobile"`
	Post    string `json:"post" binding:"required"`
}
