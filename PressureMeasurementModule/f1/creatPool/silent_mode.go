// f1/creatpool/silent_pool.go
package creatpool

import (
	"log"
	"os"
	"time"
)

// 真正静默创建 token 池（通过“伪造”命令行参数欺骗 f1 框架）
func CreateTokenPoolSilently() {
	log.Println("【静默模式】开始创建 token 池（伪造命令行参数）")

	// 保存原始参数
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }() // 结束后恢复

	// 伪造一个 f1 能识别的合法参数
	// 关键：必须包含 "run", "constant", "login" 这些关键字
	// -r 20/s 表示每秒20个并发
	// -d 10m 表示最多跑10分钟（实际会提前结束，因为你有 success >= targetTokens 就停）
	os.Args = []string{
		os.Args[0], // 程序名
		"run",
		"constant",
		"login",
		"-r", "10/s", // 20个并发，够快
		"-d", "20s", // 最多10分钟（实际几秒就够了）
	}

	// 重置全局状态（防止多次调用残留）
	success = 0
	failure = 0
	latencies = latencies[:0]
	activeUsers = append([]string(nil), realUsers...) // 重新加载所有账号

	startTime = time.Now()

	// 直接调用你原来的函数！它会：
	// - 解析上面的“假参数”
	// - 启动 20 并发疯狂登录
	// - 失败自动移除账号
	// - 成功存 token
	// - 达到 150 个自动停止
	RunPoolTasks()

	log.Printf("【静默模式】token 池创建完成！最终大小: %d 个", len(TokenPool))
	if len(TokenPool) == 0 {
		log.Fatal("静默造 token 失败！请检查账号密码、CSV、网络")
	}
}
