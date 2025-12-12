// 修改后的 login.go: 新增写入 CSV 逻辑
package creatpool

import (
	login "autoTest/API/deskApi/loginApi"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	// 全局 token 池（线程安全）
	TokenPool = make(map[string]*context.Context) // username -> token
	PoolMu    sync.RWMutex
)

func LoginTask(username, password string) error {
	log.Printf("开始登录 → %s", username) // 打印开始

	start := time.Now()

	// ←←← 关键：加 nil 检查 + 详细错误
	ctx, err := login.ReturnContextLoginY1(username, password)
	if err != nil {
		log.Printf("ReturnContextLoginY1 err → %v for %s", err, username)
		return err
	}
	if ctx == nil {
		log.Printf("ctx is nil for %s! 这就是 panic 原因！", username)
		return fmt.Errorf("ctx is nil")
	}

	// 存池子
	PoolMu.Lock()
	TokenPool[username] = ctx
	PoolMu.Unlock()

	// 新增: 提取 token 并写入 CSV
	token := extractTokenFromCtx(ctx) // 见下面函数
	if err := writeTokenToCSV(username, token); err != nil {
		log.Printf("写入 CSV 失败 for %s: %v", username, err)
	}

	elapsed := time.Since(start).Milliseconds()
	log.Printf("登录成功 → %s  耗时 %dms", username, elapsed)
	return nil
}

// 新增: 从 ctx 提取 token 字符串（根据你的实际 ctx 结构调整）
func extractTokenFromCtx(ctx *context.Context) string {
	// 假设 ctx 携带 Authorization header
	if auth, ok := (*ctx).Value("Authorization").(string); ok {
		return strings.TrimPrefix(auth, "Bearer ") // 提取纯 token
	}
	// 如果不同，调整为你的实际方式，如 (*ctx).Value("token").(string)
	return "" // 默认空，如果失败
}

// 新增: 写入 CSV (username, token)
func writeTokenToCSV(username, token string) error {
	f, err := os.OpenFile("tokens.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	return writer.Write([]string{username, token})
}

// GetToken 辅助函数：从池子取 token
func GetToken(username string) (*context.Context, bool) {
	PoolMu.RLock()
	defer PoolMu.RUnlock()
	t, ok := TokenPool[username]
	return t, ok
}

// RemoveToken 辅助函数：删除 token
func RemoveToken(username string) {
	PoolMu.Lock()
	delete(TokenPool, username)
	PoolMu.Unlock()
}

// 新增: 从 CSV 加载 token 到 TokenPool（用于压测前）
func LoadTokensFromCSV() error {
	TokenPool = make(map[string]*context.Context) // 清空并重建

	f, err := os.Open("tokens.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, row := range records {
		if len(row) != 2 {
			continue // 跳过无效行
		}
		username, token := row[0], row[1]
		// 重建 ctx
		ctx := context.WithValue(context.Background(), "Authorization", "Bearer "+token)
		PoolMu.Lock()
		TokenPool[username] = &ctx
		PoolMu.Unlock()
	}

	log.Printf("从 CSV 加载 %d 个 token 到池子", len(TokenPool))
	return nil
}
