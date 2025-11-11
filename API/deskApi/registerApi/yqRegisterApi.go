package registerapi

import (
	membermanagement "autoTest/API/adminApi/memberManagement"
	login "autoTest/API/deskApi/loginApi"
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"context"
	"fmt"
)

// 总代注册，没有上级的注册方式

type MobileAutoLoginStruct struct {
	UserName            string `json:"userName"`
	VerifyCode          string `json:"verifyCode"`
	RegisterDevice      string `json:"registerDevice"`
	RegisterFingerprint string `json:"registerFingerprint"`
	InviteCode          string `json:"inviteCode"`
	Rrack               string `json:"track"`
	PackageName         string `json:"packageName"`
	model.BaseStruct
}

// 要求注册，通过链接邀请,前台的登录

// 传入一个用户，手机号码+区号
func MobileAutoLoginFunc(userName string) (*model.BetResponse, context.Context, error) {
	api := "/api/Home/MobileAutoLogin"
	// 发送验证码
	// 获取验证码
	ctx := context.Background()
	if res, verifyCode, err := membermanagement.SendToGetVerCode(&ctx, 18, userName); err != nil {
		return model.ResponseToBetResponse(res), nil, err
	} else {
		registerFingerprint := utils.GenerateCryptoRandomString(32)
		timestamp, random, language := request.GetTimeRandom()
		payloadList := []interface{}{userName, verifyCode, "", registerFingerprint, "", "", "", random, language, "", timestamp}
		// payload的构建
		payloadStruct := &MobileAutoLoginStruct{}
		if repBoy, _, err := requstmodle.DeskTrodRegRequest2(&ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Home/MobileAutoLogin登录/注册请求失败", err)), nil, err
		} else {
			// 还需要解析出token进行保存
			token, err := model.GetJsonToken(string(repBoy))
			if err != nil {
				logger.LogError("/api/Home/MobileAutoLogin,token 解析失败", err)
				errs := fmt.Errorf("/api/Home/MobileAutoLogin,token 解析失败%s", err)
				return model.HandlerErrorRes2(errs), nil, err
			}
			ctxToken := context.WithValue(ctx, login.DeskAuthTokenKey, token)
			if betResponse, err := model.ParseResponse2(repBoy); err != nil {
				return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Home/MobileAutoLogin响应结果解析失败", err)), nil, err
			} else {
				return betResponse, ctxToken, nil
			}
		}
	}
}

// 嵌套结构体 Track
type Track struct {
	IsTrusted bool  `json:"isTrusted"`
	Vts       int64 `json:"_vts"`
}

// 主结构体
type RegisterStruct struct {
	UserName            string `json:"userName"`
	VerifyCode          string `json:"verifyCode"`
	InviteCode          string `json:"inviteCode"`
	RegisterFingerprint string `json:"registerFingerprint"`
	Track               Track  `json:"track"`
	Language            string `json:"language"`
	Random              int64  `json:"random"`
	Signature           string `json:"signature"`
	Timestamp           int64  `json:"timestamp"`
}

// userName 用户名  invitationCode邀请码
func RegisterMobileLoginFunc(userName, invitationCode string) (*model.BetResponse, context.Context, error) {
	api := "/api/Home/MobileAutoLogin"
	// 发送验证码
	// 获取验证码
	ctx := context.Background()
	if res, verifyCode, err := membermanagement.SendToGetVerCode(&ctx, 18, userName); err != nil {
		return model.ResponseToBetResponse(res), nil, err
	} else {
		registerFingerprint := utils.GenerateCryptoRandomString(32)
		timestamp, random, language := request.GetTimeRandom()
		// payload的构建
		payloadStruct := &RegisterStruct{}
		payloadList := []interface{}{userName, verifyCode, invitationCode, registerFingerprint, Track{IsTrusted: true, Vts: timestamp}, language, random, "", timestamp}
		if repBoy, _, err := requstmodle.DeskTrodRegRequest2(&ctx, api, payloadStruct, payloadList, request.InitStructToMap); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Home/MobileAutoLogin登录/注册请求失败", err)), nil, err
		} else {
			// 解析tokne并进行保存
			// 还需要解析出token进行保存
			token, err := model.GetJsonToken(string(repBoy))
			if err != nil {
				logger.LogError("/api/Home/MobileAutoLogin,token 解析失败", err)
				errs := fmt.Errorf("/api/Home/MobileAutoLogin,token 解析失败%s", err)
				return model.HandlerErrorRes2(errs), nil, err
			}
			ctxToken := context.WithValue(ctx, login.DeskAuthTokenKey, token)
			if betResponse, err := model.ParseResponse2(repBoy); err != nil {
				return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Home/MobileAutoLogin响应结果解析失败", err)), nil, err
			} else {
				return betResponse, ctxToken, nil
			}
		}
	}
}
