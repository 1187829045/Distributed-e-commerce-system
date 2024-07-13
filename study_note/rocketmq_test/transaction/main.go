package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// OrderListener 结构体用于实现事务监听器
type OrderListener struct{}

// ExecuteLocalTransaction 执行本地事务
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地逻辑")
	time.Sleep(time.Second * 3) // 模拟本地事务执行时间
	fmt.Println("执行本地逻辑失败")
	// 本地事务执行失败，例如代码异常或系统宕机
	return primitive.UnknowState // 返回未知状态，表示事务状态不确定
}

// CheckLocalTransaction 检查本地事务状态
func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("rocketmq的消息回查")
	time.Sleep(time.Second * 15)        // 模拟检查事务状态的时间
	return primitive.CommitMessageState // 返回提交消息状态，表示事务成功
}

func main() {
	// 创建一个新的事务生产者，设置NameServer地址和事务监听器
	p, err := rocketmq.NewTransactionProducer(
		&OrderListener{},
		producer.WithNameServer([]string{"192.168.128.128:9876"}),
	)
	if err != nil {
		panic("生成producer失败") // 如果创建生产者失败，终止程序并打印错误信息
	}

	// 启动生产者
	if err = p.Start(); err != nil {
		panic("启动producer失败") // 如果启动生产者失败，终止程序并打印错误信息
	}

	// 发送事务消息，并返回结果
	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("TransTopic", []byte("this is transaction message2")))
	if err != nil {
		// 如果发送失败，打印错误信息
		fmt.Printf("发送失败: %s\n", err)
	} else {
		// 如果发送成功，打印发送结果
		fmt.Printf("发送成功: %s\n", res.String())
	}

	// 为了防止主goroutine退出，使用time.Sleep让程序等待一个小时
	time.Sleep(time.Hour)

	// 关闭生产者
	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败") // 如果关闭生产者失败，终止程序并打印错误信息
	}
}
