package withdrawaltimeoutcompensation

import (
	withdrawalorders "autoTest/API/adminApi/financialManagement/withdrawalOrders"
	"autoTest/API/adminApi/login"
	memberlist "autoTest/API/adminApi/memberList/memberList"
	fundtransactionrecords "autoTest/API/adminApi/reportManagement/FundTransactionRecords"
	desklogin "autoTest/API/deskApi/loginApi"
	withdrawcash "autoTest/API/deskApi/withdrawcash"
	"autoTest/API/utils"
	"autoTest/store/config"
	"autoTest/store/logger"
	"context"
	"fmt"
	"sort"
	"sync"
)

// 提现超时赔付
// 获取昨天的提现超时赔付的账变的账号以及金额
func GetYesterdayUserAmont() {
	financialTypeList := []string{"WithdrawTimeoutCompensation"}
	fundTransactionInfo := fundtransactionrecords.RunFundtransactionrecordsApi(financialTypeList)
	result := AggregateByUserID(
		fundTransactionInfo.UserInfo,
		func(item fundtransactionrecords.FundTransactionInfo) int64 {
			return int64(item.UserID)
		},
		func(item fundtransactionrecords.FundTransactionInfo) float64 {
			return item.Moneny
		},
	)
	TotalMoneny := 0.0   // 总超时赔付的金额
	TotalOrderCount := 0 // 总超时赔付的的订单数
	logger.Logger.Info("超时赔付领取成功的数据:")
	for _, r := range result {
		fmt.Printf("用户 %d 有 %d 笔订单，明细：[", r.UserID, r.OrderCount)
		for i, a := range r.Amounts {
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("%.2f", a)
		}
		fmt.Printf("]，总金额：%.2f\n", r.TotalAmount)
		TotalMoneny += r.TotalAmount
		TotalOrderCount += r.OrderCount
	}
	peopleNubmer := len(result)
	fmt.Printf("总超时赔付的金额:%.2f,总超时赔付的人数:%d,总超时赔付的的订单数:%d\n", TotalMoneny, peopleNubmer, TotalOrderCount)
	fmt.Printf("人均赔付金额%.2f,平均每笔订单赔付金额%.2f\n", TotalMoneny/float64(peopleNubmer), TotalMoneny/float64(TotalOrderCount))
}

type UserOrderSummary struct {
	UserID      int64
	Amounts     []float64
	OrderCount  int
	TotalAmount float64
}

// 关键：这个函数接收任意包含 UserID 和 Amount 字段的结构体切片
// T 必须有一个 UserID() int64 和 Amount() float64 方法（下面会用反射或直接传）
func AggregateByUserID[T any](items []T, getUserID func(T) int64, getAmount func(T) float64) []UserOrderSummary {
	m := make(map[int64]*UserOrderSummary)

	for _, item := range items {
		userID := getUserID(item)
		amount := getAmount(item)

		if s, ok := m[userID]; ok {
			s.Amounts = append(s.Amounts, amount)
			s.TotalAmount += amount
			s.OrderCount++
		} else {
			m[userID] = &UserOrderSummary{
				UserID:      userID,
				Amounts:     []float64{amount},
				TotalAmount: amount,
				OrderCount:  1,
			}
		}
	}

	var result []UserOrderSummary
	for _, v := range m {
		result = append(result, *v)
	}

	// 按 UserID 排序
	sort.Slice(result, func(i, j int) bool {
		return result[i].UserID < result[j].UserID
	})

	return result
}

// 把会员id查找出对应的会员账号同时修改它们的登录密码,返回一个账号的队列
func GetAmountQueue(ctx *context.Context, userId int) string {
	userAmounts := ""
	// 会员id查找出对应的会员账号
	wg := &sync.WaitGroup{}
	wg.Add(2)
	// 根据id查找出用户
	go func(wg *sync.WaitGroup, ctxToUse *context.Context, userId int) {
		defer wg.Done()
		if _, userAmount, err := memberlist.GetUserAmount(ctx, userId); err != nil {
			logger.LogError("根据id查找出用户报错", err)
			return
		} else {
			//logger.Logger.Warn("根据id查找出用户的响应", resp)
			userAmounts = userAmount
		}
	}(wg, ctx, userId)

	// 修改用户密码
	go func(wg *sync.WaitGroup, userId int, ctxToUse *context.Context) {
		defer wg.Done()
		if _, err := memberlist.UpdatePassword(ctx, int64(userId), config.SUB_PWD); err != nil {
			logger.LogError("根据id查找出修改用户密码报错", err)

			return
		}
	}(wg, userId, ctx)

	wg.Wait()
	return userAmounts
}

// 执行，提现订单 + 提现锁定 返回会员列表
func ExcelQueue() []string {
	// 后台登录
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		logger.LogError("执行，提现订单 + 提现锁定 返回会员列表登录失败", err)
		return nil
	} else {
		list := withdrawalorders.GetWithdrawalAmountList()
		fmt.Printf("提现的会员id列表:%v,提现的人数:%d\n", list, len(list))
		userAmountList := []string{}
		for _, v := range list {
			st := GetAmountQueue(ctx, v)
			if st == "" {
				logger.Logger.Warn("该后台账号没有查询账号的权限")
				continue
			}
			userAmountList = append(userAmountList, st)
		}
		return userAmountList
	}
}

var (
	ToBeCollected int = 0 // 待领取
	Collected     int = 0 // 已领取
	Expired       int = 0 // 已过期
)

// task 负责登录和获取昨日的提现订单，并且返回那些触发了，昨天提现的订单数
func RunWithdrawTask(userName string) []withdrawcash.WithdrawHistoryInfo {
	var list []withdrawcash.WithdrawHistoryInfo
	//登录
	if ctx, err := desklogin.ReturnContextLoginY1(userName, "qwer1234"); err != nil {
		logger.Logger.Warn("提现订单的前台登录失败", err)
		return []withdrawcash.WithdrawHistoryInfo{}
	} else {
		_, startTime, _, endTime, _ := utils.ParseTimeRangeToTimestamp(config.StartTime, config.EndTime)
		if _, withdrawInfoList, err := withdrawcash.GetWithdrawHistoryApi(ctx, startTime, endTime); err != nil {
			logger.Logger.Warn("该用户的历史提现订单获取失败", err)
			return []withdrawcash.WithdrawHistoryInfo{}
		} else {
			// 把今天的剔除出去
			for _, v := range withdrawInfoList {
				// 并且统计这个账号触发了提现赔付的
				if v.CreateTime < endTime && v.CompensationState > 0 {
					switch v.CompensationState {
					case 1:
						ToBeCollected++
					case 2:
						Collected++
					case 3:
						Expired++
					}
					list = append(list, v)
				}
			}
			return list

		}
	}
}

// 执行提现历史的数据
func ExecelWithdrawHistoryInfo() {
	userAmountList := ExcelQueue()
	if len(userAmountList) == 0 {
		logger.Logger.Warn("没有获取到提现历史的数据")
		return
	}
	n := 0
	for _, v := range userAmountList {
		tasks := RunWithdrawTask(v)
		for _, v := range tasks {
			fmt.Println(v)
			n++
		}
	}
	fmt.Printf("提现触发了提现赔付的单子条数%d,待领取%d,已领取%d,已过期%d", n, ToBeCollected, Collected, Expired)
}
