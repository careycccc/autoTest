package topup

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
)

// 进行充值操作
type GoodsDepositRechargeStruct struct {
	RechargeCategoryId any `json:"rechargeCategoryId"`
	ReturnUrl          any `json:"returnUrl"`
	UrlInfo            any `json:"urlInfo"`
	VendorId           any `json:"vendorId"`
	RechargeGoodsId    any `json:"rechargeGoodsId"`
	model.BaseStruct
}

/*
进行充值操作
rechargeCategoryId
vendorId
rechargeGoodsId // 供应商id bankcard 0
*
*/
func GoodsDepositRechargeApi(ctx *context.Context, rechargeCategoryId, rechargeGoodsId int) (*model.Response, error) {
	api := "/api/Recharge/GoodsDepositRecharge"
	payloadStruct := &GoodsDepositRechargeStruct{}
	returnUrl := config.GoodsDeposit_URL + "#/main"
	urlInfo := config.GoodsDeposit_URL + ",status/rechargeStatus"
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{rechargeCategoryId, returnUrl, urlInfo, 0, rechargeGoodsId, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Recharge/GoodsDepositRecharge请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Recharge/GoodsDepositRecharge解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

// 本地充值
type RechargeRequest struct {
	RechargeCategoryID int64        `json:"rechargeCategoryId"`
	ReturnURL          string       `json:"returnUrl"`
	URLInfo            string       `json:"urlInfo"`
	VendorID           int          `json:"vendorId"`
	CustomerInfo       CustomerInfo `json:"customerInfo"`
	RechargeGoodsID    int          `json:"rechargeGoodsId"`
	Language           string       `json:"language"`
	Random             int64        `json:"random"`
	Signature          string       `json:"signature"`
	Timestamp          int64        `json:"timestamp"`
}

type CustomerInfo struct {
	AccountNo  string `json:"accountNo"`
	HolderName string `json:"holderName"`
}
