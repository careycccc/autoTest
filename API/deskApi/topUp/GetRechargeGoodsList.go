package topup

import (
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

type Response struct {
	Data       []DataItem `json:"data"`
	Code       int        `json:"code"`
	Msg        string     `json:"msg"`
	MsgCode    int        `json:"msgCode"`
	ServerTime int64      `json:"serverTime"`
}

type DataItem struct {
	ID                  int               `json:"id"`
	GoodsImg            string            `json:"goodsImg"`
	RechargeAmount      float64           `json:"rechargeAmount"`
	GiftAmount          float64           `json:"giftAmount"`
	SupportCategoryIds  string            `json:"supportCategoryIds"`
	SupportCategories   []SupportCategory `json:"supportCategories"`
	IsRecommendedAmount bool              `json:"isRecommendedAmount"`
}

// 单个金额下面的提现通道类型
type SupportCategory struct {
	ID                int               `json:"id"`
	Name              string            `json:"name"`
	RechargeType      string            `json:"rechargeType"`
	State             int               `json:"state"`
	Sort              int               `json:"sort"`
	IconUrl           string            `json:"iconUrl"`
	SelectedIconUrl   string            `json:"selectedIconUrl"`
	Rate              float64           `json:"rate"`
	MinAmount         float64           `json:"minAmount"`
	MaxAmount         float64           `json:"maxAmount"`
	RechargeGiftRatio RechargeGiftRatio `json:"rechargeGiftRatio"`
	QuickConfigList   []interface{}     `json:"quickConfigList"`
	GiftRatioType     int               `json:"giftRatioType"`
	GiftAmount        float64           `json:"giftAmount"`
}

type RechargeGiftRatio struct {
	GiftRatioType     int         `json:"giftRatioType"`
	ScaleType         int         `json:"scaleType"`
	UniformRatioData  interface{} `json:"uniformRatioData"`
	IntervalRatioList interface{} `json:"intervalRatioList"`
}

type GetRechargeGoods struct {
	RechargeAmount   float64           // 充值的金额
	RechargeGoodsId  int8              // 商品编号
	SupportCategorys []SupportCategory // 金额下面的提现通道的集合
}

// 获取充值的充值金额键盘和配置
// 返回GetRechargeGoods，一个金额对应下面的充值的通道
func GetRechargeGoodsListApi(ctx *context.Context) (*model.Response, *[]GetRechargeGoods, error) {
	api := "/api/Recharge/GetRechargeGoodsList"
	payloadStruct, payloadList := utils.BaseStructHandler()
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.InitStructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Recharge/GetRechargeGoodsList请求失败", err)), nil, err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Recharge/GetRechargeGoodsList解析失败", err)), nil, err
		} else {
			var res Response
			if err := json.Unmarshal(respBoy, &res); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Recharge/GetRechargeGoodsList[Response]解析失败", err)), nil, err
			}
			listRechargeList := make([]GetRechargeGoods, 0, 40)
			for _, v := range res.Data {
				listRechargeList = append(listRechargeList, GetRechargeGoods{
					RechargeAmount:   v.RechargeAmount,
					RechargeGoodsId:  int8(v.ID),
					SupportCategorys: v.SupportCategories,
				})
			}
			return resp, &listRechargeList, nil
		}
	}
}
