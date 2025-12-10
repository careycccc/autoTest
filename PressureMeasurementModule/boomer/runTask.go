package boomer

import (
	"log"
	"os"
	"time"

	"github.com/myzhan/boomer"
)

// main.go —— 用真实池大小等待
func RunTasks() {
	// 阶段1：登录（后台）
	go boomer.Run(RunLoginPhase())

	// 2. 关键！强制设置 10000 并发
	// 主线程等待真实池大小达到 100
	log.Println("等待拿到 100 个有效 token...")
	for {
		PoolMu.RLock()
		size := len(TokenPool)
		PoolMu.RUnlock()

		if size >= 100 {
			log.Printf("成功拿到 %d 个有效 token，进入下一阶段！", size)
			break
		}
		time.Sleep(1 * time.Second)
	}

	// 等待 6 秒
	log.Println("等待 6 秒后开始查询VIP...")
	time.Sleep(6 * time.Second)

	// 阶段2：200人查询VIP
	log.Println("开始查询VIP阶段...")
	boomer.Run(QueryVipTask())

	time.Sleep(60 * time.Second)
	log.Println("完整压测结束！")
	os.Exit(0)
}
