// common/universal_vegeta_with_alert.go
package common

import (
	creatpool "autoTest/PressureMeasurementModule/f1/creatPool"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var (
	// tokenCounter 使用原子操作进行安全的并发计数
	tokenCounter uint64 = 0
)

// === 配置告警阈值（你可以调这些值，根据你的设计目标）===
// AlertConfig 定义了压测监控的各项告警阈值配置
type AlertConfig struct {
	ErrorRateThreshold      float64       // 错误率超过 1% 告警
	LatencyP95ThresholdMs   int64         // P95 超过 500ms 告警
	LatencyP99ThresholdMs   int64         // P99 超过 1000ms 告警
	ThroughputDropThreshold float64       // TPS 下降超过 10% 视为拐点
	TargetTPS               float64       // 设计目标 TPS（未达到也告警）
	MonitorInterval         time.Duration // 监控间隔时间
}

// defaultAlertConfig 定义了默认的告警阈值配置
var defaultAlertConfig = AlertConfig{
	ErrorRateThreshold:      0.01, // 1%
	LatencyP95ThresholdMs:   500,
	LatencyP99ThresholdMs:   1000,
	ThroughputDropThreshold: 0.10, // 10%
	TargetTPS:               200,  // 示例设计目标：200 TPS，根据你的系统改
	MonitorInterval:         10 * time.Second,
}

// getRoundRobinToken 使用轮询方式从token池中获取token
func getRoundRobinToken() *context.Context {
	if len(creatpool.TokenPool) == 0 {
		return nil
	}
	var tokens []*context.Context
	for _, t := range creatpool.TokenPool {
		tokens = append(tokens, t)
	}
	current := atomic.AddUint64(&tokenCounter, 1)
	index := (current - 1) % uint64(len(tokens))
	return tokens[index]
}

// 万能压测函数（焦点监控：Error Rate + Latency + Throughput/TPS + 拐点检测）
func UniversalVegetaAttack(
	totalRequests uint64, // 总请求数（0=无限）
	tps uint64, // 目标 RPS（作为 TPS 目标基准）t
	duration time.Duration, // 最大运行时间
	taskName string, // 任务名
	reportFile string, // 报表文件名
	bizFunc func(tokenCtx *context.Context, args ...any) (int64, error), // 业务函数
	args ...any,
) {
	log.Printf("启动万能压测：%s | 总请求 %d | 目标RPS/TPS %d | 时长 %v", taskName, totalRequests, tps, duration)

	startTime := time.Now()
	rate := vegeta.Rate{Freq: int(tps), Per: time.Second}

	var (
		successCount uint64
		failureCount uint64
		durations    []int64
		durMu        sync.Mutex
		requestCount uint64

		// 用于 TPS 监控和拐点检测
		lastTPS      float64 = 0 // 上一个间隔的 TPS
		peakTPS      float64 = 0 // 历史峰值 TPS
		lastReqCount uint64      // 上次请求计数
		cfg          = defaultAlertConfig
	)

	// 更新目标 TPS（如果场景 RPS > 设计目标，动态用 RPS 作为基准）
	if float64(tps) > cfg.TargetTPS {
		cfg.TargetTPS = float64(tps)
	}

	// 创建 target 和 attacker
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://placeholder",
	})
	attacker := vegeta.NewAttacker()
	results := attacker.Attack(targeter, rate, duration, taskName)

	// 启动监控 goroutine：每间隔检查 Error Rate / Latency / TPS
	monitorDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(cfg.MonitorInterval)
		defer ticker.Stop()

		for {
			select {
			case <-monitorDone:
				return
			case <-ticker.C:
				currentReq := atomic.LoadUint64(&requestCount)
				currentFailure := atomic.LoadUint64(&failureCount)
				elapsedSec := time.Since(startTime).Seconds()
				if elapsedSec < float64(cfg.MonitorInterval.Seconds()) {
					continue // 跳过初始阶段
				}

				// 计算当前 TPS（Throughput = 成功请求 / 时间）
				recentReq := currentReq - lastReqCount
				currentTPS := float64(recentReq) / cfg.MonitorInterval.Seconds()
				lastReqCount = currentReq

				// 更新峰值 TPS
				if currentTPS > peakTPS {
					peakTPS = currentTPS
				}

				// 监控1: Error Rate（错误率）
				errorRate := float64(currentFailure) / float64(currentReq+1)
				if errorRate > cfg.ErrorRateThreshold {
					log.Printf("🚨 告警！%s 错误率过高: %.2f%% > 阈值 %.2f%% | 可能原因: 接口异常或负载过高",
						taskName, errorRate*100, cfg.ErrorRateThreshold*100)
				}

				// 监控2: Latency（响应时间 P95/P99）
				durMu.Lock()
				localDurs := make([]int64, len(durations))
				copy(localDurs, durations)
				durMu.Unlock()

				if len(localDurs) > 0 {
					sort.Slice(localDurs, func(i, j int) bool { return localDurs[i] < localDurs[j] })
					n := len(localDurs)
					p95 := localDurs[int(math.Min(float64(n-1)*0.95, float64(n-1)))]
					p99 := localDurs[int(math.Min(float64(n-1)*0.99, float64(n-1)))]

					if p95 > cfg.LatencyP95ThresholdMs {
						log.Printf("🚨 告警！%s P95 响应时间过长: %dms > 阈值 %dms | 可能原因: 服务瓶颈或网络延迟\n",
							taskName, p95, cfg.LatencyP95ThresholdMs)
					}
					if p99 > cfg.LatencyP99ThresholdMs {
						log.Printf("🚨 告警！%s P99 响应时间过长: %dms > 阈值 %dms | 可能原因: 极端 case 或资源争用\n",
							taskName, p99, cfg.LatencyP99ThresholdMs)
					}
				}

				// 监控3: Throughput/TPS + 拐点检测
				// 检测1: TPS 未达到设计目标
				if currentTPS < cfg.TargetTPS*(1-cfg.ThroughputDropThreshold) {
					log.Printf("🚨 告警！%s TPS 未达目标: 当前 %.2f < 目标 %.2f | 可能原因: 系统容量不足\n",
						taskName, currentTPS, cfg.TargetTPS)
				}

				// 检测2: 随着负载增加，TPS 不升反降（拐点：当前 TPS < 上次 TPS * (1 - threshold)）
				if lastTPS > 0 && currentTPS < lastTPS*(1-cfg.ThroughputDropThreshold) {
					log.Printf("🚨 拐点检测！%s TPS 不升反降: 当前 %.2f < 上次 %.2f (下降 %.2f%%) | 峰值 %.2f | 可能原因: 性能饱和或资源耗尽\n",
						taskName, currentTPS, lastTPS, (1-currentTPS/lastTPS)*100, peakTPS)
				}
				lastTPS = currentTPS

				// 额外检查: token 池问题
				if len(creatpool.TokenPool) == 0 {
					log.Println("🚨 异常！ token 池为空 | 可能原因: 认证失效")
				}
			}
		}
	}()

	// 处理结果
	for range results {
		if totalRequests > 0 && atomic.LoadUint64(&requestCount) >= totalRequests {
			break
		}

		elapsed, err := bizFunc(getRoundRobinToken(), args...)
		atomic.AddUint64(&requestCount, 1)

		durMu.Lock()
		durations = append(durations, elapsed)
		durMu.Unlock()

		if err != nil {
			atomic.AddUint64(&failureCount, 1)
		} else {
			atomic.AddUint64(&successCount, 1)
		}
	}

	// 停止监控
	close(monitorDone)

	// 生成报表
	generateReport(taskName, reportFile, startTime, successCount, failureCount, durations)
}

