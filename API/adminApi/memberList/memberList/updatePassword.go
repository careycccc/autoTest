package memberlist

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
)

type UpdataPasswordstruct struct {
	UserId   int64  `json:"userId"`
	Password string `json:"password"`
	model.BaseStruct
}

// 后台修改密码
func UpdatePassword(ctx *context.Context, userid int64, password string) (*model.Response, error) {
	api := "/api/Users/UpdatePassword"
	payloadStruct := &UpdataPasswordstruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{userid, password, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/UpdatePassword请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/UpdatePassword解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}
