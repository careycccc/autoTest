package main

import (
	invitationcarousel "autoTest/API/deskApi/invitationCarousel"
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
	invitationcarousel.RunSpinInvitedWheelWork() // 邀请转盘
	//invitecode.RunInvite()  // 多级下级邀请
	//excel.RunExcel()
	//addwallet.RunAddWallet()
	// withdrawcash.RunWithDrawCase()
	// withdrawalorders.RunWithdraw(2441424, "BankCard", 1211, 1211)
	// vip.RunVip()
}
