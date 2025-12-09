package f1

import (
	login "autoTest/API/deskApi/loginApi"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

var (
	// 全局 token 池（线程安全）
	TokenPool = make(map[string]*context.Context) // username -> token
	PoolMu    sync.RWMutex
)

// LoginTask 改造后：
// 1. 不再需要手动计算 elapsed 时间
// 2. 返回 error 即可，f1 会根据是否返回 error 判断成功/失败
// func LoginTask(username, password string) error {
// 	// 执行登录请求
// 	token, err := login.ReturnContextLoginY1(username, password)

// 	// 失败判定
// 	if err != nil {
// 		return err
// 	}
// 	if token == nil {
// 		return errors.New("login returned nil token")
// 	}

// 	// 登录成功 → 存入 token 池
// 	PoolMu.Lock()
// 	TokenPool[username] = token
// 	PoolMu.Unlock()

// 	return nil
// }

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

	// 存池子（你的原逻辑）
	PoolMu.Lock()
	TokenPool[username] = ctx
	PoolMu.Unlock()

	elapsed := time.Since(start).Milliseconds()
	log.Printf("登录成功 → %s  耗时 %dms", username, elapsed)
	return nil
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
