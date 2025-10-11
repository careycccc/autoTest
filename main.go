package main

import (
	rechargewheel "autoTest/API/deskApi/activeGIft/rechargeWheel"
	_ "autoTest/API/deskApi/loginApi"

	_ "autoTest/store/config"
	"autoTest/store/logger"
)

func init() {
	// 初始化日志 如果需要把日志写入到yaml文件中，就调用logger.InitLogger2()
	logger.InitLogger()
	// logger.Init(config.LogLevel)
	// logger.Logger.Info("这是一个信息日志",
	// 	"key", "value",
	// )
	logger.Logger.Info("logger init sucessfully....")
	// 模拟一个错误
	// err := someFunction()
	// if err != nil {
	// 	logger.LogError("报错消息", err)
	// }
}

func main() {
	rechargewheel.RunRechargeWheelCondition()
}

// // someFunction 模拟一个返回错误的函数
// func someFunction() error {
// 	return fmt.Errorf("出现错误")
// }
