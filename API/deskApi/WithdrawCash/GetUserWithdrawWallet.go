package withdrawcash

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

type GetUserWithdrawWalletStruct struct {
	WithdrawType any `json:"withdrawType"`
	model.BaseStruct
}

type getUserWithdrawResponse struct {
	Data struct {
		WalletId string `json:"walletId"`
	} `json:"data"`
}

// 获取提现的WithdrawWalletid
// 返回WithdrawWalletid
func GetUserWithdrawWallet(ctx *context.Context, withdrawType string) (*model.Response, string, error) {
	api := "/api/Withdraw/GetUserWithdrawWallet"
	payloadStruct := &GetUserWithdrawWalletStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{withdrawType, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Withdraw/GetUserWithdrawWallet请求失败", err)), "", err
	} else {
		var getWithdraw getUserWithdrawResponse
		if err := json.Unmarshal(respBoy, &getWithdraw); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Withdraw/GetUserWithdrawWallet[getWithdraw]解析失败", err)), "", err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Withdraw/GetUserWithdrawWallet解析失败", err)), "", err
		} else {
			return resp, getWithdraw.Data.WalletId, nil
		}
	}
}
