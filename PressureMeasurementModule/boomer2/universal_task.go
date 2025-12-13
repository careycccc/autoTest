package boomer2

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/myzhan/boomer"
)

var tokenCounter uint64 = 0

// 标准业务函数类型（所有压测接口必须符合这个签名）
type StandardBizFunc func(tokenCtx *context.Context) (int64, bool, string)

// 无参数版本（最常用）
func UniversalTaskNoArgs(
	totalRequests uint64,
	targetRPS int64,
	taskDuration int64,
	taskName string,
	bizFunc StandardBizFunc,
) *boomer.Task {

	return UniversalTaskWithRoundRobinToken(totalRequests, targetRPS, taskDuration, taskName,
		func(tokenCtx *context.Context, _ ...interface{}) (int64, bool, string) {
			return bizFunc(tokenCtx)
		})
}

// 原版带参数版本（保留，兼容未来）
func UniversalTaskWithRoundRobinToken(
	totalRequests uint64,
	targetRPS int64,
	taskDuration int64,
	taskName string,
	bizFunc func(tokenCtx *context.Context, args ...interface{}) (int64, bool, string),
	args ...interface{},
) *boomer.Task {

	log.Printf("启动任务: %s | 目标RPS: %d | 最大运行: %ds", taskName, targetRPS, taskDuration)

	var counter uint64 = 0
	startTime := time.Now()

	return &boomer.Task{
		Name:   taskName,
		Weight: 10,
		Fn: func() {
			if totalRequests > 0 && atomic.LoadUint64(&counter) >= totalRequests {
				return
			}

			current := atomic.AddUint64(&counter, 1)
			if totalRequests > 0 && current > totalRequests {
				return
			}

			// 精确限速
			if targetRPS > 0 {
				interval := time.Second / time.Duration(targetRPS)
				expected := startTime.Add(time.Duration(current-1) * interval)
				if sleep := time.Until(expected); sleep > 0 {
					time.Sleep(sleep)
				}
			}

			// 超时退出
			if taskDuration > 0 && time.Since(startTime) > time.Duration(taskDuration)*time.Second {
				return
			}

			// 轮询获取 token（最均匀）
			tokenCtx := getRoundRobinToken()
			if tokenCtx == nil {
				time.Sleep(500 * time.Millisecond)
				return
			}

			elapsed, success, msg := bizFunc(tokenCtx, args...)

			if success {
				boomer.RecordSuccess("https", taskName, elapsed, 0)
			} else {
				boomer.RecordFailure("https", taskName, elapsed, msg)
			}
		},
	}
}

func getRoundRobinToken() *context.Context {
	PoolMu.RLock()
	defer PoolMu.RUnlock()

	if len(TokenPool) == 0 {
		return nil
	}

	var list []*context.Context
	for _, ctx := range TokenPool {
		list = append(list, ctx)
	}

	idx := atomic.AddUint64(&tokenCounter, 1) - 1
	return list[idx%uint64(len(list))]
}
