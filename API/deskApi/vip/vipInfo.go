package vip

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
)

// 获取前台的vip信息
func GetVipLevelConfig(ctx *context.Context) (*model.Response, *context.Context, error) {
	api := "/api/VipLevel/GetVipLevelConfig"
	payloadStruct := model.BaseStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadData := []interface{}{random, language, "", timestamp}
	if respBody, _, err := requstmodle.DeskTenAuthorRequest[model.BaseStruct](ctx, api, &payloadStruct, payloadData, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("api/VipLevel/GetVipLevelConfig 请求失败", err)), nil, err
	} else {
		if response, err := model.ParseResponse(respBody); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("api/VipLevel/GetVipLevelConfig 响应结果解析失败", err)), nil, err
		} else {

			return response, ctx, nil
		}
	}
}
