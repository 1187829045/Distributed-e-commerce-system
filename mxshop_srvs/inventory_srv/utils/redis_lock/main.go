package main

import (
	"fmt"
	"sync"
	"time"

	goredislib "github.com/go-redis/redis/v8"           // 引入 go-redis 库
	"github.com/go-redsync/redsync/v4"                  // 引入 redsync 库
	"github.com/go-redsync/redsync/v4/redis/goredis/v8" // 引入 redsync 的 go-redis 适配库
)

func main() {
	// 使用 go-redis 创建一个 Redis 客户端（或使用 redigo），
	// 这个客户端将用于 redsync 与 Redis 通信。
	// 此外，也可以是任何实现了 `redis.Pool` 接口的连接池。
	client := goredislib.NewClient(&goredislib.Options{
		Addr: "192.168.128.136:6379", // 配置 Redis 地址
	})
	pool := goredis.NewPool(client) // 创建一个 Redis 连接池

	// 创建一个 redsync 实例，用于获取互斥锁。
	rs := redsync.New(pool)

	// 使用相同的名称获取一个新的互斥锁，所有需要相同锁的实例都使用这个名称。
	gNum := 2          // 定义协程数量
	mutexname := "421" // 定义互斥锁名称

	var wg sync.WaitGroup // 定义一个等待组
	wg.Add(gNum)          // 设置等待组计数

	for i := 0; i < gNum; i++ {
		go func() {
			defer wg.Done()                 // 协程完成后减少等待组计数
			mutex := rs.NewMutex(mutexname) // 创建一个新的互斥锁

			fmt.Println("开始获取锁")                 // 输出开始获取锁信息
			if err := mutex.Lock(); err != nil { // 尝试获取锁
				panic(err) // 获取锁失败时抛出错误
			}

			fmt.Println("获取锁成功") // 输出获取锁成功信息

			time.Sleep(time.Second * 8) // 睡眠8秒，模拟业务逻辑处理时间

			fmt.Println("开始释放锁")                              // 输出开始释放锁信息
			if ok, err := mutex.Unlock(); !ok || err != nil { // 尝试释放锁
				panic("unlock failed") // 释放锁失败时抛出错误
			}
			fmt.Println("释放锁成功") // 输出释放锁成功信息
		}()
	}
	wg.Wait() // 等待所有协程完成
}
