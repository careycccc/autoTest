package lotterygameapi

// 彩票投注

import (
	"autoTest/store/config"
	"autoTest/store/model"
	"autoTest/store/request"
	"fmt"
)

// BetRequest 定义请求体的结构体
type BetRequest struct {
	GameCode    string `json:"gameCode"`
	IssueNumber string `json:"issueNumber"`
	Amount      int    `json:"amount"`
	BetMultiple int    `json:"betMultiple"`
	BetContent  string `json:"betContent"`
	Language    string `json:"language"`
	Random      int64  `json:"random"`
	Signature   string `json:"signature"`
	Timestamp   int64  `json:"timestamp"`
}

type ResponseStruct struct {
	Code        int
	Msg         string
	MsgCode     int
	ServiceTime int64
}

/*
*
gameCode  彩票投注种类
amount 投注金额 = 单个金额 * 倍率
betMultiple 投注倍率
betContent 投注盘口
issueNumber 期号
token token对象
*/
func BetWingo(gameCode string, amount, betMultiple int, betContent, issueNumber, token, username string) (*model.BetResponse, error) {
	var apiArg string
	switch gameCode {
	case "TrxWinGo_10M":
		apiArg = "TrxWinGoBet"
	case "WinGo_5M":
		apiArg = "WinGoBet"
	}
	// 请求体地址
	api := "/api/Lottery/" + apiArg
	url := config.LOTTERY_H5
	// 参数化
	bet := &BetRequest{}
	timestamp, random, language := request.GetTimeRandom()
	betResultList := []interface{}{gameCode, issueNumber, amount, betMultiple, betContent, language, random, "", timestamp}
	resultMap, _ := request.InitStructToMap(bet, betResultList)
	// 获取请求头
	deskA := &BetTokenStruct{}
	url_h5 := config.WMG_H5
	desSlice := []interface{}{url_h5, url_h5, token}
	headMap, _ := request.AssignSliceToStructMap(deskA, desSlice)
	respBody, _, err := request.PostRequestCofig(resultMap, url, api, headMap)
	if err != nil {
		//logger.LogError(api+"报错消息", err)
		errs := fmt.Errorf(api+"报错消息%s", err)
		return model.HandlerErrorRes2(errs), nil
	}
	if resp, err := model.ParseResponse2(respBody); err != nil {
		//logger.LogError("/api/Home/Login 响应解析失败", err)
		errs := fmt.Errorf("/api/Home/Login 响应解析失败%s", err)
		return model.HandlerErrorRes2(errs), err
	} else {
		code := resp.Code
		msgcode := resp.MsgCode
		msg := resp.Msg
		if code == 0 && msgcode == 0 && msg == "Succeed" {
			return resp, nil
		} else {
			return resp, nil
		}
	}

}
