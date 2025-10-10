package lotterygameapi

import (
	"autoTest/store/config"
	"autoTest/store/request"
	"encoding/json"
	"time"
)

// 获取期号
type IssnumResponse struct {
	Code int `json:"code"`
	Data struct {
		StartTime      int64  `json:"startTime"`
		EndTime        int64  `json:"endTime"`
		IssueNumber    string `json:"issueNumber"`
		IntervalMinute int64  `json:"intervalMinute"`
		GameCode       string `json:"gameCode"`
		Diif           int    `json:"diif"`
		Countdown      int    `json:"countdown"`
	} `json:"data"`
	Msg        string `json:"msg"`
	MsgCode    int    `json:"msgCode"`
	ServerTime int64  `json:"serverTime"`
}

// 获取当期的期号
func GetNowBetNumber(token, betType, lotteryType string) (map[string]interface{}, error) {
	api := "/webapi/kv/issue/" + betType
	// 获取路径地址
	base_url := config.WMG_H5
	// 获取请求头

	headMap := GetIssNunmberHeaderFunc(token, betType, lotteryType)
	k := make(map[string]interface{})
	respBody, _, err := request.GetRequest(base_url, api, headMap, k)
	if err != nil {
		//fmt.Printf("%v", err)
		return nil, err
	}
	// fmt.Printf("响应期号%v", string(resp))
	var response IssnumResponse
	error := json.Unmarshal([]byte(string(respBody)), &response)
	if error != nil {
		//fmt.Printf("响应期号解析失败%v", error)
		return nil, error
	}
	nowBetNumber := map[string]interface{}{
		"startTime":      response.Data.StartTime,      // 开始时间
		"endTime":        response.Data.EndTime,        // 结束时间
		"issueNumber":    response.Data.IssueNumber,    // 期号
		"intervalMinute": response.Data.IntervalMinute, // 间隔时间
	}
	// fmt.Println(nowBetNumber)
	return nowBetNumber, nil
}

/*
判断是否可以下注,并且返回期号
betType 投注的方式 wingo 30s  wingo1min  wingo 3min  wing 5min
lotteryType 彩票类型  WinGo
*
*/
func IsBet(token, betType, lotteryType string) (bool, string) {
	nowBetNumber, err := GetNowBetNumber(token, betType, lotteryType)
	if err != nil {
		//fmt.Println("没有成功获取到期号")
		return false, "-1"
	}

	endTime := nowBetNumber["endTime"].(int64)
	issueNumber := nowBetNumber["issueNumber"]
	// 获取当前时间
	now := time.Now()
	// 获取时间戳（秒）
	secTimestamp := now.UnixMilli()
	// 结束时间 - 当前时间 >= 动画7s
	if endTime-secTimestamp >= 7000 {
		// 可以投注
		return true, issueNumber.(string)
	} else {
		// 不可以投注
		// 需要等待7s
		time.Sleep(time.Second * 7)
		return IsBet(token, betType, lotteryType)
	}
}

/*
token
betType 投注的方式 wingo 30s  wingo1min  wingo 3min  wing 5min
lotteryType 彩票类型  WinGo
*
*/
func GetIssNunmberHeaderFunc(token, betType, lotteryType string) map[string]interface{} {
	// result := "https://h5.wmgametransit.com/WinGo/"
	result := config.WMG_H5 + lotteryType
	if token == "" {
		//游客的方式
		result = result + betType
	} else {
		// token有值的情况
		r1 := "?Lang=en&Skin=Classic&SkinColor=Default&Token="
		r2 := "&RedirectUrl=" + config.PLANT_H5 + "%2Fgame%2Fcategory%3FcategoryCode%3DC202505280608510046&Beck=0"
		result = result + betType + r1 + token + r2
	}
	return map[string]interface{}{
		"Referer": result,
	}
}
