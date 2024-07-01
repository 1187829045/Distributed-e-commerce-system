package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

//普通消息发送

func main() {
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.128.128:9876"}))
	if err != nil {
		panic("NewProducer err")
	}
	if err = p.Start(); err != nil {
		panic("start err")
	}
	res, err := p.SendSync(context.Background(), primitive.NewMessage("llb-test", []byte("this is llb")))
	if err != nil {
		panic("SendSync err")
	}
	fmt.Println("res:", res)
	if err = p.Shutdown(); err != nil {
		panic("Shutdown err")
	}
}
