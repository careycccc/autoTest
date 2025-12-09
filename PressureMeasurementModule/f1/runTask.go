package f1

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
	startTime time.Time // 添加测试开始时间
)

// 关键：用 channel 做任务队列，保证顺序 + 唯一
var userQueue chan string

func init() {
	data, err := os.ReadFile(config.CSVADDR)
	if err != nil {
		log.Fatal("读取CSV失败:", err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if u := strings.TrimSpace(line); u != "" {
			realUsers = append(realUsers, u)
		}
	}
	log.Printf("加载 %d 个真实账号，每个账号只登录一次", len(realUsers))

	// 创建 channel 队列，容量就是账号数
	userQueue = make(chan string, len(realUsers))
	for _, u := range realUsers {
		userQueue <- u // 预先放入队列
	}
}

func RunTasks() {
	f := f1.New()
	startTime = time.Now() // 记录开始时间
	f.Add("login", func(t *f1testing.T) f1testing.RunFn {
		t.Logf("开始为 %d 个真实账号各登录一次", len(realUsers))

		// 结束后生成报表
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
			rps := float64(atomic.LoadUint64(&success)+atomic.LoadUint64(&failure)) / duration

			report := map[string]interface{}{
				"场景":          "登录接口压测",
				"总迭代次数":       atomic.LoadUint64(&success) + atomic.LoadUint64(&failure),
				"成功次数":        atomic.LoadUint64(&success),
				"失败次数":        atomic.LoadUint64(&failure),
				"成功率":         fmt.Sprintf("%.2f%%", float64(atomic.LoadUint64(&success))*100/float64(atomic.LoadUint64(&success)+atomic.LoadUint64(&failure))),
				"Current_RPS": fmt.Sprintf("%.2f", rps),
				"平均响应时间_ms":   fmt.Sprintf("%.2f", float64(sum(durations))/float64(n)),
				"P95_响应时间_ms": p95,
				"P99_响应时间_ms": p99,
				"token池大小":    len(TokenPool),
				"压测时长":        time.Since(startTime).String(),
			}

			data, _ := json.MarshalIndent(report, "", "  ")
			os.WriteFile("yc.json", data, 0644)
			log.Println("报表已生成：yc.json")
		})

		return func(t *f1testing.T) {
			// 从队列取一个账号（拿不到就返回）
			select {
			case username := <-userQueue:
				// 成功拿到，执行登录
				start := time.Now()
				err := LoginTask(username, "qwer1234")
				elapsed := time.Since(start).Milliseconds()

				if err != nil {
					atomic.AddUint64(&failure, 1)
					t.Errorf("登录失败 → %s: %v", username, err)
					return
				}

				latMu.Lock()
				latencies = append(latencies, elapsed)
				latMu.Unlock()
				atomic.AddUint64(&success, 1)
				t.Logf("登录成功 → %s (耗时 %dms)", username, elapsed)

			default:
				// 队列空了，说明所有账号都用完了
				return
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
