package chickenroadgame

import (
	threegameapi "autoTest/API/betApi/threeGameApi"
	desklogin "autoTest/API/deskApi/loginApi"
	"context"
	"fmt"
)

// 鸡路小游戏

// 获取游戏列表
// 传入上下文，返回游戏列表的token
func GetChickenRoadGameListToken(ctx *context.Context) (string, error) {
	// 获取鸡路小游戏的游戏列表
	resp, _, err := threegameapi.GetGameUrlCommon(ctx, "chicken-road-two", "INOUT", "https://www.mggametransit.com/game/allGames?vendorCode=INOUT&name=INOUT", 104589)
	if err != nil {
		fmt.Println("获取鸡路小游戏的游戏列表失败", err)
		return "", err
	}
	fmt.Println("获取鸡路小游戏的游戏列表成功:", resp)
	return "", nil
}

func RunChickenRoadGame() {
	loginCtx, err := desklogin.ReturnContextLoginY1("911117313425", "qwer1234")
	if err != nil {
		fmt.Println("登录失败，无法进行鸡路小游戏的测试", err)
		return
	}
	GetChickenRoadGameListToken(loginCtx)
}
