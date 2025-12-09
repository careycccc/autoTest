// logout.go
package boomer

import (
	"time"

	"github.com/myzhan/boomer"
)

func LogoutTask(username string) {
	start := time.Now()

	token, ok := GetToken(username)
	if !ok {
		boomer.RecordFailure("http", "退出接口", 0, "no token for "+username)
		return
	}

	// 模拟退出请求
	_ = token
	time.Sleep(1 * time.Second)

	// 清理 token
	RemoveToken(username)

	elapsed := time.Since(start).Milliseconds()
	boomer.RecordSuccess("http", "退出接口", elapsed, 0)
}
