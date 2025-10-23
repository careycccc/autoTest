package withdrawalorders

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

// 提现出款
// 获取可以提现出款的三方通道
type GetCanWithdrawChannelByOrder struct {
	Remark     any `json:"remark"` // 备注
	UserId     any `json:"userId"`
	OrderNo    any `json:"orderNo"`    // 订单号
	CreateTime any `json:"createTime"` // 订单创建时间
	model.BaseStruct
}

type WithdrawChannelId struct {
	Data []struct {
		ChannelId int `json:"channelId"` // 通道id
	}
}

/*
userId int, // 传入用户id
withdrawinfo // 传入提现的信息
返回三方提现的通道id
*
*/
func GetCanWithdrawChannelByOrderApi(ctx *context.Context, userId int, withdrawinfo Withdrawinfo) (*model.Response, int, error) {
	api := "/api/WithdrawOrder/GetCanWithdrawChannelByOrder"
	payloadStruct := &GetCanWithdrawChannelByOrder{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{config.Remark, userId, withdrawinfo.orderNo, withdrawinfo.createTime, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/WithdrawOrder/GetCanWithdrawChannelByOrder请求失败", err)), 0, err
	} else {
		var response WithdrawChannelId
		if err := json.Unmarshal(respBoy, &response); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/WithdrawOrder/GetCanWithdrawChannelByOrder【WithdrawChannelId】解析失败", err)), 0, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/WithdrawOrder/GetCanWithdrawChannelByOrder解析失败", err)), 0, err
		} else {
			return resp, response.Data[0].ChannelId, nil
		}
	}
}

type ThirdWithdrawStruct struct {
	OrderNo           any `json:"orderNo"`
	UserId            any `json:"userId"`
	CreateTime        any `json:"createTime"`
	Remark            any `json:"remark"`
	WithdrawChannelId any `json:"withdrawChannelId"` // 提现通道id
	model.BaseStruct
}

// 点击三方通道进行提现
func ThirdWithdrawApi(ctx *context.Context, userId int, withdrawinfo Withdrawinfo, WithdrawChannelId int) (*model.BetResponse, error) {
	api := "/api/WithdrawOrder/ThirdWithdraw"
	payloadStruct := &ThirdWithdrawStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{withdrawinfo.orderNo, userId, withdrawinfo.createTime, config.Remark, WithdrawChannelId, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes2(model.ErrorLoggerType("/api/WithdrawOrder/ThirdWithdraw请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse2(respBoy); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/WithdrawOrder/ThirdWithdraw解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}
