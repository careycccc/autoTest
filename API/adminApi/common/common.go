package common

import (
	financialmanagement "autoTest/API/adminApi/financialManagement"
	"autoTest/API/adminApi/login"
	memberlist "autoTest/API/adminApi/memberList/memberList"
	withdrawcash "autoTest/API/deskApi/WithdrawCash"
	"autoTest/store/config"
	"autoTest/store/logger"
	sutils "autoTest/store/utils"
	"context"
	"sync"
	"time"
)

// id转账号并且修改登录密码
func IdToAmountAndUpdatePassword(ctx *context.Context, id int) string {
	// 查询这个账号是不是游客账号
	_, userinfo, err := memberlist.GetUserDetail(ctx, id)
	if err != nil {
		logger.Logger.Warn("id查询详情失败")
		return ""

	}
	if userinfo.UserType == 2 {
		logger.Logger.Warn("该账号是游客账号")
		return ""
	}
	result := ""
	adminCtx := ctx
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func(wg *sync.WaitGroup, ctx *context.Context, id int) {
		defer wg.Done()
		if _, amount, err := memberlist.GetUserAmount(ctx, id); err != nil {
			logger.Logger.Warn("id转账号失败")
			return
		} else {
			result = amount
		}
	}(wg, adminCtx, id)
	go func(wg *sync.WaitGroup, ctx *context.Context, id int) {
		defer wg.Done()
		if _, err := memberlist.UpdatePassword(ctx, int64(id), config.ADMIN_PWD); err != nil {
			logger.Logger.Warn("修改密码失败")
			return
		}
	}(wg, adminCtx, id)
	wg.Wait()
	return result
}

// 随机从用户列表中选择用户进行充值，投注
func RandUserTopupGame() {
	var adminCtx *context.Context
	// 后台登录
	ctx, err := login.RunAdminSitLogin()
	if err != nil {
		logger.LogError("随机从用户列表中选择用户进行充值后台登录失败", err)
		return
	} else {
		adminCtx = ctx
	}
	UserAmount := make([]*memberlist.UserInfo, 0, 100)
	// 从后台获取用户列表
	for i := 5; i < 7; i++ {
		if _, userList, err := memberlist.GetUserVipListApi(adminCtx, 1, 20, 0, i); err != nil {
			continue
		} else {
			UserAmount = append(UserAmount, userList...)
		}
	}
	// amountList := make([]string, 0, 100)
	for _, user := range UserAmount {
		// 随机选择用户进行充值,
		if user.State != 1 {
			// 解冻操作
			if _, err := memberlist.UpdateUserStateApi(adminCtx, int(user.UserId), 1); err != nil {
				logger.Logger.Warn("解冻用户失败", user.UserId)
				continue
			}
		}

		time.Sleep(time.Second * 3)
		// 充值
		money, err := sutils.GenerateRandomInt(config.MIN_MONENY, config.MAX_MONENY)
		if err != nil {
			logger.Logger.Warn("生成随机金额失败", err)
			continue
		}
		if _, err := financialmanagement.ArtificialRechargeFunc(adminCtx, int(user.UserId), money, 2); err != nil {
			logger.Logger.Warn("人工充值失败", err)
			continue
		}
		userAmount := IdToAmountAndUpdatePassword(adminCtx, int(user.UserId))
		// 前台登录并进行投注
		// if err := lotterygameapi.BetRun(userAmount); err != nil {
		// 	logger.LogError("前台投注失败", err)
		// }
		// amountList = append(amountList, userAmount)
		// 绑定银行卡,提现
		withdrawcash.RunWithDrawCase(userAmount) // 提现
	}
	// // 保存到csv中
	// err = accounts.WriteConcurrently(amountList, 5, config.CSVADDR) // 保存到csv中
	// if err != nil {
	// 	logger.LogError("保存到csv失败", err)
	// }
	// 等待15分钟
	// time.Sleep(time.Minute * 15)
	// for _, amount := range amountList {
	// 	// 进行登录
	// 	ctx, err := desklogin.ReturnContextLoginY1(amount, config.ADMIN_PWD)
	// 	if err != nil {
	// 		logger.LogError("用户登录失败", err)
	// 		continue
	// 	}
	// 	// 加入锦标赛
	// 	if resp, err := champion.JoinChampion(ctx, 500028); err != nil {
	// 		logger.LogError("加入锦标赛失败", err)
	// 		continue
	// 	} else {
	// 		logger.Logger.Info("加入锦标赛成功", resp)
	// 		// 进行投注
	// 		time.Sleep(time.Second * 10)
	// 		// 先把投注的结果随机出来
	// 		gameCode, betContent, money, betMultiple := lotterygameapi.GetBetResult()
	// 		if err := lotterygameapi.RunBetFunc(ctx, gameCode, betContent, amount, money, betMultiple); err != nil {
	// 			continue
	// 		}
	// 	}
	// }
}
