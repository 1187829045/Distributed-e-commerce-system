package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"sale_master/study_note/stream_grpc_test/proto"
	"sync"
	"time"
)

func main() {
	//创建一个 gRPC 客户端连接
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	//服务端流模式
	//创建一个 gRPC 客户端对象，用于与 gRPC 服务器进行通信。
	c := proto.NewGreeterClient(conn)
	//调用 gRPC 客户端对象 c 的 GetStream 方法，向服务器发送请求，并获取服务器的响应
	res, _ := c.GetStream(context.Background(), &proto.StreamReqData{Data: "llb"})
	for {
		a, err := res.Recv()
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(a.Data)
	}
	//客户端流模式
	putS, _ := c.PutStream(context.Background())
	i := 0
	for {
		i++
		putS.Send(&proto.StreamReqData{
			Data: fmt.Sprintf("llb%d", i),
		})
		time.Sleep(time.Second)
		if i > 10 {
			break
		}
	}
	//////////////////////////////////////////
	//双向流模式
	server, _ := c.AllStream(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			data, _ := server.Recv()
			fmt.Println("收到服务端消息:" + data.Data)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			server.Send(&proto.StreamReqData{Data: "我是客户端"})
			time.Sleep(time.Second)
		}
	}()
	wg.Wait()
}
