package login

import (
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"context"
)

type contextKey string

const (
	DeskAuthTokenKey contextKey = "desk_auth_token"
)

type userloginY1 struct {
	UserName    string `json:"userName"`
	Password    string `json:"password"`
	LoginType   string `json:"loginType"`
	DeviceId    string `json:"deviceId"`
	BrowserId   string `json:"browserId"`
	PackageName string `json:"packageName"`
	model.BaseStruct
}

// 前台y1需要输入账号和密码登录
func LoginY1(ctx context.Context, username, password string) (*model.Response, *context.Context, error) {
	api := "/api/Home/Login"
	payloadStruct := &userloginY1{}
	browserid := utils.GenerateCryptoRandomString(32)
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{username, password, "Mobile", "", browserid, "", random, language, "", timestamp}
	header_url := config.PLANT_H5
	base_url := config.SIT_WEB_API
	headerStruct := &model.DeskHeaderTenantIdStruct{}
	headerList := []interface{}{"3003", header_url, header_url, header_url}
	respBody, _, err := request.PostGenericsFuncFlatten[userloginY1, model.DeskHeaderTenantIdStruct](base_url, api, payloadStruct, payloadList, headerStruct, headerList, request.StructToMap, request.InitStructToMap)
	if err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/Login 请求失败", err)), nil, err
	}
	token, err := model.GetJsonToken(string(respBody))
	if err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/Login,token 解析失败", err)), nil, err
	}
	ctxToken := context.WithValue(ctx, DeskAuthTokenKey, token)
	if resp, err := model.ParseResponse(respBody); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/Login 响应解析失败", err)), nil, err
	} else {
		return resp, &ctxToken, nil
	}
}

// 主要是返回上下文
func ReturnContextLoginY1(username, password string) (*context.Context, error) {
	ctx := context.Background()
	if _, ctxToken, err := LoginY1(ctx, username, password); err != nil {
		logger.LogError("LoginY1['/api/Home/Login']登录失败", err)
		return nil, err
	} else {
		return ctxToken, nil
	}
}
