package withdrawalorders

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

// 查询订单
type GetWithdrawLockPageList struct {
	UserId            any `json:"userId"`            // 用户id
	WithdrawType      any `json:"withdrawType"`      // 提现的类型
	MinWithdrawAmount any `json:"minWithdrawAmount"` // 最小金额
	MaxWithdrawAmount any `json:"maxWithdrawAmount"` // 最大的金额
	DateType          any `json:"dateType"`          // 默认0
	model.QueryPayloadStruct
}

type GetWithdrawLockPageListResponse struct {
	Data struct {
		List []struct {
			OrderNo    string `json:"orderNo"`
			CreateTime int64  `json:"createTime"`
		} `json:"list"`
	} `json:"data"`
}

type Withdrawinfo struct {
	orderNo    string // 订单号
	createTime int64  // 订单创建时间
}

/*
userId // 用户id
WithdrawType string, // 提现的类型
minWithdrawAmount,maxWithdrawAmount float64// 最小金额// 最大的金额
返回提现的订单号,订单创建时间
*
*/
func GetWithdrawLockPageListApi(ctx *context.Context, userId int, WithdrawType string, minWithdrawAmount, maxWithdrawAmount float64) (*model.Response, *Withdrawinfo, error) {
	api := "/api/WithdrawOrder/GetWithdrawLockPageList"
	payloadStruct := &GetWithdrawLockPageList{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{userId, WithdrawType, minWithdrawAmount, maxWithdrawAmount, 0, 1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/WithdrawOrder/GetWithdrawLockPageList请求失败", err)), nil, err
	} else {
		var response GetWithdrawLockPageListResponse
		if err := json.Unmarshal(respBoy, &response); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/WithdrawOrder/GetWithdrawLockPageList[GetWithdrawLockPageListResponse]解析失败", err)), nil, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("GetWithdrawLockPageListResponse解析失败", err)), nil, err
		} else {
			return resp, &Withdrawinfo{
				orderNo:    response.Data.List[0].OrderNo,
				createTime: response.Data.List[0].CreateTime,
			}, nil
		}
	}
}

// 点击锁定
type LockWithdrawOrder struct {
	UserId     any `json:"userId"`
	OrderNo    any `json:"orderNo"`
	CreateTime any `json:"createTime"`
	Remark     any `json:"remark"`
	model.BaseStruct
}

// 点击锁定
func LockWithdrawOrderApi(ctx *context.Context, userId int, withdrawinfo *Withdrawinfo) (*model.BetResponse, error) {
	api := "/api/WithdrawOrder/LockWithdrawOrder"
	payloadStruct := &LockWithdrawOrder{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{userId, withdrawinfo.orderNo, withdrawinfo.createTime, config.Remark, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes2(model.ErrorLoggerType("/api/WithdrawOrder/LockWithdrawOrder请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse2(respBoy); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/WithdrawOrder/LockWithdrawOrder解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}
