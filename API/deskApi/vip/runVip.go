package vip

import (
	registerapi "autoTest/API/deskApi/registerApi"
	"autoTest/store/logger"
)

// 执行vip
func RunVip() {
	// 进行登录
	userName := "911023199716"
	if _, ctx, err := registerapi.GeneralAgentRegister(userName); err != nil {
		logger.LogError("vip登录失败", err)
		return
	} else {
		ctxToken := ctx
		if _, vipInfo, err := GetUserVipInfo(ctxToken); err != nil {
			logger.LogError("获取用户的vip信息失败", err)
			return
		} else {
			if vipInfo.VipLevel == 0 {
				// 没有vip等级
				logger.LogError("该会员的vip等级为0", nil)
				return
			}
			// 领取vip的奖励
			// 领取升级奖励
			if resp, err := PickVipRewardApi(ctxToken, 2, int8(vipInfo.VipLevel)); err != nil {
				logger.LogError("该会员的升级奖励领取失败", err)
			} else {
				logger.Logger.Info("升级奖励领取结果：", resp)
			}
			// 领取周奖励和月奖励
			// 判断周奖励和月奖励是否开启
			if vipInfo.WeekRewardState {
				// 周奖励开启
				if resp, err := PickVipRewardApi(ctxToken, 3, int8(vipInfo.VipLevel)); err != nil {
					logger.LogError("该会员的周奖励领取失败", err)
				} else {
					logger.Logger.Info("周奖励领取结果：", resp)
				}
			} else {
				logger.LogError("该会员的周奖励未开启", nil)
			}
			if vipInfo.MonthRewardState {
				// 周奖励开启
				if resp, err := PickVipRewardApi(ctxToken, 4, int8(vipInfo.VipLevel)); err != nil {
					logger.LogError("该会员的月奖励领取失败", err)
				} else {
					logger.Logger.Info("月奖励领取结果：", resp)
				}
			} else {
				logger.LogError("该会员的月奖励未开启", nil)
			}
		}
	}
}
