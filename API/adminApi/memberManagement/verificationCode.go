package membermanagement

import (
	"autoTest/API/adminApi/login"
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// 发送和收取验证码

// 发送验证码
type SendVerifiyCodeStruct struct {
	VerifyCodeType any `json:"verifyCodeType"`
	PhoneOrEmail   any `json:"phoneOrEmail"`
	CodeType       any `json:"codeType"`
	model.BaseStruct
}

/*
需要传入上下文，手机号码
返回 响应结果，以及错误信息
发送验证码
codeType 验证码类型 18是登录验证 1是注册验证
*
*/
func SendVerificationCode(userName string, codeType int8, ch chan struct{}) (*model.BetResponse, error) {
	defer func() {
		time.Sleep(time.Second)
		ch <- struct{}{}
	}()
	api := "/api/Home/SendVerifiyCode"
	base_url := config.GoodsDeposit_URL
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &SendVerifiyCodeStruct{}
	payloadList := []interface{}{1, userName, codeType, random, language, "", timestamp}
	headerStruct := &model.DeskHeaderTenantIdStruct2{}
	plant_h5 := config.REGISTER_URL
	// heaerList := []interface{}{config.TENANTID, plant_h5, plant_h5, plant_h5}
	heaerList := []interface{}{plant_h5, plant_h5}
	if respBoy, _, err := request.PostGenericsFuncFlatten(base_url, api, payloadStruct, payloadList, headerStruct, heaerList, request.StructToMap, request.InitStructToMap); err != nil {
		errs := fmt.Errorf("/api/Home/SendVerifiyCode发送请求失败%s", err)
		return model.HandlerErrorRes2(errs), err
	} else {
		// 将返回值进行解析
		if betRes, err := model.ParseResponse2(respBoy); err != nil {
			errs := fmt.Errorf("/api/Home/SendVerifiyCode响应解析失败%s", err)
			return model.HandlerErrorRes2(errs), err
		} else {
			return betRes, nil
		}
	}
}

// 查询验证码
type QueryTifyStruct struct {
	MobileOrEmail string `json:"mobileOrEmail"`
	model.QueryPayloadStruct
}

// 验证码类型
type TifyResponse struct {
	Data struct {
		List       []ListItem `json:"list"`
		PageNo     int        `json:"pageNo"`
		TotalPage  int        `json:"totalPage"`
		TotalCount int        `json:"totalCount"`
	} `json:"data"`
	MsgParameters interface{} `json:"msgParameters"`
	Code          int         `json:"code"`
	Msg           string      `json:"msg"`
	MsgCode       int         `json:"msgCode"`
}

// 定义结构体来映射JSON中的list项
type ListItem struct {
	ID             int64  `json:"id"`
	Category       int    `json:"category"`
	CategoryText   string `json:"categoryText"`
	CodeType       int    `json:"codeType"`
	CodeTypeText   string `json:"codeTypeText"`
	UserName       string `json:"userName"`
	Number         string `json:"number"`
	State          int    `json:"state"`
	StateText      string `json:"stateText"`
	Remark         string `json:"remark"`
	SendIP         string `json:"sendIP"`
	UserID         int64  `json:"userId"`
	ExpirationTime int64  `json:"expirationTime"`
	Creator        string `json:"creator"`
	CreateTime     int64  `json:"createTime"`
	LastUpdateMan  string `json:"lastUpdateMan"`
	LastUpdateTime int64  `json:"lastUpdateTime"`
}

/*
// 获取验证码
// 需要传入上下文，需要获取的验证码的手机号码
返回响应和验证码
*
*/
func GetVerificationCode(ctx *context.Context, userName string) (*model.Response, string, error) {
	api := "/api/Users/GetVerifyCodePageList"
	payloadStruct := &QueryTifyStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{userName, 1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Users/GetVerifyCodePageList的请求失败", err)), "", err
	} else {
		var tifyResponse TifyResponse
		if err := json.Unmarshal(respBoy, &tifyResponse); err != nil {
			return &model.Response{
				Code:    -1,
				Msg:     "/api/Users/GetVerifyCodePageList" + err.Error(),
				MsgCode: -1,
			}, "", err
		} else {
			return &model.Response{
				Code:    0,
				Msg:     "Succeed",
				MsgCode: 0,
			}, tifyResponse.Data.List[0].Number, nil
		}

	}

}

// 发送要验证码到接收验证码
// codeType 验证码类型 18是登录验证 1是注册验证
func SendToGetVerCode(ctx *context.Context, codeType int8, userName string) (*model.Response, string, error) {
	ch := make(chan struct{}, 1)
	// 1.发送验证码
	SendVerificationCode(userName, codeType, ch)
	// 2.获取验证码
	// 2.1需要后端登录
	if resp, ctxToken, err := login.AdminSitLogin(ctx); err != nil {
		return resp, "", err
	} else {
		<-ch
		if res, verficationCode, err := GetVerificationCode(ctxToken, userName); err != nil {
			return res, "", err
		} else {
			return res, verficationCode, err
		}
	}
}
