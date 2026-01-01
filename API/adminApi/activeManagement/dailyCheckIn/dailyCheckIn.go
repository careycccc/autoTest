package dailycheckin

import (
	"autoTest/API/adminApi/common"
	rechargeorders "autoTest/API/adminApi/financialManagement/rechargeOrders"
	"autoTest/API/adminApi/login"
	fundtransactionrecords "autoTest/API/adminApi/reportManagement/FundTransactionRecords"
	memberreports "autoTest/API/adminApi/reportManagement/memberReports"
	"autoTest/API/deskApi/active/everydayCheckin"
	desklogin "autoTest/API/deskApi/loginApi"
	"autoTest/API/utils"
	"autoTest/PressureMeasurementModule/accounts"
	"autoTest/store/config"
	"autoTest/store/logger"
	"context"
	"fmt"
	"time"
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
	// 当日充值金额和人数
	CompleteRechargeTaskList []SummaryItem
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
	TouchNumber int = 85
	// 完成充值任务的人数
	ActiveRechargeNumber []int
	// 派发总奖励的金额
	AwardAmount float64
	// 手动领取人数
	ManualReceiveNumber []int
	//  手动领取的金额
	ManualReceiveAmount float64
	// 自动领取人数
	AutoReceiveNumber []int
	// 自动领取的金额
	AutoReceiveAmount float64
)

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
		fmt.Println("昨日有充值行为的人的去重前的数据：", len(list))
		summarList := SummarizeOrders(list)
		CompleteRechargeTaskList = summarList
		// 去重充值会员ID  经过合并就已经去重了
		RechargeNumber = len(summarList) // 充值人数(有充值行为的人)
	}
}

// 查询领取记录
func GetReceiveRecord() {
	// RunDailyCheckInUserList()
	//RunDailyCheckInActivityList()
	RunDailyCheckInUserList()
}

// 查询昨天的数据，主要是数据报表
func RunCheckinDataValidation() {
	activeId := 27 // 活动id
	activeName := ""
	//1，获取每日登录人数
	//2，获取每日充值人数，充值金额
	GetDailyCheckInInfo()
	for _, item := range CompleteRechargeTaskList {
		// 3 查询这个会员是否参与了本轮活动
		// 3.1 id转账号进行前台登录
		amount := common.IdToAmountAndUpdatePassword(AdminCtx, item.UserId)
		if amount == "" {
			logger.Logger.Warn("会员账号转换失败/或者是游客账号", item.UserId)
			continue
		}
		// 3.2 登录
		ctx, err := desklogin.ReturnContextLoginY1(amount, "qwer1234")
		if err != nil {
			return
		}
		// 1.获取用户签到信息
		res, respData, err := everydayCheckin.GetUserCheckInActivityData(ctx)
		if err != nil {
			return
		}
		if res.Msg != "Succeed" {
			// 表示这个账号没有参与过本轮活动，直接跳过
			// 记录这个账号的下标
			fmt.Println("没有参与过活动", item.UserId)
			continue
		}
		//4，获取这个会员的当前签到天数，如果这个天数=1，表示今天第一条参与，不做统计，只有是连续签到才做统计，
		if respData.Data.ActivityId == activeId {
			activeName = respData.Data.ActivityName
			// 只统计活动id为27好的
			if respData.Data.CurrentCheckInDays > 1 {
				// 获取这个会员当期的配置选项
				reward := respData.Data.RewardDetail[respData.Data.CurrentCheckInDays-2]
				// 查找这个会员的昨天的充值和任务充值金额是否满足
				if item.TotalActualAmount >= float64(reward.RechargeAmount) {
					// 满足，派发总金额的累加
					AwardAmount += float64(reward.RewardAmount)
					// 完成人数累加
					ActiveRechargeNumber = append(ActiveRechargeNumber, item.UserId)
				}
			}
		}

	}
	//5，满足了要统计一下，总共发放的金额
	//6，查询昨天的该会员的账变记录，如果有数据，说明是昨天手动领取的，剩余的就是自动发的金额
	financialTypeList := []string{"DailyCheckInReward"}
	for _, item := range ActiveRechargeNumber {
		//fmt.Println("完成充值任务的人数的id", item)
		// 查找手动派发的人数
		_, start, _, end, _ := utils.ParseTimeRangeToTimestamp(config.StartTime, config.EndTime)
		if _, fundtrans, err := fundtransactionrecords.GetFinancialTypeById(AdminCtx, item, financialTypeList, start, end); err != nil {
			logger.Logger.Error("查询账变记录失败", err)
			continue
		} else {

			if len(fundtrans.Data.List) == 0 {
				// 表示没有查询到昨日的数据，没有手动领取 // 自动领取的人数
				AutoReceiveNumber = append(AutoReceiveNumber, item)
			} else {
				ManualReceiveNumber = append(ManualReceiveNumber, item)
				ManualReceiveAmount += float64(fundtrans.Data.List[0].Amount)
			}
		}
	}
	yesterdayFishPerson := len(ActiveRechargeNumber)
	yesterdayManualReceiveNumber := len(ManualReceiveNumber)
	yesterdayAutoReceiveNumber := len(AutoReceiveNumber)
	AutoReceiveAmount = AwardAmount - ManualReceiveAmount
	fmt.Printf("昨日登录人数:%d,昨日参与的人数：%d\n", LoginNumber, RechargeNumber)
	fmt.Printf("活动id:%d,活动名称:%s,昨日派发总金额:%.2f,昨日完成任务人数:%d\n", activeId, activeName, AwardAmount, yesterdayFishPerson)
	fmt.Printf("手动领取人数:%d,手动领取金额:%.2f\n", yesterdayManualReceiveNumber, ManualReceiveAmount)
	fmt.Printf("自动领取人数:%d,自动领取金额:%.2f\n", yesterdayAutoReceiveNumber, AutoReceiveAmount)
	fmt.Printf("触达人数:%d,触达率：%.2f%%,参与率：%.2f%%\n", TouchNumber,
		float64(TouchNumber)/float64(LoginNumber)*100,
		float64(RechargeNumber)/float64(TouchNumber)*100)
	fmt.Printf("完成率：%.2f%%,手动领取率：%.2f%%,自动领取率：%.2f%%\n",
		float64(yesterdayFishPerson)/float64(RechargeNumber)*100,
		float64(yesterdayManualReceiveNumber)/float64(yesterdayFishPerson)*100,
		float64(yesterdayAutoReceiveNumber)/float64(yesterdayFishPerson)*100)
	fmt.Printf("手动领取成本率:%.2f%%,自动领取成本率：%.2f%%\n",
		float64(ManualReceiveAmount)/float64(AwardAmount)*100,
		float64(AutoReceiveAmount)/float64(AwardAmount)*100)

}

// 每日的执行，从csv中读取数据，进行登录，进行充值，进行签到，或者加入黑名单
func PrepareDataByCsv() []string {
	// 从csv中读取数据
	if amountlist, err := accounts.ProcessCSVConcurrently(CSVADDR, 4); err != nil {
		logger.Logger.Error("从csv中读取数据失败", err)
		return nil
	} else {
		return amountlist
	}
}

// 准备数据和签到活动
func RunDailyCheckInActivity() {
	// 数据准备
	// PrepareData()
	//time.Sleep(time.Second * 5)
	//从csv里面读取数据出来
	list := PrepareDataByCsv()
	//SingleCheckinTask("911030331131")
	// ExcelEverDayCheckIn(list)
	for _, item := range list {
		time.Sleep(time.Second * 5)
		res ,_ := SingleCheckinTask(item)
		if res == "5479" {
			continue
		}
	}
}
