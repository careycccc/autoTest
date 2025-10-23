package withdrawalorders

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
)

// 提现订单
type ConfirmWithdrawOrder struct {
	OrderNo    any `json:"orderNo"`
	UserId     any `json:"userId"`
	CreateTime any `json:"createTime"`
	Remark     any `json:"remark"`
	model.BaseStruct
}

// 点击确认出款
func ConfirmWithdrawOrderApi(ctx *context.Context, userId int, withdrawinfo Withdrawinfo) (*model.BetResponse, error) {
	api := "/api/WithdrawOrder/ConfirmWithdrawOrder"
	payloadStruct := &ConfirmWithdrawOrder{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{withdrawinfo.orderNo, userId, withdrawinfo.createTime, config.Remark, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes2(model.ErrorLoggerType("/api/WithdrawOrder/ConfirmWithdrawOrder请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse2(respBoy); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/WithdrawOrder/ConfirmWithdrawOrder解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}
