package topup

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
)

// arpay充值方式

type ArpayRequest struct {
	RechargeCategoryID int    `json:"rechargeCategoryId"` //700055
	ReturnURL          string `json:"returnUrl"`
	URLInfo            string `json:"urlInfo"`
	VendorID           int    `json:"vendorId"` // 默认0
	Amount             int    `json:"amount"`   // 充值金额
	model.BaseStruct
}

/*
arp充值
amount 充值金额
*
*/
func ArpayRequestApi(ctx *context.Context, amount int) (*model.Response, error) {
	api := "/api/Recharge/DepositRecharge"
	payloadStruct := &ArpayRequest{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{700055, "https://arplatsaassit1.club/#/main", "https://arplatsaassit1.club,status/rechargeStatus", 0, amount, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Recharge/DepositRecharge 请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Recharge/DepositRecharge 解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}
