package getuserinfo

import (
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

// 获取前台的用户信息
type UserInfoDataResponse struct {
	UserID           int     `json:"userId"`
	NickName         string  `json:"nickName"`
	UserPhoto        string  `json:"userPhoto"`
	LastLoginTime    int64   `json:"lastLoginTime"`
	IsOpenVip        bool    `json:"isOpenVip"`
	VipLevel         int     `json:"vipLevel"`
	RechargeLevel    int     `json:"rechargeLevel"`
	WalletBalance    float64 `json:"walletBalance"`
	SafeBoxAmount    float64 `json:"safeBoxAmount"`
	BoolAttr         int     `json:"boolAttr"`
	HasNoReadMessage bool    `json:"hasNoReadMessage"`
	RegisterType     int     `json:"registerType"`
	VerifyMethods    struct {
		Email  string `json:"email"`
		Phone  string `json:"phone"`
		Google string `json:"google"`
	} `json:"verifyMethods"`
	BindGoogleVerifyMethod         int     `json:"bindGoogleVerifyMethod"`
	LastLoginSysLanguage           string  `json:"lastLoginSysLanguage"`
	InviteCode                     string  `json:"inviteCode"`
	YesterdayRebateAmount          float64 `json:"yesterdayRebateAmount"`
	UserUnGrandMsgCount            int     `json:"userUnGrandMsgCount"`
	UserUnreadInmailCount          int     `json:"userUnreadInmailCount"`
	UserUnreceiveInmailRewardCount int     `json:"userUnreceiveInmailRewardCount"`
	CanSetPassword                 bool    `json:"canSetPassword"`
}

type UserInfoResponse struct {
	Data UserInfoDataResponse `json:"data"`
}

// 获取前台用户的信息
func GetUserInfo(ctx *context.Context) (*model.Response, *UserInfoDataResponse, error) {
	api := "/api/User/GetUserInfo"
	payloadStruct, payloadList := utils.BaseStructHandler()
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/User/GetUserInfo请求失败", err)), &UserInfoDataResponse{}, err
	} else {
		var respdata UserInfoResponse
		if err := json.Unmarshal(respBoy, &respdata); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/User/GetUserInfo[respdata]解析失败", err)), &UserInfoDataResponse{}, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/User/GetUserInfo解析失败", err)), &UserInfoDataResponse{}, err
		} else {
			return resp, &respdata.Data, nil
		}
	}
}
