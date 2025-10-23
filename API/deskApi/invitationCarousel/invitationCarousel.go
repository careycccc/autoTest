package invitationcarousel

import (
	financialmanagement "autoTest/API/adminApi/financialManagement"
	adminLogin "autoTest/API/adminApi/login"
	memberlist "autoTest/API/adminApi/memberList/memberList"
	lotterygameapi "autoTest/API/betApi/LotteryGameApi"
	login "autoTest/API/deskApi/loginApi"
	registerapi "autoTest/API/deskApi/registerApi"
	utils "autoTest/API/utils"
	requstmodle "autoTest/requstModle"
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	sutils "autoTest/store/utils"
	"context"
	"encoding/json"
	"math/rand"
	"sync"
	"time"
)

// 邀请转盘
// 定义结构体来映射 JSON 数据
type ClickResponse struct {
	Data struct {
		IsFirstInvitedWheel bool `json:"isFirstInvitedWheel"`
	} `json:"data"`
}

/*
开启邀请转盘的活动点击4个礼物盒
返回响应和点击响应的结果 true表示礼物盒没有开，点击了，false表示活动是开启的
*
*/
func ClickSpinInvitedWheel(ctx *context.Context) (*model.Response, bool, error) {
	api := "/api/Activity/SpinInvitedWheel"
	payloadStruct, payloadList := utils.BaseStructHandler()
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.InitStructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/SpinInvitedWheel请求失败", err)), false, err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/SpinInvitedWheel解析失败", err)), false, err
		} else {
			var response ClickResponse
			if err := json.Unmarshal([]byte(string(respBoy)), &response); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/SpinInvitedWheel['ClickResponse']解析失败", err)), false, err
			}
			return resp, response.Data.IsFirstInvitedWheel, err
		}
	}

}

// 定义结构体来映射 JSON 数据
type ClickShareResponse struct {
	Data struct {
		InviteCode string `json:"inviteCode"`
	} `json:"data"`
}

// 点击分享链接
// 返回响应和邀请码
func ClickShareLink(ctx *context.Context) (*model.Response, string, error) {
	api := "/api/Activity/GetUserInviteLinkAddress"
	payloadStruct, payloadList := utils.BaseStructHandler()
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.InitStructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/GetUserInviteLinkAddress请求失败", err)), "", err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/GetUserInviteLinkAddress解析失败", err)), "", err
		} else {
			var response ClickShareResponse
			if err := json.Unmarshal([]byte(string(respBoy)), &response); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/GetUserInviteLinkAddress['ClickShareResponse']解析失败", err)), "", err
			}
			return resp, response.Data.InviteCode, err
		}
	}

}

// 定义结构体来映射 JSON 数据
type ClickSpinningTurntableResponse struct {
	Data struct {
		PrizeAmount float64 `json:"prizeAmount"` // 每次旋转的金额
		IsWin       bool    `json:"isWin"`
	} `json:"data"`
}

/*
点击选择转盘
返回响应和本次旋转的金额，和是否已经中奖
*
*/
func ClickSpinningTurntable(ctx *context.Context) (*model.Response, float64, bool, error) {
	api := "/api/Activity/SpinInvitedWheel"
	payloadStruct, payloadList := utils.BaseStructHandler()
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.InitStructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/GetUserInviteLinkAddress请求失败", err)), -1.1, false, err
	} else {
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/GetUserInviteLinkAddress解析失败", err)), -1.1, false, err
		} else {
			var response ClickSpinningTurntableResponse
			if err := json.Unmarshal([]byte(string(respBoy)), &response); err != nil {
				return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/GetUserInviteLinkAddress['ClickShareResponse']解析失败", err)), -1.1, false, err
			}
			return resp, response.Data.PrizeAmount, response.Data.IsWin, err
		}
	}
}

// 运行邀请转盘
func RunSpinInvitedWheel() error {
	userName := "911006976778"
	if ctx, err := login.ReturnContextLoginY1(userName, "qwer1234"); err != nil {
		return err
	} else {
		if _, result, err := ClickSpinInvitedWheel(ctx); err != nil {
			return err
		} else {
			if result {
				logger.Logger.Info("邀请转盘的4个礼物盒点击成功", userName, "结果:", result)
			}

			// 点击分享链接，返回邀请码
			if _, InviteCode, err := ClickShareLink(ctx); err != nil {
				return err
			} else {
				// 使用重试机制
				// 旋转转盘
				if result, err := request.RetryWrapper(ClickSpinningTurntable, ctx); err != nil {
					return err
				} else {
					logger.Logger.Info("本次旋转的金额", result[1], InviteCode)
					return nil
				}
			}
		}
	}
}

