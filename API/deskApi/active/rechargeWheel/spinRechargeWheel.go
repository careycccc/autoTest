package rechargewheel

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
)

// 旋转充值转盘
type SpinRechargeWheel struct {
	RechargeWheelType any `json:"rechargeWheelType"` // 充值转盘类型
	model.BaseStruct
}

//	type SpinRechargeWheelResponse struct {
//		Data struct {
//			Id           int     `json:"id"`         // 奖励类型 1，积分 2，代金券 3，实物奖励
//			RewardType   int     `json:"rewardType"` // 奖励数量
//			RewardAmount float64 `json:"rewardAmount"`
//		} `json:"data"`
//	}
//
// 旋转转盘
func SpinRechargeWheelApi(ctx *context.Context, rechargeWheelType int8) (*model.Response, error) {
	api := "/api/Activity/SpinRechargeWheel"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &SpinRechargeWheel{}
	payloadList := []interface{}{rechargeWheelType, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/SpinRechargeWheel 请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/SpinRechargeWheel 解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}
