package f1

import (
	creatpool "autoTest/PressureMeasurementModule/f1/creatPool"
	f1 "autoTest/PressureMeasurementModule/f1/vip"
	"log"
	"os"
)

func RunF1() {
	if len(os.Args) > 1 {
		log.Println("检测到参数 → 执行用户自定义登录压测（只造 token）")
		creatpool.RunPoolTasks()
		return
	}

	log.Println("无参数启动 → 执行完整 F1 压测（自动创建 token 池）")

	creatpool.CreateTokenPoolSilently() // 伪造参数，完美运行你全部智能逻辑

	log.Printf("token 池已就绪，大小: %d，开始业务压测...", len(creatpool.TokenPool))

	f1.RunQueryViptasks()

	log.Println("所有压测场景执行完毕！报表在 reports/ 目录")
}

// f1怎么结合vegeta
