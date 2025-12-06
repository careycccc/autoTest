package memberreports

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

// 会员登录报表

type MemberLoginReport struct {
	MemberIdType int    `json:"memberIdType"` // 会员ID类型
	StartTime    string `json:"startTime"`    // 开始时间
	EndTime      string `json:"endTime"`      // 结束时间
	model.QueryPayloadStruct
}

type MemberLoginReportResponse struct {
	Data struct {
		TotalCount int `json:"totalCount"` // 会员登录人数总数
	} `json:"data"`
}

/*
会员登录报表
返回会员登录人数总数
*
*/
func GetUserLoginLogPageListApi(ctx *context.Context, startTime, endTime string) (*model.Response, int, error) {
	api := "/api/RptUserInfo/GetUserLoginLogPageList"
	payloadStruct := &MemberLoginReport{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{1, startTime, endTime, 1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/RptUserInfo/GetUserLoginLogPageList请求失败", err)), 0, err
	} else {
		var memberLoginReportResponse MemberLoginReportResponse
		if err := json.Unmarshal(respBoy, &memberLoginReportResponse); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/RptUserInfo/GetUserLoginLogPageList【memberLoginReportResponse解析失败】", err)), 0, err
		} else {
			if resp, err := model.ParseResponse(respBoy); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/RptUserInfo/GetUserLoginLogPageList 解析失败", err)), 0, err
			} else {
				return resp, memberLoginReportResponse.Data.TotalCount, nil
			}
		}
	}
}
