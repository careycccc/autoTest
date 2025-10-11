package sitemessage

import (
	"autoTest/API/adminApi/login"
	uploadfile "autoTest/API/uploadFile"
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

// 站内信
// Translation represents the translations array in the JSON
type Translation struct {
	Language  string `json:"language"`
	Content   string `json:"content"`
	Title     string `json:"title"`
	Thumbnail string `json:"thumbnail"`
}

// FreeReward represents the freeReward object in rewardConfig
type FreeReward struct {
	RewardAmount         int    `json:"rewardAmount"`
	AmountCodingMultiple int    `json:"amountCodingMultiple"`
	CouponIds            string `json:"couponIds"`
}

// RechargeReward represents the rechargeReward object in rewardConfig
type RechargeReward struct {
	RechargeAmount       int    `json:"rechargeAmount"`
	RechargeCount        int    `json:"rechargeCount"`
	RewardAmount         int    `json:"rewardAmount"`
	AmountCodingMultiple int    `json:"amountCodingMultiple"`
	CouponIds            string `json:"couponIds"`
}

// RewardConfig represents the rewardConfig object in the JSON
type RewardConfig struct {
	FreeReward     FreeReward     `json:"freeReward"`
	RewardTypes    []int          `json:"rewardTypes"`
	RechargeReward RechargeReward `json:"rechargeReward"`
	ExpireType     int            `json:"expireType"`
}

// Message represents the entire JSON structure
type Message struct {
	BackstageDisplayName string        `json:"backstageDisplayName"`
	ValidType            int           `json:"validType"`
	Title                string        `json:"title"`
	JumpType             int           `json:"jumpType"`
	JumpPage             int           `json:"jumpPage"`
	JumpButtonText       string        `json:"jumpButtonText"`
	TargetType           int           `json:"targetType"`
	Translations         []Translation `json:"translations"`
	SendType             int           `json:"sendType"`
	IsHasReward          bool          `json:"isHasReward"`
	RewardConfig         RewardConfig  `json:"rewardConfig"`
	Random               int64         `json:"random"`
	Language             string        `json:"language"`
	Signature            string        `json:"signature"`
	Timestamp            int64         `json:"timestamp"`
}

/*
backstageDisplayName string,   站内信的名字
validType int,     默认值1
jumpType int,  跳转类型
jumpPage int,  跳转页面
jumpButtonText string  跳转的按钮文字
targetType int, 跳转目标 ，接收对象
content string,  站内信的内容
sendType int,  发送类型
thumbnail 图片上传的地址
*
*/
func CreateMessage(
	backstageDisplayName string,
	validType int,
	jumpType int,
	jumpPage int,
	jumpButtonText string,
	targetType int,
	content string,
	sendType int,
	random int64,
	timestamp int64,
	title string,
	thumbnail string,
) map[string]interface{} {
	// Create the Message struct with provided parameters and default values from JSON
	message := Message{
		BackstageDisplayName: backstageDisplayName,
		ValidType:            validType,
		JumpType:             jumpType,
		JumpPage:             jumpPage,
		JumpButtonText:       jumpButtonText,
		TargetType:           targetType,
		SendType:             sendType,
		Random:               random,
		Timestamp:            timestamp,
		Language:             "en", // Default from JSON
		Signature:            "",   // Default from JSON
		IsHasReward:          true, // Default from JSON
		Translations: []Translation{
			{
				Language:  "hi",
				Content:   content, // Use provided content for hi
				Title:     title,
				Thumbnail: thumbnail,
			},
			{
				Language:  "en",
				Content:   content, // Use provided content for en
				Title:     title,
				Thumbnail: thumbnail,
			},
			{
				Language:  "zh",
				Content:   content, // Use provided content for zh
				Title:     title,
				Thumbnail: thumbnail,
			},
		},
		RewardConfig: RewardConfig{
			FreeReward: FreeReward{
				RewardAmount:         10,
				AmountCodingMultiple: 1,
				CouponIds:            "400007",
			},
			RewardTypes: []int{1, 2},
			RechargeReward: RechargeReward{
				RechargeAmount:       1000,
				RechargeCount:        1,
				RewardAmount:         100,
				AmountCodingMultiple: 11,
				CouponIds:            "400007",
			},
			ExpireType: 1,
		},
	}

	// Convert Message struct to map[string]interface{}
	result := map[string]interface{}{
		"backstageDisplayName": message.BackstageDisplayName,
		"validType":            message.ValidType,
		"jumpType":             message.JumpType,
		"jumpPage":             message.JumpPage,
		"jumpButtonText":       message.JumpButtonText,
		"targetType":           message.TargetType,
		"translations":         message.Translations,
		"sendType":             message.SendType,
		"isHasReward":          message.IsHasReward,
		"rewardConfig":         message.RewardConfig,
		"random":               message.Random,
		"language":             message.Language,
		"signature":            message.Signature,
		"timestamp":            message.Timestamp,
		"title":                message.Title,
	}

	return result
}

type QuerySiteMessageStruct struct {
	PageNo  any `json:"pageNo"`
	PgeSize any `json:"pageSize"`
	OrderBy any `json:"orderBy"`
	model.BaseStruct
}

type Root struct {
	Data Data `json:"data"`
}

type Data struct {
	List []Notification `json:"list"`
}

type Notification struct {
	ID int64 `json:"id"`
}

// 返回所有的站内信的id
func QuerySiteMessage(ctx *context.Context) (*model.Response, []int64, error) {
	api := "/api/Inmail/GetPageList"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &QuerySiteMessageStruct{}
	payloadList := []interface{}{1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Inmail/GetPageList请求失败", err)), nil, err
	} else {
		var root Root
		if err := json.Unmarshal(respBoy, &root); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Inmail/GetPageList 【Root】解析失败", err)), nil, err
		}
		// 提取所有 id
		var ids []int64
		for _, notification := range root.Data.List {
			ids = append(ids, notification.ID)
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Inmail/GetPageList解析失败", err)), nil, err
		} else {
			return resp, ids, nil
		}
	}
}

