package memberlist

import (
	"autoTest/API/adminApi/login"
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
)

//修改用户状态

type UpdateUserState struct {
	UserId int `json:"userId"`
	State int `json:"state"`
	model.BaseStruct
}

func UpdateUserStateApi(ctx *context.Context,userId,state int) (*model.Response,error) {
	api := "/api/Users/UpdateState"
	payloadStruct := &UpdateUserState{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{userId, state, random,language,"",timestamp}
	if respBoy,_,err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList,request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/UpdateState请求失败",err)), err
	} else {
		if string(respBoy) == ""{
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/UpdateState响应为空",err)), err
		}else {
			if resp,err := model.ParseResponse(respBoy); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/UpdateState响应解析失败",err)), err
			}else {
				return resp, nil
			}
		}
	}
}	

// 解冻
func UnfreezeUser(userId int)  {
	if ctx ,err := login.RunAdminSitLogin(); err != nil {
		logger.Logger.Warn("解冻用户，后台登录失败")
	    return 
	}else {
		resp,err := UpdateUserStateApi(ctx,userId,1)
		if err != nil {
			logger.Logger.Warn("解冻用户失败",userId)
		}

		logger.Logger.Info("解冻用户成功",userId,resp)
	}
}