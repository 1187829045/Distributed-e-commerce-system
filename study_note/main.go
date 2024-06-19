package main

import (
	"fmt"
	"time"
)

func main() {

	timestamp := int64(1624067115) // 一个Unix时间戳
	fmt.Println(timestamp)
	t := time.Unix(timestamp, 0) // 将Unix时间戳转换为时间对象
	fmt.Println(t)
}
