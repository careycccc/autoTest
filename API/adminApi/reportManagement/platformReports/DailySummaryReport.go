package platformreports

import (
	"autoTest/API/adminApi/login"
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
	"fmt"
)

// 平台报表-> 每日汇总报表

type GetPlatRptStatisticPage struct {
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	DateType  int    `json:"dateType"` // 默认值1
	model.QueryPayloadStruct
}

type Report struct {
	ReportDate        string  `json:"reportDate"`
	TenantID          int     `json:"tenantId"`
	TenantName        string  `json:"tenantName"`
	RegisterCount     int     `json:"registerCount"`
	LoginCount        int     `json:"loginCount"`
	BetUserCount      int     `json:"betUserCount"`
	BetAmount         float64 `json:"betAmount"`
	FeeAmount         float64 `json:"feeAmount"`
	WinAmount         float64 `json:"winAmount"`
	ActivityAmount    float64 `json:"activityAmount"`
	RechargeCount     int     `json:"rechargeCount"`
	RechargeAmount    float64 `json:"rechargeAmount"`
	WithdrawCount     int     `json:"withdrawCount"`
	WithdrawAmount    float64 `json:"withdrawAmount"`
	PlatNetDeposit    float64 `json:"platNetDeposit"`
	PlatWinLoseAmount float64 `json:"platWinLoseAmount"`
	PlatNetProfit     float64 `json:"platNetProfit"`
	BetCount          int     `json:"betCount"`
}

type ReportResponse struct {
	Data struct {
		List []Report `json:"list"`
	} `json:"data"`
}

/*
每日汇总报表
*
*/
func GetPlatRptStatisticPageList(ctx *context.Context) (*model.Response, *Report, error) {
	api := "/api/RptUserActivity/GetPlatRptStatisticPageList"
	starttime, endtime := utils.GetDayStartEnd()
	payloadStruct := &GetPlatRptStatisticPage{}
	timestamp, random, language := request.GetTimeRandom()
	// 临时处理总条数为2000条
	payloadList := []interface{}{starttime, endtime, 1, 1, 2000, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/RptUserActivity/GetPlatRptStatisticPageList请求失败", err)), &Report{}, err
	} else {
		var reportResponse ReportResponse
		if err := json.Unmarshal(respBoy, &reportResponse); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/RptUserActivity/GetPlatRptStatisticPageList【reportResponse解析失败】", err)), &Report{}, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/RptUserActivity/GetPlatRptStatisticPageList 解析失败", err)), &Report{}, err
		} else {
			return resp, &reportResponse.Data.List[0], nil
		}
	}
}

// 运行获取每日汇总报表
func RunGetPlatRptStatisticPageList() {
	ctx := context.Background()
	_, ctxToken, err := login.AdminSitLogin(&ctx)
	if err != nil {
		logger.LogError("Login error:%s", err)
		return
	}
	if _, report, err := GetPlatRptStatisticPageList(ctxToken); err != nil {
		logger.LogError("GetPlatRptStatisticPageList运行的错误信息 error:%s", err)
		return
	} else {
		fmt.Println("GetPlatRptStatisticPageList运行的结果 response:", report)
	}
}
