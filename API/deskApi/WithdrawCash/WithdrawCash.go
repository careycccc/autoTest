package withdrawcash

import (
	withdrawalorders "autoTest/API/adminApi/financialManagement/withdrawalOrders"
	addwallet "autoTest/API/adminApi/memberList/addWallet"
	recoversaasbalance "autoTest/API/deskApi/WithdrawCash/RecoverSaasBalance"
	getuserinfo "autoTest/API/deskApi/getUserinfo"
	registerapi "autoTest/API/deskApi/registerApi"
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"context"
	"encoding/json"
	"strconv"
	"sync"
)

// 后台提现的结构体
type withDrawaInfo struct {
	withDrawaAmont float64 // 提现的金额
	withDrawaType  string  // 提现的类型
}

// 提现
func RunWithDrawCase() {
	userName := "911204199711"

	_, deskCtx, err := registerapi.GeneralAgentRegister(userName)
	if err != nil {
		logger.LogError("提现登录失败", err)
		return
	}
	if deskCtx == nil {
		logger.LogError("提现登录返回空 context", nil)
		return
	}

	// 三个 channel
	moneyChan := make(chan *float64, 1)
	allWithdrawChan := make(chan *recoversaasbalance.AllWithdraw, 1)
	userIdChan := make(chan int, 1)

	var wg sync.WaitGroup
	wg.Add(3)

	// 全部改成真正的 goroutine（关键！）
	go func() {
		defer wg.Done()
		if _, allWithdraw, err := recoversaasbalance.GetWithdrawBasicInfo(deskCtx); err != nil {
			logger.LogError("获取提现基础信息失败", err)
			return
		} else {

			allWithdrawChan <- allWithdraw
		}
	}()

	go func() {
		defer wg.Done()
		if _, userInfo, err := getuserinfo.GetUserInfo(deskCtx); err != nil {
			logger.LogError("获取用户信息失败", err)
			return
		} else {
			userIdChan <- userInfo.UserID
		}
	}()

	go func() {
		defer wg.Done()
		if _, amount, err := recoversaasbalance.RecoverSaasBalance(deskCtx); err != nil {
			logger.LogError("恢复余额失败", err)
			return
		} else {

			moneyChan <- &amount
		}
	}()

	wg.Wait() // 现在这三个是真的并行，最慢的决定时间

	// 获取用户ID
	userid := <-userIdChan

	// 添加钱包后台异步执行（不阻塞！）
	addwallet.RunAddWallet(strconv.Itoa(userid))

	// 设置提现密码（同步执行，没问题）
	if _, err := SetWithdrawPasswordApi(deskCtx); err != nil {
		logger.LogError("设置提现密码失败", err)
		return
	} else {
		// 提现逻辑
		withDrawaChan := make(chan *withDrawaInfo, 1)
		WithDrawCase(deskCtx, <-moneyChan, <-allWithdrawChan, withDrawaChan)
		withDrawa := <-withDrawaChan
		// 下单
		withdrawalorders.RunWithdraw(userid, withDrawa.withDrawaType, withDrawa.withDrawaAmont, withDrawa.withDrawaAmont)
	}

}

