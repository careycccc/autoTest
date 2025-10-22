package recoversaasbalance

import (
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

// 提现大类的类型
type WithdrawCategory struct {
	ID                int     `json:"id"`
	TenantID          int     `json:"tenantId"`
	WithdrawType      string  `json:"withdrawType"`
	Name              string  `json:"name"`
	IconURL           string  `json:"iconUrl"`
	SelectedIconURL   string  `json:"selectedIconUrl"`
	UserMaxBindCount  int     `json:"userMaxBindCount"`
	MaxWithdrawTimes  int     `json:"maxWithdrawTimes"`
	MinAmount         float64 `json:"minAmount"`
	MaxAmount         float64 `json:"maxAmount"`
	FeeAmountRangeMin float64 `json:"feeAmountRangeMin"`
	FeeAmountRangeMax float64 `json:"feeAmountRangeMax"`
	FeeType           int     `json:"feeType"`
	Fee               float64 `json:"fee"`
	AllowStartTime    string  `json:"allowStartTime"`
	AllowEndTime      string  `json:"allowEndTime"`
	Sort              int     `json:"sort"`
	State             int     `json:"state"`
}

type WithdrawBasicInfoResp struct {
	Data struct {
		FixedWithdrawAmountList    []float64          `json:"fixedWithdrawAmountList"`    // 提现的金额的list表
		UserTodayWithdrawAmount    float64            `json:"userTodayWithdrawAmount"`    // 今日可提现的总金额
		UserTodayWithdrawCount     int                `json:"userTodayWithdrawCount"`     // 今日可提现的总次数
		UserTodayWithdrawFreeCount float64            `json:"userTodayWithdrawFreeCount"` // 打码量
		WithdrawCategoryList       []WithdrawCategory `json:"withdrawCategoryList"`       // 提现大类
	} `json:"data"`
}

type AllWithdraw struct {
	WithdrawAmountList         []float64
	UserTodayWithdrawAmount    float64
	UserTodayWithdrawCount     int
	UserTodayWithdrawFreeCount float64
	WithdrawCategoryList       []WithdrawCategory
}

// 获取提现的提现的金额列表，就是固定提现的那个list
func GetWithdrawBasicInfo(ctx *context.Context) (*model.Response, *AllWithdraw, error) {
	api := "/api/Withdraw/GetWithdrawBasicInfo"
	payloadStruct, payloadList := utils.BaseStructHandler()
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.InitStructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Withdraw/GetWithdrawBasicInfo请求失败", err)), nil, err
	} else {
		var balanceResp WithdrawBasicInfoResp
		if err := json.Unmarshal(respBoy, &balanceResp); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Withdraw/GetWithdrawBasicInfo[WithdrawBasicInfoResp]解析失败", err)), nil, err
		}

		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Withdraw/GetWithdrawBasicInfo解析失败", err)), nil, err
		} else {
			// 把返回的信息整合一下
			return resp, &AllWithdraw{
				WithdrawAmountList:         balanceResp.Data.FixedWithdrawAmountList,
				UserTodayWithdrawAmount:    balanceResp.Data.UserTodayWithdrawAmount,
				UserTodayWithdrawCount:     balanceResp.Data.UserTodayWithdrawCount,
				UserTodayWithdrawFreeCount: balanceResp.Data.UserTodayWithdrawFreeCount,
				WithdrawCategoryList:       balanceResp.Data.WithdrawCategoryList,
			}, nil
		}
	}
}