// 点击启用
type ClickSiteMessageStruct struct {
	State any `json:"state"`
	Id    any `json:"id"`
	model.BaseStruct
}

// 站内信的启用
func ClickSiteMessage(ctx *context.Context, id int64) (*model.Response, error) {
	api := "/api/Inmail/UpdateState"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &ClickSiteMessageStruct{}
	payloadList := []interface{}{1, id, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Inmail/UpdateState请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Inmail/UpdateState解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

// 获取请求头的map
func GetHeaderMap(ctx *context.Context) (map[string]interface{}, string) {
	// 请求头
	headerStruct := &model.AdminHeaderStruct{}
	header_url := config.ADMIN_SYSTEM_URL
	token := (*ctx).Value(login.AuthTokenKey)
	desSlice := []interface{}{header_url, header_url, header_url, token}
	headMap, err := request.AssignSliceToStructMap(headerStruct, desSlice)
	if err != nil {
		logger.LogError("headerMap获取失败", err)
		return nil, ""
	}
	return headMap, header_url
}

// 需要提供跳转类型，和跳转文字,上传地址
func SendSiteMessage(ctx *context.Context, jumpNumber int, jumpText, thumbnail string) (*model.Response, error) {
	timestamp, random, _ := request.GetTimeRandom()
	znxTitle := "自动化生成的站内信" + strconv.FormatInt(timestamp, 10)
	result := CreateMessage(znxTitle, 1, 1, jumpNumber, jumpText, 1, "这是内容", 1, random, timestamp, znxTitle, thumbnail)
	headMap, header_url := GetHeaderMap(ctx)
	api := "/api/Inmail/Add"
	respBoy, _, err := request.PostRequestCofig(result, header_url, api, headMap)
	if err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("站内信的请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("站内信的响应解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

// 运行站内信
func RunSendSiteMessage() error {
	idlist := utils.RandmoUserId(config.SiteMessageNumber)
	results, err := MoreSendSiteMessage(idlist, config.SiteMessageConcurrent)
	if err != nil {
		return err
	}
	// 打印结果
	for _, result := range results {
		fmt.Println(result)
	}
	return nil
}

/*
多并发的发送站内信
包含并发逻辑，保证 a -> b -> c 串行
inputs需要提供第一个函数的入参
concurrency 并发量
*
*/
func MoreSendSiteMessage(inputs []string, concurrency int) ([]string, error) {
	ctx, err := login.RunAdminSitLogin()
	if err != nil {
		logger.LogError("站内信的【后台登录失败】", err)
		return nil, err
	}
	idChan := make(chan struct{}, 1)
	sem := make(chan struct{}, concurrency) // 信号量通道，控制并发量
	var wg sync.WaitGroup                   // 等待所有任务完成

	// 处理每个输入
	for i, input := range inputs {
		wg.Add(1)
		go func(idx int, input string, dChan chan struct{}) {

			defer wg.Done()
			sem <- struct{}{}        // 占用一个并发槽
			defer func() { <-sem }() // 释放并发槽
			// 文件上传
			fileName := fmt.Sprintf("./assert/workerOder/%d.png", i+1)
			ch := make(chan struct{}, 1)
			_, thumbnail := uploadfile.RunWorkerOderActiveZx(ctx, fileName, ch)
			<-ch
			// 串行执行 a -> b -> c
			// 按顺序执行 a -> b -> c
			_, err := SendSiteMessage(ctx, i, "跳转"+strconv.Itoa(i), thumbnail)
			if err != nil {
				logger.LogError("站内信的创建失败", err)
				return
			}
			logger.Logger.Info("站内信创建成功")
			idChan <- struct{}{}

		}(i, input, idChan)
	}
	// 上面的步骤结束了，进行收集所有的站内信的id

	<-idChan
	_, znxIdList, err := QuerySiteMessage(ctx)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(znxIdList); i++ {
		if _, err := ClickSiteMessage(ctx, znxIdList[i]); err != nil {
			logger.Logger.Warn("站内信的启用失败", err)
			continue
		}
	}

	// 关闭结果通道
	go func() {
		wg.Wait()
		//close(results)
	}()

	// 收集结果并按索引排序
	finalResults := make([]string, len(inputs))
	// for result := range results {
	// 	finalResults[result.index] = result.result
	// }

	return finalResults, err
}
