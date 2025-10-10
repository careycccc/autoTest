package financialmanagement

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"sync"
)

// 人工充值接口M
type ManualRecharge struct {
	ArtificialRechargeType int8   `json:"artificialRechargeType"`
	RechargeAmount         int64  `json:"rechargeAmount"` // 充值金额
	Remark                 string `json:"remark"`         // 备注
	AmountOfCode           int8   `json:"amountOfCode"`   // 打码量 null表示默认，数字表示倍数
	UserId                 int64  `json:"userId"`
	model.BaseStruct
}

/*
*
userid 用户id
rechargeAmount 充值金额
amountOfCode 打码量
*/
func ArtificialRechargeFunc(ctx *context.Context, userid, rechargeAmount int64, amountOfCode int8, wg *sync.WaitGroup) (*model.Response, error) {
	wg.Add(1)
	defer wg.Done()
	api := "/api/ArtificialRechargeRecord/ArtificialRecharge"
	payloadStruct := &ManualRecharge{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{3, rechargeAmount, "carey3003", amountOfCode, userid, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/ArtificialRechargeRecord/ArtificialRecharge请求失败", err)), err
	} else {
		logger.Logger.Info("充值成功的金额", rechargeAmount)
		if res, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/ArtificialRechargeRecord/ArtificialRecharge解析失败", err)), err
		} else {
			return res, nil
		}
	}
}
