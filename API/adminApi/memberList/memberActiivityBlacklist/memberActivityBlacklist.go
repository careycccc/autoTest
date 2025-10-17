package memberactiivityblacklist

import (
	"autoTest/API/adminApi/login"
	memberlist "autoTest/API/adminApi/memberList/memberList"
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"context"
	"strconv"
)

// 添加会员黑名单
type UserActivityBlockStruct struct {
	Remark            any `json:"remark"`  // 备注
	UserIds           any `json:"userIds"` // 会员IDs  多个会员用\n
	ActivityBlockType any `json:"activityBlockType"`
	model.BaseStruct      // 黑名单类型 1 活动黑名单 2 投注黑名单
}

/*
activityBlockType  11 三层代理返佣，13 会员排行榜奖金，14 代理排行榜奖金
userIds 1245784\n1245786
*
*/
func AddMemberActivityBlacklistApi(ctx *context.Context, activityBlockType int8, userIds string) (*model.Response, error) {
	api := "/api/UserActivityBlock/Add"
	payloadStruct := &UserActivityBlockStruct{}
	remark := "测试11111"
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{remark, userIds, activityBlockType, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/UserActivityBlock/Add 请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/UserActivityBlock/Add 响应解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

// 从yaml文件中读取用户信息，batchAddMemberActivityBlacklist 批量添加会员黑名单
func RunAddBlacklist() {
	if list, err := utils.ReadYAMLByLine(config.SUBUSERYAML); err != nil {
		logger.LogError("读取yaml文件失败", err)
		return
	} else {
		ctxToken, err := login.RunAdminSitLogin()
		if err != nil {
			logger.LogError("登录失败", err)
			return
		}
		for _, v := range list {
			// fmt.Print("开始添加会员黑名单，类型为11，用户为", v, "\n")
			// time.Sleep(time.Second * 1)
			// 根据账号获取用户id
			if _, userId, err := memberlist.GetUserIdApi(ctxToken, v); err != nil {
				logger.LogError("获取用户id失败", err)
				break
			} else {
				str := strconv.FormatInt(userId, 10)
				if resp, err := AddMemberActivityBlacklistApi(ctxToken, 14, str); err != nil {
					logger.LogError("添加会员黑名单失败", err)
					break
				} else {
					logger.Logger.Info("添加会员黑名单成功", "响应结果", resp)
				}
			}

		}
	}
	// if list, err := utils.ReadYAMLByLine(config.SUBUSERYAML); err != nil {
	// 	logger.LogError("读取yaml文件失败", err)
	// 	return
	// } else {
	// 	ctxToken, err := login.RunAdminSitLogin()
	// 	if err != nil {
	// 		logger.LogError("登录失败", err)
	// 		return
	// 	}
	// 	li := []int8{11, 13, 14}
	// 	for _, v := range list {
	// 		time.Sleep(time.Second * 1)
	// 		for _, j := range li {
	// 			fmt.Print("开始添加会员黑名单，类型为", j, "用户为", v, "\n")
	// 			if resp, err := AddMemberActivityBlacklistApi(ctxToken, j, v); err != nil {
	// 				logger.LogError("添加会员黑名单失败", err)
	// 				break
	// 			} else {
	// 				logger.Logger.Info("添加会员黑名单成功", "响应结果", resp)
	// 			}
	// 		}

	// 	}
	// }

}
