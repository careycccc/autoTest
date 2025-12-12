package creatpool

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"autoTest/store/config"

	"github.com/form3tech-oss/f1/v2/pkg/f1"
	f1testing "github.com/form3tech-oss/f1/v2/pkg/f1/testing"
)

var (
	realUsers []string
	latencies []int64
	latMu     sync.Mutex
	success   uint64
	failure   uint64
	startTime time.Time

	// 关键：可尝试账号队列（动态移除失败账号）
	activeUsers []string
	userMu      sync.Mutex

	// 你想要的最大 token 数（可自由设置）
	targetTokens uint64 = 150 // ←←← 你想要的 token 上限
)

func init() {
	data, err := os.ReadFile(config.CSVADDR)
	if err != nil {
		log.Fatal("读取CSV失败:", err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if u := strings.TrimSpace(line); u != "" {
			realUsers = append(realUsers, u)
			activeUsers = append(activeUsers, u) // 初始全部可用
		}
	}
	log.Printf("加载 %d 个真实账号，失败即移除，直到拿到 %d 个有效 token（或账号耗尽）", len(realUsers), targetTokens)
}

// token池的执行
func RunPoolTasks() {
	f := f1.New()
	startTime = time.Now()

	f.Add("login", func(t *f1testing.T) f1testing.RunFn {
		t.Logf("开始登录压测：%d 个真实账号，失败即移除，目标 %d 个有效 token", len(realUsers), targetTokens)

		// 压测结束自动生成报表
		t.Cleanup(func() {
			latMu.Lock()
			durations := make([]int64, len(latencies))
			copy(durations, latencies)
			latMu.Unlock()

			sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
			n := len(durations)
			p95 := int64(0)
			p99 := int64(0)
			if n > 0 {
				p95 = durations[int(float64(n)*0.95)]
				p99 = durations[min(int(float64(n)*0.99), n-1)]
			}

			duration := time.Since(startTime).Seconds()
			rps := 0.0
			total := atomic.LoadUint64(&success) + atomic.LoadUint64(&failure)
			if duration > 0 {
				rps = float64(total) / duration
			}

			report := map[string]interface{}{
				"场景":          "登录接口压测",
				"总尝试次数":       total,
				"成功次数":        atomic.LoadUint64(&success),
				"失败次数":        atomic.LoadUint64(&failure),
				"成功率":         fmt.Sprintf("%.2f%%", float64(atomic.LoadUint64(&success))*100/float64(total+1)),
				"Current_RPS": fmt.Sprintf("%.2f", rps),
				"平均响应时间_ms":   fmt.Sprintf("%.2f", float64(sum(durations))/float64(n+1)),
				"P95_响应时间_ms": p95,
				"P99_响应时间_ms": p99,
				"token池大小":    len(TokenPool),
				"压测时长":        time.Since(startTime).String(),
			}

			data, _ := json.MarshalIndent(report, "", "  ")
			os.WriteFile("./reports/createPool.json", data, 0644)
			log.Println("报表已生成：createPool.json")
		})

		return func(t *f1testing.T) {
			// 条件1：达到目标 token 数 → 停止
			if atomic.LoadUint64(&success) >= targetTokens {
				//time.Sleep(1 * time.Hour)
				return
			}

			// 条件2：账号已耗尽 → 停止（永不死循环！）
			userMu.Lock()
			if len(activeUsers) == 0 {
				userMu.Unlock()
				//time.Sleep(1 * time.Hour)
				return
			}

			// 取一个账号并立即移除（只尝试一次）
			username := activeUsers[0]
			activeUsers = activeUsers[1:] // 移除
			userMu.Unlock()

			// 执行登录
			start := time.Now()
			err := LoginTask(username, "qwer1234")
			elapsed := time.Since(start).Milliseconds()

			if err != nil {
				atomic.AddUint64(&failure, 1)
				t.Errorf("登录失败（已移除） → %s: %v", username, err)
				log.Printf("账号失败并移除 → %s", username)
				return
			}

			// 成功
			latMu.Lock()
			latencies = append(latencies, elapsed)
			latMu.Unlock()
			atomic.AddUint64(&success, 1)
			t.Logf("第 %d 个有效 token 已缓存 → %s (耗时 %dms)", atomic.LoadUint64(&success), username, elapsed)

			// 达到目标或账号耗尽
			if atomic.LoadUint64(&success) >= targetTokens || len(activeUsers) == 0 {
				finalSize := len(TokenPool)
				log.Printf("登录阶段结束！最终获取 %d 个有效 token（目标 %d，剩余账号 %d）", finalSize, targetTokens, len(activeUsers))
			}
		}
	})

	f.Execute()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sum(arr []int64) int64 {
	var s int64
	for _, v := range arr {
		s += v
	}
	return s
}
