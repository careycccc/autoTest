package request

import (
	"autoTest/store/config"
	"autoTest/store/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

/*
注册的泛型
P,H是需要传入的两个结构体
p 表示的是 payload的结构体
H 表示的是 header的结构体
api 表示接口地址
payloadList 表示需要赋值的 payload
headerList 表示需要赋值的 header
payloadFunc 表示需要处理的pay的func
headerFunc 表示需要处理的header的func
token 需要传入的token值
*
*/
func PostGenericsFunc[P, H any](api string, payload *P, payloadList []interface{}, headerStruct *H,
	headerList []interface{}, payloadFunc func(structType interface{}, slice []interface{}) (map[string]interface{}, error), headerFunc func(structType interface{}, slice []interface{}) (map[string]interface{}, error)) ([]byte, *http.Response, error) {
	// 结构体转 Map
	payloadMap, err := payloadFunc(payload, payloadList)
	if err != nil {
		return nil, nil, errors.New("failed to convert payloadMap struct to map")
	}
	// 获取 token 和 base_url
	base_url := config.SIT_WEB_API
	// 确保 headerList 包含必要参数
	headerMap, err := headerFunc(headerStruct, headerList)
	if err != nil {
		return nil, nil, errors.New("failed to convert headerMap struct to map")
	}

	respBody, req, err := PostRequestCofig(payloadMap, base_url, api, headerMap)
	if err != nil {
		return nil, nil, fmt.Errorf("请求失败:%s", err)
	}
	return respBody, req, nil
}

type structSlice func(structType interface{}, slice []interface{}) (map[string]interface{}, error)

/*
PostGenericsFuncFlatten 的需要进行结构体平铺的版本
P,H是需要传入的两个结构体
p 表示的是 payload的结构体
H 表示的是 header的结构体
api 表示接口地址
payloadList 表示需要赋值的 payload
headerList 表示需要赋值的 header
payloadFunc 表示需要处理的pay的func
headerFunc 表示需要处理的header的func
token 需要传入的token值
*
*/
func PostGenericsFuncFlatten[P, H any](base_url, api string, payload *P, payloadList []interface{}, headerStruct *H,
	headerList []interface{}, payloadFunc structSlice, headerFunc structSlice) ([]byte, *http.Response, error) {
	// 结构体转 Map
	payloadMap, err := payloadFunc(payload, payloadList)
	if err != nil {
		return nil, nil, errors.New("failed to convert payloadMap struct to map")
	}
	// 获取 token 和 base_url
	// base_url := common.SIT_WEB_API
	// 确保 headerList 包含必要参数
	headerMap, err := headerFunc(headerStruct, headerList)
	if err != nil {
		return nil, nil, errors.New("failed to convert headerMap struct to map")
	}
	// 将嵌套map进行平铺
	FlattendMap := FlattenMap(payloadMap)

	respBody, req, err := PostRequestCofig(FlattendMap, base_url, api, headerMap)
	if err != nil {
		return nil, nil, fmt.Errorf("请求失败:%s", err)
	}
	var result model.Response
	err = json.Unmarshal([]byte(string(respBody)), &result)
	if err != nil {
		return nil, nil, fmt.Errorf("错误代码反序列化:%s", err)
	}
	return respBody, req, nil
}
