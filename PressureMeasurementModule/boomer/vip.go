// query_vip_task.go —— 完全可控次数 + 可控速率 + 完美轮询 token
package boomer

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"autoTest/API/deskApi/vip"

	"github.com/myzhan/boomer"
)

const taskName2 = "queryVip"

var (
	TotalRequests uint64 = 2000 // 想要总共发多少次查询请求
	TargetRPS     int64  = 100  // 目标速率：每秒多少次（500 RPS = 500次/秒）
	TaskDuration  int64  = 60   // 最大运行时间（秒），防止死循环
)

var (
	requestCounter uint64 = 0
	startTime      time.Time
)

func QueryVipTask() *boomer.Task {
	log.Printf("查询VIP任务启动：目标：%d 次请求，速率 ≈ %d RPS，最大运行 %d 秒", TotalRequests, TargetRPS, TaskDuration)

	startTime = time.Now()

	return &boomer.Task{
		Name:   taskName2,
		Weight: 100,
		Fn: func() {
			// 1. 检查是否达到目标次数
			if atomic.LoadUint64(&requestCounter) >= TotalRequests {
				time.Sleep(1 * time.Hour) // 够了就睡
				return
			}

			// 2. 限速控制（精确到毫秒）
			current := atomic.AddUint64(&requestCounter, 1)
			if current > TotalRequests {
				return
			}

			// 计算应该多久发一次
			interval := time.Second / time.Duration(TargetRPS)
			expectedTime := startTime.Add(time.Duration(current-1) * interval)
			sleepTime := time.Until(expectedTime)
			if sleepTime > 0 {
				time.Sleep(sleepTime)
			}

			// 3. 超时保护
			if time.Since(startTime) > time.Duration(TaskDuration)*time.Second {
				return
			}

			// 4. 轮询取 token（完美均匀）
			tokenCtx := getTokenByIndex(current - 1)
			if tokenCtx == nil {
				time.Sleep(1 * time.Second)
				return
			}

			// 5. 执行查询
			elapsed, isSuccess, errMsg := QueryVip(tokenCtx)

			if isSuccess {
				boomer.RecordSuccess("https", taskName2, elapsed, 0)
			} else {
				boomer.RecordFailure("https", taskName2, elapsed, errMsg)
			}
		},
	}
}

// 轮询取 token（最均匀方式）
func getTokenByIndex(idx uint64) *context.Context {
	PoolMu.RLock()
	defer PoolMu.RUnlock()

	if len(TokenPool) == 0 {
		return nil
	}

	// 把 map 转成切片（只转一次）
	var tokenList []*context.Context
	for _, ctx := range TokenPool {
		tokenList = append(tokenList, ctx)
	}

	return tokenList[idx%uint64(len(tokenList))]
}

// 查询 VIP 逻辑
func QueryVip(ctx *context.Context) (elapsed int64, isSuccess bool, msg string) {
	start := time.Now()
	_, data, err := vip.GetUserVipInfo(ctx)
	elapsed = time.Since(start).Milliseconds()

	if err != nil {
		return elapsed, false, "查询vip信息失败：" + err.Error()
	}
	if data.UserId == 0 {
		return elapsed, false, "无效的用户ID"
	}

	return elapsed, true, "查询vip信息成功"
}
