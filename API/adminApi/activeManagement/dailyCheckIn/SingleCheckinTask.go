package dailycheckin

import (
	financialmanagement "autoTest/API/adminApi/financialManagement"
	memberactivityblacklist "autoTest/API/adminApi/memberList/MemberActivityBlacklist"
	"autoTest/API/deskApi/active/everydayCheckin"
	getuserinfo "autoTest/API/deskApi/getUserinfo"
	login "autoTest/API/deskApi/loginApi"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	util "autoTest/store/utils"
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
func SingleCheckinTask(userAccount string) (string, error) {
	// 1.用户进行前台登录
	if ctx, err := login.ReturnContextLoginY1(userAccount, config.SUB_PWD); err != nil {
		logger.Logger.Warn(userAccount, "前台登录失败", err)
		return "5479", err
	} else {
		// 2.获取会员信息
		if _, userinfo, err := getuserinfo.GetUserInfo(ctx); err != nil {
			logger.Logger.Warn(userAccount, "获取会员信息失败", err)
			return "", err
		} else {
			// 查询这个会员在不在黑名单中
			_,isBlackList, _ := memberactivityblacklist.UserActivityBlockIsInBlockApi(ctx, strconv.Itoa(userinfo.UserID), 21)
			// 3.判断会员是否在黑名单中
			if isBlackList {
				// 在黑名单中,50%的几率解除黑名单
				if RandomIntAndCompare(100, 40) {
					RemoveUserFromBlackList(userinfo.UserID)
				}
			} else {
				// 不在黑名单中,40%的几率加入黑名单
				if RandomIntAndCompare(100, 60) {
					AddUserToBlackList(userinfo.UserID)
				}
			}
			// 等待5秒，进行充值操作
			time.Sleep(5 * time.Second)
			moneny, _ := util.GenerateRandomInt(2000, 3000)
			if resp, err := financialmanagement.ArtificialRechargeFunc(AdminCtx, userinfo.UserID, moneny, 1); err != nil {
				logger.Logger.Warn(userAccount, "会员充值失败", err)
				return resp.Msg, err
			} else {
				isSuccess := model.IsSuccess(resp)
				if !isSuccess {
					logger.Logger.Warn(userAccount, "isSuccess会员充值失败", err)
					return "", err
				} else {
					// 等待5秒进行签到操作
					time.Sleep(5 * time.Second)
					// 50%的几率签到
					if RandomIntAndCompare(100, 80) {
						// 1.获取用户签到信息
						resp, respData, err := everydayCheckin.GetUserCheckInActivityData(ctx)
						if err != nil {
							return resp.Msg, err
						}
						if resp.Msg == "" {
							return "该会员不满足本轮活动的签到条件", nil
						}
						//logger.Logger.Info("每日签到信息", res, respData)
						id := respData.Data.ActivityId
						if id == 0 {
							logger.Logger.Warn("没有获取到用户签到信息")
							return "", err
						}
						// 2.点击签到按钮
						resp, err = everydayCheckin.ReceiveDailyCheckInRewardApi(ctx, id, 0)
						if err != nil {
							return resp.Msg, err
						}
						if resp.Msg != "Succeed" {
							logger.Logger.Warn("每日签到失败", resp.Msg)
							return resp.Msg, err
						}
						logger.Logger.Info("每日签到信息", resp.Msg)
					}
					return "", nil
				}
			}
			
		}
	}
}

// 主动将该会员加入黑名单中 有20的概率加入黑名单
func AddUserToBlackList(userId int) bool {
	result := RandomIntAndCompare(100, 20)
	if result {
		// 加入黑名单
		if resp, err := memberactivityblacklist.UserActivityBlockAddApi(AdminCtx, strconv.Itoa(userId), 21); err != nil {
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
		if resp, err := memberactivityblacklist.UserActivityBlockDeleteApi(AdminCtx, userId, 21); err != nil {
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
