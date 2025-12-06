package rechargeorders

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

// 充值订单

type RechargeOrder struct {
	RechargeState string `json:"rechargeState"` // 充值状态 Payed=已支付，Pending=待支付，Failed=支付失败
	StartTime     int64  `json:"startTime"`     // 开始时间
	EndTime       int64  `json:"endTime"`       // 结束时间
	DateType      int    `json:"dateType"`      // 日期类型  1
	model.QueryPayloadStruct
}

// 定义顶层响应结构体
type RechargeOrderResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Data   `json:"data"`
}

// Data 结构体包含了列表和总计数
type Data struct {
	List       []OrderItem `json:"list"`
	Summary    Summary     `json:"summary"`
	TotalCount int         `json:"totalCount"`
}

// OrderItem 结构体用于提取列表中的订单详情
type OrderItem struct {
	UserId              int     `json:"userId"`
	ActualAmount        float64 `json:"actualAmount"`
	RechargeChannelName string  `json:"rechargeChannelName"`
	RechargeChannelType string  `json:"rechargeChannelType"`
	RechargeType        string  `json:"rechargeType"`
}

// Summary 结构体用于提取汇总信息
type Summary struct {
	TotalAmount float64 `json:"totalAmount"`
}

/*
查询的是已支付的订单
返回充值订单列表 RechargeOrderResponse
*
*/
func GetRechargeOrderPageListApi(ctx *context.Context, startTime, endTime int64) (*model.Response, *RechargeOrderResponse, error) {
	api := "/api/RechargeOrder/GetRechargeOrderPageList"
	payloadStruct := &RechargeOrder{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{"Payed", startTime, endTime, 1, 1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/RechargeOrder/GetRechargeOrderPageList请求失败", err)), &RechargeOrderResponse{}, err
	} else {
		if res, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/RechargeOrder/GetRechargeOrderPageList解析失败", err)), &RechargeOrderResponse{}, err
		} else {
			var rechargeOrderResponse RechargeOrderResponse
			if err := json.Unmarshal(respBoy, &rechargeOrderResponse); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/RechargeOrder/GetRechargeOrderPageList【rechargeOrderResponse解析失败】", err)), &RechargeOrderResponse{}, err
			} else {
				return res, &rechargeOrderResponse, nil
			}
		}
	}
}
