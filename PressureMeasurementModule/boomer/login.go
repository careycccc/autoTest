// login.go —— v1.6.1 兼容版：失败不计数，成功才计数
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
	targetLogins int    = 100 // 目标：100 个有效 token
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

	log.Printf("加载真实账号 %d 个，失败自动重试，直到拿到 %d 个有效 token", len(realUsers), targetLogins)
}

func LoginTask(username, password string) (elapsed int64, isSuccess bool, msg string) {
	start := time.Now()
	token, err := login.ReturnContextLoginY1(username, password)
	elapsed = time.Since(start).Milliseconds()

	if err != nil || token == nil {
		return elapsed, false, "登录失败: " + err.Error()
	}

	PoolMu.Lock()
	TokenPool[username] = token
	PoolMu.Unlock()

	atomic.AddUint64(&SuccessCount, 1)
	return elapsed, true, ""
}

// 修复版 RunLoginPhase：失败不计入 SuccessCount
// login_task.go —— 终极修复：用真实池大小判断结束
func RunLoginPhase() *boomer.Task {
	return &boomer.Task{
		Name:   taskName,
		Weight: 100,
		Fn: func() {
			// 关键修改：用 len(TokenPool) 判断，而不是 SuccessCount！
			PoolMu.RLock()
			currentPoolSize := len(TokenPool)
			PoolMu.RUnlock()

			if currentPoolSize >= targetLogins {
				return
			}

			// 轮询账号
			i := atomic.AddUint64(&idx, 1) - 1
			username := realUsers[i%uint64(len(realUsers))]
			password := "qwer1234"

			elapsed, isSuccess, errMsg := LoginTask(username, password)

			if isSuccess {
				boomer.RecordSuccess("http", taskName, elapsed, 0)
				// 实时打印当前池大小
				PoolMu.RLock()
				log.Printf("第 %d 个有效 token 已缓存 → %s (当前池大小=%d)",
					len(TokenPool), username, len(TokenPool))
				PoolMu.RUnlock()
			} else {
				atomic.AddUint64(&FailureCount, 1)
				boomer.RecordFailure("http", taskName, elapsed, errMsg)
			}
		},
	}
}
