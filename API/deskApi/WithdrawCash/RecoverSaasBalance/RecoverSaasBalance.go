package recoversaasbalance

import (
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

type RecoverSaasBalanceResponse struct {
	Data struct {
		Balance float64 `json:"balance"` // 用户的金额
	} `json:"data"`
}

// 获取当前用户的金额，就是充值和提现页面的刷新按钮
// 返回当前用户的金额
func RecoverSaasBalance(ctx *context.Context) (*model.Response, float64, error) {
	api := "/api/ThirdGame/RecoverSaasBalance"
	payloadStruct, payloadList := utils.BaseStructHandler()
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.InitStructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("", err)), 0, err
	} else {
		var balanceResp RecoverSaasBalanceResponse
		if err := json.Unmarshal(respBoy, &balanceResp); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("", err)), 0, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("", err)), 0, err
		} else {
			return resp, balanceResp.Data.Balance, nil
		}
	}
}
