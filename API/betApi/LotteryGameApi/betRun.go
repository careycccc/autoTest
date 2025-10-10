package lotterygameapi

import (
	login "autoTest/API/deskApi/loginApi"
	"autoTest/store/logger"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 负责组装投注
func BetRun(userName string) error {
	//userName := "919111997678"
	// 进行前台登录
	ctx := context.Background()
	if _, tokenCtx, err := login.LoginY1(ctx, userName, "qwer1234"); err != nil {
		fmt.Println(err)
		return err
	} else {
		// 进行投注
		// 先把投注的结果随机出来
		gameCode, betContent, amount, betMultiple := GetBetResult()
		if err := RunBetFunc(tokenCtx, gameCode, betContent, userName, amount, betMultiple); err != nil {
			return err
		} else {
			return nil
		}
	}
}

/*
gameCode, 投注的彩票类型"WinGo_5M", "TrxWinGo_10M"
betContent, 投注彩票的颜色或者大小
amount, 投注彩票的金额
betMultiple 投注彩票的倍率
*
*/
func GetBetResult() (gameCode, betContent string, amount, betMultiple int) {
	num := RandomInt(2)
	gameCodeList := []string{"WinGo_5M", "TrxWinGo_10M"}
	gameCode = gameCodeList[num]
	num1 := RandomInt(5)
	betContentList := []string{"Color_Green", "Color_Violet", "Color_Red", "BigSmall_Big", "BigSmall_Small"}
	betContent = betContentList[num1]
	num2 := RandomInt(4)
	amountList := []int{10, 20, 50, 100}
	amount = amountList[num2]
	num3 := RandomInt(4)
	betMultipleList := []int{10, 20, 50, 100}
	betMultiple = betMultipleList[num3]
	return
}

// RandomNumber 为每个 goroutine 生成 [0, n) 范围内的随机整数
func RandomInt(n int) int {
	if n <= 0 {
		return 0 // 处理无效输入
	}
	// 使用 sync.Once 确保每个 goroutine 初始化一次随机源
	var src *rand.Rand
	var once sync.Once
	once.Do(func() {
		src = rand.New(rand.NewSource(time.Now().UnixNano())) // 为 goroutine 创建随机源
	})
	return src.Intn(n) // 生成 [0, n) 范围的随机数
}

/*
投注的函数
gameCode, 投注的彩票类型"WinGo_5M", "TrxWinGo_10M"
betContent, 投注彩票的颜色或者大小
amount, 投注彩票的金额
betMultiple 投注彩票的倍率
**/

func RunBetFunc(ctx *context.Context, gameCode, betContent, userName string, amount, betMultiple int) error {
	BalanceToken, balance, err := GetBalanceInfoFunc(ctx, gameCode)
	if err != nil {
		return err
	}
	if balance == 0.0 {
		//fmt.Println("------------------------余额为0,不可以投注------------------------")
		logger.Logger.Warn("余额为0,不可以投注")
		return fmt.Errorf("------------------------余额为0,不可以投注------------------------")
	} else {
		// 是否可以投注
		isBet, issNumber := IsBet("", gameCode, "WinGo")
		if isBet {
			// 可以投注
			resp, err := BetWingo(gameCode, amount, betMultiple, betContent, issNumber, BalanceToken, userName)
			if err != nil {
				logger.Logger.Error("BetWingo报错信息", err)
			} else {
				logger.Logger.Info("投注成功", resp)
			}
		} else {
			//fmt.Println("不可以投注")
			logger.Logger.Warn("当期期号,不可以投注")
			return fmt.Errorf("不可以投注")
		}
		return nil
	}
}
