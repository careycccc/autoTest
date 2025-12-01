// main.go  ← 完整替换成这个版本
package pressuremeasurementmodule

import (
	"autoTest/PressureMeasurementModule/accounts"
	"autoTest/store/config"
)

// 运行压测脚本
func RunPressusreModle() {
	accounts.LoadFromCSV(config.CSVADDR)

	// 两个完全一样的任务，只是权重分时控制
	// wave1 := &boomer.Task{Name: "wave1 (0-5s)", Weight: 100, Fn: task.UserTask}
	// wave2 := &boomer.Task{Name: "wave2 (10-15s)", Weight: 0, Fn: task.UserTask}

	// boomer.Run(wave1, wave2)

	// // 精确时间轴控制（2025 年最强写法）
	// go func() {
	// 	// 第 5 秒：第一波结束（关闭 wave1）
	// 	time.Sleep(5 * time.Second)
	// 	wave1.Weight = 0
	// 	log.Println("第一波 1 万用户登录完成（5 秒结束）")

	// 	// 等待 5 秒（5~10 秒什么都不干）
	// 	time.Sleep(5 * time.Second)
	// 	log.Println("等待 5 秒结束，开始第二波")

	// 	// 第 10 秒：开启第二波
	// 	wave2.Weight = 100
	// 	log.Println("第二波开始（10~15 秒冲 1 万用户）")

	// 	// 第 15 秒后：第二波也结束，只剩保活
	// 	time.Sleep(5 * time.Second)
	// 	wave2.Weight = 0
	// 	log.Println("第二波完成，总 2 万用户全部上线，开始每 5 秒保活")

	// 	// 可选：2 小时后自动退出
	// 	// time.Sleep(7200 * time.Second)
	// 	// log.Println("压测结束")
	// 	// os.Exit(0)
	// }()

	// log.Println("Boomer 已启动，前 5 秒冲 1 万用户...")
}
