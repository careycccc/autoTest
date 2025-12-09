package boomer

import (
	login "autoTest/API/deskApi/loginApi"
	"context"
	"sync"
	"sync/atomic"
	"time"
)

var (
	// 全局 token 池（线程安全）
	TokenPool = make(map[string]*context.Context) // username -> token
	PoolMu    sync.RWMutex

	// 统计（可选，方便 main.go 打印）
	SuccessCount uint64
	FailureCount uint64
)

// 改造后的 LoginTask —— 所有逻辑都在这里！
// login.go - 改造后 (只执行登录和缓存，不负责统计上报)
func LoginTask(username, password string) (elapsed int64, isSuccess bool, msg string) {
	start := time.Now()
	token, err := login.ReturnContextLoginY1(username, password)
	elapsed = time.Since(start).Milliseconds()

	if err != nil || token == nil {
		return elapsed, false, "登录失败: " + err.Error()
	}

	// 登录成功 → 存入 token 池
	PoolMu.Lock()
	TokenPool[username] = token
	PoolMu.Unlock()

	atomic.AddUint64(&SuccessCount, 1)
	return elapsed, true, "" // 成功不带消息
}

// 辅助函数：从池子取 token（供 recharge.go 等使用）
func GetToken(username string) (*context.Context, bool) {
	PoolMu.RLock()
	defer PoolMu.RUnlock()
	t, ok := TokenPool[username]
	return t, ok
}

// 辅助函数：删除 token（退出时用）
func RemoveToken(username string) {
	PoolMu.Lock()
	delete(TokenPool, username)
	PoolMu.Unlock()
}
