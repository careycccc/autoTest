// task/user.go  ← 完整替换成这个版本（2025 年最强写法）
package task

import (
	"autoTest/PressureMeasurementModule/accounts"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/myzhan/boomer"
)

// 模拟真实接口：登录返回一个随机 token
func fakeLogin(username string) string {
	// 模拟网络延迟 + 接口处理
	time.Sleep(1000 * time.Millisecond)

	// 生成一个随机 token（模拟后端返回）
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// 模拟真实保活接口：需要带 token
func fakeKeepAlive(token string) {
	time.Sleep(500 * time.Millisecond)
	// 这里可以模拟校验 token、刷新超时等逻辑
	// fmt.Printf("保活成功，token: %s\n", token[:8])
}

func UserTask() {
	username := accounts.Next()
	// 第一步：登录获取 token（每个用户独立持有）
	start := time.Now()
	token := fakeLogin(username)
	elapsed := time.Since(start).Milliseconds()
	boomer.RecordSuccess("http", "login", elapsed, int64(len(token)))

	// 第二步：登录成功后，每 5 秒带 token 保活一次
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		start = time.Now()
		fakeKeepAlive(token)
		elapsed = time.Since(start).Milliseconds()
		boomer.RecordSuccess("http", "keepalive", elapsed, 0)
	}
}