// 提现
func WithDrawCase(ctx *context.Context, money *float64, allwithdraw *recoversaasbalance.AllWithdraw, ch chan<- *withDrawaInfo) {
	// 判断用户是否有钱，每日提现金额是否有值，提现是否有次数，打码量是否满足
	if *money <= 0.0 {
		logger.Logger.Warnln("提现获取用户金额小于等于0", nil)
		ch <- &withDrawaInfo{
			withDrawaAmont: 0,
			withDrawaType:  "",
		}
		return
	}
	if allwithdraw.UserTodayWithdrawCount == 0 || allwithdraw.AmountCoding != 0 {
		logger.Logger.Warn("用户的提现次数等于0,或者 用户的打码量不等于0", nil)
		ch <- &withDrawaInfo{
			withDrawaAmont: 0,
			withDrawaType:  "",
		}
		return
	}
	// 要保证提现金额要有大于整个提现list里面的值
	canWithDrawCaseList := filterGreaterOrEqual(*money, allwithdraw.WithdrawAmountList)
PT:
	canWithDrawCaseListLen := len(canWithDrawCaseList)
	i := 0
	if canWithDrawCaseListLen == 1 {
		i = 0
	} else {
		i = utils.RandInt(0, canWithDrawCaseListLen-1)
	}
	// 随机出来的值 大于 今日可提现的总金额
	if canWithDrawCaseList[i] > allwithdraw.UserTodayWithdrawAmount {
		goto PT
	}
	// 筛选出了可以提现的金额
	// 随机提现的大类
	WithdrawCategoryListLen := len(allwithdraw.WithdrawCategoryList)
PT2:
	j := 0
	if WithdrawCategoryListLen == 1 {
		j = 0
	} else {
		j = utils.RandInt(0, WithdrawCategoryListLen-1)
	}
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

	ch <- &withDrawaInfo{
		withDrawaAmont: canWithDrawCaseList[i],
		withDrawaType:  withdrawType,
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

type TaskYestarDayWithdrawaNumber struct {
	TotalWithdrawaNumber int // 昨日总的订单数
	IsWithdrawa          int // 触发了多少条提现补偿
}

type WithdrawalRequest struct {
	// 提现类型
	WithdrawType string `json:"withdrawType"`
	// 提现状态
	WithdrawState string `json:"withdrawState"`
	// 开始时间 (毫秒时间戳)
	StartTime int64 `json:"startTime"`
	// 结束时间 (毫秒时间戳)
	EndTime  int64 `json:"endTime"`
	PageNo   int   `json:"pageNo"`
	PageSize int   `json:"pageSize"`
	model.BaseStruct
}

type WithdrawHistoryInfo struct {
	CreateTime        int64   `json:"createTime"`        // 订单创建时间
	OrderNo           string  `json:"orderNo"`           // 订单号
	WithdrawType      string  `json:"withdrawType"`      // 提现通道类型
	Amount            float64 `json:"amount"`            // 订单金额
	CompensationState int     `json:"compensationState"` // 提现赔付的标志位 0 不需要赔付  1 需要赔付 2 待领取 3 已过期
}

// 内部列表项结构：只包含您需要的字段
type ListItem struct {
	OrderNo           string  `json:"orderNo"`
	WithdrawType      string  `json:"withdrawType"`
	Amount            float64 `json:"amount"`
	CreateTime        int64   `json:"createTime"`
	CompensationState int     `json:"compensationState"`
}

// Data 层结构
type Data struct {
	List       []ListItem `json:"list"`
	PageNo     int        `json:"pageNo"`
	TotalPage  int        `json:"totalPage"`
	TotalCount int        `json:"totalCount"`
}

// 完整的根结构
type Response struct {
	Data       Data   `json:"data"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	MsgCode    int    `json:"msgCode"`
	ServerTime int64  `json:"serverTime"`
}

// 提现历史 昨日
func GetWithdrawHistoryApi(ctx *context.Context, startTime, endTime int64) (*model.Response, []WithdrawHistoryInfo, error) {
	var WithdrawHistoryInfoList []WithdrawHistoryInfo
	api := "/api/Withdraw/GetWithdrawHistory"
	payloadStruct := &WithdrawalRequest{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{"", "", startTime, endTime, 1, 20, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Withdraw/GetWithdrawHistory请求报错", err)), nil, err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Withdraw/GetWithdrawHistory解析失败", err)), nil, err
		} else {
			// 解析想要的数据
			var res Response
			if err := json.Unmarshal(respBoy, &res); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Withdraw/GetWithdrawHistory【Response】解析失败", err)), nil, err
			} else {
				for _, v := range res.Data.List {
					WithdrawHistoryInfoList = append(WithdrawHistoryInfoList, WithdrawHistoryInfo{
						CreateTime:        v.CreateTime,
						OrderNo:           v.OrderNo,
						WithdrawType:      v.WithdrawType,
						Amount:            v.Amount,
						CompensationState: v.CompensationState,
					})
				}
			}
			return resp, WithdrawHistoryInfoList, nil
		}
	}
}
