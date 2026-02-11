package chickenroadgame

import (
	threegameapi "autoTest/API/betApi/threeGameApi"
	desklogin "autoTest/API/deskApi/loginApi"
	"autoTest/store/logger"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// 鸡路小游戏

// 获取游戏列表
// 传入上下文，返回游戏列表的token
func GetChickenRoadGameListToken(ctx *context.Context) (*threegameapi.ThirdGameType, error) {
	// 获取鸡路小游戏的游戏列表
	_, thirdGameType, err := threegameapi.GetGameUrlCommon(ctx, "chicken-road-two", "INOUT", "https://www.mggametransit.com/game/allGames?vendorCode=INOUT&name=INOUT", 104589)
	if err != nil {
		logger.LogError("获取鸡路小游戏的游戏列表失败", err)
		return &threegameapi.ThirdGameType{}, err
	}
	//fmt.Println("获取鸡路小游戏的游戏列表成功:", resp)
	return thirdGameType, nil
}

type ChickenRoadGameEnterStruct struct {
	Operator   string `json:"operator"`
	Auth_token string `json:"auth_token"`
	Currency   string `json:"currency"`
	Game_mode  string `json:"game_mode"`
}

// AuthInOutGames 完成 /api/auth 认证请求
// 参数 authToken: 你从抓包或上一步得到的 auth_token 字符串
// 返回: 响应体的字符串 + error
func EnterChickenRoadGame(authToken, operatorId string) (string, error) {
	fmt.Println("进入鸡路小游戏，使用的authToken:", authToken)
	url := "https://api.inout.games/api/auth"
	//OperatorID, _ := threegameapi.GenerateOperatorID()
	// 请求体（只改 auth_token 部分）
	payload := []byte(fmt.Sprintf(`{
  "operator": "%s",
  "auth_token": "%s",
  "currency": "INR",
  "game_mode": "chicken-road-two"
}`, operatorId, authToken))

	// 创建请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	// 完整还原你抓包的所有关键 Header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br, zstd")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Mobile Safari/537.36")
	req.Header.Set("Origin", "https://chicken-road-two.inout.games")
	req.Header.Set("Referer", "https://chicken-road-two.inout.games/")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?1")
	req.Header.Set("Sec-CH-UA-Platform", `"Android"`)
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")

	// 创建客户端（带超时）
	client := &http.Client{
		Timeout: 20 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 自动处理 gzip 压缩（服务端可能返回压缩）
	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gz, err := gzip.NewReader(resp.Body)
		if err != nil {
			return "", err
		}
		defer gz.Close()
		reader = gz
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

/*
// 玩鸡路小游戏，1-10步
authToken  token
operatorId  厂商id
betAmount   下注金额
loopForever 是否无限循环 ture false
*
*/
func AutoPlayChicken(authToken, operatorId string, betAmount int64, loopForever bool) {
	if betAmount == 0 {
		betAmount = 100
	}
	rand.Seed(time.Now().UnixNano())

	// 完整URL模板
	baseURL := "wss://api.inout.games/io/"
	query := fmt.Sprintf("gameMode=chicken-road-two&operatorId=%s&Authorization=%s&EIO=4&transport=websocket",
		operatorId, authToken)
	wsURL := baseURL + "?" + query

	log.Printf("启动小鸡过马路自动机器人")
	log.Printf("OperatorID : %s", operatorId)
	log.Printf("下注金额   : %d", betAmount)
	log.Printf("无限循环   : %v", loopForever)

	for {
		if !playOneRound(wsURL, betAmount) {
			log.Println("连接失败（JWT过期或网络问题），5秒后重试...")
			time.Sleep(5 * time.Second)
			continue
		}

		if !loopForever {
			log.Println("单局结束，程序退出")
			return
		}
		time.Sleep(2 * time.Second) // 局间休息
	}
}

// 一局完整流程
func playOneRound(wsURL string, betAmount int64) bool {
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Printf("WebSocket连接失败: %v", err)
		return false
	}
	defer conn.Close()
	log.Println("WebSocket 连接成功！")

	stepCount := 0
	stopAt := rand.Intn(10) + 1 // 随机1~10步停止
	log.Printf("本局目标：第 %d 步自动停止", stopAt)

	// 下注
	send(conn, "bet", map[string]any{
		"amount": betAmount,
	})
	time.Sleep(600 * time.Millisecond)

	// 开始游戏
	send(conn, "start", nil)

	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			s := string(msg)

			if !strings.HasPrefix(s, "42") {
				continue
			}

			jsonPart := strings.TrimPrefix(s, "42")
			var parts []json.RawMessage
			if json.Unmarshal([]byte(jsonPart), &parts) != nil {
				continue
			}
			if len(parts) < 2 {
				continue
			}

			var event string
			json.Unmarshal(parts[0], &event)
			if event != "gameService" {
				continue
			}

			var data map[string]any
			json.Unmarshal(parts[1], &data)
			action := data["action"].(string)

			if action == "step" {
				stepCount++
				line := data["payload"].(map[string]any)["lineNumber"].(float64)
				log.Printf("第 %d 步 → 当前车道 %.0f", stepCount, line)

				if stepCount >= stopAt {
					log.Printf("达到目标！第 %d 步执行停止", stepCount)
					send(conn, "stop", nil)
					time.Sleep(1500 * time.Millisecond) // 等待结算
					return
				}
			}

			if action == "gameOver" || action == "crashed" || action == "result" {
				log.Printf("本局结束：%s", action)
				time.Sleep(1 * time.Second)
				return
			}
		}
	}()

	// 超时保护
	select {
	case <-done:
	case <-time.After(35 * time.Second):
		log.Println("本局超时，强制结束")
	}
	return true
}

