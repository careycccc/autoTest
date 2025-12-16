package dailycheckin

import (
	rechargeorders "autoTest/API/adminApi/financialManagement/rechargeOrders"
	memberlist "autoTest/API/adminApi/memberList/memberList"
	"autoTest/PressureMeasurementModule/accounts"
	"autoTest/store/config"
	"autoTest/store/logger"
	"fmt"
	"sync"
	"time"

	"github.com/panjf2000/ants/v2"
)

// SummaryItem 结构体用于存放我们按 UserId 汇总后的新数据 (输出结果)
type SummaryItem struct {
	UserId              int     `json:"userId"`
	TotalActualAmount   float64 `json:"totalAmount"` // 汇总后的金额
	RechargeChannelName string  `json:"rechargeChannelName"`
	RechargeChannelType string  `json:"rechargeChannelType"`
	RechargeType        string  `json:"rechargeType"`
}

// 签到通用方法

// ExcelEverDayCheckIn 封装了并发安全处理 []string 切片的逻辑，返回 []UserDailyCheckInInfo。
func ExcelEverDayCheckIn(input []string) []string {
	poolSize := 10 // 协程池大小

	p, err := ants.NewPool(poolSize)
	if err != nil {
		logger.LogError("创建协程池失败", err)
		return nil
	}
	defer p.Release()

	resultChan := make(chan string, len(input))
	var wg sync.WaitGroup

	fmt.Printf("启动 ants 协程池，最大并发 Goroutine 数量: %d\n", poolSize)
	fmt.Printf("总任务数: %d\n", len(input))

	for i, rawString := range input {
		time.Sleep(time.Second * 2)
		wg.Add(1)

		taskFunc := func(data string, index int) func() {
			return func() {
				defer wg.Done()
				processedInfo, err := SingleCheckinTask(data)
				if err != nil {
					// 任务执行失败！
					// 记录错误信息，并通过 return 跳过该任务的结果收集。
					fmt.Printf("  [Task %d] 处理失败，跳过结果收集: %s, 错误: %v\n", index, data, err)
					return // <-- 关键：跳过任务的剩余部分
				}

				// 任务成功，将结果安全地写入 resultChan
				resultChan <- processedInfo

			}
		}(rawString, i)

		err := p.Submit(taskFunc)
		if err != nil {
			logger.LogError("提交任务失败", err)
			wg.Done()
		}
	}

	// 等待所有任务完成并关闭结果 Channel
	go func() {
		wg.Wait()
		close(resultChan)
		fmt.Println("\n所有任务已执行完毕，开始收集结果...")
	}()

	// 收集结果
	var output []string
	for result := range resultChan {
		output = append(output, result)
	}

	return output
}

// SummarizeOrders 接收 OrderItem 切片，按 UserId 汇总 ActualAmount
func SummarizeOrders(items []rechargeorders.OrderItem) []SummaryItem {
	summaryMap := make(map[int]SummaryItem)

	for _, item := range items {
		if summary, exists := summaryMap[item.UserId]; exists {
			summary.TotalActualAmount += item.ActualAmount
			summaryMap[item.UserId] = summary
		} else {
			summaryMap[item.UserId] = SummaryItem{
				UserId:              item.UserId,
				TotalActualAmount:   item.ActualAmount,
				RechargeChannelName: item.RechargeChannelName,
				RechargeChannelType: item.RechargeChannelType,
				RechargeType:        item.RechargeType,
			}
		}
	}

	var summarizedList []SummaryItem
	for _, summary := range summaryMap {
		summarizedList = append(summarizedList, summary)
	}

	return summarizedList
}

