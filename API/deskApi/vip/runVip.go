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

		}
	}
}
