package vip

import (
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

type getVipLevelConfigStruct struct {
	Level                   int     // vip等级
	Recharges               float64 // 充值金额
	BetCode                 float64 // 打码量
	UpgradeRewards          float64 // 升级奖励
	MothRelegationRecharges float64 // 月保级充值金额
	MothRelegationBetCode   float64 // 月保级投注金额
	MothRewardRecharges     float64 // 月充值金额达到这个要求才能领取奖励
	MothRewardBetCode       float64 // 月打码量
	MothReward              float64 // 月充值奖励
	WeekRewardRecharges     float64 // 周充值金额达到这个要求才能领取奖励
	WeekRewardBetCode       float64 // 周打码量
	WeekReward              float64 // 周奖励
	DayWithdrawAmount       float64 // 每天提现的金额
	DayWithdrawCount        int     //每天可以提现的次数
	DayWithdrawNoFeeCount   int     // 每日免手续费次数
}

type GetVipLevelConfigResponse struct {
	Data []getVipLevelConfigStruct `json:"data"`
}

// 获取前台的vip信息
func GetVipLevelConfig(ctx *context.Context) (*model.Response, []getVipLevelConfigStruct, error) {
	api := "/api/VipLevel/GetVipLevelConfig"
	payloadStruct := model.BaseStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadData := []interface{}{random, language, "", timestamp}
	if respBody, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, &payloadStruct, payloadData, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("api/VipLevel/GetVipLevelConfig 请求失败", err)), nil, err
	} else {
		if resp, err := model.ParseResponse(respBody); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("api/VipLevel/GetVipLevelConfig 响应结果解析失败", err)), nil, err
		} else {
			var response GetVipLevelConfigResponse
			if err := json.Unmarshal(respBody, &response); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("api/VipLevel/GetVipLevelConfig[GetVipLevelConfigResponse] 响应结果解析失败", err)), nil, err
			}
			return resp, response.Data, nil
		}
	}
}

type GetUserVipInfoResponse struct {
	Data GetUserVipInfoData
}

type GetUserVipInfoData struct {
	UserId                int64
	VipLevel              int // vip等级
	DaysLeft              int
	UpLevelBetAmount      float64
	UpLevelRechargeAmount float64
	WeekBetAmount         float64
	WeekRechargeAmount    float64
	WeekRewardState       bool
	MonthBetAmount        float64
	MonthRechargeAmount   float64
	MonthRewardState      bool
	ReceivedLevels        string
	VipAmountOfCode       int
}

// 获取会员的vip信息
func GetUserVipInfo(ctx *context.Context) (*model.Response, GetUserVipInfoData, error) {
	api := "/api/VipLevel/GetUserVipInfo"
	payloadStruct, payloadList := utils.BaseStructHandler()
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.InitStructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/VipLevel/GetUserVipInfo请求报错", err)), GetUserVipInfoData{}, err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/VipLevel/GetUserVipInfo解析报错", err)), GetUserVipInfoData{}, err
		} else {
			var response GetUserVipInfoResponse
			if err := json.Unmarshal(respBoy, &response); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/VipLevel/GetUserVipInfo【GetUserVipInfoResponse】解析报错", err)), GetUserVipInfoData{}, err
			} else {
				return resp, response.Data, nil
			}
		}
	}
}

// 领取vip奖励
type PickVipRewardstruct struct {
	RewardType  any `json:"rewardType"`  // 3周奖励 4月奖励 2是升级奖励
	RewardLevel any `json:"rewardLevel"` // vip等级
	model.BaseStruct
}

/*
RewardType 领取奖励的类型 // 3周奖励 4月奖励 2是升级奖励
RewardLevel  vip等级
*
*/
func PickVipRewardApi(ctx *context.Context, rewardType, RewardLevel int8) (*model.Response, error) {
	api := "/api/VipLevel/PickVipReward"
	payloadStruct := &PickVipRewardstruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{rewardType, RewardLevel, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/VipLevel/PickVipReward请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/VipLevel/PickVipReward解析失败", err)), err
		} else {
			return resp, err
		}
	}
}
