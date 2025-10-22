package withdrawcash

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
)

// 设置提现密码的接口
type SetWithdrawPasswordStruct struct {
	WithdrawPassword any `json:"withdrawPassword"` // 设置提现密码
	model.BaseStruct
}

// 设置提现密码的请求
func SetWithdrawPasswordApi(ctx *context.Context) (*model.BetResponse, error) {
	api := "/api/User/SetWithdrawPassword"
	payloadStruct := &SetWithdrawPasswordStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{config.WithdrawPassword, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes2(model.ErrorLoggerType("/api/User/SetWithdrawPassword设置提现密码请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse2(respBoy); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/User/SetWithdrawPassword设置提现密码解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}
