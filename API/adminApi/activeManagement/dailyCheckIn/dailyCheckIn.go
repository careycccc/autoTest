package dailycheckin

import (
	rechargeorders "autoTest/API/adminApi/financialManagement/rechargeOrders"
	"autoTest/API/adminApi/login"
	memberreports "autoTest/API/adminApi/reportManagement/memberReports"
	"autoTest/API/utils"
	"autoTest/PressureMeasurementModule/accounts"
	"autoTest/store/config"
	"autoTest/store/logger"
	"context"
	"fmt"
)

// 每日签到活动

// 每个用户的当日签到信息
type UserDailyCheckInInfo struct {
	UserAccount        string  `json:"userAccount"`        // 会员账号
	UserId             int     `json:"userId"`             // 会员ID
	VipLevel           int     `json:"vipLevel"`           // 会员当前VIP等级
	RechargeAmount     float64 `json:"rechargeAmount"`     // 当天充值金额
	CheckinNumberDay   int     `json:"checkinNumberDay"`   // 连续签到天数
	CheckinAward       float64 `json:"checkinAward"`       // 当天签到奖励金额
	IsAutoReceiveAward bool    `json:"isAutoReceiveAward"` // 是否自动领取奖励
	IsBlacklist        bool    `json:"isBlacklist"`        // 是否黑名单
	IsDetailsPage      bool    `json:"isDetailsPage"`      // 是否进入活动详情页
	AlarmInformation   string  `json:"alarmInformation"`   // 警告信息
}

// // 活动配置项
// type CheckInInfo struct {
// 	ActiveName           string  `json:"activeName"`     // 活动名称
// 	ActiveType           string  `json:"activeType"`     // 活动类型
// 	ShowObject           []int   `json:"showObject"`     // 活动展示对象
// 	AcitveRechargeAmount float64 `json:"rechargeAmount"` // 活动充值金额
// }

var (
	//充值人数(有充值行为的人)
	RechargeNumber int
	// 完成充值任务的人数
	CompleteRechargeTaskNumber int
	//每日登录人数
	LoginNumber int
	// 后台登录的ctx
	AdminCtx *context.Context
	// 每日用户的签到信息列表
	UserDailyCheckInInfoList []UserDailyCheckInInfo
	// 昨日用户的签到信息列表
	LastDayUserDailyCheckInInfoList []UserDailyCheckInInfo
	// 首次生成账号的csv文件路径
	CSVADDR string = "./API/adminApi/activeManagement/dailyCheckIn/initUserAccount.csv"
	// 配置项的信息 活动名称
	ActiveName string
	// 活动类型
	ActiveType string
	// 活动展示对象
	ShowObject []int
	// 自动签到的人
	AutoCheckInNumber []int
	// 活动充值金额
	AcitveRechargeAmount float64
	// 当日触达人数
	TouchNumber []int
	// 完成充值任务的人数
	ActiveRechargeNumber []int
	// 派发总奖励的金额
	AwardAmount float64
	// 手动领取人数
	ManualReceiveNumber []int
	//  手动领取的金额
	ManualReceiveAmount float64
)

// 获取当期活动配置
func GetActiveInfo() {
	ActiveName = "每日签到"
	ActiveType = "每日签到"
	ShowObject = []int{1, 2, 3}
	AcitveRechargeAmount = 100

}

func init() {
	// 后台登录
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		logger.Logger.Warn("每日签到后台登录失败", err)
		return
	} else {
		AdminCtx = ctx
	}
}

