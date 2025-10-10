package activeinformation

import (
	"autoTest/API/adminApi/login"
	messagemanagement "autoTest/API/adminApi/operationsManagement/messageManagement"
	uploadfile "autoTest/API/uploadFile"
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"fmt"
)

// 活动资讯

type ActiveInformationStruct struct {
	ImgUrl          any `json:"imgUrl"`          // 图片展示的地址
	InformationType any `json:"informationType"` // 0 链接  1.富文本，2，界面  3.工单  新建的类型选择
	Title           any `json:"title"`           // 活动标题
	Sort            any `json:"sort"`            // 排序
	DisplayTarget   any `json:"displayTarget"`   // 展示目标 1表示全平台
	PageId          any `json:"pageId"`          // 工单类型  1.外部链接   2. 一对一客服
	SysLanguage     any `json:"sysLanguage"`     // 系统语言
	Content         any `json:"content"`         // 空的字符串
	model.BaseStruct
}

/*
添加活动资讯的工单跳转
ImgUrl  图片展示的地址
InformationType  0 链接  1.富文本，2，界面  3.工单  新建的类型选择
Title 活动标题
Sort 排序
DisplayTarget 展示目标 1表示全平台
**/

func AddActiveInformationFunc(ctx *context.Context, imgUrl, title string, InformationType, sort, displayTarget, pageId int) (*model.Response, error) {
	api := "/api/ActivityInformation/Add"
	payloadStruct := &ActiveInformationStruct{}
	timestamp, random, language := request.GetTimeRandom()
	paloadList := []interface{}{imgUrl, InformationType, title, sort, displayTarget, pageId, language, "", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest[ActiveInformationStruct](ctx, api, payloadStruct, paloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/ActivityInformation/Add 请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/ActivityInformation/Add 响应解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

// 执行全量的活动资讯工单
func RunAddActiveInformation() error {
	// 进行后台登录
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		return err
	} else {
		// 图片上传
		worderList := messagemanagement.WordOrderList
		for i := range worderList {
			fileName := fmt.Sprintf("./assert/workerOder/%d.png", i+1)
			ch := make(chan struct{}, 1)
			resp, str := uploadfile.RunWorkerOderActiveZx(ctx, fileName, ch)
			<-ch
			if res, err := AddActiveInformationFunc(ctx, str, worderList[i], 3, 10+i, 1, i+1); err != nil {
				logger.LogError("AddActiveInformationFunc调用失败", err)
				return err
			} else {
				fmt.Println(resp)
				fmt.Println(res)
			}
		}
		return nil
	}
	// 创建 活动资讯的工单
}