// 注册的方式，
func RunSpinInvitedWheelWork() error {
	// 随机生成账号
	userList := utils.RandmoUserId(config.GeneralAgentNumber)
	// 把生成的账号进行总代的方式进行注册
	for i := range userList {
		// 把生成的下级账号写入到yaml中
		go sutils.WriteYAML(userList[i])
		//进行注册
		if _, ctxToken, err := registerapi.GeneralAgentRegister(userList[i]); err != nil {
			return err
		} else {
			logger.Logger.Info("注册成功并且成功登录", userList[i])
			ch := make(chan struct{}, 1)
			err := NewRound(ctxToken, ch)
			if err != nil {
				return err
			}
			<-ch
			// 进行下一轮的判别
			if config.WHEELNUMBER == 1 {
				return nil
			}
			newRoundCh := make(chan struct{}, 1)
			for i := 1; i <= config.WHEELNUMBER; i++ {
				err := NewRound(ctxToken, newRoundCh)
				if err != nil {
					return err
				}
				<-newRoundCh
			}
		}
	}
	return nil

}

// 新的一轮
func NewRound(ctx *context.Context, ch chan<- struct{}) error {
	// 点击4个礼物盒子
	if _, result, err := ClickSpinInvitedWheel(ctx); err != nil {
		return err
	} else {
		logger.Logger.Info("邀请转盘的活动开启", result)
		// 旋转邀请转盘	旋转免费次数
		if result, err := request.RetryWrapper(ClickSpinningTurntable, ctx); err != nil {
			return err
		} else {
			logger.Logger.Info("本次旋转的金额", result[1])
			// 分享邀请链接
			if _, InviteCode, err := ClickShareLink(ctx); err != nil {
				return err
			} else {
				logger.Logger.Info("分享链接以点击", InviteCode)
				// 邀请下一级
				if moneny, err := sutils.GenerateRandomInt(config.MIN_MONENY, config.MAX_MONENY); err != nil {
					return err
				} else {
					RunTaskWhille(InviteCode, moneny, ctx)
				}
			}
		}
	}
	ch <- struct{}{}
	return nil
}

// 定义结构体来映射 JSON 数据
type GetUserInvitedWheelInfoResponse struct {
	Data struct {
		InvitedWheelTotalPrizeAmount float64 `json:"invitedWheelTotalPrizeAmount"` // 当前用户的旋转转盘的金额
	} `json:"data"`
}

// 获取用户的邀请转盘的提款金额
func GetUserInvitedWheelInfo(ctx *context.Context) (*model.Response, float64, error) {
	api := "/api/Activity/GetUserInvitedWheelInfo"
	payloadStruct, payloadList := utils.BaseStructHandler()
	if respBoy, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.InitStructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("", err)), -1.0, err
	} else {
		var response GetUserInvitedWheelInfoResponse
		err := json.Unmarshal(respBoy, &response)
		if err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("", err)), -1.0, err
		}
		if resp, err := model.ParseResponse(respBoy); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("", err)), -1.0, err
		} else {
			// 把金额给提取出来
			return resp, response.Data.InvitedWheelTotalPrizeAmount, nil
		}
	}
}

type ClickWheelWithdraw struct {
	Amount int64 `json:"amount"` // 提现金额
	model.BaseStruct
}

// 点击提现按钮
// 传入提现的金额，token
func ClickWheelWithdrawFunc(amount float64, ctx *context.Context) (*model.Response, error) {
	api := "/api/Activity/SumitInvitedWheelWithdraw"
	payloadStruct := &ClickWheelWithdraw{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{amount, random, language, "", timestamp}
	if respBody, _, err := requstmodle.DeskTenAuthorRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/SumitInvitedWheelWithdraw请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse(respBody); err != nil {
			return model.HandlerErrorRes(model.ErrorLoggerType("/api/Activity/SumitInvitedWheelWithdraw解析失败", err)), err
		} else {
			return resp, nil
		}
	}
}

