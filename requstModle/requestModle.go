package requstmodle

import (
	login "autoTest/API/adminApi/login"
	desklogin "autoTest/API/deskApi/loginApi"
	"autoTest/store/config"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type payloadMapfunc func(structType interface{}, slice []interface{}) (map[string]interface{}, error)

/*
主要是抽象出各种请求模板
前台的带TenantId  Referer  Origin  Domainurl Authorization
需要传入api  payload的结构体需要在外面进行实例化后
payloadData 数据切片
payloadFunc 处理payload和payloadData的函数
*
*/
func DeskTenAuthorRequest[P any](ctx *context.Context, api string, payload *P, payloadData []any, payloadFunc payloadMapfunc) ([]byte, *http.Response, error) {
	base_url := config.SIT_WEB_API
	// 请求头的设定
	header_struct := &model.DeskHeaderAuthorizationStruct{}
	plant_h5 := config.PLANT_H5
	token := (*ctx).Value(desklogin.DeskAuthTokenKey)
	//fmt.Println("前台登陆后的token", token)
	header_list := []interface{}{config.TENANTID, plant_h5, plant_h5, plant_h5, token}
	if headerMap, err := request.AssignSliceToStructMap(header_struct, header_list); err != nil {
		return nil, nil, errors.New("failed to convert headerMap struct to map")

	} else {
		// 设置请求体的
		payloadMap, err := payloadFunc(payload, payloadData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert payloadMap struct to map:%s", err)
		}
		//请求体和请求头的map
		// 将请求体map进行平铺
		FlattendMap := request.FlattenMap(payloadMap)
		respBody, req, err := request.PostRequestCofig(FlattendMap, base_url, api, headerMap)
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
}

/*
主要是抽象出各种请求模板
前台的带TenantId  Referer  Origin  Domainurl
需要传入api  payload的结构体需要在外面进行实例化后
payloadData 数据切片
payloadFunc 处理payload和payloadData的函数
注册的封装，前台手机号+验证码登录的封装
*
*/
func DeskTrodRegRequest[P any](ctx *context.Context, api string, payload *P, payloadData []any, payloadFunc payloadMapfunc) ([]byte, *http.Response, error) {
	base_url := config.SIT_WEB_API
	// 请求头的设定
	header_struct := &model.DeskHeaderTenantIdStruct{}
	plant_h5 := config.PLANT_H5
	header_list := []interface{}{config.TENANTID, plant_h5, plant_h5, plant_h5}
	if headerMap, err := request.InitStructToMap(header_struct, header_list); err != nil {
		return nil, nil, errors.New("failed to convert headerMap struct to map")

	} else {
		// 设置请求体的
		payloadMap, err := payloadFunc(payload, payloadData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert payloadMap struct to map:%s", err)
		}
		//请求体和请求头的map
		// 将请求体map进行平铺
		FlattendMap := request.FlattenMap(payloadMap)
		respBody, req, err := request.PostRequestCofig(FlattendMap, base_url, api, headerMap)
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
}

/*
商户后台的请求封装
Referer  Origin  Domainurl Authorization
*
*/
func AdminRodAutRequest[P any](ctx *context.Context, api string, payload *P, payloadData []any, payloadFunc payloadMapfunc) ([]byte, *http.Response, error) {
	base_url := config.ADMIN_SYSTEM_URL
	// 请求头的设定
	header_struct := &model.AdminHeaderStruct{}
	token := (*ctx).Value(login.AuthTokenKey)
	header_list := []interface{}{base_url, base_url, base_url, token}
	if headerMap, err := request.AssignSliceToStructMap(header_struct, header_list); err != nil {
		return nil, nil, errors.New("failed to convert headerMap struct to map")

	} else {
		// 设置请求体的
		payloadMap, err := payloadFunc(payload, payloadData)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to convert payloadMap struct to map:%s", err)
		}
		//请求体和请求头的map
		// 将请求体map进行平铺
		FlattendMap := request.FlattenMap(payloadMap)
		respBody, req, err := request.PostRequestCofig(FlattendMap, base_url, api, headerMap)
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
}
