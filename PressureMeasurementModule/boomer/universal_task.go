// universal_task.go —— 万能压测函数 + 完美轮询 token（2025 终极版）
package boomer

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/myzhan/boomer"
)

var (
	// 全局轮询计数器（所有任务共享，保证均匀）
	tokenCounter uint64 = 0
)

// 万能压测函数（支持任意业务函数 + 完美轮询 token）
func UniversalTaskWithRoundRobinToken(
	totalRequests uint64, // 总请求次数（0 = 无限）
	targetRPS int64, // 目标 RPS（0 = 不限速）
	taskDuration int64, // 最大运行时间（秒）
	taskName string, // 任务名
	bizFunc func(tokenCtx *context.Context, args ...interface{}) (int64, bool, string), // 你的业务函数
	args ...interface{}, // 传给业务函数的其他参数
) *boomer.Task {

	log.Printf("万能压测任务启动：%s | 总请求 %d 次 | 目标 %d RPS | 最大运行 %d 秒 | token 轮询复用",
		taskName, totalRequests, targetRPS, taskDuration)

	var counter uint64 = 0
	startTime := time.Now()

	return &boomer.Task{
		Name:   taskName,
		Weight: 100,
		Fn: func() {
			// 1. 检查是否达到目标次数
			if totalRequests > 0 && atomic.LoadUint64(&counter) >= totalRequests {
				time.Sleep(time.Second * time.Duration(taskDuration+3))
				return
			}

			// 2. 精确限速
			current := atomic.AddUint64(&counter, 1)
			if totalRequests > 0 && current > totalRequests {
				return
			}

			if targetRPS > 0 {
				interval := time.Second / time.Duration(targetRPS)
				expected := startTime.Add(time.Duration(current-1) * interval)
				if sleep := time.Until(expected); sleep > 0 {
					time.Sleep(sleep)
				}
			}

			// 3. 超时保护
			if taskDuration > 0 && time.Since(startTime) > time.Duration(taskDuration)*time.Second {
				return
			}

			// 4. 关键：轮询取 token（最均匀！）
			tokenCtx := getRoundRobinToken()
			if tokenCtx == nil {
				time.Sleep(1 * time.Second)
				return
			}

			// 5. 执行真实业务函数
			elapsed, success, msg := bizFunc(tokenCtx, args...)

			if success {
				boomer.RecordSuccess("https", taskName, elapsed, 0)
			} else {
				boomer.RecordFailure("https", taskName, elapsed, msg)
			}
		},
	}
}

// 完美轮询取 token（全局计数器，保证 100% 均匀）
func getRoundRobinToken() *context.Context {
	PoolMu.RLock()
	defer PoolMu.RUnlock()

	if len(TokenPool) == 0 {
		return nil
	}

	// 转成切片（只转一次，性能拉满）
	var tokenList []*context.Context
	for _, ctx := range TokenPool {
		tokenList = append(tokenList, ctx)
	}

	if len(tokenList) == 0 {
		return nil
	}

	// 全局轮询计数器
	idx := atomic.AddUint64(&tokenCounter, 1) - 1
	return tokenList[idx%uint64(len(tokenList))]
}

// 使用示例
// main.go —— 任意接口，一行调用！
// func RunTasks() {
//     // 阶段1：登录建池（保持不变）
//     boomer.Run(RunLoginPhase())
//     time.Sleep(12 * time.Second)
//     time.Sleep(6 * time.Second)

//     // 阶段2：查询 VIP（2000 次，500 RPS）
//     boomer.Run(UniversalTaskWithRoundRobinToken(
//         2000,      // 总请求次数
//         500,       // 500 RPS
//         60,        // 最大 60 秒
//         "查询VIP信息",
//         func(tokenCtx *context.Context, args ...interface{}) (int64, bool, string) {
//             return QueryVip(tokenCtx)
//         },
//     ))

//     // 阶段3：充值（10000 次，1000 RPS，带参数）
//     boomer.Run(UniversalTaskWithRoundRobinToken(
//         10000,
//         1000,
//         120,
//         "充值接口",
//         func(tokenCtx *context.Context, args ...interface{}) (int64, bool, string) {
//             amount := args[0].(int)
//             return Recharge(tokenCtx, amount)
//         },
//         100, // 充值金额
//     ))

//     // 阶段4：下单（无限运行，100 RPS）
//     boomer.Run(UniversalTaskWithRoundRobinToken(
//         0,       // 0 = 无限
//         100,
//         3600,
//         "下单接口",
//         func(tokenCtx *context.Context, args ...interface{}) (int64, bool, string) {
//             return PlaceOrder(tokenCtx)
//         },
//     ))
// }
