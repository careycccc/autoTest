package dailycheckin

import (
	financialmanagement "autoTest/API/adminApi/financialManagement"
	memberactivityblacklist "autoTest/API/adminApi/memberList/MemberActivityBlacklist"
	getuserinfo "autoTest/API/deskApi/getUserinfo"
	login "autoTest/API/deskApi/loginApi"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	sutils "autoTest/store/utils"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

/*
单个会员的签到流程
返回会员今日的签到信息
*
*/
func SingleCheckinTask(userAccount string) (UserDailyCheckInInfo, error) {
	// 0.登录前获取到上一天的该用户的状态
	yesterdayUserinfo := RecoverYesterdayDataByAccount(userAccount)
	// 1.用户进行前台登录
	if ctx, err := login.ReturnContextLoginY1(userAccount, config.SUB_PWD); err != nil {
		logger.Logger.Warn(userAccount, "前台登录失败", err)
		info := NewUserDailyCheckInInfoBuilder().WithUserAccount(userAccount).WithAlarmInformation("前台登录失败").WithCheckinNumberDay(0).Build()
		return info, err
	} else {
		// 2.获取用户的vip等级和金额，进行会员等级的判断，该活动会员是否可见
		if _, userinfo, err := getuserinfo.GetUserInfo(ctx); err != nil {
			logger.Logger.Warn(userAccount, "获取用户信息失败", err)
			info := NewUserDailyCheckInInfoBuilder().WithUserAccount(userAccount).WithAlarmInformation("获取用户信息失败").WithCheckinNumberDay(0).Build()
			return info, err
		} else {
			userLevel := userinfo.VipLevel
			userId := userinfo.UserID
			addBlack := yesterdayUserinfo.IsBlacklist // 这个起始值是昨天的该会员的黑名单的状态
			// 主动将该会员加入黑名单中
			if addBlack {
				// 解除黑名单 有40的概率解除黑名单
				result := RemoveUserFromBlackList(userId)
				if result {
					// 解除黑名单成功
					addBlack = !result
				}
			} else {
				// 加入黑名单 有20的概率加入黑名单
				result := AddUserToBlackList(userId)
				if result {
					// 加入黑名单成功
					addBlack = result
				}
			}
			//判断，该活动会员是否可见
			if userLevel < ShowObject[0] {
				info := NewUserDailyCheckInInfoBuilder().WithUserAccount(userAccount).WithVipLevel(userLevel).WithUserId(userId).WithAlarmInformation("会员vip等级不满足参与本次签到活动").WithCheckinNumberDay(0).Build()
				return info, nil
			} else {
				// 3.是否进入详情页 80%的概率进入活动详情页
				pageInfo := RandomIntAndCompare(100, 80)
				if pageInfo {
					// 进入详情页 发起详情页的请求
				} else {
					// 不进入详情页
					pageInfo = false
				}
				// 4.充值金额是否满足解锁活动的条件
				if money, err := sutils.GenerateRandomInt(config.MIN_MONENY, config.MAX_MONENY); err != nil {
					logger.Logger.Warn(userAccount, "生成充值金额失败", err)
					info := NewUserDailyCheckInInfoBuilder().WithUserAccount(userAccount).WithAlarmInformation("生成充值金额失败").WithIsDetailsPage(pageInfo).WithVipLevel(userLevel).WithUserId(userId).WithCheckinNumberDay(0).Build()
					return info, err
				} else {
					// 4.充值金额是否满足解锁活动的条件
					if resp, err := financialmanagement.ArtificialRechargeFunc(AdminCtx, userId, float64(money), 0); err != nil {
						logger.Logger.Warn(userAccount, "后台充值请求失败", err)
						info := NewUserDailyCheckInInfoBuilder().WithUserAccount(userAccount).WithAlarmInformation("后台充值请求失败").WithIsDetailsPage(pageInfo).WithVipLevel(userLevel).WithUserId(userId).WithCheckinNumberDay(0).Build()
						return info, err
					} else {
						isSuccess := model.IsSuccess(resp)
						if !isSuccess {
							// 充值失败
							logger.Logger.Warn(userAccount, "充值失败", err)
							info := NewUserDailyCheckInInfoBuilder().WithUserAccount(userAccount).WithIsDetailsPage(pageInfo).WithAlarmInformation("充值失败").WithVipLevel(userLevel).WithUserId(userId).WithCheckinNumberDay(0).Build()
							return info, err
						} else {
							// 充值成功
							if money < AcitveRechargeAmount {
								logger.Logger.Warn(userAccount, "充值金额小于参与活动的金额", err)
								info := NewUserDailyCheckInInfoBuilder().WithUserAccount(userAccount).WithIsDetailsPage(pageInfo).WithVipLevel(userLevel).WithUserId(userId).WithRechargeAmount(money).WithAlarmInformation("充值金额小于参与活动的金额").WithCheckinNumberDay(0).Build()
								return info, err
							} else {
								// 4.充值金额满足解锁活动的条件
								// 5.判断当前用户是否在黑名单中
								if addBlack {
									// 在黑名单中,
									info := NewUserDailyCheckInInfoBuilder().WithUserAccount(userAccount).WithIsDetailsPage(pageInfo).WithVipLevel(userLevel).WithUserId(userId).WithRechargeAmount(money).WithIsBlacklist(addBlack).WithCheckinNumberDay(yesterdayUserinfo.CheckinNumberDay + 1).Build()
									return info, nil
								} else {
									// 不在黑名单 50%发送签到请求
									sendChickein := RandomIntAndCompare(100, 50)
									if sendChickein {
										// 发起签到请求   *******  需要获取连续签到的天数和奖励
										info := NewUserDailyCheckInInfoBuilder().WithUserAccount(userAccount).WithIsDetailsPage(pageInfo).WithVipLevel(userLevel).WithUserId(userId).WithRechargeAmount(money).WithIsBlacklist(addBlack).WithIsAutoReceiveAward(sendChickein).WithCheckinNumberDay(yesterdayUserinfo.CheckinNumberDay + 1).WithCheckinAward(AcitveRechargeAmount).Build()
										return info, nil
									} else {
										AutoCheckInNumber = append(AutoCheckInNumber, userId)
										// 不发起 签到请求,自动领取
										info := NewUserDailyCheckInInfoBuilder().WithUserAccount(userAccount).WithIsDetailsPage(pageInfo).WithVipLevel(userLevel).WithUserId(userId).WithRechargeAmount(money).WithIsBlacklist(addBlack).WithIsAutoReceiveAward(sendChickein).WithCheckinNumberDay(yesterdayUserinfo.CheckinNumberDay + 1).Build()
										return info, nil
									}
								}

							}
						}
					}
				}
			}
		}
	}
	// 6.点击签到 手动签到，不点击就是自动签到，明天会自动派发
	// 7.查看余额是否增加 账变的金额
}

// 主动将该会员加入黑名单中 有20的概率加入黑名单
func AddUserToBlackList(userId int) bool {
	result := RandomIntAndCompare(100, 20)
	if result {
		// 加入黑名单
		if resp, err := memberactivityblacklist.UserActivityBlockAddApi(AdminCtx, strconv.Itoa(userId), 1); err != nil {
			logger.Logger.Warn(userId, "加入黑名单失败", err)
			return false
		} else {
			isSuccess := model.IsSuccess(resp)
			if !isSuccess {
				logger.Logger.Warn(userId, "加入黑名单失败", err)
				return false
			} else {
				return true
			}
		}
	}
	return false
}

// 主动将该会员加入黑名单中 有40的概率解除黑名单
func RemoveUserFromBlackList(userId int) bool {
	result := RandomIntAndCompare(100, 40)
	if result {
		// 解除黑名单
		if resp, err := memberactivityblacklist.UserActivityBlockDeleteApi(AdminCtx, userId, 1); err != nil {
			logger.Logger.Warn(userId, "解除黑名单失败", err)
			return false
		} else {
			isSuccess := model.IsSuccess(resp)
			if !isSuccess {
				logger.Logger.Warn(userId, "解除黑名单失败", err)
				return false
			} else {
				return true
			}
		}
	}
	return false
}

// 1. 定义一个全局或包级别的随机源和锁
// 确保它只被初始化一次，并且可以被多次调用共享。
// 使用 sync.Once 确保全局随机源只初始化一次
var (
	globalRand *rand.Rand
	randMutex  sync.Mutex
	once       sync.Once
)

// initGlobalRand 初始化全局随机源
func initGlobalRand() {
	// 使用高精度时间作为种子
	seed := time.Now().UnixNano()
	globalRand = rand.New(rand.NewSource(seed))
}

// RandomIntAndCompare 生成 [0, n) 范围的随机整数，并与 compareValue 进行比较
// n: 随机数的上限（不包含）
// compareValue: 用于比较的值
// 返回值: 生成的随机数 <= compareValue 则为 true，否则为 false
func RandomIntAndCompare(n int, compareValue int) bool {
	if n <= 0 {
		// 无效输入，返回 false
		return false
	}

	// 确保全局随机源只初始化一次
	once.Do(initGlobalRand)

	// 锁定，确保在并发调用时，Intn() 是原子操作，不会被中断
	randMutex.Lock()
	defer randMutex.Unlock()

	// 生成 [0, n) 范围的随机数
	randomValue := globalRand.Intn(n)

	return randomValue <= compareValue
}
