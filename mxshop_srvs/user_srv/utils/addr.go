package utils

import (
	"net"
)

// GetFreePort 函数用于获取一个空闲的端口号。
// 返回获取的空闲端口号和可能发生的错误。
func GetFreePort() (int, error) {
	// 使用 net 包进行操作

	// 解析地址 "localhost:0"，0 表示让系统自动分配一个未使用的端口号
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	// 在指定地址上监听 TCP 连接
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close() // 延迟关闭监听

	// 返回监听的地址的端口号作为空闲端口号
	return l.Addr().(*net.TCPAddr).Port, nil
}
