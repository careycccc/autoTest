package desklogintest

import (
	lotterygameapi "autoTest/API/betApi/LotteryGameApi"
	login "autoTest/API/deskApi/loginApi"
	"autoTest/testcase/common"
	"context"
	"log"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

// 测试投注，
func TestBetWinGo(t *testing.T) {
	resultsDir := common.SetupAllureResultsDir()
	runner.Run(t, "登录 /api/Home/Login 测试", func(t provider.T) {
		t.Title("测试登录 /api/Home/Login、充值和玩游戏")
		t.Description("验证登录 /api/Home/Login 返回正确的响应和token,查询金额接口正常处理,下注接口返回正确结果")
		t.Labels(common.CommonLabels("/api/Home/Login", "critical")...)
		ctx := context.Background()

		// Declare variables in the outer scope
		var balance float64
		var rechargeResp string // Replace with actual type
		var gameCode string
		var betContent string
		var amount int
		var betMultiple int
		userName := "919111997678"
		// Step 1: Test login
		resp, tokenCtx, err := login.LoginY1(ctx, userName, "qwer1234")
		if err != nil {
			t.Errorf("LoginY1 failed: %v", err)
			t.WithNewStep("发送登录请求", func(s provider.StepCtx) {
				s.Assert().False(true, "login.LoginY1: %v", err.Error())
			})
			t.Fail()
			return
		}
		t.WithNewStep("发送登录请求", func(s provider.StepCtx) {
			common.VerifyLoginResponse(s, resp, "TestLogin")
		})

		// Step 2: Query balance
		t.WithNewStep("查询金额Balance", func(s provider.StepCtx) {
			gameCode, betContent, amount, betMultiple = lotterygameapi.GetBetResult()
			rechargeResp, balance, err = lotterygameapi.GetBalanceInfoFunc(tokenCtx, gameCode)
			if err != nil {
				s.Assert().False(true, "lotterygameapi.GetBalanceInfoFunc 调用失败: %v", err.Error())
				s.Fail()
				return
			}
			s.Assert().NotNil(rechargeResp, "rechargeResp is nil")
			//s.Assert().NotNil(tokenCtx, "tokenCtx is nil")
			s.Assert().True(balance >= 0, "金额应大于或等于0")
		})

		// Step 3: Test betting
		t.WithNewStep("发起投注", func(s provider.StepCtx) {
			if balance == 0 {
				s.Assert().False(true, "balance的金额为0")
				s.Fail()
				return
			}

			// Check if betting is allowed
			isBet, issNumber := lotterygameapi.IsBet("", gameCode, "WinGo")
			if !isBet {
				log.Printf("当期期号,不可以投注")
				s.Assert().False(true, "当期期号,不可以投注")
				s.Fail()
				return
			}

			// Add nil checks for critical variables
			// if tokenCtx == nil {
			// 	s.Assert().False(true, "tokenCtx is nil")
			// 	s.Fail()
			// 	return
			// }
			// if rechargeResp == "" {
			// 	s.Assert().False(true, "rechargeResp is nil")
			// 	s.Fail()
			// 	return
			// }
			// if gameCode == "" {
			// 	s.Assert().False(true, "gameCode is nil")
			// 	s.Fail()
			// 	return
			// }
			// if amount == 0 {
			// 	s.Assert().False(true, "amount is nil")
			// 	s.Fail()
			// 	return
			// }
			// if betMultiple == 0 {
			// 	s.Assert().False(true, "betMultiple is nil")
			// 	s.Fail()
			// 	return
			// }
			// if betContent == "" {
			// 	s.Assert().False(true, "betContent is nil")
			// 	s.Fail()
			// 	return
			// }
			// if issNumber == "" {
			// 	s.Assert().False(true, "issNumber is nil")
			// 	s.Fail()
			// 	return
			// }
			// if userName == "" {
			// 	s.Assert().False(true, "userName is nil")
			// 	s.Fail()
			// 	return
			// }

			// Perform betting
			resp, err := lotterygameapi.BetWingo(gameCode, amount, betMultiple, betContent, issNumber, rechargeResp, userName)
			if err != nil {
				s.Assert().False(true, "lotterygameapi.BetWingo 调用失败: %v", err.Error())
				s.Fail()
				return
			}
			s.Assert().NotNil(resp, "BetWingo response is nil")
			s.Assert().NotNil(amount*betMultiple, "投注金额")
			s.Assert().NotNil(gameCode, "游戏名称")
			common.VerifyLoginResponse2(s, resp, "TestBetWinGo")
		})
	})
	common.FlushLogs()
	common.CheckAllureResultsDir(resultsDir)
}
