package rechargewheel

import (
	financialmanagement "autoTest/API/adminApi/financialManagement"
	"autoTest/API/adminApi/login"
	getgiftinfo "autoTest/API/deskApi/activeGIft/GetGiftInfo"
	getuserinfo "autoTest/API/deskApi/getUserinfo"
	registerapi "autoTest/API/deskApi/registerApi"
	"autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"context"
	"encoding/json"
	"sync"
	"time"
)

// 充值转盘的相关的逻辑

type UserRechargeWheelInfo struct {
	isOpenRechargeWheel          bool    // 是否开启了充值转盘 是否显示充值转盘
	RechargeWheelRemainSpinCount float64 // 充值轮盘剩余旋转次数
}

// 获取当前用户充值转盘的，开启信息，剩余旋转次数
func GetUserRechargeWheelInfo(ctx *context.Context) (UserRechargeWheelInfo, error) {
	if _, rechargeWheelInfo, err := getgiftinfo.GetGiftInfoApi(ctx); err != nil {
		return UserRechargeWheelInfo{}, err
	} else {
		info := UserRechargeWheelInfo{
			isOpenRechargeWheel:          rechargeWheelInfo.IsOpenRechargeWheel.(bool),
			RechargeWheelRemainSpinCount: rechargeWheelInfo.RechargeWheelRemainSpinCount.(float64),
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

// 返回充值转盘是否开启，剩余旋转次数,充值金额，上下文，用户id
type getRechargeWheelInfo struct {
	isShow      bool             // 是否开启充值转盘
	wheelNumber int              // 剩余旋转次数
	amount      float64          // 充值金额
	ctx         *context.Context // 上下文
	userId      int              // 用户id
}

// 执行充值转盘的逻辑
func execRechargeWheel(value1 int8) {
	rechargeWheelInfo, _ := CallRechargeWheelCondition(value1)
	if rechargeWheelInfo.isShow {
		logger.Logger.Info("充值转盘是否开启", rechargeWheelInfo.isShow)
		// 在进行充值
		if _, err := financialmanagement.ArtificialRechargeFunc(rechargeWheelInfo.ctx, rechargeWheelInfo.userId, rechargeWheelInfo.amount, 2); err != nil {
			logger.LogError("充值转盘设置充值金额失败", err)
			return
		}
		time.Sleep(1 * time.Second) // 等待1秒，确保后台处理完成
		// 旋转充值转盘
		if resp, err := SpinRechargeWheelApi(rechargeWheelInfo.ctx, 1); err != nil {
			logger.LogError("旋转充值转盘失败", err)
			return
		} else {
			if resp.Data != nil {
				logger.Logger.Info("旋转充值转盘成功", "resp.Data", resp.Data)
			}
		}
	}
}

/*
随机一个账号，设置充值转盘的条件
value1  0表示不需要充值   1，首充 2，二充  3，三充
返回充值转盘是否开启，剩余旋转次数,充值金额，后台的上下文
*
*/
func CallRechargeWheelCondition(value1 int8) (*getRechargeWheelInfo, *context.Context) {
	// 第一步后台登录
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		logger.LogError("充值转盘的后台登录失败", err)
		return &getRechargeWheelInfo{isShow: false, wheelNumber: 0, amount: 0.0}, ctx
	} else {
		// 复制 ctx 避免作用域问题
		adminCtx := ctx
		// 获取需要充值的金额
		amount := ReturnRechargeAmount(ctx)

		// 随机生成一个用户
		userCount, err := utils.RandmoUserCount() // 注意：应该是 RandUserCount？
		if err != nil {
			logger.LogError("充值转盘随机生成用户失败", err)
			return &getRechargeWheelInfo{isShow: false, wheelNumber: 0, amount: 0.0}, ctx
		}
		logger.Logger.Info("充值转盘随机生成的用户", userCount)

		// 获取进行注册登录总代的方式
		_, ctxToken, err := registerapi.GeneralAgentRegister(userCount)
		if err != nil {
			logger.LogError("充值转盘随机生成用户注册登录失败", err)
			return &getRechargeWheelInfo{isShow: false, wheelNumber: 0, amount: 0.0}, ctx
		}

		_, userInfo, err := getuserinfo.GetUserInfo(ctxToken)
		if err != nil {
			logger.LogError("充值转盘随机生成用户注册登录后获取用户信息失败", err)
			return &getRechargeWheelInfo{isShow: false, wheelNumber: 0, amount: 0.0}, ctx
		}
		//logger.Logger.Info("充值转盘随机生成用户注册登录后获取用户信息", resp, userInfo)
		userId := int(userInfo.UserID)
		// 并发处理
		wg := &sync.WaitGroup{}
		wg.Add(2)
		// 设置充值转盘的条件
		go func(wg *sync.WaitGroup, ctxToUse *context.Context, val int8) {
			defer wg.Done()
			if _, err := SetRechargeWheelCondition(ctxToUse, val); err != nil {
				logger.LogError("设置充值转盘条件失败", err)
				return
			}
			//logger.Logger.Info("充值转盘条件设置成功")
		}(wg, adminCtx, value1)

		// 设置充值金额
		go func(wg *sync.WaitGroup, ctxToUse *context.Context) {
			defer wg.Done()
			if _, err := financialmanagement.ArtificialRechargeFunc(ctxToUse, userId, amount, 2); err != nil {
				logger.LogError("call充值转盘设置充值金额失败", err)
				return
			}
		}(wg, adminCtx)

		wg.Wait()
		time.Sleep(1 * time.Second) // 等待1秒，确保后台处理完成
		// 最后获取用户的充值转盘信息
		if info, err := GetUserRechargeWheelInfo(ctxToken); err != nil {
			logger.LogError("获取用户充值转盘信息失败", err)
			return &getRechargeWheelInfo{isShow: false, wheelNumber: 0, amount: 0.0}, ctx
		} else {
			return &getRechargeWheelInfo{isShow: info.isOpenRechargeWheel, wheelNumber: int(info.RechargeWheelRemainSpinCount), amount: amount, ctx: ctxToken, userId: userId}, ctx
		}
	}
}

// 运行充值转盘的任务
// 0 表示不需要充值   1，首充 2，二充  3，三充
func RunRechargeWheelCondition(value1 int8) {
	switch value1 {
	case 0:
		execRechargeWheel(value1)
	case 1:
		// 首充
		rechargeWheelInfo, ctx := CallRechargeWheelCondition(value1)
		// 判断充值转盘是否开启
		if rechargeWheelInfo.isShow {
			logger.Logger.Info("充值转盘是否开启", rechargeWheelInfo.isShow)
			// 判断是否有剩余的旋转次数
			if rechargeWheelInfo.wheelNumber > 0 {
				logger.LogError("需要首充判断，只进行了首充却有了旋转次数", nil)
				return
			} else {
				time.Sleep(1 * time.Second) // 等待1秒，确保后台处理完成
				// 在进行充值
				userId := int(rechargeWheelInfo.userId)
				if resp, err := financialmanagement.ArtificialRechargeFunc(ctx, userId, rechargeWheelInfo.amount, 2); err != nil {
					logger.LogError("首充充值转盘设置充值金额失败", err)
					return
				} else {
					logger.Logger.Info("首充充值转盘充值成功", resp)
				}
				time.Sleep(1 * time.Second) // 等待1秒，确保后台处理完成
				// 旋转充值转盘
				if resp, err := SpinRechargeWheelApi(rechargeWheelInfo.ctx, 1); err != nil {
					logger.LogError("首充旋转充值转盘失败", err)
					return
				} else {
					logger.Logger.Info("首充旋转充值转盘成功", resp)
				}
			}
		}
	case 2:
		// 二充
		// 第一次充值
		rechargeWheelInfo, ctx := CallRechargeWheelCondition(value1)
		if rechargeWheelInfo.isShow {
			logger.LogError("二充，第一次充值就开启了充值转盘", nil)
			return
		}
		// 第二次充值
		if _, err := financialmanagement.ArtificialRechargeFunc(ctx, rechargeWheelInfo.userId, rechargeWheelInfo.amount, 2); err != nil {
			logger.LogError("充值转盘设置充值金额失败", err)
			return
		}
		// 查看是否有旋转次数
		time.Sleep(1 * time.Second) // 等待1秒，确保后台处理完成
		if info, err := GetUserRechargeWheelInfo(rechargeWheelInfo.ctx); err != nil {
			logger.LogError("获取用户充值转盘信息失败", err)
			return
		} else {
			if info.isOpenRechargeWheel {
				logger.Logger.Info("二充第二次充值转盘是否开启", info.isOpenRechargeWheel)
				// 开启就是正常的
				//再次充值
				if _, err := financialmanagement.ArtificialRechargeFunc(ctx, rechargeWheelInfo.userId, rechargeWheelInfo.amount, 2); err != nil {
					logger.LogError("二充充值转盘设置充值金额失败", err)
					return
				}
				time.Sleep(1 * time.Second) // 等待1秒，确保后台处理完成
				// 旋转充值转盘
				if resp, err := SpinRechargeWheelApi(rechargeWheelInfo.ctx, 1); err != nil {
					logger.LogError("二充第二次充值旋转充值转盘失败", err)
					return
				} else {
					if resp.Data != nil {
						logger.Logger.Info("二充第二次充值旋转充值转盘成功", "resp.Data", resp.Data)
					}
				}
			} else {
				logger.LogError("二充第二次充值转盘没有开启", nil)
				return
			}
		}
	case 3:
		// 三充
		// 第一次充值
		rechargeWheelInfo, ctx := CallRechargeWheelCondition(value1)
		if rechargeWheelInfo.isShow {
			logger.LogError("三充，第一次充值就开启了充值转盘", nil)
			return
		}
		// 第二次充值
		if _, err := financialmanagement.ArtificialRechargeFunc(ctx, rechargeWheelInfo.userId, rechargeWheelInfo.amount, 2); err != nil {
			logger.LogError("充值转盘设置充值金额失败", err)
			return
		}

		// 查看前台是否开启了充值转盘
		time.Sleep(1 * time.Second) // 等待1秒，确保后台处理完成
		if info, err := GetUserRechargeWheelInfo(rechargeWheelInfo.ctx); err != nil {
			logger.LogError("获取用户充值转盘信息失败", err)
			return
		} else {
			if info.isOpenRechargeWheel {
				logger.LogError("三充，第二次充值,就开启了充值转盘", nil)
				return
			} else {
				// 第三次充值
				if _, err := financialmanagement.ArtificialRechargeFunc(ctx, rechargeWheelInfo.userId, rechargeWheelInfo.amount, 2); err != nil {
					logger.LogError("三充充值转盘设置充值金额失败", err)
					return
				}
				time.Sleep(1 * time.Second) // 等待1秒，确保后台处理完成
				// 查看前台是否开启了充值转盘
				if info, err := GetUserRechargeWheelInfo(rechargeWheelInfo.ctx); err != nil {
					logger.LogError("三充获取用户充值转盘信息失败", err)
					return
				} else {
					if info.isOpenRechargeWheel {
						logger.Logger.Info("三充第三次充值转盘是否开启", info.isOpenRechargeWheel)
						//再次充值
						time.Sleep(1 * time.Second) // 等待1秒，确保后台处理完成
						if _, err := financialmanagement.ArtificialRechargeFunc(ctx, rechargeWheelInfo.userId, rechargeWheelInfo.amount, 2); err != nil {
							logger.LogError("充值转盘设置充值金额失败", err)
							return
						}
						time.Sleep(1 * time.Second) // 等待1秒，确保后台处理完成
						// 旋转充值转盘
						if resp, err := SpinRechargeWheelApi(rechargeWheelInfo.ctx, 1); err != nil {
							logger.LogError("三充第三次充值旋转充值转盘失败", err)
							return
						} else {
							if resp.Data != nil {
								logger.Logger.Info("三充第三次充值旋转充值转盘成功", "resp.Data", resp.Data)
							}
						}
					} else {
						logger.LogError("三充第三次充值转盘没有开启", nil)
						return
					}
				}

			}
		}
	default:
		logger.LogError("输入的参数不正确,只能是0,1,2,3", nil)
		return
	}

}
