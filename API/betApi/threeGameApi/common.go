package threegameapi

import (
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"context"
	"encoding/json"
	"net/url"

	"github.com/gofrs/uuid"
)

// 三方游戏的公共的结构体和函数
type ThirdGameStruct struct {
	GameCode     any `json:"gameCode"`
	VendorCode   any `json:"vendorCode"`
	GameID       any `json:"gameId"`
	ReturnURL    any `json:"returnUrl"`
	DeviceType   any `json:"deviceType"`
	DeviceTypeID any `json:"deviceTypeId"`
	model.BaseStruct
}

// token, operatorId 厂商id
type ThirdGameType struct {
	AuthToken  string
	OperatorId string
}

/*
*
获取游戏的url
gameCode 游戏code
vendorCode 厂商code
returnUrl 返回的url
gameId 游戏id
*
*/
func GetGameUrlCommon(ctx *context.Context, gameCode, vendorCode, returnUrl string, gameId int) (*model.Response, *ThirdGameType, error) {
	api := "/api/ThirdGame/GetGameUrl"
	payloadStruct := &ThirdGameStruct{}
	deviceTypeId := utils.GenerateCryptoRandomString(32)
	timestamp, random, language := request.GetTimeRandom()
	payloadData := []any{gameCode, vendorCode, gameId, returnUrl, "H5", deviceTypeId, random, language, "", timestamp}
	if respBody, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadData, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("获取三方游戏的/api/ThirdGame/GetGameUrl请求失败", err)), &ThirdGameType{}, err
	} else {
		// 解析出token值
		var resp struct {
			Data struct {
				URL string `json:"url"`
			} `json:"data"`
		}

		if err := json.Unmarshal([]byte(respBody), &resp); err != nil {
			logger.LogError("解析三方游戏的/api/ThirdGame/GetGameUrl响应失败", err)
			return model.HandlerErrorRes(model.ErrorLoggerType("解析三方游戏的/api/ThirdGame/GetGameUrl响应失败", err)), &ThirdGameType{}, err
		}

		u, err := url.Parse(resp.Data.URL)
		if err != nil {
			logger.LogError("解析三方游戏的url失败", err)
			return model.HandlerErrorRes(model.ErrorLoggerType("解析三方游戏的url失败", err)), &ThirdGameType{}, err
		}

		authToken := u.Query().Get("authToken")
		operatorId := u.Query().Get("operatorId")

		if resp, err := model.ParseResponse(respBody); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("解析三方游戏的/api/ThirdGame/GetGameUrl响应失败", err)), &ThirdGameType{}, err
		} else {
			return resp, &ThirdGameType{
				AuthToken:  authToken,
				OperatorId: operatorId,
			}, nil
		}
	}
}

// 生成三方游戏的operatorID
func GenerateOperatorID() (string, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