// 初始化 有充值行为的人 完成充值任务的人数 每日登录人数  主要是数据报表的查看和对比的时候需要执行
func GetDailyCheckInInfo() {
	// 获取每日登录人数
	if _, number, err := memberreports.GetUserLoginLogPageListApi(AdminCtx, config.StartTime, config.EndTime); err != nil {
		logger.LogError("获取每日签到活动的登录人数失败", err)
		return
	} else {
		LoginNumber = number
	}
	// 获取每日充值人数
	_, startTime, _, endTime, _ := utils.ParseTimeRangeToTimestamp(config.StartTime, config.EndTime)
	if _, rechargeOrderResponse, err := rechargeorders.GetRechargeOrderPageListApi(AdminCtx, startTime, endTime); err != nil {
		logger.LogError("获取每日签到活动的充值人数失败", err)
		return
	} else {
		list := rechargeOrderResponse.Data.List
		summarList := SummarizeOrders(list)
		// 去重充值会员ID  经过合并就已经去重了
		RechargeNumber = len(summarList) // 充值人数(有充值行为的人)
		CompleteRechargeTaskList := make([]SummaryItem, 0, len(summarList))
		// 获取每日完成充值任务的人数
		for _, item := range summarList {
			if item.TotalActualAmount >= 1000 {
				CompleteRechargeTaskList = append(CompleteRechargeTaskList, item)
			}
		}
		CompleteRechargeTaskNumber = len(CompleteRechargeTaskList) // 完成充值任务的人数
	}
}

// 每日的执行，从csv中读取数据，进行登录，进行充值，进行签到，或者加入黑名单
func PrepareDataByCsv() {
	// 从csv中读取数据
	if amountlist, err := accounts.ProcessCSVConcurrently(CSVADDR, 4); err != nil {
		logger.Logger.Error("从csv中读取数据失败", err)
		return
	} else {
		UserDailyCheckInInfoList = make([]UserDailyCheckInInfo, 0, len(amountlist))
		LastDayUserDailyCheckInInfoList = make([]UserDailyCheckInInfo, 0, len(amountlist))
		// 需要上一天的数据恢复
		LastDayUserDailyCheckInInfoList = RecoverYesterdayData()
		// 进行登录，进行充值，进行签到，或者加入黑名单
		UserDailyCheckInInfoList = ExcelEverDayCheckIn(amountlist)
		ReportStatistics(UserDailyCheckInInfoList)
		// 保存今日所有参与的数据
		SaveToDayData(UserDailyCheckInInfoList)
	}
}

func RunDailyCheckInActivity() {
	//PrepareData() // 只能运行一次，把初始化的用户账号写入到csv中
	// fmt.Println("每日签到活动开始运行...", CompleteRechargeTaskNumber)
	PrepareDataByCsv()
}

// 报表统计
func ReportStatistics(userinfo []UserDailyCheckInInfo) {
	if len(userinfo) == 0 {
		logger.Logger.Warn("每日用户的签到信息列表数据为空")
		return
	}
	for _, info := range userinfo {
		if info.IsDetailsPage {
			TouchNumber = append(TouchNumber, info.UserId) // 触达人数统计
		}
		if info.RechargeAmount >= AcitveRechargeAmount {
			ActiveRechargeNumber = append(ActiveRechargeNumber, info.UserId) // 完成充值任务人数统计
		}
		// 统计今日总的金额
		AwardAmount += info.RechargeAmount
		// 手动领取人数统计
		if info.IsAutoReceiveAward {
			ManualReceiveNumber = append(ManualReceiveNumber, info.UserId)
			// 手动领取的金额
			ManualReceiveAmount += info.RechargeAmount
		}
	}
	touchNumber := len(TouchNumber)
	fmt.Printf("当日触达人数:%d,触达率:%.2f\n", touchNumber, float64(touchNumber/LoginNumber)*100)
	fmt.Printf("当日参与人数:%d,参与率:%.2f\n", RechargeNumber, float64(RechargeNumber/touchNumber)*100)
	activeRechargeNumber := len(ActiveRechargeNumber)
	fmt.Printf("当日完成人数：%d,完成率:%.2f\n", activeRechargeNumber, float64(activeRechargeNumber/RechargeNumber)*100)
	fmt.Printf("当日派发奖励金额:%f\n", AwardAmount)
	manualReceiveNumber := len(ManualReceiveNumber)
	fmt.Printf("当日手动领取人数:%d,手动领取率:%.2f,手动领取金额:%.2f,手动领取成本率:%.2f\n", manualReceiveNumber, float64(manualReceiveNumber/activeRechargeNumber)*100, ManualReceiveAmount, float64(ManualReceiveAmount/AwardAmount)*100)
}
