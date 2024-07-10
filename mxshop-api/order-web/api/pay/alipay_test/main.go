package main

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
)

func main() {
	// 创建支付宝客户端
	// 参数分别为：支付宝应用的 App ID、应用私钥、是否为生产环境（false表示沙箱环境）
	client, err := alipay.New("your-app-id", "your-private-key", false) //自己生产的私钥
	if err != nil {
		fmt.Println("Error creating Alipay client:", err)
		return
	}

	// 设置支付宝的公钥，用于验证支付宝返回的数据的签名
	err = client.LoadAliPayPublicKey("your-alipay-public-key") //支付宝生成的公钥
	if err != nil {
		fmt.Println("Error loading Alipay public key:", err)
		return
	}

	// 构造WAP支付请求
	var p = alipay.TradeWapPay{}
	// 支付结果异步通知URL
	p.NotifyURL = "http://example.com/notify" //哲云西生成的
	// 支付完成后返回的URL
	p.ReturnURL = "http://example.com/return" // 127.0.0.1
	// 订单标题
	p.Subject = "购买商品"
	// 商户订单号，需保证在商户端唯一
	p.OutTradeNo = "order-001"
	// 订单总金额，以元为单位，必须精确到小数点后两位
	p.TotalAmount = "10.00"
	// 产品码，QUICK_WAP_WAY表示移动网页支付
	p.ProductCode = "QUICK_WAP_WAY"

	// 生成支付URL
	payURL, err := client.TradeWapPay(p)
	if err != nil {
		fmt.Println("Error generating Alipay WAP URL:", err)
		return
	}

	// 将支付URL转换为字符串并输出
	finalURL := payURL.String()
	fmt.Println("Generated Alipay WAP payment URL:", finalURL)
}
