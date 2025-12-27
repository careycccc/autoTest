package rechargewheel

import (
	"autoTest/API/adminApi/login"
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
)

//后台充值转盘新增
type RechargeWheelConfig struct {
    RechargeWheelType     int          `json:"rechargeWheelType"`
    SpecialWheelUnlockCond int         `json:"specialWheelUnlockCond"`
    TaskConfig            []TaskConfig `json:"taskConfig"`
    RewardConfig          []RewardConfig `json:"rewardConfig"`
    Random                int64        `json:"random"`
    Language              string       `json:"language"`
    Signature             string       `json:"signature"`
    Timestamp             int64        `json:"timestamp"`
}

type TaskConfig struct {
    ID             int `json:"id"`
    RechargeType   int `json:"rechargeType"`
    RechargeAmount int `json:"rechargeAmount"`
    SpinCount      int `json:"spinCount"`
}

type RewardConfig struct {
    ID           int    `json:"id"`
    RewardType   int    `json:"rewardType"`
    RewardAmount int    `json:"rewardAmount"`
    WashCode     int    `json:"washCode"`
    Weight       int    `json:"weight"`
    Icon         string `json:"icon"`
}
var (
	timestamp, random, language = request.GetTimeRandom()
	// 白银盘的增加
	payloadListType1 = []interface{}{
		1, // rechargeWheelType
		0, // specialWheelUnlockCond

		// taskConfig []TaskConfig
		[]TaskConfig{
			{ID: 2, RechargeType: 1, RechargeAmount: 300, SpinCount: 1},
			{ID: 3, RechargeType: 2, RechargeAmount: 500, SpinCount: 1},
			{ID: 13, RechargeType: 1, RechargeAmount: 800, SpinCount: 0},
		},

		// rewardConfig []RewardConfig
		[]RewardConfig{
			{ID: 4, RewardType: 1, RewardAmount: 17, WashCode: 2, Weight: 30, Icon: "3004/banner/123432408-3064-file_20251124123432407.webp"},
			{ID: 5, RewardType: 1, RewardAmount: 57, WashCode: 2, Weight: 30, Icon: "3004/banner/123753240-3067-file_20251124123753236.webp"},
			{ID: 6, RewardType: 1, RewardAmount: 77, WashCode: 2, Weight: 20, Icon: "3004/banner/123744166-3066-file_20251124123744165.webp"},
			{ID: 7, RewardType: 1, RewardAmount: 177, WashCode: 2, Weight: 20, Icon: "3004/banner/124125067-3071-file_20251124124124971.webp"},
			{ID: 8, RewardType: 1, RewardAmount: 277, WashCode: 2, Weight: 10, Icon: "3004/banner/124130579-3072-file_20251124124130575.webp"},
			{ID: 9, RewardType: 1, RewardAmount: 377, WashCode: 2, Weight: 10, Icon: "3004/banner/124135709-3073-file_20251124124135705.webp"},
			{ID: 11, RewardType: 3, RewardAmount: 1, WashCode: 1, Weight: 10, Icon: "3004/banner/051701003-3135-file_20251125051701002.webp"},
			{ID: 12, RewardType: 2, RewardAmount: 1, WashCode: 1, Weight: 20, Icon: "3004/banner/051707876-3136-file_20251125051707875.webp"},
		},

		random, // random
		language,                // language
		"", // signature
		timestamp,   // timestamp
	}
	// 黄金盘的增加
	payloadListType2 = []interface{}{
		2, // rechargeWheelType
		0, // specialWheelUnlockCond
	
		// taskConfig []TaskConfig
		[]TaskConfig{
			{ID: 2, RechargeType: 1, RechargeAmount: 600, SpinCount: 1},
			{ID: 3, RechargeType: 2, RechargeAmount: 1200, SpinCount: 1},
		},
	
		// rewardConfig []RewardConfig
		[]RewardConfig{
			{ID: 4, RewardType: 1, RewardAmount: 77, WashCode: 2, Weight: 30, Icon: "3004/banner/033305278-3074-file_20251125033305277.webp"},
			{ID: 5, RewardType: 1, RewardAmount: 177, WashCode: 2, Weight: 30, Icon: "3004/banner/033313850-3075-file_20251125033313849.webp"},
			{ID: 6, RewardType: 1, RewardAmount: 277, WashCode: 2, Weight: 30, Icon: "3004/banner/033321065-3076-file_20251125033321063.webp"},
			{ID: 7, RewardType: 1, RewardAmount: 377, WashCode: 2, Weight: 20, Icon: "3004/banner/033339561-3077-file_20251125033339559.webp"},
			{ID: 8, RewardType: 1, RewardAmount: 577, WashCode: 2, Weight: 20, Icon: "3004/banner/033348037-3078-file_20251125033348035.webp"},
			{ID: 9, RewardType: 1, RewardAmount: 777, WashCode: 2, Weight: 10, Icon: "3004/banner/033354191-3079-file_20251125033354189.webp"},
			{ID: 10, RewardType: 3, RewardAmount: 1, WashCode: 1, Weight: 20, Icon: "3004/banner/051725376-3137-file_20251125051725374.webp"},
			{ID: 11, RewardType: 4, RewardAmount: 1, WashCode: 1, Weight: 20, Icon: "3004/banner/051731684-3138-file_20251125051731680.webp"},
		},
	
		random, // random
		language,                // language
		"", // signature
		timestamp,   // timestamp
	}
	// 钻石盘的增加
	payloadListType3 = []interface{}{
		3, // rechargeWheelType
		0, // specialWheelUnlockCond
	
		// taskConfig []TaskConfig
		[]TaskConfig{
			{ID: 2, RechargeType: 1, RechargeAmount: 800, SpinCount: 1},
			{ID: 3, RechargeType: 1, RechargeAmount: 1500, SpinCount: 1},
			{ID: 4, RechargeType: 2, RechargeAmount: 2000, SpinCount: 1},
		},
	
		// rewardConfig []RewardConfig
		[]RewardConfig{
			{ID: 5, RewardType: 1, RewardAmount: 177, WashCode: 2, Weight: 30, Icon: "3004/banner/033925509-3082-file_20251125033925499.webp"},
			{ID: 6, RewardType: 1, RewardAmount: 277, WashCode: 2, Weight: 30, Icon: "3004/banner/033930789-3083-file_20251125033930788.webp"},
			{ID: 7, RewardType: 1, RewardAmount: 377, WashCode: 2, Weight: 30, Icon: "3004/banner/033936984-3084-file_20251125033936982.webp"},
			{ID: 8, RewardType: 4, RewardAmount: 1, WashCode: 1, Weight: 20, Icon: "3004/banner/051743558-3139-file_20251125051743552.webp"},
			{ID: 9, RewardType: 1, RewardAmount: 577, WashCode: 2, Weight: 20, Icon: "3004/banner/033951311-3086-file_20251125033951310.webp"},
			{ID: 10, RewardType: 1, RewardAmount: 777, WashCode: 2, Weight: 20, Icon: "3004/banner/033958361-3087-file_20251125033958359.webp"},
			{ID: 11, RewardType: 1, RewardAmount: 1777, WashCode: 2, Weight: 20, Icon: "3004/banner/034005724-3088-file_20251125034005723.webp"},
			{ID: 12, RewardType: 4, RewardAmount: 2, WashCode: 1, Weight: 10, Icon: "3004/banner/051749781-3140-file_20251125051749780.webp"},
		},
	
		random, // random
		language,                // language
		"", // signature
		timestamp,   // timestamp
	}
	payloadListType4 = []interface{}{
		4, // rechargeWheelType
		3777, // specialWheelUnlockCond（注意：这里不是0，是3777）
	
		// taskConfig []TaskConfig
		[]TaskConfig{
			{ID: 2, RechargeType: 1, RechargeAmount: 800, SpinCount: 1},
			{ID: 3, RechargeType: 1, RechargeAmount: 15000, SpinCount: 1},
			{ID: 4, RechargeType: 1, RechargeAmount: 2000, SpinCount: 1},
			{ID: 5, RechargeType: 2, RechargeAmount: 4000, SpinCount: 1},
		},
	
		// rewardConfig []RewardConfig
		[]RewardConfig{
			{ID: 6, RewardType: 1, RewardAmount: 277, WashCode: 2, Weight: 30, Icon: "3004/banner/034603404-3097-file_20251125034603403.webp"},
			{ID: 7, RewardType: 3, RewardAmount: 2, WashCode: 1, Weight: 30, Icon: "3004/banner/051801006-3141-file_20251125051801005.webp"},
			{ID: 8, RewardType: 1, RewardAmount: 577, WashCode: 2, Weight: 30, Icon: "3004/banner/034550007-3095-file_20251125034550005.webp"},
			{ID: 9, RewardType: 3, RewardAmount: 4, WashCode: 1, Weight: 30, Icon: "3004/banner/051808665-3142-file_20251125051808664.webp"},
			{ID: 10, RewardType: 1, RewardAmount: 1777, WashCode: 2, Weight: 20, Icon: "3004/banner/034532090-3093-file_20251125034532089.webp"},
			{ID: 11, RewardType: 4, RewardAmount: 2, WashCode: 1, Weight: 20, Icon: "3004/banner/051815323-3143-file_20251125051815322.webp"},
			{ID: 12, RewardType: 1, RewardAmount: 2777, WashCode: 2, Weight: 20, Icon: "3004/banner/034519853-3091-file_20251125034519852.webp"},
			{ID: 13, RewardType: 5, RewardAmount: 1, WashCode: 1, Weight: 10, Icon: "3004/banner/052205348-3146-file_20251125052205346.webp"},
		},
	
		random, // random
		language,                // language
		"", // signature
		timestamp,   // timestamp
	}
)

func AddWheel(ctx *context.Context,payloadList []interface{})(*model.Response,error) {
	api := "/api/RechargeWheel/Update"
	payloadStruct := &RechargeWheelConfig{}
	if respBody,_,err := requstmodle.AdminRodAutRequest(ctx,api,payloadStruct,payloadList,request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/RechargeWheel/Update请求失败",err)),err
	} else {
		if string(respBody) == ""{
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/RechargeWheel/Update的响应为空",err)),err
		}
		if resp,err := model.ParseResponse(respBody); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/RechargeWheel/Updatt解析为空",err)),err
		} else {
			return resp,nil
		}
	}
}

// 添加充值转盘
// 一键为后台配置4个充值转盘的配置
func RunAddWheel() {
	if ctx,err := login.RunAdminSitLogin();err != nil {
		logger.Logger.Warn("添加充值转盘的后台登录失败",err)
		return
	}else {
		list := [][]interface{}{payloadListType1, payloadListType2, payloadListType3, payloadListType4}
		for index, payload := range list {
			if resp, err := AddWheel(ctx, payload); err != nil {
				logger.Logger.Warn("添加充值转盘失败",index, err)
			} else {
				logger.Logger.Info("添加充值转盘成功",index, resp)
			}
		}

	}
}