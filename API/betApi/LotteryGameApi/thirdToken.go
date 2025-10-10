package lotterygameapi

import (
	login "autoTest/API/deskApi/loginApi"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"context"
	"encoding/json"
	"strings"
)

// 前台获取下注的token
type BetTokenStruct struct {
	Referer       interface{}
	Origin        interface{}
	Authorization interface{}
}

// 主要是针对那三个get请求获取token，为后面的投注的token
//  gameCode=WinGo_5M&language=en&random=131601285634&signature=C9A95C24297A0345B9DF3FC970BBB766&timestamp=1756999407

type GetGameInfoStruct struct {
	GameCode  string `json:"gameCode"`
	Language  string `json:"language"`
	Random    int64  `json:"random"`
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
}

// 请求GetGameInfo并且返回token
func ThridTokenFunc(token, gameCode string) (string, error) {
	// 初始化结构体并且赋值
	GetGameInfo := &GetGameInfoStruct{}
	timestamp, random, language := request.GetTimeRandom()
	values := []interface{}{gameCode, language, random, "", timestamp}
	paramsMap, _ := request.InitStructToMap(GetGameInfo, values)
	// 获取签名
	verifyPwd := ""
	signatureStr := utils.GetSignature(paramsMap, &verifyPwd)
	paramsMap["signature"] = signatureStr

	api := "/api/Lottery/GetGameInfo"
	baseUrl := config.LOTTERY_H5
	// 获取请求头
	deskA := &BetTokenStruct{}
	url_h5 := config.WMG_H5
	ThirdGametoken, err := ThirdGameFunc(token, gameCode)
	if err != nil {
		logger.LogError("/api/ThirdGame/GetGameUrl报错", err)
		return "", err
	}
	desSlice := []interface{}{url_h5, url_h5, ThirdGametoken}
	headMap, _ := request.AssignSliceToStructMap(deskA, desSlice)
	_, resp, err := request.GetRequest(baseUrl, api, headMap, paramsMap)
	if err != nil {
		logger.LogError("/api/Lottery/GetGameInfo请求报错", err)
		return "", err
	}
	authorization := resp.Header.Get("Authorization")
	// 去掉前缀 "Bearer "
	cleanToken := strings.TrimPrefix(authorization, "Bearer ")
	// fmt.Println("响应头的token", cleanToken)
	return cleanToken, nil
}

type GetBalanceInfoStruct struct {
	Language  string `json:"language"`
	Random    int64  `json:"random"`
	Signature string `json:"signature"`
	Timestamp int64  `json:"timestamp"`
}

// 获取banlance的值
type BalanceStruct struct {
	Data struct {
		Balance float64 `json:"balance"`
	} `json:"data"`
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	MsgCode     int    `json:"msgCode"`
	ServiceTime int64  `json:"serviceTime"`
}

// 请求GetBalance
func GetBalanceInfoFunc(ctx *context.Context, gameCode string) (string, float64, error) {
	// 初始化结构体并且赋值
	GetGameInfo := &model.BaseStruct{}
	timestamp, random, language := request.GetTimeRandom()
	values := []interface{}{random, language, "", timestamp}
	paramsMap, _ := request.InitStructToMap(GetGameInfo, values)

	// 获取签名
	verifyPwd := ""
	signatureStr := utils.GetSignature(paramsMap, &verifyPwd)
	paramsMap["signature"] = signatureStr
	api := "/api/Lottery/GetBalance"
	baseUrl := config.LOTTERY_H5
	// 获取请求头
	deskA := &BetTokenStruct{}
	url_h5 := config.WMG_H5
	token := (*ctx).Value(login.DeskAuthTokenKey)

	ThridTokenToken, err := ThridTokenFunc(token.(string), gameCode)
	if err != nil {
		logger.LogError("/api/ThirdGame/GetGameUrl报错", err)
		return "", 0.0, err
	}

	desSlice := []interface{}{url_h5, url_h5, ThridTokenToken}
	headMap, _ := request.AssignSliceToStructMap(deskA, desSlice)
	respBoy, resp, err := request.GetRequest(baseUrl, api, headMap, paramsMap)
	if err != nil {
		logger.LogError("/api/Lottery/GetBalance请求报错消息", err)
		return "", 0.0, err
	}
	// fmt.Println("GetBalanceInfoFunc的响应结果", string(respBoy))
	var balance BalanceStruct
	err = json.Unmarshal([]byte(string(respBoy)), &balance)
	if err != nil {
		logger.LogError("GetBalanceInfoFunc的反序列失败", err)
		return "", 0.0, err
	}

	authorization := resp.Header.Get("Authorization")
	// 去掉前缀 "Bearer "
	cleanToken := strings.TrimPrefix(authorization, "Bearer ")

	return cleanToken, balance.Data.Balance, nil
}
