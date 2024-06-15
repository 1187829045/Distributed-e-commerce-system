package iniitialize

import "go.uber.org/zap"

func InitLogger() {
	// 创建一个新的生产环境的 logger 实例
	logger, _ := zap.NewProduction()
	// 将全局 logger 替换为新创建的 logger
	//全局 logger 替换为自定义的 logger，可以确保你的应用程序中所有的日志记录操作都使用同一个 logger 实例
	zap.ReplaceGlobals(logger)
}