// 运行并行的任务
// 只要填写邀请码，自动邀请下级，并且充值
func RunTaskWhille(yqCode string, monenyCount float64, ctx *context.Context) {
	rand.Seed(time.Now().UnixNano()) // 初始化随机种子

	// 定义单个锁保护所有公共数据
	var lock sync.Mutex

	var wg sync.WaitGroup
	// 并发控制，使用通道限制最大并发数为
	semaphore := make(chan struct{}, config.SUB_CONCURRENT)
	// 随机生成下级数量
	number, _ := sutils.GenerateRandomInt(config.SUB_MINNUMBER, config.SUB_MAXMUMBER)
	// 执行 几 次任务
	for i := 1; i <= int(number); i++ {
		wg.Add(1)
		semaphore <- struct{}{} // 获取一个令牌，控制并发数

		go func(id int) {
			defer func() { <-semaphore }() // 释放令牌
			// 将三个公共数据及其锁作为参数传入
			TaskWhille(id, &wg, &yqCode, &monenyCount, &lock)
		}(i)
	}

	// 等待所有任务完成
	wg.Wait()
	// 充值结束后等待5s后，进行点击提现
	_, moneny, err := GetUserInvitedWheelInfo(ctx)
	if err != nil {
		logger.LogError("获取用户的转盘信息失败", err)
		return
	}
	logger.Logger.Info("当前用户的转盘总额", moneny)
	// 还差提现金额的获取
	time.Sleep(time.Second * 5)
	_, errs := ClickWheelWithdrawFunc(moneny, ctx) // 点击转盘提现
	if errs != nil {
		//fmt.Println("点击转盘提现失败", err)
		logger.LogError("点击转盘提现失败", errs)
		return
	}

}

// 并行邀请人
// 任务函数，接收三个公共数据及其对应的锁
// 任务函数，调用 RunWhille 并更新公共数据
func TaskWhille(id int, wg *sync.WaitGroup, yqCode *string, monenyCount *float64, lock *sync.Mutex) {
	defer wg.Done()
	// 随机生成账号
	userAmount, _ := utils.RandmoUserCount()
	// 使用单个锁保护对三个公共数据的联合操作
	lock.Lock()
	// 调用 RunWhille，传入当前公共数据值
	RunWhille(userAmount, *yqCode, *monenyCount)
	lock.Unlock()
}

// 邀请转盘邀请下一级
/*
userAmount  邀请下级的账号
yqCode  邀请人的邀请码
monenyCount 邀请转盘的充值金额
**/
func RunWhille(userAmount string, yqCode string, monenyCount float64) error {
	// 发送注册
	if _, ctxToken, err := registerapi.RegisterMobileLoginFunc(userAmount, yqCode); err != nil {
		return err
	} else {
		logger.Logger.Info("注册成功", userAmount)
		// 后台登录
		ctxAdminToken, err := adminLogin.RunAdminSitLogin()
		if err != nil {
			logger.LogError("后台登录失败", err)
			return err
		}
		time.Sleep(time.Second * 1)
		// 后台获取用户id
		_, userId, err := memberlist.GetUserIdApi(ctxAdminToken, userAmount)
		if err != nil {
			return err
		}
		logger.Logger.Info("userid的值", userId)
		adminToken := *ctxAdminToken
		wg := &sync.WaitGroup{}
		wg.Add(2)
		// 根据id进行充值
		go func(wg *sync.WaitGroup, ctxToUse *context.Context) {
			defer wg.Done()
			financialmanagement.ArtificialRechargeFunc(ctxToUse, int(userId), monenyCount, 2)
		}(wg, &adminToken)

		// 修改用户密码
		go func(wg *sync.WaitGroup, ctxToUse *context.Context) {
			defer wg.Done()
			memberlist.UpdatePassword(ctxToUse, userId, config.SUB_PWD)
		}(wg, &adminToken)

		wg.Wait()
		logger.Logger.Info("充值金额", monenyCount)
		logger.Logger.Info("成功修改密码", config.SUB_PWD)
		time.Sleep(time.Second * 3)
		// 充值结束后进行投注
		gameCode, betContent, amount, betMultiple := lotterygameapi.GetBetResult()
		if err := lotterygameapi.RunBetFunc(&ctxToken, gameCode, betContent, userAmount, amount, betMultiple); err != nil {
			return err
		} else {
			return nil
		}
	}

	// 初始化随机数种子,有些下级充值有些下级不充值
	// rand.Seed(time.Now().UnixNano())
	// // 生成一个[0, 1]范围内的随机数
	// randomNumber := rand.Intn(2)
	// // 检查随机数以确定是否触发50%几率的事件
	// if randomNumber == 0 {
	// 	// 后台进行登录和人工充值
	// 	adminRun(userAmount, monenyCount)
	// }
}
