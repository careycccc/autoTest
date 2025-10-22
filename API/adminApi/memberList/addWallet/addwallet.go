package addwallet

import (
	"autoTest/API/adminApi/login"
	requstmodle "autoTest/requstModle"
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"autoTest/store/utils"
	"context"
	"fmt"
	"sync"
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
	fmt.Println("upi-----", upiId)
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
func RunAddWallet() {
	userId := "2441424"
	// 后台登录
	ctx := context.Background()
	if _, ctxToken, err := login.AdminSitLogin(&ctx); err != nil {
		logger.LogError("报错消息添加银行信息的后台登录失败", err)
		return
	} else {
		adminToken := ctxToken
		wg := &sync.WaitGroup{}
		wg.Add(4)
		go func(wg *sync.WaitGroup, ctxToUse *context.Context, userId string) {
			defer wg.Done()
			if resp, err := AddUserBank(adminToken, userId); err != nil {
				logger.LogError("添加银行信息的异步报错", err)
				return
			} else {
				logger.Logger.Info("添加银行信息的异步", resp)
			}
		}(wg, adminToken, userId)
		go func(wg *sync.WaitGroup, ctxToUse *context.Context, userId string) {
			defer wg.Done()
			if resp, err := AddUserPix(adminToken, userId); err != nil {
				logger.LogError("添加pix信息的异步报错", err)
				return
			} else {
				logger.Logger.Info("添加pix信息的异步", resp)
			}
		}(wg, adminToken, userId)
		go func(wg *sync.WaitGroup, ctxToUse *context.Context, userId string) {
			defer wg.Done()
			if resp, err := AddUserUsdt(adminToken, userId); err != nil {
				logger.LogError("添加usdt信息的异步报错", err)
				return
			} else {
				logger.Logger.Info("添加usdt信息的异步", resp)
			}
		}(wg, adminToken, userId)
		go func(wg *sync.WaitGroup, ctxToUse *context.Context, userId string) {
			defer wg.Done()
			if resp, err := AddUserWallet(adminToken, userId); err != nil {
				logger.LogError("添加电子钱包信息的异步报错", err)
				return
			} else {
				logger.Logger.Info("添加电子钱包信息的异步", resp)
			}
		}(wg, adminToken, userId)

		wg.Wait()
	}
}
