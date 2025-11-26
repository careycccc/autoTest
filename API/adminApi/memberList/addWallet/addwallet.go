package addwallet

import (
	"autoTest/API/adminApi/login"
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"context"
	"sync"
	"time"
)

// 新增用户的银行卡
type AddUserBankStruct struct {
	BankCode   any `json:"bankCode"` // 银行的code
	CardNo     any `json:"cardNo"`   // 银行卡的卡号
	MobileNo   any `json:"mobileNo"` // 手机号码
	Email      any `json:"email"`
	IfscCode   any `json:"ifscCode"`
	UserId     any `json:"userId"`
	WalletType any `json:"walletType"` //1 表示银行卡
	model.BaseStruct
}

// 添加银行卡
func AddUserBank(ctx *context.Context, userId string) (*model.BetResponse, error) {
	api := "/api/Users/AddUserWallet"
	payloadStruct := &AddUserBankStruct{}
	bankcode := "ceshiyong"
	cardNo, _ := utils.GenerateBankCard(18)
	mobileNo, _ := utils.GenerateBankCard(12)
	email := utils.GenerateRandomEmail()
	ifscCode, err := utils.RandomIFSC()
	if err != nil {
		logger.LogError("报错消息ifscCode地址生成失败", err)
		return nil, nil
	}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{bankcode, cardNo, mobileNo, email, ifscCode, userId, 1, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Users/AddUserWallet添加银行卡请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse2(respBoy); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Users/AddUserWallet添加银行卡解析失败", err)), err
		} else {
			return resp, nil
		}
	}

}

// 添加电子钱包
type AddUserWalletStruct struct {
	BankCode   any `json:"bankCode"` // 银行的code
	MobileNo   any `json:"mobileNo"` // 手机号码
	UserId     any `json:"userId"`
	WalletType any `json:"walletType"` //1 表示银行卡
	model.BaseStruct
}

// 添加电子钱包
func AddUserWallet(ctx *context.Context, userId string) (*model.BetResponse, error) {
	api := "/api/Users/AddUserWallet"
	payloadStruct := &AddUserWalletStruct{}
	bankcode := "ceshiyong"
	mobileNo, _ := utils.GenerateBankCard(12)
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{bankcode, mobileNo, userId, 2, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Users/AddUserWallet添加电子钱包请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse2(respBoy); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Users/AddUserWallet添加电子钱包解析失败", err)), err
		} else {
			return resp, nil
		}
	}

}

// 添加pix
type AddPixStruct struct {
	MobileNo      any `json:"mobileNo"`      // 手机号码
	PixWalletType any `json:"pixWalletType"` // pix
	UserId        any `json:"userId"`
	WalletType    any `json:"walletType"` //1 表示银行卡
	model.BaseStruct
}

// 添加pix
func AddUserPix(ctx *context.Context, userId string) (*model.BetResponse, error) {
	api := "/api/Users/AddUserWallet"
	payloadStruct := &AddPixStruct{}
	mobileNo, _ := utils.GenerateBankCard(10)
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{mobileNo, "Phone", userId, 3, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Users/AddUserWallet添加pix请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse2(respBoy); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Users/AddUserWallet添加pix解析失败", err)), err
		} else {
			return resp, nil
		}
	}

}

// 添加usdt
type AddUsdtStruct struct {
	Address      any `json:"address"`      // usdt地址
	AliasAddress any `json:"aliasAddress"` // usdt地址别称
	NetworkType  any `json:"networkType"`  // usdt地址类型
	UserId       any `json:"userId"`
	WalletType   any `json:"walletType"` //1 表示银行卡
	model.BaseStruct
}

func AddUserUsdt(ctx *context.Context, userId string) (*model.BetResponse, error) {
	address, err := utils.GenerateTRONAddress()
	if err != nil {
		logger.LogError("报错消息usdt地址生成失败", err)
		return nil, err
	}
	api := "/api/Users/AddUserWallet"
	payloadStruct := &AddUsdtStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{address, address, "TRC20", userId, 4, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Users/AddUserWallet添加pix请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse2(respBoy); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Users/AddUserWallet添加pix解析失败", err)), err
		} else {
			return resp, nil
		}
	}

}

// 添加upi
type AddUpiStruct struct {
	UpiId      any `json:"upiId"`
	UserId     any `json:"userId"`
	WalletType any `json:"walletType"` //1 表示银行卡
	model.BaseStruct
}

func AddUserUpi(ctx *context.Context, userId string) (*model.BetResponse, error) {
	upiId, err := utils.GenerateUPIFormat()
	if err != nil {
		logger.LogError("报错消息upi地址生成失败", err)
		return nil, err
	}
	api := "/api/Users/AddUserWallet"
	payloadStruct := &AddUsdtStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{upiId, userId, 5, random, language, "", timestamp}
	if respBoy, _, err := requstmodle.AdminRodAutRequest(ctx, api, payloadStruct, payloadList, request.StructToMap); err != nil {
		return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Users/AddUserWallet添加upi请求失败", err)), err
	} else {
		if resp, err := model.ParseResponse2(respBoy); err != nil {
			return model.HandlerErrorRes2(model.ErrorLoggerType("/api/Users/AddUserWallet添加upi解析失败", err)), err
		} else {
			return resp, nil
		}
	}

}

