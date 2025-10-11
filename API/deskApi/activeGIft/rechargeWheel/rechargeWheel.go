package rechargewheel

import (
	"autoTest/API/adminApi/login"
	getgiftinfo "autoTest/API/deskApi/activeGIft/GetGiftInfo"
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
	"fmt"
)

// 充值转盘的相关的逻辑

type UserRechargeWheelInfo struct {
	isOpenRechargeWheel          bool // 是否开启了充值转盘 是否显示充值转盘
	RechargeWheelRemainSpinCount int  // 充值轮盘剩余旋转次数
}

// 获取当前用户充值转盘的，开启信息，剩余旋转次数
func GetUserRechargeWheelInfo(ctx *context.Context) (UserRechargeWheelInfo, error) {
	if _, rechargeWheelInfo, err := getgiftinfo.GetGiftInfoApi(ctx); err != nil {
		return UserRechargeWheelInfo{}, err
	} else {
		info := UserRechargeWheelInfo{
			isOpenRechargeWheel:          rechargeWheelInfo.IsOpenRechargeWheel.(bool),
			RechargeWheelRemainSpinCount: rechargeWheelInfo.RechargeWheelRemainSpinCount.(int),
		}
		return info, nil
	}
}

type SetRechargeWheelConditionStruct struct {
	SettingKey any `json:"settingKey"`
	Value1     any `json:"value1"` // 0表示不需要充值   1，首充 2，二充  3，三充
	model.BaseStruct
}

// 设置充值转盘的条件 无需首充，需首充 ，二充，三充
func SetRechargeWheelCondition(ctx *context.Context, value1 int8) (*model.Response, error) {
	api := "/api/RechargeWheel/UpdateConfig"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &SetRechargeWheelConditionStruct{}
	payloadList := []interface{}{"RechargeWheelNeedFirstRechargeSwitch", value1, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/RechargeWheel/UpdateConfig 请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/RechargeWheel/UpdateConfig 解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

// 充值转盘第一个转盘的充值配置
type GetFirstRechargeWheel struct {
	RechargeWheelType any `json:"rechargeWheelType"` // 默认值1
	model.BaseStruct
}

// 充值转盘的任务配置
type TaskConfig struct {
	Id             any `json:"id"`
	RechargeType   any `json:"rechargeType"`   // 1 表示累计充值  2表示循环充值
	RechargeAmount any `json:"rechargeAmount"` // 充值金额
	SpinCount      any `json:"spinCount"`      // 奖励转盘的次数
}

// 定义结构体来映射 JSON 数据
type TaskConfigResponse struct {
	Data struct {
		TaskConfig []TaskConfig `json:"taskConfig"`
	} `json:"data"`
}

// 获取充值转盘第一个转盘的充值配置
func GetFirstRechargeWheelInfo(ctx *context.Context) (*model.Response, []TaskConfig, error) {
	api := "/api/RechargeWheel/Get"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &GetFirstRechargeWheel{}
	payloadList := []interface{}{1, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("", err)), nil, err
	} else {
		var task TaskConfigResponse
		if err := json.Unmarshal(respBoy, &task); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("", err)), nil, err
		} else {
			if resp, err := model.ParseResponse(respBoy); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("", err)), nil, err
			} else {
				return resp, task.Data.TaskConfig, nil
			}
		}
	}
}

// 进行比较，保证充值的金额至少要满足有充值保存有旋转的次数产生
func ReturnRechargeAmount(ctx *context.Context) (amount float64) {
	if _, list, err := GetFirstRechargeWheelInfo(ctx); err != nil {
		return
	} else {
		if len(list) > 1 {
			for i := 0; i < len(list)-1; i++ {
				if list[i].RechargeAmount.(float64) >= list[i+1].RechargeAmount.(float64) {
					amount = list[i].RechargeAmount.(float64)
				} else {
					amount = list[i+1].RechargeAmount.(float64)
				}
			}
		} else if len(list) == 1 {
			amount = list[0].RechargeAmount.(float64)
		} else {
			logger.LogError("没有获取到充值转盘的配置项", err)
			return
		}
	}
	return
}

type GetPageListRewardRecordStruct struct {
	UserId any `json:"userId"`
	model.QueryPayloadStruct
}

// 定义结构体来映射 JSON 数据
type GetPageListRewardRecordResponse struct {
	Data struct {
		List []GetPageListRewardRecordList `json:"list"`
	} `json:"data"`
}

type GetPageListRewardRecordList struct {
	UserId            any `json:"userId"`
	RechargeWheelType any `json:"rechargeWheelType"`
	RewardType        any `json:"rewardType"`
	RewardAmount      any `json:"rewardAmount"`
	CreateTime        any `json:"createTime"`
}

// 转盘奖励记录
// 返回用户id，
func GetPageListRewardRecord(ctx *context.Context, userId int) (*model.Response, []GetPageListRewardRecordList, error) {
	api := "/api/RechargeWheel/GetPageListRewardRecord"
	timestamp, random, language := request.GetTimeRandom()
	payloadStruct := &GetPageListRewardRecordStruct{}
	payloadList := []interface{}{userId, 1, 20, "Desc", random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("", err)), nil, err
	} else {
		var record GetPageListRewardRecordResponse
		if err := json.Unmarshal(respBoy, &record); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("", err)), nil, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("", err)), nil, err
		} else {
			return resp, record.Data.List, nil
		}
	}
}

// 运行充值转盘的任务
func RunRechargeWheelCondition() {
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		logger.LogError("充值转盘的后台登录失败", err)
		return
	} else {
		amount := ReturnRechargeAmount(ctx) // 要充值的金额
		fmt.Println(amount)
		resp, list, _ := GetPageListRewardRecord(ctx, 2440414)
		fmt.Println(resp)
		fmt.Println(list)
	}
}

func CallRechargeWheelCondition(value1 int8) {
	// 第一步后台登录
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		logger.LogError("充值转盘的后台登录失败", err)
		return
	} else {
		//var wg *sync.WaitGroup
		go func() {
			// 设置条件
			if _, err := SetRechargeWheelCondition(ctx, value1); err != nil {
				return
			}
		}()

	}
}
