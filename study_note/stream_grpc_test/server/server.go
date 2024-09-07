package main

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"sale_master/study_note/stream_grpc_test/proto"
	"sync"
	"time"
)

const PORT = ":50052"

type Server struct {
	proto.UnimplementedGreeterServer // 嵌入未实现的GreeterServer，确保我们遵循接口的所有方法
}

func (s *Server) GetStream(req *proto.StreamReqData, res proto.Greeter_GetStreamServer) error {
	i := 0
	for {
		i++
		_ = res.Send(&proto.StreamResData{
			Data: fmt.Sprintf("%v", time.Now().Unix()),
		})
		if i > 10 {
			fmt.Println("i>10,over")
			break
		}
	}
	return nil
}

func (s *Server) PutStream(client proto.Greeter_PutStreamServer) error {
	for {
		if a, err := client.Recv(); err != nil {
			fmt.Println(a.Data)
			break
		} else {
			fmt.Println(a.Data)
		}
	}
	return nil
}

// AllStream 方法需要实现 GreeterServer 接口的 AllStream 方法
func (s *Server) AllStream(server proto.Greeter_AllStreamServer) error {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			data, _ := server.Recv()
			fmt.Println("收到客户端消息:" + data.Data)
		}
	}()
	go func() {
		defer wg.Done()
		for {
			server.Send(&proto.StreamResData{Data: "我是服务器"})
			time.Sleep(time.Second)
		}
	}()
	wg.Wait()
	return nil
}

func main() {
	lis, err := net.Listen("tcp", PORT)
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &Server{})
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}
