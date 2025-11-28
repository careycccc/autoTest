package withdrawalorders

import (
	querycommonfunc "autoTest/API/adminApi/financialManagement/withdrawalOrders/QueryCommonFunc"
	"autoTest/API/adminApi/login"
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/logger"
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

// 查询昨天所有的提现了的会员,返回提现订单所有的会员id列表
func QueryWithdrawaAmount(ctx *context.Context, startTime, endTime int64) (*model.Response, []int, error) {
	api := "/api/WithdrawOrder/GetWithdrawOrderPageList"
	if resp, list, err := querycommonfunc.QueryCommonFuncApi(ctx, api, startTime, endTime); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/WithdrawOrder/GetWithdrawOrderPageList报错", err)), nil, err
	} else {
		return resp, list, nil
	}

}

// 运行查询提现订单的函数
func RunWithDrawCase() []int {
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		logger.LogError("提现订单的后台登录报错", err)
		return nil
	} else {
		start, end := utils.GetYesterdayStartEndMilli()
		if _, userList, err := QueryWithdrawaAmount(ctx, start, end); err != nil {
			logger.LogError("提现订单的查询的报错信息", err)
			return nil
		} else {
			// logger.Logger.Info("提现订单的查询结果", userList)
			return userList
		}
	}
}
