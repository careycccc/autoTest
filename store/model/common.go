package model

import (
	"autoTest/store/logger"
	json "encoding/json"
	"fmt"
	"time"

	easyjson "github.com/mailru/easyjson"
)

// 解析token
//
//easyjson:json
type ResponseToken struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

// 解析响应
//
//easyjson:json
type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	MsgCode int         `json:"msgCode"`
	Data    interface{} `json:"data,omitempty"`
}

// 投注类型的响应
type BetResponse struct {
	Code        int    `json:"code"`
	Msg         string `json:"msg"`
	MsgCode     int    `json:"msgCode"`
	ServiceTime int64  `json:"serviceTime"`
}

// 解析响应中的token
func GetJsonToken(jsonStr string) (string, error) {
	// 解析 JSON
	var response ResponseToken
	err := easyjson.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return "", fmt.Errorf("JSON 解析错误: %v", err)
	}

	// 提取 token
	token := response.Data.Token
	return token, nil
}

// 解析响应中的response
func ParseResponse(respBody []byte) (*Response, error) {
	var resp Response
	err := easyjson.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, fmt.Errorf("JSON 解析错误: %v", err)
	}
	return &resp, nil
}

// 解析响应中的投注类型的
func ParseResponse2(respBody []byte) (*BetResponse, error) {
	var resp BetResponse
	err := json.Unmarshal(respBody, &resp)
	if err != nil {
		return nil, fmt.Errorf("JSON 解析错误: %v", err)
	}
	return &resp, nil
}

// 传入一个error 返回一个Response
// 主要把错误信息写入到响应里面
func HandlerErrorRes(err error) *Response {
	return &Response{
		Code:    0,
		Msg:     err.Error(),
		MsgCode: 0,
		Data:    "",
	}
}

func HandlerErrorRes2(err error) *BetResponse {
	return &BetResponse{
		Code:        0,
		Msg:         err.Error(),
		MsgCode:     0,
		ServiceTime: 0,
	}
}

// 两种类型的转换
// 函数将 Response 转换为 BetResponse
func ResponseToBetResponse(resp *Response) *BetResponse {
	// 初始化 BetResponse
	betResp := BetResponse{
		Code:    resp.Code,    // 赋值相同字段
		Msg:     resp.Msg,     // 赋值相同字段
		MsgCode: resp.MsgCode, // 赋值相同字段
	}

	// 处理 Data 字段到 ServiceTime
	if resp.Data != nil {
		// 假设 Data 是 map[string]int64 格式，包含 serviceTime
		if dataMap, ok := resp.Data.(map[string]int64); ok {
			if serviceTime, exists := dataMap["serviceTime"]; exists {
				betResp.ServiceTime = serviceTime
			} else {
				// 如果 Data 中没有 serviceTime，使用默认值（当前时间戳）
				betResp.ServiceTime = time.Now().Unix()
			}
		} else if serviceTime, ok := resp.Data.(int64); ok {
			// 如果 Data 直接是 int64 类型
			betResp.ServiceTime = serviceTime
		} else {
			// 其他类型，使用默认值
			betResp.ServiceTime = time.Now().Unix()
		}
	} else {
		// 如果 Data 为空，使用默认值
		betResp.ServiceTime = time.Now().Unix()
	}

	return &betResp
}

// betResponseToResponse 将 BetResponse 转换为 Response
func BetResponseToResponse(betResp BetResponse) Response {
	return Response{
		Code:    betResp.Code,                                         // 赋值相同字段
		Msg:     betResp.Msg,                                          // 赋值相同字段
		MsgCode: betResp.MsgCode,                                      // 赋值相同字段
		Data:    map[string]int64{"serviceTime": betResp.ServiceTime}, // ServiceTime 放入 Data
	}
}

// 错误的情况下日志 + 错误类型的返回
func ErrorLoggerType(s string, err error) error {
	logger.LogError(fmt.Sprintf("%s", s), err)
	errs := fmt.Errorf("%s-%s", s, err)
	return errs
}
