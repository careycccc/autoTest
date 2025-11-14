package GameBetOrders

import (
	"autoTest/API/adminApi/login"
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// 游戏投注订单
type QueryBetRecordPageList struct {
	QueryTimeType string `json:"queryTimeType"`
	CategoryType  int    `json:"categoryType"`
	BeginTimeUnix int64  `json:"beginTimeUnix"`
	EndTimeUnix   int64  `json:"endTimeUnix"`
	SortField     string `json:"sortField"`
	model.QueryPayloadStruct
}

type BetSummary struct {
	BetAmountSum   float64 `json:"betAmountSum"`   // 总投注金额
	ValidAmountSum float64 `json:"validAmountSum"` // 有效投注金额
	WinAmountSum   float64 `json:"winAmountSum"`   // 总赢金额
	WinLoseAmount  float64 `json:"winLoseAmount"`  // 盈亏金额（负数表示亏损）
	FeeAmountSum   float64 `json:"feeAmountSum"`   // 总手续费
}

type SumResponse struct {
	Data struct {
		Data struct {
			Sum BetSummary `json:"sum"`
		} `json:"data"`
		List []struct {
			UserId int `json:"userId"`
		} `json:"list"`
		TotalCount int `json:"totalCount"` // 总下注条数
	} `json:"data"`
}

// 返回的结构体
type BetRecordPageList struct {
	BetSummary
	TotalCount int   // 总条数
	UserIdList []int // 用户id
}

/*
categoryType  // 0表示电子游戏  1表示真人视讯 2 表示体育竞技  3 彩票 4 棋牌
startTime // 开始时间
endTime int64  // 结束时间
sortField string // BetTime 投注时间
返回该游戏大类下的投注和税收的数据和总条数
*
*/
func QueryGameBetOrders(ctx *context.Context, categoryType int8, startTime, endTime int64, sortField string) (*model.Response, *BetRecordPageList, error) {
	api := "/api/ThirdGame/GetBetRecordPageList"
	payloadStruct := &QueryBetRecordPageList{}
	timestamp, random, language := request.GetTimeRandom()
	// 临时处理总条数为2000条
	payloadList := []interface{}{sortField, categoryType, startTime, endTime, sortField, 1, 2000, "Desc", random, language, "", timestamp}
	if respBody, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("", err)), &BetRecordPageList{}, err
	} else {
		var sumResponse SumResponse
		if err := json.Unmarshal(respBody, &sumResponse); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("", err)), &BetRecordPageList{}, err
		}
		if resp, err := model.ParseResponse(respBody); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("", err)), &BetRecordPageList{}, err
		} else {
			// 获取用户id列表
			userIdList := make([]int, len(sumResponse.Data.List))
			for i, v := range sumResponse.Data.List {
				userIdList[i] = v.UserId
			}
			return resp, &BetRecordPageList{
				BetSummary: sumResponse.Data.Data.Sum,
				TotalCount: sumResponse.Data.TotalCount,
				UserIdList: userIdList,
			}, nil
		}
	}
}

// 游戏大类的数量
const GameTypeNumber = 5

// 游戏大类的名称
var GameTypeName = []string{"电子游戏", "真人视讯", "体育竞技", "彩票", "棋牌"}

type GameBetOrders struct {
	Info     *BetRecordPageList
	Name     string
	UserList []int
}

// 游戏的查询的详细信息ch
var GameBetOrdersCh = make(chan *GameBetOrders, GameTypeNumber)

// 后台游戏订单
func RunGameBetOrders() {
	// 后台登录
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		logger.LogError("后台游戏订单查询的后台登录失败", err)
		return
	} else {
		startTime, endTime := utils.GetTodayStartAndEnd()
		var wg sync.WaitGroup
		wg.Add(GameTypeNumber)
		for i := 0; i < GameTypeNumber; i++ {
			go func(i int, chBetOrders chan<- *GameBetOrders) {
				defer wg.Done()
				if _, BetRecordPageList, err := QueryGameBetOrders(ctx, int8(i), startTime, endTime, "BetTime"); err != nil {
					str := fmt.Sprintf("%s后台游戏订单查询失败", GameTypeName[i])
					logger.LogError(str, err)
					return
				} else {
					// 计算单个游戏大类的人数
					userList := GameOrderPersonNumber(BetRecordPageList.UserIdList)
					chBetOrders <- &GameBetOrders{
						Info:     BetRecordPageList,
						Name:     GameTypeName[i],
						UserList: userList,
					}
				}
			}(i, GameBetOrdersCh)
		}
		wg.Wait()
		close(GameBetOrdersCh)
		AnalysisBetRecordPageList(GameBetOrdersCh)

	}
}

// 辅助函数，传入一个ch，给我进行汇总和分析
func AnalysisBetRecordPageList(ch <-chan *GameBetOrders) {
	var BetAmountSum, ValidAmountSum, WinAmountSum, WinLoseAmount, FeeAmountSum float64
	var totalUserList []int
	for betRecord := range ch {
		if betRecord.Info.BetAmountSum > 0 {
			// 有投注金额才进行打印
			logger.Logger.Info("游戏名称:", betRecord.Name, "\n", "投注金额:", betRecord.Info.BetAmountSum, "\n", "有效投注:", betRecord.Info.ValidAmountSum, "\n", "派奖金额:", betRecord.Info.WinAmountSum, "\n", "盈亏:", betRecord.Info.WinLoseAmount, "\n", "手续费:", betRecord.Info.FeeAmountSum, "\n", "游戏人数:", len(betRecord.UserList))
			//汇总计算
			BetAmountSum += betRecord.Info.BetAmountSum
			ValidAmountSum += betRecord.Info.ValidAmountSum
			WinAmountSum += betRecord.Info.WinAmountSum
			WinLoseAmount += betRecord.Info.WinLoseAmount
			FeeAmountSum += betRecord.Info.FeeAmountSum
			// 把用户列表进行合并
			totalUserList = append(totalUserList, betRecord.UserList...)
		}
	}
	// 计算总人数
	list := GameOrderPersonNumber(totalUserList)
	logger.Logger.Info("总计", "\n", "投注金额:", BetAmountSum, "\n", "有效投注:", ValidAmountSum, "\n", "派奖金额:", WinAmountSum, "\n", "盈亏:", WinLoseAmount, "\n", "手续费:", FeeAmountSum, "\n", "游戏人数:", len(list))
}

// 计算单个游戏大类的人数
func GameOrderPersonNumber(list []int) []int {
	// 获取用户id列表
	seen := make(map[int]bool)
	result := make([]int, 0, len(list))

	for _, item := range list {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	return result
}