// 第一次的数据准备  10个vip0的新用户，10个vip1的老用户，10个vip2的老用户
func PrepareData() {
	// 随机10个vip0的新用户
	// userList := utils.RandmoUserId(10)
	// for _, v := range userList {
	// 	// 进行注册
	// 	_, _, err := registerapi.NewGeneralAgentRegister(v)
	// 	if err != nil {
	// 		logger.Logger.Error("签到活动的数据准备-注册新用户失败", err)
	// 		return
	// 	}
	// }

	// vip的老用户
	vipList := make([]string, 0, 20)
	// 随机10个vip1,vip2的老用户
	for i := 1; i <= 4; i++ {
		if i == 2 || i == 3 {
			continue
		}
		if _, userinfoList, err := memberlist.GetUserVipListApi(AdminCtx, 5, 20, 0, i); err != nil {
			logger.Logger.Warn("GetUserVipListApi请求失败", err)
			continue
		} else {
			// 转成账号
			for _, v := range userinfoList {
				// id -> 账号
				if _, amount, err := memberlist.GetUserAmount(AdminCtx, int(v.UserId)); err != nil {
					logger.Logger.Warn("转账号失败", err)
					continue
				} else {
					// 修改账号的密码
					if _, err := memberlist.UpdatePassword(AdminCtx, v.UserId, config.SUB_PWD); err != nil {
						logger.Logger.Warn("修改密码失败", err)
						continue
					} else {
						vipList = append(vipList, amount)
					}
				}

			}
		}
	}
	// 所有账号准备完毕
	//userList = append(userList, vipList...)
	//把账号写入到csv中
	accounts.WriteConcurrently(vipList, 5, CSVADDR)
}

// ====================================================================
// 构建器 (Builder) 结构体
// ====================================================================

// UserDailyCheckInInfoBuilder 用于暂存和构建 UserDailyCheckInInfo 实例
type UserDailyCheckInInfoBuilder struct {
	info *UserDailyCheckInInfo
}

// 🚀 NewUserDailyCheckInInfoBuilder
// 构造函数：创建一个新的构建器实例。结构体初始化为所有字段的零值。
func NewUserDailyCheckInInfoBuilder() *UserDailyCheckInInfoBuilder {
	return &UserDailyCheckInInfoBuilder{
		// 初始化一个所有字段都是零值的结构体指针
		info: &UserDailyCheckInInfo{},
	}
}

// ====================================================================
// 可选参数设置方法 (链式调用)
// 每个方法都返回构建器本身的指针，支持链式调用。
// ====================================================================

func (b *UserDailyCheckInInfoBuilder) WithUserAccount(userAccount string) *UserDailyCheckInInfoBuilder {
	b.info.UserAccount = userAccount
	return b
}

func (b *UserDailyCheckInInfoBuilder) WithUserId(id int) *UserDailyCheckInInfoBuilder {
	b.info.UserId = id
	return b
}

func (b *UserDailyCheckInInfoBuilder) WithVipLevel(level int) *UserDailyCheckInInfoBuilder {
	b.info.VipLevel = level
	return b
}

// 设置充值金额
func (b *UserDailyCheckInInfoBuilder) WithRechargeAmount(amount float64) *UserDailyCheckInInfoBuilder {
	b.info.RechargeAmount = amount
	return b
}

// 设置签到天数
func (b *UserDailyCheckInInfoBuilder) WithCheckinNumberDay(day int) *UserDailyCheckInInfoBuilder {
	b.info.CheckinNumberDay = day
	return b
}

// 设置签到奖励
func (b *UserDailyCheckInInfoBuilder) WithCheckinAward(award float64) *UserDailyCheckInInfoBuilder {
	b.info.CheckinAward = award
	return b
}

// 设置是否自动领取
func (b *UserDailyCheckInInfoBuilder) WithIsAutoReceiveAward(isAuto bool) *UserDailyCheckInInfoBuilder {
	b.info.IsAutoReceiveAward = isAuto
	return b
}

// 设置是否是黑名单
func (b *UserDailyCheckInInfoBuilder) WithIsBlacklist(isBlacklist bool) *UserDailyCheckInInfoBuilder {
	b.info.IsBlacklist = isBlacklist
	return b
}

// 设置是否是详情
func (b *UserDailyCheckInInfoBuilder) WithIsDetailsPage(isDetails bool) *UserDailyCheckInInfoBuilder {
	b.info.IsDetailsPage = isDetails
	return b
}

// 设置警告信息
func (b *UserDailyCheckInInfoBuilder) WithAlarmInformation(alarm string) *UserDailyCheckInInfoBuilder {
	b.info.AlarmInformation = alarm
	return b
}

// ====================================================================
// 构建完成方法
// ====================================================================

// Build 最终返回构建好的 UserDailyCheckInInfo 结构体
func (b *UserDailyCheckInInfoBuilder) Build() UserDailyCheckInInfo {
	// 返回构建好的结构体值 (非指针)
	return *b.info
}
