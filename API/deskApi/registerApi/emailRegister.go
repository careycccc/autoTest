package registerapi

import (
	membermanagement "autoTest/API/adminApi/memberManagement"
	login "autoTest/API/deskApi/loginApi"
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"context"
	"fmt"
)

type EmailRegister struct {
	UserName       string `json:"userName"`   // 电子邮箱地址
	InviteCode     string `json:"inviteCode"` // 邀请码
	LoginType      string `json:"loginType"`  // "Email"
	TurnstileToken string `json:"turnstileToken"`
	Password       string `json:"password"` // 密码
	Code           string `json:"code"`     // 验证吗
	model.BaseStruct
}

// 邮箱的方式进行注册，需要邀请去吗
func EmailRegisterApi(userName, inviteCode string) (*model.Response, context.Context, error) {
	// 发送验证码，获取验证码
	ctx := context.Background()
	if res, verifyCode, err := membermanagement.SendToGetVerCode(&ctx, 2, int8(2), userName); err != nil {
		logger.Logger.Warn("验证码发送信息", res)
		logger.LogError("验证码发送信息", err)
		return res, nil, err
	} else {
		api := "/api/Home/Register"
		payload := &EmailRegister{}
		timestamp, random, language := request.GetTimeRandom()
		payloadList := []any{userName, inviteCode, "Email", "", config.ADMIN_PWD, verifyCode, random, language, "", timestamp}
		if respBody, _, err := requstmodle.DeskTrodRegRequest2(&ctx, api, payload, payloadList, request.StructToMap); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/Register 请求失败", err)), nil, err
		} else {
			if string(respBody) == "" {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/Register respBoy为空", err)), nil, err
			} else {
				// 解析token出来
				token, err := model.GetJsonToken(string(respBody))
				if err != nil {
					return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/Register,token获取token失败", err)), nil, err
				}
				ctxToken := context.WithValue(ctx, login.DeskAuthTokenKey, token)
				if resp, err := model.ParseResponse(respBody); err != nil {
					return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/Register解析失败", err)), nil, err
				} else {
					return resp, ctxToken, nil
				}
			}
		}
	}
}

// RunEmailregeister 是一个执行邮箱注册的函数
// 该函数生成随机邮箱地址，并调用注册API进行注册
func RunEmailregeister() {
	// 初始化邀请码为空字符串
	inviteCode := "PBZXTZW"
	// 生成随机邮箱地址作为用户名
	userName := utils.GenerateRandomEmail()

	// 调用邮箱注册API，传入邀请码和用户名
	// 处理可能的错误情况
	if resp, ctx, err := EmailRegisterApi(inviteCode, userName); err != nil {
		// 记录注册失败的警告日志，包含API响应
		logger.Logger.Warn("注册失败", resp)
		// 记录详细的错误日志
		logger.LogError("注册失败", err)
	} else {
		// 注册成功时打印响应和上下文信息
		fmt.Println(resp, ctx)
		fmt.Println("注册的邮箱用户名", userName)
	}
}
