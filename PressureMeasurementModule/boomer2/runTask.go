package boomer2

import (
	"log"
	"os"
	"time"

	"github.com/myzhan/boomer"
)

func main() {
	RunTasks()
}

func RunTasks() {
	// 阶段0：后台启动登录建池
	go boomer.Run(RunLoginPhase())

	// 等待 token 池就绪
	log.Println("等待获取 100 个有效 token...")
	for {
		PoolMu.RLock()
		size := len(TokenPool)
		PoolMu.RUnlock()

		if size >= targetLogins {
			log.Printf("成功获取 %d 个有效 token，开始压测！", size)
			break
		}
		time.Sleep(1 * time.Second)
	}

	time.Sleep(6 * time.Second) // 稳定一下

	// 开始完整 5 阶段压测
	RunVipPressureTest()

	log.Println("所有压测阶段结束，程序退出")
	os.Exit(0)
}

// 公共函数：VIP 接口完整 5 阶段压测
func RunVipPressureTest() {
	log.Println("==================== 开始 VIP 查询接口压测 ====================")

	bizFunc := VipQueryUserInfo // 你的业务函数

	// 阶段1：预热
	log.Println("【阶段1】预热 - 2000 RPS，持续 5 分钟")
	boomer.Run(UniversalTaskNoArgs(0, 2000, 300, "VIP查询-预热", bizFunc))
	time.Sleep(10 * time.Second)

	// 阶段2：线性加压（找拐点）
	log.Println("【阶段2】线性加压 - 逐步提升 RPS")
	steps := []int64{5000, 10000, 20000, 30000, 40000, 50000, 60000, 70000, 80000}
	for _, rps := range steps {
		log.Printf("   ---> 加压到 %d RPS，持续 3 分钟", rps)
		boomer.Run(UniversalTaskNoArgs(0, rps, 180, "VIP查询-加压"+formatRPS(rps), bizFunc))
		time.Sleep(10 * time.Second)
	}

	// 阶段3：稳定持压（根据阶段2实测调整此值）
	log.Println("【阶段3】稳定持压 - 70000 RPS，持续 30 分钟")
	boomer.Run(UniversalTaskNoArgs(0, 70000, 1800, "VIP查询-稳定持压7w", bizFunc))
	time.Sleep(15 * time.Second)

	// 阶段4：峰值冲击
	log.Println("【阶段4】峰值冲击 - 100000 RPS，持续 2 分钟")
	boomer.Run(UniversalTaskNoArgs(0, 100000, 120, "VIP查询-峰值冲击10w", bizFunc))
	time.Sleep(15 * time.Second)

	// 阶段5：降压恢复
	log.Println("【阶段5】降压恢复 - 5000 RPS，持续 5 分钟")
	boomer.Run(UniversalTaskNoArgs(0, 5000, 300, "VIP查询-降压恢复", bizFunc))

	log.Println("==================== VIP 查询压测全部完成 ====================")
}

func formatRPS(rps int64) string {
	if rps >= 10000 {
		return string(rune(rps/10000+'0')) + "w"
	}
	return string(rune(rps/1000+'0')) + "k"
}
