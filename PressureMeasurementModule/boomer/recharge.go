// recharge.go
package boomer

import (
	"time"

	"github.com/myzhan/boomer"
)

func RechargeTask(username string) {
	start := time.Now()

	token, ok := GetToken(username)
	if !ok {
		boomer.RecordFailure("http", "充值", 0, "no token for "+username)
		return
	}

	// 模拟充值请求（带 token）
	_ = token                   // 在真实场景这里加到 Header: Authorization: Bearer xxx
	time.Sleep(1 * time.Second) // 模拟业务 1s

	elapsed := time.Since(start).Milliseconds()
	boomer.RecordSuccess("http", "充值", elapsed, 0)
}