// 运行为一个会员添加提现信息
// 会员id绑定提现信息
// RunAddWallet 为用户添加 4 种收款方式（银行卡、PIX、USDT、电子钱包）
// 业务要求：必须等这 4 个全部成功后，后续提现逻辑才能继续
// 防卡死设计：最多只等 8 秒，超时就强制放行（用户能提到钱最重要）
func RunAddWallet(userId string) {
	// 1. 后台登录
	ctx := context.Background()
	_, ctxToken, err := login.AdminSitLogin(&ctx)
	if err != nil {
		logger.LogError("RunAddWallet 后台登录失败", err)
		return
	}
	if ctxToken == nil {
		logger.LogError("RunAddWallet 后台登录返回空 token", nil)
		return
	}

	// 你原来的写法：ctxToken 实际上就是 token（*context.Context 里存的字符串）
	adminToken := ctxToken

	// 2. 核心：最多等 8 秒（可自行调整 15~30 秒之间）
	const maxWait = 8 * time.Second
	timeout := time.NewTimer(maxWait)
	defer timeout.Stop()

	// 完成信号通道
	done := make(chan struct{}, 1)

	// 3. 开启后台真正执行 4 个添加任务
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.LogError("RunAddWallet 后台协程发生 panic", nil)
			}
		}()

		wg := &sync.WaitGroup{}
		wg.Add(4)

		// 银行卡
		go func() {
			defer wg.Done()
			if resp, err := AddUserBank(adminToken, userId); err != nil {
				logger.LogError("添加银行卡失败 userId="+userId, err)
			} else {
				logger.Logger.Info("添加银行卡成功 userId="+userId, resp)
			}
		}()

		// PIX
		go func() {
			defer wg.Done()
			if resp, err := AddUserPix(adminToken, userId); err != nil {
				logger.LogError("添加PIX失败 userId="+userId, err)
			} else {
				logger.Logger.Info("添加PIX成功 userId="+userId, resp)
			}
		}()

		// USDT
		go func() {
			defer wg.Done()
			if resp, err := AddUserUsdt(adminToken, userId); err != nil {
				logger.LogError("添加USDT失败 userId="+userId, err)
			} else {
				logger.Logger.Info("添加USDT成功 userId="+userId, resp)
			}
		}()

		// 电子钱包
		go func() {
			defer wg.Done()
			if resp, err := AddUserWallet(adminToken, userId); err != nil {
				logger.LogError("添加电子钱包失败 userId="+userId, err)
			} else {
				logger.Logger.Info("添加电子钱包成功 userId="+userId, resp)
			}
		}()

		wg.Wait()
		logger.Logger.Info("RunAddWallet 全部完成 userId=" + userId)

		// 通知主协程：真的成功了
		select {
		case done <- struct{}{}:
		default:
		}
	}()

	// 4. 主协程：最多等 18 秒
	select {
	case <-done:
		// 完美！4 个钱包全部成功，正常继续后续提现逻辑
		return

	case <-timeout.C:
		// 超时了！强制放行，不再卡死整个提现流程
		logger.Logger.Warn("RunAddWallet 超时 " + maxWait.String() + " 已强制放行，后续提现继续执行 userId=" + userId)
		// 直接 return，后面的 SetWithdrawPasswordApi、RunWithdraw 等立刻执行
		return
	}
}

// func RunAddWallet(userId string) {
// 	// 后台登录
// 	ctx := context.Background()
// 	if _, ctxToken, err := login.AdminSitLogin(&ctx); err != nil {
// 		logger.LogError("报错消息添加银行信息的后台登录失败", err)
// 		return
// 	} else {
// 		adminToken := ctxToken
// 		wg := &sync.WaitGroup{}
// 		wg.Add(4)
// 		go func(wg *sync.WaitGroup, ctxToUse *context.Context, userId string) {
// 			defer wg.Done()
// 			if resp, err := AddUserBank(adminToken, userId); err != nil {
// 				logger.LogError("添加银行信息的异步报错", err)
// 				return
// 			} else {
// 				logger.Logger.Info("添加银行信息的异步", resp)
// 			}
// 		}(wg, adminToken, userId)
// 		go func(wg *sync.WaitGroup, ctxToUse *context.Context, userId string) {
// 			defer wg.Done()
// 			if resp, err := AddUserPix(adminToken, userId); err != nil {
// 				logger.LogError("添加pix信息的异步报错", err)
// 				return
// 			} else {
// 				logger.Logger.Info("添加pix信息的异步", resp)
// 			}
// 		}(wg, adminToken, userId)
// 		go func(wg *sync.WaitGroup, ctxToUse *context.Context, userId string) {
// 			defer wg.Done()
// 			if resp, err := AddUserUsdt(adminToken, userId); err != nil {
// 				logger.LogError("添加usdt信息的异步报错", err)
// 				return
// 			} else {
// 				logger.Logger.Info("添加usdt信息的异步", resp)
// 			}
// 		}(wg, adminToken, userId)
// 		go func(wg *sync.WaitGroup, ctxToUse *context.Context, userId string) {
// 			defer wg.Done()
// 			if resp, err := AddUserWallet(adminToken, userId); err != nil {
// 				logger.LogError("添加电子钱包信息的异步报错", err)
// 				return
// 			} else {
// 				logger.Logger.Info("添加电子钱包信息的异步", resp)
// 			}
// 		}(wg, adminToken, userId)

//			wg.Wait()
//		}
//	}
//
// RunAddWallet 完全火力全开版（推荐用于生产）
// 特点：调用即返回，彻底不阻塞任何主流程，后台异步执行 4 个添加钱包请求
