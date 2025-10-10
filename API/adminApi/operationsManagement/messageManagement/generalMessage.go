package messagemanagement

import (
	"autoTest/API/adminApi/login"
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"fmt"
	"time"
)

type AddCarouselStruct struct {
	Type            any `json:"type"`            // 消息类型  4表示轮播图  5定制化轮播图
	Sort            any `json:"sort"`            // 排序数字越大越靠前
	MessageJumpType any `json:"messageJumpType"` // 消息跳转目标类型 5 表示工单
	ImageUrl        any `json:"imageUrl"`        // 上传的图片地址
	SysLanguage     any `json:"sysLanguage"`     // 系统语言
	CustomPopupId   any `josn:"customPopupId"`   // 工单的目标  1 外部链接
	TargetType      any `json:"targetType"`      // 消息跳转目标  1 全平台会员
	model.BaseStruct
}

// 通用消息
// 新增轮播图
/*
messageType 消息类型  4表示轮播图  5定制化轮播图
MessageJumpType  消息跳转目标类型 5 表示工单
TargetType 消息跳转目标  1 全平台会员
customPopupId 工单类型
**/
func AddCarousel(ctx *context.Context, sort int, messageType, messageJumpType, customPopupId, targetType int8) (*model.Response, error) {
	api := "/api/Message/Add"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &AddCarouselStruct{}
	payloadList := []interface{}{messageType, sort, messageJumpType, "3003/other/082900768-1997-3.webp", "en", customPopupId, targetType, random, language, "", timestamp}
	if respBody, _, err := requstmodle.AdminRodAutRequest[AddCarouselStruct](ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Message/Add请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBody); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Message/Add解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

type AddCustomizedCarouselStruct struct {
	Type            any `json:"type"`            // 消息类型  4表示轮播图  5定制化轮播图
	Sort            any `json:"sort"`            // 排序数字越大越靠前
	MessageJumpType any `json:"messageJumpType"` // 消息跳转目标类型 5 表示工单
	ImageUrl        any `json:"imageUrl"`        // 上传的图片地址
	SysLanguage     any `json:"sysLanguage"`     // 系统语言
	CustomPopupId   any `josn:"customPopupId"`   // 工单的目标  1 外部链接
	ButtonTxt       any `json:"buttonTxt"`       // 按钮文字
	TargetType      any `json:"targetType"`      // 消息跳转目标  1 全平台会员
	model.BaseStruct
}

/*
// 定制化弹窗轮播图 跳工单
messageType 消息类型  4表示轮播图  5定制化轮播图
MessageJumpType  消息跳转目标类型 5 表示工单
TargetType 消息跳转目标  1 全平台会员
customPopupId 工单类型
*
*/
func AddCustomizedCarousel(ctx *context.Context, sort int, messageType, messageJumpType, customPopupId, targetType int8, buttonTxt string) (*model.Response, error) {
	api := "/api/Message/Add"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &AddCustomizedCarouselStruct{}
	payloadList := []interface{}{messageType, sort, messageJumpType, "3003/other/082900768-1997-3.webp", "en", customPopupId, buttonTxt, targetType, random, language, "", timestamp}
	if respBody, _, err := requstmodle.AdminRodAutRequest[AddCustomizedCarouselStruct](ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Message/Add请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBody); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Message/Add解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

// 全量轮播图的的工单的
func AllWorkOrderCreate() {
	len := len(WordOrderList)
	ctx := context.Background()
	_, ctxT, err := login.AdminSitLogin(&ctx)
	if err != nil {
		fmt.Println("Login error:", err)
		return
	}
	for i := 1; i <= len; i++ {
		AddCarousel(ctxT, i+10, 4, 5, int8(i), 1)
	}
}

// 全量定制化轮播图的的工单的
func AllCustomizedCarousel() {
	len := len(WordOrderList)

	ctx := context.Background()
	_, ctxT, err := login.AdminSitLogin(&ctx)
	if err != nil {
		fmt.Println("Login error:", err)
		return
	}

	for i := 1; i <= len; i++ {
		time.Sleep(500 * time.Millisecond)
		AddCustomizedCarousel(ctxT, i+10, 5, 5, int8(i), 1, WordOrderList[i-1])
	}
}
