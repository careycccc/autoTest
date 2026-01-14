package champion

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
)

type ChampionStruct struct {
	ChampionId any `json:"championId"`
	model.BaseStruct
}

// 加入竞标赛
// championId 竞标赛活动id
func JoinChampion(ctx *context.Context, championId int) (*model.Response, error) {
	api := "/api/Activity/AddChampion"
	payloadStruc := &ChampionStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{championId, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruc, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/AddChampion请求失败", err)), err
	} else {
		if string(respBoy) == "" {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/AddChampion返回空", err)), err
		}
		// 处理正常返回
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/AddChampion解析返回失败", err)), err
		} else {
			return resp, nil
		}
	}
}
