package main

import (
	"context" // 导入上下文包，用于控制多个goroutine之间的信号传递
	"fmt"     // 导入格式化包，用于格式化输出

	"github.com/apache/rocketmq-client-go/v2"           // 导入RocketMQ客户端包
	"github.com/apache/rocketmq-client-go/v2/primitive" // 导入RocketMQ原始数据类型包
	"github.com/apache/rocketmq-client-go/v2/producer"  // 导入RocketMQ生产者包
)

func main() {
	// 创建一个新的Producer（生产者），设置NameServer地址
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.128.128:9876"}))
	if err != nil {
		panic("生成producer失败") // 如果创建生产者失败，终止程序并打印错误信息
	}

	// 启动生产者
	if err = p.Start(); err != nil {
		panic("启动producer失败") // 如果启动生产者失败，终止程序并打印错误信息
	}

	// 创建一条新的消息，设置主题为llb1，消息内容为"this is delay message"
	msg := primitive.NewMessage("llb1", []byte("this is delay message"))
	// 设置消息的延迟时间级别
	msg.WithDelayTimeLevel(3)
	// 同步发送消息，并返回结果
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		// 如果发送失败，打印错误信息
		fmt.Printf("发送失败: %s\n", err)
	} else {
		// 如果发送成功，打印发送结果
		fmt.Printf("发送成功: %s\n", res.String())
	}

	// 关闭生产者
	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败") // 如果关闭生产者失败，终止程序并打印错误信息
	}

	// 在支付场景中，例如淘宝或12306购票，超时归还需要定时执行某些逻辑
	// 可以通过轮询的方式实现，但轮询有一些问题：
	// 1. 需要确定轮询的频率，例如每30分钟执行一次
	// 2. 在12:00执行过一次，下一次执行是在12:30，但如果用户在12:01下单，超时应该在12:31处理
	// 如果轮询频率设置为1分钟一次，订单量不大的情况下，多数查询都是无用的，而且还会频繁查询数据库
	// 使用RocketMQ的延迟消息，可以在指定时间执行，并且消息中包含了订单编号，只查询相关订单编号
}
