package boomer

import (
	"log"
	"os"
	"time"

	"github.com/myzhan/boomer"
)

var (
	// 全局信号：登录阶段完成
	loginDone = make(chan struct{})
)

func RunTasks() {
	log.Println("完整压测启动：100个真实账号各登录2次 → 等全部完成 → 200人查询VIP")

	// ==================== 阶段1：登录建 token 池 ====================
	go func() {
		boomer.Run(RunLoginPhase()) // 登录阶段后台运行
		// 登录阶段自然结束（所有 goroutine 都 sleep 了）
		log.Println("登录阶段彻底结束！100 个 token 已准备好")
		loginDone <- struct{}{} // 发送信号：第一阶段完成
	}()

	// 主 goroutine 等待第一阶段彻底结束
	log.Println("主线程等待登录阶段完成...")
	<-loginDone // 卡在这里，直到登录阶段发信号
	log.Println("登录阶段已确认完成，开始进入等待倒计时...")

	// 等待 2 秒
	log.Println("等待 2 秒后开始查询VIP...")
	time.Sleep(2 * time.Second)

	// ==================== 阶段2：200 人查询VIP ====================
	log.Println("开始查询VIP阶段：200 个虚拟用户（100 token 复用）")
	boomer.Run(QueryVipTask()) // 这行会阻塞，直到压测结束

	log.Println("查询VIP压测结束！完整流程结束！")
	os.Exit(0)
}
