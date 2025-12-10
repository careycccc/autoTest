package rechargeorders

import (
	"autoTest/API/adminApi/login"
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
	"fmt"
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
	CreateTime          int64   `json:"createTime"`
	OrderNo             string  `json:"orderNo"`
	RechargeState       string  `json:"rechargeState"`
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

// 根据id获取订单详情
type GetRechargeOrderById struct {
	UserID          int64  `json:"userId"`
	StartTime       int64  `json:"startTime"`
	EndTime         int64  `json:"endTime"`
	MinActualAmount *int64 `json:"minActualAmount"`
	MaxActualAmount *int64 `json:"maxActualAmount"`
	DateType        int    `json:"dateType"`
	model.QueryPayloadStruct
}

/*
根据id获取订单
**/

func GetRechargeOrderByIdApi(ctx *context.Context, userid int, startTime, endTime int64) (*model.Response, *RechargeOrderResponse, error) {
	api := "/api/RechargeOrder/GetRechargeOrderPageList"
	payloadStruct := &GetRechargeOrderById{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{userid, startTime, endTime, nil, nil, 1, 1, 20, "Desc", random, language, "", timestamp}
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

type OrderDetail struct {
	UserID       int    `json:"userId"`
	OrderNo      string `json:"orderNo"`
	ActualAmount int64  `json:"actualAmount"`
	CreateTime   int64  `json:"createTime"`
	Remark       string `json:"remark"`
	model.BaseStruct
}

/*
点击确认补单
actualAmount  充值金额
createTime 订单创建时间
orderNo 订单号
*
*/
func ClickRechargeOrderApi(ctx *context.Context, userId int, actualAmount float64, createTime int64, orderNo string) (*model.Response, error) {
	api := "/api/RechargeOrder/ManualAuditRechargeOrder"
	payloadStruct := &OrderDetail{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{userId, orderNo, actualAmount, createTime, config.Remark, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/RechargeOrder/ManualAuditRechargeOrder请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/RechargeOrder/ManualAuditRechargeOrder解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

// 只要传入一个userid就可以点击确认补单
func RunToUpArpay(userid int) error {
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		logger.LogError("充值订单的手动补单后台登录失败", err)
		return err
	} else {
		_, start, _, end, err := utils.ParseTimeRangeToTimestamp(config.StartTime, config.EndTime)
		if err != nil {
			logger.LogError("充值订单的手动补单时间转换失败", err)
			return err
		}
		if _, oderInfo, err := GetRechargeOrderByIdApi(ctx, userid, start, end); err != nil {
			logger.LogError("充值订单的手动补单获取订单信息失败", err)
			return err
		} else {
			if len(oderInfo.Data.List) > 0 {
				for _, v := range oderInfo.Data.List {
					// 已支付状态的订单不需要补单
					if v.RechargeState == "Payed" {
						continue
					}
					if _, err := ClickRechargeOrderApi(ctx, v.UserId, v.ActualAmount, v.CreateTime, v.OrderNo); err != nil {
						logger.LogError("充值订单的手动补单失败", err)
						return err
					} else {
						return nil
					}
				}
			}
			return fmt.Errorf("没有查询到该账号的订单")
		}

	}
}
