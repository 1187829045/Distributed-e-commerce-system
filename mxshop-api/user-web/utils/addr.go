package utils

import (
	"net"
)

// GetFreePort 查找本地机器上的一个可用端口并返回。
// 如果成功，返回端口号；如果出错，返回错误信息。
func GetFreePort() (int, error) {
	// 使用 "tcp" 网络类型和地址 "localhost:0" 解析一个 TCP 地址。
	// 其中 "0" 作为端口号，告诉系统查找一个可用端口。
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	//0xc00001c4b0
	if err != nil {
		return 0, err // 如果解析地址失败，返回 0 和错误信息。
	}

	// 在解析的 TCP 地址上开始监听。
	// 由于地址中指定了 "0" 端口，这将查找到一个可用端口。
	//*net.TCPListener 类型的指针和一个错误值
	l, err := net.ListenTCP("tcp", addr)
	//0xc000066540
	if err != nil {
		return 0, err // 如果监听地址失败，返回 0 和错误信息。
	}
	defer l.Close() // 确保在函数退出时关闭监听器，防止资源泄漏。

	// 从监听器的地址中获取端口号并返回。
	//55520
	return l.Addr().(*net.TCPAddr).Port, nil
}
