package fundtransactionrecords

import (
	"autoTest/API/adminApi/login"
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
)

// 资金账变查询

type QueryRequest struct {
	// 财务类型列表
	FinancialTypeList []string `json:"financialTypeList"`
	// 查询开始时间（毫秒时间戳）
	StartTime int64 `json:"startTime"`
	// 查询结束时间（毫秒时间戳）
	EndTime int64 `json:"endTime"`
	model.QueryPayloadStruct
}

type QueryResponse struct {
	Data struct {
		List []struct {
			ID            string  `json:"id"`
			SysCurrency   string  `json:"sysCurrency"`
			UserID        int64   `json:"userId"` // 这里就是你想要的 userId
			VendorCode    string  `json:"vendorCode"`
			TenantAccount string  `json:"tenantAccount"`
			Type          string  `json:"type"`
			BeforeAmount  float64 `json:"beforeAmount"`
			Amount        float64 `json:"amount"` // 这里就是你想要的 amount
			BackAmount    float64 `json:"backAmount"`
			GameRate      float64 `json:"gameRate"`
			OrderNo       string  `json:"orderNo"`
			Remark        string  `json:"remark"`
			CreateTime    int64   `json:"createTime"`
		} `json:"list"`
		Summary struct {
			TotalAmount float64 `json:"totalAmount"`
		} `json:"summary"`
		PageNo     int `json:"pageNo"`
		TotalPage  int `json:"totalPage"`
		TotalCount int `json:"totalCount"`
	} `json:"data"`
	MsgParameters interface{} `json:"msgParameters"`
	Code          int         `json:"code"`
	Msg           string      `json:"msg"`
	MsgCode       int         `json:"msgCode"`
}

type FundTransactionInfo struct {
	UserID int     // userid
	Moneny float64 // 账变金额
}

// 只需要会员id和账变的金额以及总条数
type FundTransactionInfoList struct {
	UserInfo   []FundTransactionInfo
	TotalCount int
}

/*
financialTypeList 财务类型列表
startTime,endTime 查询开始时间（毫秒时间戳）查询结束时间（毫秒时间戳）
*
*/
func FundtransactionrecordsApi(ctx *context.Context, financialTypeList []string, startTime, endTime int64) (*model.Response, *FundTransactionInfoList, error) {
	api := "/api/Financial/GetPageList"
	payloadStruct := &QueryRequest{}
	timestamp, random, language := request.GetTimeRandom()
	// 暂时处理获取2000条数据
	payloadList := []interface{}{financialTypeList, startTime, endTime, 1, 2000, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Financial/GetPageList请求失败", err)), &FundTransactionInfoList{}, err
	} else {
		var resp QueryResponse
		if err := json.Unmarshal(respBoy, &resp); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Financial/GetPageList【QueryResponse】反序列化失败", err)), &FundTransactionInfoList{}, err
		} else {
			// 将响应结果进行反序列化
			if res, err := model.ParseResponse(respBoy); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Financial/GetPageList反序列化失败", err)), &FundTransactionInfoList{}, err
			} else {
				var funlist []FundTransactionInfo
				// 遍历这个响应list
				for _, v := range resp.Data.List {
					funlist = append(funlist, FundTransactionInfo{
						UserID: int(v.UserID),
						Moneny: v.Amount,
					})
				}
				return res, &FundTransactionInfoList{
					UserInfo:   funlist,
					TotalCount: resp.Data.TotalCount,
				}, nil
			}

		}

	}
}

// financialTypeList 传入类型 是一个[]string
func RunFundtransactionrecordsApi(financialTypeList []string) *FundTransactionInfoList {
	// 后台登录
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		logger.LogError("资金账变的查询后台登录失败", err)
		return &FundTransactionInfoList{}
	} else {
		_, startTime, _, endTime, _ := utils.ParseTimeRangeToTimestamp(config.StartTime, config.EndTime)
		if _, p, err := FundtransactionrecordsApi(ctx, financialTypeList, startTime, endTime); err != nil {
			logger.Logger.Warn("提现赔付账变的错误信息", err)
			return &FundTransactionInfoList{}
		} else {
			return p
		}
	}
}
