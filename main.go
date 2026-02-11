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
	// ctxToken, err := login.RunAdminSitLogin()
	// if err != nil {
	// 	logger.LogError("登录失败", err)
	// 	return
	// } else {
	// 	fmt.Println("-------", ctxToken)
	// }
	invitationcarousel.RunSpinInvitedWheelWork() // 邀请转盘
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
	// accounts.RunWirteCsv() // 会员列表的数据写入到csv中
	// dailycheckin.RunDailyCheckInActivity() // 每日签到活动
	// checksign.RunCheckSign()
	// f1.RunF1()
	// financialmanagement.RunArtificialRechargeFunc() // 人工充值
	// boomer.RunTasks()
	// everydayCheckin.RunEverydayCheckIn()    // 每日签到活动
	// dailycheckin.GetReceiveRecord() // 每日签到活动报表查看
	// dailycheckin.RunDailyCheckInActivity() // 每日签到活动
	// rechargewheel.RunAddWheel()  // 一键为后台配置4个充值转盘的配置
	// common.RandUserTopupGame() // 随机从用户列表中选择用户进行充值，投注
	// common.ToupBet()

	// 验证邮箱注册
	// registerapi.RunEmailregeister()
	// uidList := utils.RandmoUserId(100)

	// for i := 0; i < len(uidList); i++ {
	// 	registerapi.GeneralAgentRegister1(uidList[i])
	// }
	// for i := 0; i < 1000; i++ {
	// 	registerapi.GeneralAgentRegister1("915645645874")
	// }
}
