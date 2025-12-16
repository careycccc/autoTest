package common

import (
	memberlist "autoTest/API/adminApi/memberList/memberList"
	"autoTest/store/config"
	"autoTest/store/logger"
	"context"
	"sync"
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
