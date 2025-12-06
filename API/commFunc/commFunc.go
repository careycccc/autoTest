package commfunc

import (
	financialmanagement "autoTest/API/adminApi/financialManagement"
	memberlist "autoTest/API/adminApi/memberList/memberList"
	"autoTest/store/config"
	"autoTest/store/logger"
	"context"
	"sync"
)

// 放一些前后台都需要方法
/*
传入userid 进行充值和修改密码
传入后台的token，充值金额
**/
func UpdatePasswordAndToUp(ctxAdminToken *context.Context, userid int64, monenyCount float64) error {
	adminToken := *ctxAdminToken
	wg := &sync.WaitGroup{}
	wg.Add(2)
	// 根据id进行充值
	go func(wg *sync.WaitGroup, ctxToUse *context.Context) {
		defer wg.Done()
		financialmanagement.ArtificialRechargeFunc(ctxToUse, int(userid), monenyCount, 2)
	}(wg, &adminToken)

	// 修改用户密码
	go func(wg *sync.WaitGroup, ctxToUse *context.Context) {
		defer wg.Done()
		memberlist.UpdatePassword(ctxToUse, userid, config.SUB_PWD)
	}(wg, &adminToken)

	wg.Wait()
	logger.Logger.Info("充值金额", monenyCount)
	logger.Logger.Info("成功修改密码", config.SUB_PWD)
	return nil
}
