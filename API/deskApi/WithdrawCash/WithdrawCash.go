package withdrawcash

import (
	recoversaasbalance "autoTest/API/deskApi/WithdrawCash/RecoverSaasBalance"
	registerapi "autoTest/API/deskApi/registerApi"
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"context"
	"sync"
)

// 提现
func RunWithDrawCase() {
	// 用户的手机号码
	userName := "911022199712"
	// 判断当前用户是否有钱
	if _, ctx, err := registerapi.GeneralAgentRegister(userName); err != nil {
		logger.LogError("提现登录失败", err)
		return
	} else {
		deskToken := ctx
		// 设置提现密码
		resp, err := SetWithdrawPasswordApi(deskToken)
		if err != nil {
			logger.LogError("提现密码设置失败", err)
			return
		}
		logger.Logger.Info("提现密码设置的结果", resp)
		wg := &sync.WaitGroup{}
		moneyChan := make(chan *float64, 1)
		allWithdrawChan := make(chan *recoversaasbalance.AllWithdraw, 1)
		wg.Add(2)
		go func(ctx *context.Context, ch chan<- *float64) {
			defer wg.Done()
			if _, amount, err := recoversaasbalance.RecoverSaasBalance(ctx); err != nil {
				logger.LogError("提现获取用户金额失败", err)
				return
			} else {
				ch <- &amount
			}
		}(deskToken, moneyChan)
		go func(ctx *context.Context, ch chan<- *recoversaasbalance.AllWithdraw) {
			defer wg.Done()
			if _, allWithdraw, err := recoversaasbalance.GetWithdrawBasicInfo(ctx); err != nil {
				logger.LogError("提现获取用户金额失败", err)
				return
			} else {
				ch <- allWithdraw
			}
		}(deskToken, allWithdrawChan)
		wg.Wait()
		WithDrawCase(deskToken, <-moneyChan, <-allWithdrawChan)
	}
}

// 提现
func WithDrawCase(ctx *context.Context, money *float64, allwithdraw *recoversaasbalance.AllWithdraw) {
	// 判断用户是否有钱，每日提现金额是否有值，提现是否有次数，打码量是否满足
	if *money <= 0.0 {
		logger.LogError("提现获取用户金额小于等于0", nil)
		return
	}
	if allwithdraw.UserTodayWithdrawCount <= 0 || allwithdraw.UserTodayWithdrawFreeCount != 0 {
		logger.LogError("用户的提现次数小于等于0,或者 用户的打码量不等于0", nil)
		return
	}
	// 要保证提现金额要有大于整个提现list里面的值
	canWithDrawCaseList := filterGreaterOrEqual(*money, allwithdraw.WithdrawAmountList)
PT:
	canWithDrawCaseListLen := len(canWithDrawCaseList)
	i := utils.RandInt(0, canWithDrawCaseListLen)
	// 随机出来的值 大于 今日可提现的总金额
	if canWithDrawCaseList[i] > allwithdraw.UserTodayWithdrawAmount {
		goto PT
	}
	// 筛选出了可以提现的金额
	// 随机提现的大类
	WithdrawCategoryListLen := len(allwithdraw.WithdrawCategoryList) - 1
PT2:
	j := utils.RandInt(0, WithdrawCategoryListLen)
	if allwithdraw.WithdrawCategoryList[j].WithdrawType == "UPI" {
		// 提现类型目前不支持upi
		goto PT2
	}
	withdrawType := allwithdraw.WithdrawCategoryList[j].WithdrawType
	withdrawId := allwithdraw.WithdrawCategoryList[j].ID
	logger.Logger.Info("提现通道", withdrawType)
	// 进行提现
	if resp, err := WithdrawApplyApi(ctx, canWithDrawCaseList[i], withdrawId, withdrawType); err != nil {
		logger.LogError("提现失败", err)
		return
	} else {
		logger.Logger.Info("提现结果", resp)
	}
}

// 返回可以提现的list列表
func filterGreaterOrEqual(threshold float64, numbers []float64) []float64 {
	result := []float64{}
	for _, num := range numbers {
		if num <= threshold {
			result = append(result, num)
		}
	}
	return result
}

// 提现请求
type WithdrawApplyStruct struct {
	Amount             any `json:"amount"`             // 提现的金额
	WalletId           any `json:"walletId"`           // 提现的随机号
	WithdrawCategoryId any `json:"withdrawCategoryId"` // 提现通道的id
	WithdrawType       any `json:"withdrawType"`       // 提现通道类型
	WithdrawPassword   any `json:"withdrawPassword"`   // 提现密码
	model.BaseStruct
}

/*
提现请求
需要传入Amount float64, 提现金额
WithdrawCategoryId int, 提现通道的id
WithdrawType string 提现通道类型
*
*/
func WithdrawApplyApi(ctx *context.Context, Amount float64, WithdrawCategoryId int, WithdrawType string) (*model.BetResponse, error) {
	api := "/api/Withdraw/WithdrawApply"
	payloadStruct := &WithdrawApplyStruct{}
	_, walletId, err := GetUserWithdrawWallet(ctx, WithdrawType)
	if err != nil {
		return &model.BetResponse{}, nil
	}
	// fmt.Println("提现的walletId", walletId)
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{Amount, walletId, WithdrawCategoryId, WithdrawType, config.WithdrawPassword, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Withdraw/WithdrawApply请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse2(respBoy); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Withdraw/WithdrawApply解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}
