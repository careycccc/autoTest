package memberactivityblacklist

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
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
	payloadList := []interface{}{config.Remark, userid, activityBlockType, random, language, "", timestamp}
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
	payloadList := []interface{}{config.Remark, userid, activityBlockType, random, language, "", timestamp}
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
