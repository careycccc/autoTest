package getgiftinfo

import (
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

// 响应中data的struct
type GetGiftInfoDataResponse struct {
	IsOpenInvitedWheel           any `json:"isOpenInvitedWheel"`           // 是否开启邀请转盘
	InvitedWheelTotalPrizeAmount any `json:"invitedWheelTotalPrizeAmount"` // 邀请转盘的总金额
	UserInvitedWheelAmount       any `json:"userInvitedWheelAmount"`       // 用户在邀请转盘的当前旋转金额
	IsFirstInvitedWheel          any `json:"isFirstInvitedWheel"`          // 用户是否开启了邀请转盘
	IsOpenRechargeWheel          any `json:"isOpenRechargeWheel"`          // 是否开启了充值转盘 是否显示充值转盘
	RechargeWheelRemainSpinCount any `json:"rechargeWheelRemainSpinCount"` // 充值轮盘剩余旋转次数
	HasUnreceivedGiftPack        any `json:"hasUnreceivedGiftPack"`        // 是否有未收到的活动礼包
	UnreceivedGiftPackCount      any `json:"unreceivedGiftPackCount"`      // 有未收到的活动礼包数量
	UserInmailCount              any `json:"userInmailCount"`              // 站内信未领取礼包的数量
	GiftPackageCount             any `json:"giftPackageCount"`             // 活动礼包数量
}

// 定义结构体来映射 JSON 数据
type GetGiftInfoResponse struct {
	Data GetGiftInfoDataResponse
}

// 获取前台用户的礼包状态信息
func GetGiftInfoApi(ctx *context.Context) (*model.Response, GetGiftInfoDataResponse, error) {
	api := "/api/Home/GetGiftInfo"
	payloadStruct, payloadList := utils.BaseStructHandler()
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.InitStructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("", err)), GetGiftInfoDataResponse{}, err
	} else {
		var dataResponse GetGiftInfoResponse
		if err := json.Unmarshal(respBoy, &dataResponse); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("", err)), GetGiftInfoDataResponse{}, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("", err)), GetGiftInfoDataResponse{}, err
		} else {
			return resp, dataResponse.Data, nil
		}
	}
}