// 统一发送 gameService 消息
func send(conn *websocket.Conn, action string, payload any) {
	msg := map[string]any{
		"action":  action,
		"payload": payload,
	}
	data := []any{"gameService", msg}
	b, _ := json.Marshal(data)
	full := "42" + string(b)
	conn.WriteMessage(websocket.TextMessage, []byte(full))
}

// 根据实际返回的 JSON 结构定义对应的结构体
type GameLoginResponse struct {
	Success            bool          `json:"success"`
	Result             string        `json:"result"`     // 就是我们要的 token
	Data               string        `json:"data"`       // 和 result 内容一样，很多平台都双发
	GameConfig         any   `json:"gameConfig"` // null
	Bonuses            []any `json:"bonuses"`
	IsLobbyEnabled     bool          `json:"isLobbyEnabled"`
	IsPromoCodeEnabled bool          `json:"isPromoCodeEnabled"`
	IsSoundEnabled     bool          `json:"isSoundEnabled"`
	IsMusicEnabled     bool          `json:"isMusicEnabled"`
}

func RunChickenRoadGame() {
	//operatorId := "bbf6be0a-71c3-4e97-a858-798b01f50fc3" // 会返回这个值/api/ThirdGame/GetGameUrl
	loginCtx, err := desklogin.ReturnContextLoginY1("911117501433", "qwer1234")
	if err != nil {
		fmt.Println("登录失败，无法进行鸡路小游戏的测试", err)
		return
	}
	thirdGameType, err := GetChickenRoadGameListToken(loginCtx)
	if err != nil {
		fmt.Println("获取鸡路小游戏的游戏列表token失败，无法进行后续测试", err)
		return
	}
	chickenAuthToken := thirdGameType.AuthToken
	operatorId := thirdGameType.OperatorId
	// 进入鸡路小游戏
	ChickenRoadGameResponse, err := EnterChickenRoadGame(chickenAuthToken, operatorId)
	if err != nil {
		logger.LogError("进入鸡路小游戏失败", err)
		return
	}
	var resp GameLoginResponse
	if err := json.Unmarshal([]byte(ChickenRoadGameResponse), &resp); err != nil {
		logger.LogError("解析进入鸡路小游戏响应失败", err)
		return
	}
	fmt.Println("进入鸡路小游戏成功，返回的token:", resp.Result)
	fmt.Println("开始自动玩鸡路小游戏...", operatorId)

	// 无限自动玩（推荐）
	//AutoPlayChicken(authToken, operatorId, 100, true)

	// 或者只玩一局测试
	//AutoPlayChicken(resp.Result, operatorId, 100, false)

}
