package querycommonfunc

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
	"fmt"
)

type QueryWithdrawaAmountStruct struct {
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	DateType  int    `json:"dateType"` //0申请时间  1最后更新时间
	SortField string `json:"sortField"`
	model.QueryPayloadStruct
}

type QueryWithdrawaResponse struct {
	Data struct {
		List []struct {
			UserId int `json:"userId"`
		} `json:"list"`
		Summary struct {
			TotalActualAmount float64 `json:"totalActualAmount"` // 总实际金额
		} `json:"summary"`
		TotalCount int `json:"totalCount"` // 总订单条数
	} `json:"data"`
}

type QueryWithdrawaSummary struct {
	UserIdList        []int
	TotalActualAmount float64
	TotalCount        int
}

// 查询的基础函数
func QueryCommonFuncApi(ctx *context.Context, api string, startTime, endTime int64) (*model.Response, []int, error) {
	payloadStruct := &QueryWithdrawaAmountStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{startTime, endTime, 1, "", 1, 2000, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType(api, err)), nil, err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType(api, err)), nil, err
		} else {
			var res QueryWithdrawaResponse
			if err := json.Unmarshal(respBoy, &res); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType(api, err)), nil, err
			} else {
				userList := []int{}
				for _, v := range res.Data.List {
					userList = append(userList, v.UserId)
				}
				fmt.Printf("提现总订单数:%d,提现实际总金额：%.2f\n", res.Data.TotalCount, res.Data.Summary.TotalActualAmount)
				return resp, userList, nil
			}
		}
	}
}