// 生成完整 JSON 报表
func generateReport(taskName, fileName string, startTime time.Time, success, failure uint64, durations []int64) {
	total := success + failure
	duration := time.Since(startTime).Seconds()
	rps := 0.0
	if duration > 0 {
		rps = float64(total) / duration
	}

	sort.Slice(durations, func(i, j int) bool { return durations[i] < durations[j] })
	n := len(durations)
	avg := 0.0
	p95 := int64(0)
	p99 := int64(0)
	if n > 0 {
		avg = float64(sum(durations)) / float64(n)
		p95 = durations[int(math.Min(float64(n-1)*0.95, float64(n-1)))]
		p99 = durations[int(math.Min(float64(n-1)*0.99, float64(n-1)))]
	}

	report := map[string]any{
		"场景":              taskName,
		"总请求数":            total,
		"成功次数":            success,
		"失败次数":            failure,
		"成功率":             fmt.Sprintf("%.2f%%", float64(success)*100/float64(total+1)),
		"Current_RPS/TPS": fmt.Sprintf("%.2f", rps),
		"平均响应时间_ms":       fmt.Sprintf("%.2f", avg),
		"P95_响应时间_ms":     p95,
		"P99_响应时间_ms":     p99,
		"token池大小":        len(creatpool.TokenPool),
		"压测时长":            time.Since(startTime).String(),
	}

	data, _ := json.MarshalIndent(report, "", "  ")
	os.WriteFile(fileName, data, 0644)
	log.Printf("报表已生成：%s", fileName)
	log.Printf("压测总结：成功 %d,失败 %d,TPS %.2f,P95 %dms,P99 %dms", success, failure, rps, p95, p99)
}

func sum(arr []int64) int64 {
	var s int64
	for _, v := range arr {
		s += v
	}
	return s
}
