package dailycheckin

import (
	"autoTest/API/adminApi/common"
	memberlist "autoTest/API/adminApi/memberList/memberList"
	"autoTest/PressureMeasurementModule/accounts"
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
	"fmt"
)

type GetDailyCheckInUserListStruct struct {
	ActivityName string `json:"activityName"` // 活动名称
	StartDate    string `json:"startDate"`
	EndDate      string `json:"endDate"`
	model.QueryPayloadStruct
}

type Activity struct {
	ActivityID     int     `json:"activityId"`     // 活动ID
	ActivityName   string  `json:"activityName"`   // 活动名称
	ActivityType   int     `json:"activityType"`   // 活动类型
	CodingMultiple int     `json:"codingMultiple"` // 编码倍数
	DayIndex       int     `json:"dayIndex"`       // 天数索引
	OrderNo        string  `json:"orderNo"`        // 订单号
	ReceiveAmount  float64 `json:"receiveAmount"`  // 接收金额
	ReceiveMode    int     `json:"receiveMode"`    // 接收方式
	ReceiveTime    int64   `json:"receiveTime"`    // 接收时间
	RechargeAmount float64 `json:"rechargeAmount"` // 充值金额
	StatisticsType int     `json:"statisticsType"` // 统计类型
	UserID         int     `json:"userId"`         // 用户ID
}

type DailyCheckInUserResponse struct {
	Data struct {
		List []Activity `json:"list"`
	} `json:"data"`
}

/*
GetDailyCheckInUserList 获取每日签到用户奖励领取记录
activityName 活动名称
startDate 开始时间
endDate 结束时间
返回用户ID列表

	*
*/
func GetDailyCheckInUserList(ctx *context.Context, activityName, startDate, endDate string) (*model.Response, []int, error) {
	api := "/api/DailyCheckIn/GetDailyCheckInUserList"
	payloadStruct := &GetDailyCheckInUserListStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{activityName, startDate, endDate, 1, 200, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("", err)), nil, err
	} else {
		if string(respBoy) == "" {
			return model.HandlerErrorRes(model.ErrorLoggerType("respBoy is nil", nil)), nil, nil
		} else {
			if resp, err := model.ParseResponse(respBoy); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/DailyCheckIn/GetDailyCheckInUserList解析失败", err)), nil, err
			} else {
				var res DailyCheckInUserResponse
				if err := json.Unmarshal(respBoy, &res); err != nil {
					return model.HandlerErrorRes(model.ErrorLoggerType("/api/DailyCheckIn/GetDailyCheckInUserList【DailyCheckInUserResponse】解析失败", err)), nil, err
				} else {
					userIds := make([]int, 0, len(res.Data.List))
					for _, v := range res.Data.List {
						userIds = append(userIds, v.UserID)
					}
					return resp, userIds, nil
				}
			}
		}
	}
}

// CountElements 统计切片中每个元素出现的次数
func CountElements[T comparable](slice []T) map[T]int {
	countMap := make(map[T]int)
	for _, element := range slice {
		countMap[element]++
	}
	return countMap
}
// 查询奖励领取记录表
func RunDailyCheckInUserList() {
	// 获取第一天的签到用户
	_, userIds1, err := GetDailyCheckInUserList(AdminCtx, "vip6-7每日签到", "2025-12-24 00:00:00", "2025-12-30 23:59:59")
	if err != nil {
		logger.LogError("获取每日签到用户奖励领取记录失败", err)
		return
	}
	// 循环切片找出相同的id并且记录次数
	// 初始化一个 map，key 为元素内容，value 为出现次数
	counts := make(map[int]int)

	for _, item := range userIds1 {
		// 在 Go 中，访问 map 中不存在的 key 会返回该类型的零值（int 的零值是 0）
		// 所以直接 ++ 即可实现计数
		counts[item]++
	}
	for key, value := range counts {
		fmt.Printf("用户 %d 签到了 %d 次\n", key, value)
	}
	amountList := make([]string, 0, len(userIds1))
	// 进行
	for _,userid := range userIds1 {
		if userid == 131560 {
			continue
		}
	    // 转成账号并且修改密码
		amount := common.IdToAmountAndUpdatePassword(AdminCtx,userid)
		amountList = append(amountList, amount)
		// 进行解冻
		memberlist.UpdateUserStateApi(AdminCtx,userid,1)
	}
	accounts.WriteConcurrently(amountList, 5, CSVADDR) // 保存到csv中

	// // 获取第二天的签到用户
	// _, userIds2, err := GetDailyCheckInUserList(AdminCtx, "3-2的新的不连续签到活动", "2025-12-18 00:00:00", "2025-12-18 23:59:59")
	// if err != nil {
	// 	logger.LogError("获取每日签到用户奖励领取记录失败", err)
	// 	return
	// }

	// 使用map来提高查找效率
	// userMap := make(map[int]bool)
	// for _, id := range userIds1 {
	// 	userMap[id] = true
	// }

	// // 找出两天都签到的用户
	// commonUsers := make([]int, 0)
	// for _, id := range userIds2 {
	// 	if userMap[id] {
	// 		commonUsers = append(commonUsers, id)
	// 		fmt.Printf("用户 %d 在两天都签到了\n", id)
	// 	}
	// }

	// // 使用CountElements统计每个用户签到的次数
	// allUsers := append(userIds1, userIds2...)
	// checkInCount := CountElements(allUsers)
	// fmt.Println("\n签到统计：")
	// for userID, count := range checkInCount {
	// 	fmt.Printf("用户 %d 签到了 %d 次\n", userID, count)
	// }

}


func RunDailyCheckInActivityList() {
    // 获取第一天的签到用户
	_, userIds1, err := GetDailyCheckInUserList(AdminCtx, "3-2的新的循环签到", "2025-12-23 00:00:00", "2025-12-24 23:59:59")
	if err != nil {
		logger.LogError("获取每日签到用户奖励领取记录失败", err)
		return
	}
	userlist := make([]string, 0, len(userIds1))
	// 循环遍历用户ID列表，进行修改密码，并且保存到csv中
	for _, userId := range userIds1 {
		if userId == 131560 { continue }
		if _,err := memberlist.UpdatePassword(AdminCtx,int64(userId),config.ADMIN_PWD);err != nil {
			logger.LogError("修改密码失败", err)
			continue
		}else {
			if _,amount,err := memberlist.GetUserAmount(AdminCtx,userId);err != nil {
				logger.LogError("获取用户账号失败", err)
				continue
			}else{
				//保存到csv中
				userlist = append(userlist, amount)
			}
		}
		
	}

	accounts.WriteConcurrently(userlist, 5, CSVADDR) // 保存到csv中
}
