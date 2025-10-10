package login

import (
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
)

type contextKey string

const (
	AuthTokenKey contextKey = "admin_auth_token"
)

// 登录请求结构体
type AdminLogin struct {
	UserName string `json:"userName"` // 账号
	Pwd      string `json:"pwd"`      // 密码
	model.BaseStruct
}

// 后台sit3003登录
func AdminSitLogin(ctx *context.Context) (*model.Response, *context.Context, error) {
	api := "/api/Login/Login"
	playStruct := &AdminLogin{}
	timestamp, random, language := request.GetTimeRandom()
	playList := []interface{}{config.ADMIN_UERNAME, config.ADMIN_PWD, random, language, "", timestamp}
	headerStruct := &model.BaseHeaderStruct{}
	headerUrl := config.ADMIN_SYSTEM_URL
	headerList := []interface{}{headerUrl, headerUrl, headerUrl}
	if respBody, _, err := request.PostGenericsFuncFlatten[AdminLogin, model.BaseHeaderStruct](headerUrl, api, playStruct, playList, headerStruct, headerList, request.StructToMap, request.InitStructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/Login 请求失败", err)), nil, err
	} else {
		token, err := model.GetJsonToken(string(respBody))
		if err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("后台token 解析失败", err)), nil, err
		}
		ctxToken := context.WithValue(*ctx, AuthTokenKey, token)
		if resp, err := model.ParseResponse(respBody); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/Login 响应解析失败", err)), nil, err
		} else {
			return resp, &ctxToken, nil
		}

	}

}

// 后台登录返回ctx,error
func RunAdminSitLogin() (*context.Context, error) {
	ctx := context.Background()
	_, ctxT, err := AdminSitLogin(&ctx)
	if err != nil {
		//fmt.Println("Login error:", err)
		logger.LogError("Admin Login error:%s", err)
		return nil, err
	}
	return ctxT, nil
}
