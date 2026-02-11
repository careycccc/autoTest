package memberactivityblacklist

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

// 会员列表的黑名单

// 添加黑名单
type UserActivityBlockAdd struct {
	Remark            string `json:"remark"`            // 备注
	UserIds           string `json:"userIds"`           // 会员ID列表
	ActivityBlockType int    `json:"activityBlockType"` // 黑名单类型
	model.BaseStruct
}

// 删除黑名单
type UserActivityBlockDelete struct {
	Reason            string `json:"reason"`            // 删除原因
	UserId            int    `json:"userId"`            // 会员ID列表
	ActivityBlockType int    `json:"activityBlockType"` // 黑名单类型
	model.BaseStruct
}

/*
添加会员黑名单的请求
userid: 会员ID
activityBlockType: 黑名单类型
*
*/
func UserActivityBlockAddApi(ctx *context.Context, userid string, activityBlockType int) (*model.Response, error) {
	api := "/api/UserActivityBlock/Add"
	payloadStruct := &UserActivityBlockAdd{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []any{config.Remark, userid, activityBlockType, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/UserActivityBlock/Add 请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/UserActivityBlock/Add 响应解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

/*
删除会员黑名单的请求
userid: 会员ID
activityBlockType: 黑名单类型
*
*/
func UserActivityBlockDeleteApi(ctx *context.Context, userid int, activityBlockType int) (*model.Response, error) {
	api := "/api/UserActivityBlock/Delete"
	payloadStruct := &UserActivityBlockDelete{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []any{config.Remark, userid, activityBlockType, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/UserActivityBlock/Delete 请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/UserActivityBlock/Delete 响应解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

// 查询会员是否在某类型的黑名单中
type UserActivityBlockIsInBlock struct {
	UserId            string `json:"userId"`            // 会员ID
	ActivityBlockType int    `json:"activityBlockType"` // 黑名单类型
	model.QueryPayloadStruct
}

type UserActivityBlockIsInBlockResList struct {
	ActivityBlockType int    `json:"activityBlockType"`
	CreateTime        int64  `json:"createTime"`
	Creator           string `json:"creator"`
	Remark            string `json:"remark"`
	UserId            int    `json:"userId"`
}

type UserActivityBlockIsInBlockRes struct {
	Data struct {
		List []UserActivityBlockIsInBlockResList `json:"list"`
	} `json:"data"`
}

// 查询会员是否在某类型的黑名单中
func UserActivityBlockIsInBlockApi(ctx *context.Context, userid string, activityBlockType int) (*model.Response, bool, error) {
	api := "/api/UserActivityBlock/GetPageList"
	payloadStruct := &UserActivityBlockIsInBlock{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []any{userid, activityBlockType, 1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/UserActivityBlock/GetPageList 请求失败", err)), false, err
	} else {
		if string(respBoy) == "" {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/UserActivityBlock/GetPageList respBoy 为空", err)), false, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/UserActivityBlock/GetPageList 响应解析失败", err)), false, err
		} else {
			var res UserActivityBlockIsInBlockRes
			if err := json.Unmarshal(respBoy, &res); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/UserActivityBlock/GetPageList[UserActivityBlockIsInBlockRes] 响应解析失败", err)), false, err
			}
			if len(res.Data.List) == 0 {
				return resp, false, nil
			}
			return resp, true, nil
		}
	}
}

// 调用会员黑名单的运行方法
func RunBlackList() {
	// 后台登录
	// if ctx, err := login.RunAdminSitLogin(); err != nil {
	// 	logger.Logger.Warn("会员黑名单后台登录失败", err)
	// 	return
	// } else {
	// 	// 添加黑名单
	// 	if _, conunt, err := memberreports.GetUserLoginLogPageListApi(ctx); err != nil {
	// 		logger.LogError("获取会员登录报表失败", err)
	// 		return
	// 	} else {
	// 		logger.Logger.Info("会员登录总人数", conunt)
	// 	}
	// }
}
