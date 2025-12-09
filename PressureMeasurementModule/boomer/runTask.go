package boomer

import (
	"autoTest/store/config"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/myzhan/boomer"
)

var (
	realUsers   []string     // 从 users.csv 读出的真实账号
	idx         uint64   = 0 // 原子计数器，用于轮询账号
	totalLogins uint64   = 0 // 已完成的登录次数
	// 目标：假设 users.csv 有 50 个账号，登录 2 次 = 100 次登录
	targetLogins uint64 = 120
	taskName            = "登录接口"
)

// init 函数：程序启动前自动执行，加载用户数据
func init() {
	// 假设 config.CSVADDR 已正确配置为你的 CSV 文件地址
	data, err := os.ReadFile(config.CSVADDR)
	if err != nil {
		log.Fatal("读取 users.csv 失败:", err)
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	for _, line := range lines {
		u := strings.TrimSpace(line)
		if u != "" {
			realUsers = append(realUsers, u)
		}
	}

	// 根据实际读取到的账号数量设置目标登录次数
	if len(realUsers) > 0 {
		targetLogins = uint64(len(realUsers)) * 2
	} else {
		log.Fatal("users.csv 里没有账号！")
	}

	log.Printf("加载真实账号 %d 个，将每个账号登录 2 次 → 目标完成 %d 次登录", len(realUsers), targetLogins)
}

// RunTasks 函数：定义并执行 Boomer 任务
func RunTasks() {
	log.Printf("开始：尝试在 5 秒内完成 %d 次登录，建立 Token 池...", targetLogins)

	task := &boomer.Task{
		Name:   taskName, // 任务名称用于 Locust 统计
		Weight: 100,
		Fn: func() {
			// 1. 达到目标次数后，让 Goroutine 睡眠，防止退出
			if atomic.LoadUint64(&totalLogins) >= targetLogins {
				time.Sleep(1 * time.Hour)
				return
			}

			// 2. 轮询真实账号
			i := atomic.AddUint64(&idx, 1) - 1
			username := realUsers[i%uint64(len(realUsers))]
			// 假设密码是 qwer1234
			password := "qwer1234"

			// 3. 调用 LoginTask 执行业务逻辑
			// ⚠️ 确保 LoginTask 已经修改为返回 (elapsed, isSuccess, errMsg)
			elapsed, isSuccess, errMsg := LoginTask(username, password)

			// 4. 【核心修复点】 在 Boomer 任务函数中进行统计上报
			if isSuccess {
				boomer.RecordSuccess("http", taskName, elapsed, 0)
			} else {
				boomer.RecordFailure("http", taskName, elapsed, errMsg)
			}

			// 5. 计数与日志
			current := atomic.AddUint64(&totalLogins, 1)

			// 注意：这里 TokenPool 的长度可能小于 totalLogins，因为 TokenPool 是去重后的
			log.Printf("第 %d 次登录完成 → %s (耗时:%dms)", current, username, elapsed)

			// 达到目标后打印完成信息
			if current == targetLogins {
				log.Printf("目标达成！完成 %d 次登录，Token 池已准备就绪", targetLogins)
			}
		},
	}

	// 启动 Boomer 客户端连接 Locust Master
	boomer.Run(task)
	// 主线程阻塞 6 秒，给 Boomer 充足的时间打完 100 次请求
	time.Sleep(6 * time.Second)

	log.Println("--- 登录初始化阶段结束 ---")
	log.Printf("Token 池大小: %d (可用 Token 数量)", len(TokenPool))
	os.Exit(0)
}
