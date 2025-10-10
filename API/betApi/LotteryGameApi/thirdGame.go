package lotterygameapi

import (
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"encoding/json"
	"regexp"
)

type ThirdGameStruct struct {
	GameCode     any `json:"gameCode"`
	VendorCode   any `json:"vendorCode"`
	GameId       any `json:"gameId"`
	ReturnUrl    any `json:"returnUrl"`
	DeviceType   any `json:"deviceType"`
	DeviceTypeId any `json:"deviceTypeId"`
	model.BaseStruct
}

// 定义与 JSON 结构对应的结构体
type Response struct {
	Data          interface{} `json:"data"`          // 使用 interface{} 处理 null
	MsgParameters interface{} `json:"msgParameters"` // 使用 interface{} 处理 null
	Code          int         `json:"code"`
	Msg           string      `json:"msg"`
	MsgCode       int         `json:"msgCode"`
}

/*
传入token,游戏code
返回token
*
*/
func ThirdGameFunc(token, gameCode string) (string, error) {
	api := "/api/ThirdGame/GetGameUrl"
	base_url := config.SIT_WEB_API
	// 获取payload的实例
	payloadStruct := &ThirdGameStruct{}
	returnRurl := config.PLANT_H5 + "/game?categoryCode=C202505280608510046"
	deviceTypeId := utils.GenerateCryptoRandomString(32)
	timestamp, random, language := request.GetTimeRandom()
	payloadData := []interface{}{gameCode, "ARLottery", 10003, returnRurl, "PC", deviceTypeId, random, language, "", timestamp}
	// 请求头的实例
	headerStruct := &model.AdminHeaderStruct{}
	h5_y1 := config.PLANT_H5
	headerData := []interface{}{h5_y1, h5_y1, h5_y1, token}
	if resp, _, err := request.PostGenericsFuncFlatten[ThirdGameStruct, model.AdminHeaderStruct](base_url, api, payloadStruct, payloadData, headerStruct, headerData, request.StructToMap, request.AssignSliceToStructMap); err != nil {
		logger.LogError("/api/ThirdGame/GetGameUrl报错消息", err)
		return "", err
	} else {
		var response Response
		err = json.Unmarshal([]byte(string(resp)), &response)
		if err != nil {
			logger.LogError("/api/ThirdGame/GetGameUrl响应结果反序列化失败", err)
			return "", err
		}
		result := response.Data.(map[string]interface{})["url"]
		// fmt.Println(result)
		// 寻找token
		// 查找第一个匹配
		res := result.(string)
		// 用正则匹配 Token 的值
		re := regexp.MustCompile(`Token=([^&]+)`)
		matches := re.FindStringSubmatch(res)

		if len(matches) > 1 {
			// fmt.Println("Token:", matches[1])
			return matches[1], nil
		} else {
			logger.LogError("/api/ThirdGame/GetGameUrl,token not font", err)
			return "", err
		}
	}
}
