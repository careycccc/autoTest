// login.go —— 终极修复版：失败不计数，成功才计数，直到拿到 100 个有效 token
package boomer

import (
	login "autoTest/API/deskApi/loginApi"
	"autoTest/store/config"
	"context"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/myzhan/boomer"
)

var (
	realUsers    []string
	idx          uint64 = 0
	totalLogins  uint64 = 0   // 总尝试次数（包括失败）
	targetLogins uint64 = 100 // 目标：100 个有效 token
	TokenPool           = make(map[string]*context.Context)
	PoolMu       sync.RWMutex
	SuccessCount uint64
	FailureCount uint64
)

const taskName = "login"

func init() {
	data, err := os.ReadFile(config.CSVADDR)
	if err != nil {
		log.Fatal("读取 users.csv 失败:", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, line := range lines {
		u := strings.TrimSpace(line)
		if u != "" {
			realUsers = append(realUsers, u)
		}
	}

	if len(realUsers) == 0 {
		log.Fatal("users.csv 里没有账号！")
	}

	// 每个账号最多登录 2 次，但失败不计数，直到拿到 100 个有效 token
	log.Printf("加载真实账号 %d 个，失败自动重试，直到拿到 %d 个有效 token", len(realUsers), targetLogins)
}

// login.go —— 永不崩溃版：任何情况下都不 panic
func LoginTask(username, password string) (elapsed int64, isSuccess bool, msg string) {
	start := time.Now()

	// ←←← 关键：先调用，再判断
	ctx, err := login.ReturnContextLoginY1(username, password)
	elapsed = time.Since(start).Milliseconds()

	// 情况1：err 不为 nil
	if err != nil {
		return elapsed, false, "登录失败: " + err.Error()
	}

	// 情况2：ctx 为 nil（即使 err == nil）
	if ctx == nil {
		return elapsed, false, "登录失败: 返回的 context 为 nil"
	}
	// 3. token 值为空 → 失败跳过
	tokenVal := (*ctx).Value(login.DeskAuthTokenKey)
	if tokenVal == nil {
		return elapsed, false, "登录失败: token 值为空"
	}

	token, ok := tokenVal.(string)
	if !ok || token == "" {
		return elapsed, false, "登录失败: token 类型错误或为空"
	}

	// 全部通过 → 存池子
	PoolMu.Lock()
	TokenPool[username] = ctx
	PoolMu.Unlock()

	atomic.AddUint64(&SuccessCount, 1)
	return elapsed, true, ""
}

// 关键修复：失败不计数！成功才计数！
func RunLoginPhase() *boomer.Task {
	return &boomer.Task{
		Name:   taskName,
		Weight: 100,
		Fn: func() {
			// 已经拿到 100 个有效 token 了，停止发压
			if atomic.LoadUint64(&SuccessCount) >= targetLogins {
				time.Sleep(1 * time.Hour)
				return
			}

			// 轮询账号
			i := atomic.AddUint64(&idx, 1) - 1
			username := realUsers[i%uint64(len(realUsers))]
			password := "qwer1234"

			elapsed, isSuccess, errMsg := LoginTask(username, password)

			// 总尝试次数 +1
			currentAttempt := atomic.AddUint64(&totalLogins, 1)

			if isSuccess {
				// 只有成功才算有效登录
				boomer.RecordSuccess("http", taskName, elapsed, 0)
				currentSuccess := atomic.LoadUint64(&SuccessCount)
				log.Printf("第 %d 个有效 token 已缓存 → %s (总尝试 %d 次)", currentSuccess, username, currentAttempt)
			} else {
				// 失败 → 不计数，继续让其他 goroutine 补
				atomic.AddUint64(&FailureCount, 1)
				boomer.RecordFailure("http", taskName, elapsed, errMsg)
				// 关键：不增加 totalLogins 的“有效”计数
			}

			// 目标达成
			if atomic.LoadUint64(&SuccessCount) >= targetLogins {
				log.Printf("登录阶段完成！成功获取 %d 个有效 token，准备充值", atomic.LoadUint64(&SuccessCount))
			}
		},
	}
}
