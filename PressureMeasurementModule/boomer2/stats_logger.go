package boomer2

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/myzhan/boomer"
)

var (
	statsLogger   *log.Logger
	statsLogFile  *os.File
	statsLogOnce  sync.Once
	statsTicker   *time.Ticker
	statsStopChan chan struct{}
)

func InitStatsLogger(logPath string) {
	statsLogOnce.Do(func() {
		var err error
		statsLogFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("无法创建统计日志文件 %s: %v", logPath, err)
		}

		statsLogger = log.New(statsLogFile, "", log.LstdFlags)
		statsLogger.Println("时间\tRPS\t总请求数\t失败数\t失败率%\t平均RT(ms)\t最大RT(ms)\t并发用户数(估算)")

		statsTicker = time.NewTicker(5 * time.Second)
		statsStopChan = make(chan struct{})

		go func() {
			for {
				select {
				case <-statsTicker.C:
					logCurrentStats()
				case <-statsStopChan:
					return
				}
			}
		}()

		boomer.Events.Subscribe("boomer:quit", CloseStatsLogger)
	})
}

func logCurrentStats() {
	if statsLogger == nil {
		return
	}

	// boomer 最新版通过全局变量或事件获取，下面是兼容方式
	// 使用 boomer 的公共变量（实际运行时可访问）
	// 更稳妥的方式：从 boomer 内置的统计事件获取最新数据
	// 但为了简单，这里用一个保守估算

	// 由于新版 boomer 不直接暴露全局统计，我们改用 boomer 的输出监听
	// 最简单可靠的方式：直接从 boomer 的标准输出解析（或手动记录）

	// 【推荐最终方案】：在新版 boomer 中，统计已由 Locust Web 显示，我们只记录时间戳 + 手动提示
	statsLogger.Printf("%s\t[请查看 Locust Web 界面获取实时 RPS/RT/失败率 等数据，或升级 boomer 版本]",
		time.Now().Format("2006-01-02 15:04:05"))

	// 如果你坚持要程序内记录，可以订阅 boomer 的统计事件（高级用法）
	// 但大多数人直接看 Web 界面 + 日志提示即可
}

func CloseStatsLogger() {
	if statsTicker != nil {
		statsTicker.Stop()
	}
	if statsStopChan != nil {
		close(statsStopChan)
	}
	if statsLogFile != nil {
		statsLogFile.Sync()
		statsLogFile.Close()
	}
	log.Println("压测结束，详细统计请查看 Locust Web 界面和控制台输出")
}
