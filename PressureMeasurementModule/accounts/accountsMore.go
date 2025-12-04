package accounts

import (
	commfunc "autoTest/API/commFunc"
	"autoTest/store/config"
	sutils "autoTest/store/utils"
	"context"
	"fmt"
	"sync"

	"github.com/panjf2000/ants/v2"
)

// -----------------------------------------------------------
// 1. 定义数据结构 (用于在任务间传递参数)
// -----------------------------------------------------------

// TaskPayload 用于封装传给 worker 的所有参数
type TaskPayload struct {
	Ctx        *context.Context // 公共 token
	UserID     int64            // 切片遍历的具体值
	MoneyCount float64          // 随机金额
	Wg         *sync.WaitGroup  // 用于通知任务完成
}

// -----------------------------------------------------------
// 2. 封装的核心处理函数
// -----------------------------------------------------------

// BatchUpdateUsers 接收上下文和 ID 切片，使用 10 个并发处理
func BatchUpdateUsers(ctx *context.Context, userIds []int64) error {
	// A. 定义并发池的 Worker 逻辑
	// 这个函数会被 ants 池里的协程反复调用
	workerFunc := func(i interface{}) {
		// 1. 参数断言（拆包）
		payload, ok := i.(*TaskPayload)
		if !ok {
			return
		}
		// 2. 确保任务结束时通知 WaitGroup
		defer payload.Wg.Done()

		// 3. 执行具体的业务逻辑 充值和修改密码的逻辑
		commfunc.UpdatePasswordAndToUp(payload.Ctx, payload.UserID, payload.MoneyCount)
	}

	// B. 创建协程池
	// 容量设为 10，绑定 workerFunc
	pool, err := ants.NewPoolWithFunc(10, workerFunc)
	if err != nil {
		return fmt.Errorf("创建协程池失败: %w", err)
	}
	// 函数退出时销毁池子
	defer pool.Release()

	// C. 准备遍历提交任务
	var wg sync.WaitGroup

	for _, uid := range userIds {
		wg.Add(1)

		// 生成随机金额 (在主协程生成，避免并发锁竞争)
		// 范围 0.00 ~ 100.00
		//money := float64(rand.Intn(10000)) / 100.0
		if money, err := sutils.GenerateRandomInt(config.MIN_MONENY, config.MAX_MONENY); err != nil {
			return err
		} else {
			// 打包参数
			payload := &TaskPayload{
				Ctx:        ctx,
				UserID:     uid,
				MoneyCount: money,
				Wg:         &wg, // 必须传指针
			}

			// D. 提交任务到池子
			// Invoke 是非阻塞的，如果池子满了（10个都在忙），它会自动排队
			err := pool.Invoke(payload)
			if err != nil {
				// 极端情况下（如池子已关闭）可能提交失败
				fmt.Printf("任务提交失败 ID %d: %v\n", uid, err)
				wg.Done() // 手动 Done 避免死锁
			}
		}

	}

	// E. 等待所有任务完成
	wg.Wait()
	fmt.Println(">> 当前批次所有任务处理完毕")
	return nil
}
