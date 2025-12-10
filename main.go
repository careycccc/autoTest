package main

import (
	"autoTest/PressureMeasurementModule/boomer"
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
	//invitationcarousel.RunSpinInvitedWheelWork() // 邀请转盘
	// invitationcarousel.RunSpinInvitedWheel() // 当前用户邀请转盘自动邀请下级
	// invitecode.RunInvite() // 多级下级邀请
	// withdrawcash.RunWithDrawCase() // 提现
	// topup.RunRechargeGoods()
	// GameBetOrders.RunGameBetOrders()              // 后台投注订单的查询
	// platformreports.RunGetPlatRptStatisticPageList() // 平台报表-> 每日汇总报表
	// chickenroadgame.RunChickenRoadGame() // 鸡路小游戏
	// withdrawalorders.RunWithDrawCase() // 查询订单数据
	// withdrawalorders.RunWithLockDrawCase() // 查询提现锁定的数据
	// withdrawaltimeoutcompensation.GetYesterdayUserAmont()     //获取昨天的提现超时赔付的账变的账号以及金额
	// withdrawaltimeoutcompensation.ExecelWithdrawHistoryInfo() // 提现超时赔付
	// pressuremeasurementmodule.RunPressusreModle() // 运行压测
	//accounts.RunWirteCsv() // 会员列表的数据写入到csv中
	// dailycheckin.RunDailyCheckInActivity() // 每日签到活动
	//checksign.RunCheckSign()
	boomer.RunTasks()
}
