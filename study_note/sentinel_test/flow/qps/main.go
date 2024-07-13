package main

import (
	"fmt"
	"log"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
)

// QPS 限流（Queries Per Second 限流）是指对系统在一秒钟内处理的请求数量进行限制，以确保系统的稳定性和性能
func main() {
	// 初始化 Sentinel 默认配置
	err := sentinel.InitDefault()
	if err != nil {
		// 如果初始化失败，记录日志并终止程序
		log.Fatalf("初始化sentinel 异常: %v", err)
	}

	// 配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test",     // 资源名为 "some-test"
			TokenCalculateStrategy: flow.Direct,     // 直接计算策略
			ControlBehavior:        flow.Throttling, // 匀速通过策略
			Threshold:              100,             // 阈值为 100
			StatIntervalInMs:       1000,            // 统计时间间隔为 1000 毫秒
		},
		{
			Resource:               "some-test2", // 资源名为 "some-test2"
			TokenCalculateStrategy: flow.Direct,  // 直接计算策略
			ControlBehavior:        flow.Reject,  // 直接拒绝策略
			Threshold:              10,           // 阈值为 10
			StatIntervalInMs:       1000,         // 统计时间间隔为 1000 毫秒
		},
	})

	// 如果加载规则失败，记录日志并终止程序
	if err != nil {
		log.Fatalf("加载规则失败: %v", err)
	}

	// 循环 12 次，模拟流量
	for i := 0; i < 12; i++ {
		// 尝试进入资源 "some-test"
		e, b := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
		if b != nil {
			// 如果被限流，打印限流信息
			fmt.Println("限流了")
		} else {
			// 如果通过，打印检查通过信息并退出资源
			fmt.Println("检查通过")
			e.Exit()
		}
		// 每次循环后等待 11 毫秒
		time.Sleep(11 * time.Millisecond)
	}
}
