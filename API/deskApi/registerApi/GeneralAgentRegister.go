package registerapi

import (
	membermanagement "autoTest/API/adminApi/memberManagement"
	login "autoTest/API/deskApi/loginApi"
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"context"
)

// 总代注册，没有上级的注册方式

type GeneralRegiterStruct struct {
	UserName            string `json:"userName"`
	VerifyCode          string `json:"verifyCode"`
	RegisterDevice      string `json:"registerDevice"`
	RegisterFingerprint string `json:"registerFingerprint"`
	InviteCode          string `json:"inviteCode"`
	Rrack               string `json:"track"`
	PackageName         string `json:"packageName"`
	model.BaseStruct
}

// 总代的方式进行前台注册
func GeneralAgentRegister(userName string) (*model.Response, *context.Context, error) {
	api := "/api/Home/MobileAutoLogin"
	// 发送验证码，获取验证码
	ctx := context.Background()
	if res, verifyCode, err := membermanagement.SendToGetVerCode(&ctx, userName); err != nil {
		return res, nil, err
	} else {
		// 随机浏览器指纹
		Fingerprint := utils.GenerateCryptoRandomString(32)
		devices := utils.GenerateCryptoRandomString(16)
		payloadStruct := &GeneralRegiterStruct{}
		timestamp, random, language := request.GetTimeRandom()
		payloadList := []interface{}{userName, verifyCode, devices, Fingerprint, "", "", "", random, language, "", timestamp}
		if respBody, _, err := requstmodle.DeskTrodRegRequest(&ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/MobileAutoLogin请求失败", err)), nil, err
		} else {
			// 解析tokne并进行保存
			// 还需要解析出token进行保存
			token, err := model.GetJsonToken(string(respBody))
			if err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/MobileAutoLogin,token获取token失败", err)), nil, err
			}
			ctxToken := context.WithValue(ctx, login.DeskAuthTokenKey, token)
			if resp, err := model.ParseResponse(respBody); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Home/MobileAutoLogin解析失败", err)), nil, err
			} else {
				return resp, &ctxToken, nil
			}
		}
	}
}
