package memberlist

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

// 会员列表

type GetUserApistruct struct {
	Account  string `json:"account"`  // 用户账号
	PageNo   int    `json:"pageNo"`   // 页码
	PageSize int    `json:"pageSize"` // 每页多少数据
	OrderBy  string `json:"orderBy"`  // 排序
	model.BaseStruct
}

// 提取userid
type Useridstruct struct {
	Data struct {
		List []struct {
			UserId int64 `json:"userId"`
		} `json:"list"`
	} `json:"data"`
}

/*
输入用户电话号码，返回用户id
*
*/
func GetUserIdApi(ctx *context.Context, account string) (*model.Response, int64, error) {
	api := "/api/Users/GetPageList"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &GetUserApistruct{}
	payloadList := []interface{}{account, 1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList请求失败", err)), -1, err
	} else {
		// 提取用户id
		var response Useridstruct
		if err := json.Unmarshal(respBoy, &response); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList提取用户id解析失败", err)), -1, err
		} else {
			if resp, err := model.ParseResponse(respBoy); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetPageList响应解析失败", err)), -1, err
			} else {
				return resp, response.Data.List[0].UserId, nil
			}
		}

	}
}
