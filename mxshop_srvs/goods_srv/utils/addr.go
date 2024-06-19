package utils

import (
	"net"
)

// GetFreePort 获取一个空闲的 TCP 端口号。
func GetFreePort() (int, error) {
	// 解析本地地址和端口，端口号为 0 表示系统随机分配一个空闲端口。
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	// 在指定地址上监听 TCP 连接。
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close() // 延迟关闭监听器，确保函数返回时关闭连接。

	// 返回监听的地址的端口号。
	return l.Addr().(*net.TCPAddr).Port, nil
}
