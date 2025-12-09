// query_vip_task.go —— 200 人查询 VIP（100 token 完美复用）
package boomer

import (
	"autoTest/API/deskApi/vip"
	"context"
	"log"
	"sync/atomic"
	"time"

	"github.com/myzhan/boomer"
)

const taskName2 = "queryVip"

// 全局计数器，用于轮询取 token（比 rand 更均匀）
var globalCounter uint64 = 0

func QueryVipTask() *boomer.Task {
	log.Println("查询VIP任务启动：200 个虚拟用户（100 token 轮询复用）")

	return &boomer.Task{
		Name:   taskName2,
		Weight: 100,
		Fn: func() {
			// 已经执行完 200 次，直接睡
			if atomic.LoadUint64(&globalCounter) >= 200 {
				time.Sleep(1 * time.Hour)
				return
			}

			// 计算当前是第几次请求（0~199）
			idx := atomic.AddUint64(&globalCounter, 1) - 1

			// 轮询取 token（0~99 → 第一个 token，100~199 → 第二个 token）
			tokenList := getTokenList() // 把 map 转成切片
			if len(tokenList) == 0 {
				return
			}
			tokenCtx := tokenList[idx%uint64(len(tokenList))]

			// 执行查询 VIP
			elapsed, isSuccess, errMsg := QueryVip(tokenCtx)

			if isSuccess {
				boomer.RecordSuccess("https", taskName2, elapsed, 0)
			} else {
				boomer.RecordFailure("https", taskName2, elapsed, errMsg)
			}
		},
	}
}

// 真实查询 VIP 逻辑（
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

// 把 TokenPool 转成切片（只转一次）
func getTokenList() []*context.Context {
	PoolMu.RLock()
	defer PoolMu.RUnlock()

	if len(TokenPool) == 0 {
		return nil
	}

	list := make([]*context.Context, 0, len(TokenPool))
	for _, ctx := range TokenPool {
		list = append(list, ctx)
	}
	return list
}
