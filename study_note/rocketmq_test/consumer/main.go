package main

import (
	"context"                                           // 导入上下文包，用于控制多个goroutine之间的信号传递
	"fmt"                                               // 导入格式化包，用于格式化输出
	"github.com/apache/rocketmq-client-go/v2"           // 导入RocketMQ客户端包
	"github.com/apache/rocketmq-client-go/v2/consumer"  // 导入RocketMQ消费包
	"github.com/apache/rocketmq-client-go/v2/primitive" // 导入RocketMQ原始数据类型包
	"time"                                              // 导入时间包，用于时间相关操作
)

func main() {
	// 创建一个新的PushConsumer（推送消费者），设置NameServer地址和消费者组名
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.128.128:9876"}), // NameServer地址
		consumer.WithGroupName("shop"),                            // 消费者组名
	)

	// 订阅主题llb1，并设置消息选择器和消息处理函数
	if err := c.Subscribe("llb1", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		// 遍历接收到的消息，并打印
		for i := range msgs {
			fmt.Printf("获取到值： %v \n", msgs[i]) // 打印消息内容
		}
		return consumer.ConsumeSuccess, nil // 返回消费成功结果
	}); err != nil {
		// 如果订阅失败，打印错误信息
		fmt.Println("读取消息失败")
	}

	// 启动消费者
	_ = c.Start()

	// 为了防止主goroutine退出，使用time.Sleep让程序等待一个小时
	time.Sleep(time.Hour)

	// 关闭消费者
	_ = c.Shutdown()
}
