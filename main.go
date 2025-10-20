package main

import (
	invitecode "autoTest/API/deskApi/inviteCode"
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

	// ctx := context.Background()
	// _, ctxToken, err := login.AdminSitLogin(&ctx)
	// if err != nil {
	// 	logger.LogError("登录失败", err)
	// 	return
	// }
	invitecode.RunInvite()
	//excel.RunExcel()

}
