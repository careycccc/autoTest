package f1

import (
	"autoTest/API/deskApi/vip"
	"autoTest/PressureMeasurementModule/f1/common"
	creatpool "autoTest/PressureMeasurementModule/f1/creatPool"
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

// f1 业务函数 —— 正确处理三返回值
func getUserVipInfoBizFunc(tokenCtx *context.Context, args ...interface{}) (int64, error) {
	if tokenCtx == nil {
		return 0, fmt.Errorf("tokenCtx is nil")
	}

	start := time.Now()

	// 正确接收 3 个返回值
	resp, _, err := vip.GetUserVipInfo(tokenCtx)

	elapsed := time.Since(start).Milliseconds()

	// 关键判断逻辑：根据业务定义“成功”
	if err != nil {
		return elapsed, fmt.Errorf("请求失败: %w", err)
	}
	if resp == nil {
		return elapsed, fmt.Errorf("resp is nil")
	}
	if resp.Code != 200 && resp.Code != 0 { // 根据你们接口实际成功码调整
		return elapsed, fmt.Errorf("业务失败 code=%d msg=%s", resp.Code, resp.Msg)
	}

	// 可选：打印 data 方便调试
	// log.Printf("VIP查询成功 userID=%s vipLevel=%d", data.UserID, data.VipLevel)

	return elapsed, nil // 成功
}

// 最终版的：支持多场景 + 真实业务 + 自动从 token 池取 context
func RunQueryViptasks() {
	// 1. 先创建并预热 token 池
	creatpool.RunPoolTasks()
	// 判断如果token池为空，直接退出
	if len(creatpool.TokenPool) == 0 {
		log.Println("token为空,退出", len(creatpool.TokenPool))
		return
	}
	log.Println("token 池创建完成，等待 4s 稳定...")
	time.Sleep(4 * time.Second)

	// 2. 定义你想要的复杂压测场景
	scenarios := []struct {
		taskName   string
		reportFile string
		totalReq   uint64 // 0 = 无限，直到 duration 结束
		rps        uint64
		duration   time.Duration
	}{
		{
			taskName:   "vip查询总共发送100次查询,持续10s",
			reportFile: "reports/f1_warm_10rps.json",
			totalReq:   100,
			rps:        10,
			duration:   5 * time.Second,
		},
		{
			taskName:   "vip查询每秒20次,持续1分钟",
			reportFile: "reports/f1_medium_50rps.json",
			totalReq:   0,
			rps:        20,
			duration:   1 * time.Minute,
		},
	}

	// 3. 确保报表目录存在
	if err := os.MkdirAll("reports", 0755); err != nil {
		log.Fatal("创建 reports 目录失败:", err)
	}

	// 4. 依次执行所有场景
	for _, sc := range scenarios {
		log.Printf("🚀 开始执行压测场景：%s | RPS=%d | 时长=%v | 总请求=%d",
			sc.taskName, sc.rps, sc.duration, sc.totalReq)
		time.Sleep(5 * time.Second) // 稍微等待一下，防止日志错乱
		common.UniversalVegetaAttack(
			sc.totalReq,
			sc.rps,
			sc.duration,
			sc.taskName,
			sc.reportFile,
			getUserVipInfoBizFunc, // 正确传入真实业务函数
			// 这里不需要额外参数，直接传 nil 或空
		)

		log.Printf("✅ 场景完成：%s → 报表已生成 %s\n", sc.taskName, sc.reportFile)
	}

	log.Println("🎉 所有 F1 压测场景执行完毕！")
}
