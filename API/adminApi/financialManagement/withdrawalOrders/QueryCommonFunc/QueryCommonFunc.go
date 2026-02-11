package querycommonfunc

import (
	"autoTest/API/utils"
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
			UserId         int   `json:"userId"`
			CreateTime     int64 `json:"createTime"`     // 订单创建时间
			LastUpdateTime int64 `json:"lastUpdateTime"` //订单最后更新时间
		} `json:"list"`
		Summary struct {
			TotalActualAmount float64 `json:"totalActualAmount"` // 总实际金额
		} `json:"summary"`
		TotalCount int `json:"totalCount"` // 总订单条数
	} `json:"data"`
}

// 查询的基础函数
func QueryCommonFuncApi(ctx *context.Context, api string, startTime, endTime int64) (*model.Response, []int, error) {
	payloadStruct := &QueryWithdrawaAmountStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []any{startTime, endTime, 1, "", 1, 2000, "Desc", random, language, "", timestamp}
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
				var totalTimeConsuming float64
				for _, v := range res.Data.List {
					userList = append(userList, v.UserId)
					// 计算每笔订单的耗时
					totalTimeConsuming += utils.CalculateMinutesFromTimestamps(v.CreateTime, v.LastUpdateTime)
				}
				count := res.Data.TotalCount
				fmt.Printf("提现总订单数:%d,提现实际总金额：%.2f,平均每笔订单的耗时分钟:%2.f\n", count, res.Data.Summary.TotalActualAmount, totalTimeConsuming/float64(count))
				return resp, userList, nil
			}
		}
	}
}
