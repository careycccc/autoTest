package withdrawalorders

import (
	"autoTest/API/adminApi/login"
	"autoTest/store/logger"
	"time"
)

// 解决后台的提现，锁定->选择三方提现通道->点击提现 -> 确认出款
func RunWithdraw(userId int, WithdrawType string, minWithdrawAmount, maxWithdrawAmount float64) {
	// 后台登录
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		logger.LogError("后台确定提现商户登录失败", err)
		return
	} else {
		// 查询订单
		if _, withdrawInfo, err := GetWithdrawLockPageListApi(ctx, userId, WithdrawType, minWithdrawAmount, maxWithdrawAmount); err != nil {
			logger.LogError("查询订单失败", err)
			return
		} else {
			//logger.Logger.Info("查询订单", resp)
			// 点击锁定
			if _, err := LockWithdrawOrderApi(ctx, userId, withdrawInfo); err != nil {
				logger.LogError("锁定订单失败", err)
				return
			} else {
				//logger.Logger.Info("锁定订单", resp)
				// 获取可以提现的三方通道
				if _, channleId, err := GetCanWithdrawChannelByOrderApi(ctx, userId, *withdrawInfo); err != nil {
					logger.LogError("获取提现的三方通道失败", err)
					return
				} else {
					time.Sleep(time.Second)
					//logger.Logger.Info("提现的三方通道", resp, channleId)
					// 点击确认提现
					if resp, err := ThirdWithdrawApi(ctx, userId, *withdrawInfo, channleId); err != nil {
						logger.LogError("点击提现失败", err)
						return
					} else {
						logger.Logger.Info("点击提现", resp)
						// 点击确认出款
						if resp, err := ConfirmWithdrawOrderApi(ctx, userId, *withdrawInfo); err != nil {
							logger.LogError("点击确认出款失败", err)
							return
						} else {
							logger.Logger.Info("出款结果", resp)
						}
					}
				}
			}

		}
	}
}
