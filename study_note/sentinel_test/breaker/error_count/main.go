package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/util"
)

// 定义状态变更监听器结构体
type stateChangeTestListener struct {
}

// 当熔断器从其他状态转变为关闭状态时的回调函数

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.strategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

// 当熔断器从其他状态转变为打开状态时的回调函数

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.strategy: %+v, From %s to Open, snapshot: %d, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

// 当熔断器从其他状态转变为半打开状态时的回调函数

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.strategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func main() {
	total := 0                        // 总请求数
	totalPass := 0                    // 通过的请求数
	totalBlock := 0                   // 被阻断的请求数
	totalErr := 0                     // 记录的错误数
	conf := config.NewDefaultConfig() // 创建默认配置
	// 将日志输出到控制台
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf) // 初始化 Sentinel
	if err != nil {
		log.Fatal(err) // 如果初始化失败，记录错误并退出
	}
	ch := make(chan struct{}) // 创建一个阻塞的通道
	// 注册状态变更监听器，用于观察内部熔断器的状态变化
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	// 加载熔断规则
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// 统计时间窗口=5秒，恢复超时时间=3秒，最大错误计数=50
		{
			Resource:         "abc",                     // 资源名
			Strategy:         circuitbreaker.ErrorCount, // 错误计数策略
			RetryTimeoutMs:   3000,                      // 3秒后尝试恢复
			MinRequestAmount: 10,                        // 静默请求数量
			StatIntervalMs:   5000,                      // 统计时间窗口
			Threshold:        50,                        // 错误阈值
		},
	})
	if err != nil {
		log.Fatal(err) // 如果规则加载失败，记录错误并退出
	}

	// 输出日志信息，表示程序正在运行
	logging.Info("[CircuitBreaker ErrorCount] Sentinel Go circuit breaking demo is running. You may see the pass/block metric in the metric log.")

	// 第一个模拟的服务调用
	go func() {
		for {
			total++
			e, b := sentinel.Entry("abc") // 尝试进入资源
			if b != nil {
				// 请求被阻断
				totalBlock++
				fmt.Println("协程熔断了")
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond) // 随机等待一段时间
			} else {
				// 请求通过
				totalPass++
				if rand.Uint64()%20 > 9 {
					totalErr++
					// 记录当前调用为错误
					sentinel.TraceError(e, errors.New("biz error"))
				}
				time.Sleep(time.Duration(rand.Uint64()%20+10) * time.Millisecond) // 随机等待一段时间
				e.Exit()                                                          // 退出资源
			}
		}
	}()

	// 第二个模拟的服务调用
	go func() {
		for {
			total++
			e, b := sentinel.Entry("abc") // 尝试进入资源
			if b != nil {
				// 请求被阻断
				totalBlock++
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond) // 随机等待一段时间
			} else {
				// 请求通过
				totalPass++
				time.Sleep(time.Duration(rand.Uint64()%80) * time.Millisecond) // 随机等待一段时间
				e.Exit()                                                       // 退出资源
			}
		}
	}()

	// 定期输出错误总数
	go func() {
		for {
			time.Sleep(time.Second) // 每秒输出一次
			fmt.Println(totalErr)   // 输出错误总数
		}
	}()
	<-ch // 阻塞主线程，防止程序退出
}
