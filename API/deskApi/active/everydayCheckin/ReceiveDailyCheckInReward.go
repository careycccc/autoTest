package everydayCheckin

import (
	desklogin "autoTest/API/deskApi/loginApi"
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

// 每日签到接口

type ReceiveDailyCheckInReward struct {
	ActivityId int `json:"activityId"`
	RwardType  int `json:"rewardType"`
	model.BaseStruct
}

/*
每日点击签到接口
activityId : 活动id
rewardType: 奖励类型
*
*/
func ReceiveDailyCheckInRewardApi(ctx *context.Context, activityId, rewardType int) (*model.Response, error) {
	api := "/api/Activity/ReceiveDailyCheckInReward"
	payloadStruct := &ReceiveDailyCheckInReward{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{activityId, rewardType, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/ReceiveDailyCheckInReward请求失败", err)), err
	} else {
		if string(respBoy) == "" {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/ReceiveDailyCheckInReward返回值为空", err)), err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/ReceiveDailyCheckInReward返回值解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

// 获取每个用户签到信息
type GetUserCheckInActivityinfo struct {
	Data struct {
		ActivityId         int    `json:"activityId"`         // 活动id
		ActivityName       string `json:"activityName"`       // 活动名称
		CurrentCheckInDays int    `json:"currentCheckInDays"` // 当前签到天数
		RewardDetail       []struct {
			DayIndex       int     `json:"dayIndeAx"`      // 签到第几天
			RechargeAmount float64 `json:"rechargeAmount"` // 充值金额
			RewardAmount   float64 `json:"rewardAmount"`   // 奖励金额
			RewardType     int     `json:"rewardType"`     // 奖励类型
			Status         int     `json:"status"`         // 状态
		} `json:"rewardDetail"`
	} `json:"data"`
}

// 获取到每个用户的签到信息
func GetUserCheckInActivityData(ctx *context.Context) (*model.Response, *GetUserCheckInActivityinfo, error) {
	api := "/api/Activity/GetUserCheckInActivityData"
	payloadStruct, payloadList := utils.BaseStructHandler()
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.InitStructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/GetUserCheckInActivityData请求失败", err)), nil, err
	} else {
		if string(respBoy) == "" {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/GetUserCheckInActivityData返回值为空", err)), nil, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/GetUserCheckInActivityData返回值解析失败", err)), nil, err
		} else {
			var respData *GetUserCheckInActivityinfo
			if err := json.Unmarshal(respBoy, &respData); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/GetUserCheckInActivityData[GetUserCheckInActivityinfo]返回值解析失败", err)), nil, err
			} else {
				return resp, respData, nil
			}
		}
	}
}

// 运行每日签到
func RunEverydayCheckIn() {
	ctx, err := desklogin.ReturnContextLoginY1("911166699876", "qwer1234")
	if err != nil {
		return
	}
	// 1.获取用户签到信息
	res, respData, err := GetUserCheckInActivityData(ctx)
	if err != nil {
		return
	}
	logger.Logger.Info("每日签到信息", res, respData)
	id := respData.Data.ActivityId
	if id == 0 {
		logger.Logger.Warn("没有获取到用户签到信息")
		return
	}
	// 2.点击签到按钮
	resp, err := ReceiveDailyCheckInRewardApi(ctx, id, 0)
	if err != nil {
		return
	}
	if resp.Msg != "Succeed" {
		logger.Logger.Warn("每日签到失败", resp.Msg)
		return
	}
	logger.Logger.Info("每日签到信息", resp)
}
