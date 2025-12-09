package accounts

import (
	"autoTest/API/adminApi/login"
	memberlist "autoTest/API/adminApi/memberList/memberList"
	"autoTest/store/config"
	"autoTest/store/logger"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

var accounts []string
var idx uint64

// SummaryItem 结构体用于存放处理后的单行数据（可以根据实际需求替换为 OrderItem 或其他结构体）
type ProcessedRecord struct {
	Record   []string
	WorkerID int
}

// ConcurrentCSVReader 封装了标准的 csv.Reader 并添加了互斥锁，确保并发安全读取
type ConcurrentCSVReader struct {
	reader *csv.Reader
	mu     sync.Mutex // 互斥锁，用于保护读取状态
}

// NewConcurrentCSVReader 构造函数
func NewConcurrentCSVReader(r io.Reader) *ConcurrentCSVReader {
	return &ConcurrentCSVReader{
		reader: csv.NewReader(r),
	}
}

// Read 线程安全地从 CSV 中读取一行记录
func (c *ConcurrentCSVReader) Read() (record []string, err error) {
	// 锁定互斥锁，开始临界区
	c.mu.Lock()
	defer c.mu.Unlock() // 确保在函数返回时释放锁

	return c.reader.Read()
}

// ProcessCSVConcurrently 是最终封装的函数。
// 它接收 CSV 文件路径和希望启动的工作 Goroutine 数量。
// 返回处理后的记录切片 (这里使用 [][]string，方便演示) 和错误信息。
func ProcessCSVConcurrently(filePath string, numWorkers int) ([]string, error) {
	// 1. 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件 %s: %w", filePath, err)
	}
	defer file.Close() // 确保文件关闭

	// 2. 创建并发安全的 Reader
	concurrentReader := NewConcurrentCSVReader(file)

	// 可选：读取并忽略 Header
	if _, err := concurrentReader.Read(); err != nil && err != io.EOF {
		return nil, fmt.Errorf("读取 CSV Header 失败: %w", err)
	}

	var wg sync.WaitGroup
	// 使用 Channel 收集结果，Channel 的容量应预估
	resultChan := make(chan []string, 100)

	// 3. 启动工作 Goroutine
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				// Mutex 保证读取安全
				record, err := concurrentReader.Read()

				if err == io.EOF {
					// 读取完毕
					return
				}
				if err != nil {
					// 记录读取错误并退出当前 Goroutine
					fmt.Printf("[Worker %d] 读取出错: %v\n", workerID, err)
					return
				}

				// --- 4. 数据处理和发送结果 ---
				// 模拟数据处理逻辑 (例如：类型转换、业务计算等)
				//processedRecord := append([]string{fmt.Sprintf("Worker-%d", workerID)}, record...)

				// 将处理后的结果安全地发送到 Channel
				resultChan <- record
			}
		}(i)
	}

	// 5. 启动一个 Goroutine 等待所有工作完成并关闭 Channel
	go func() {
		wg.Wait()
		close(resultChan) // 确保在所有数据发送完毕后关闭 Channel
	}()

	// 6. 从 Channel 中收集所有结果到切片中
	var allFlatRecords []string // 声明为一维切片
	for record := range resultChan {
		// 遍历每一行记录，将其所有字段附加到一维切片中
		allFlatRecords = append(allFlatRecords, record...)
	}

	return allFlatRecords, nil
}

// =============================================================================
// SafeCSVWriter 定义了一个线程安全的 CSV 写入器
// =============================================================================

// SafeCSVWriter 结构体
type SafeCSVWriter struct {
	file  *os.File    // 目标文件句柄
	csvW  *csv.Writer // 标准 csv.Writer
	mu    sync.Mutex  // 用于保护对 csvW 的写入操作
	errCh chan error  // 错误通道，用于接收异步写入的错误
}

// NewSafeCSVWriter 创建并初始化 SafeCSVWriter
func NewSafeCSVWriter(filename string, errorChanSize int) (*SafeCSVWriter, error) {
	// 创建或截断文件
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)
	writer.Comma = ',' // 默认使用逗号作为分隔符

	return &SafeCSVWriter{
		file:  file,
		csvW:  writer,
		errCh: make(chan error, errorChanSize),
	}, nil
}

// Write 是线程安全的写入方法。它只写入一行。
func (w *SafeCSVWriter) Write(record []string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 写入记录
	if err := w.csvW.Write(record); err != nil {
		// 如果写入失败，将错误发送到错误通道
		select {
		case w.errCh <- err:
			// 错误发送成功
		default:
			// 错误通道已满，避免阻塞，打印日志
			log.Printf("警告：错误通道已满，丢弃写入错误: %v", err)
		}
		return err // 同时返回错误给调用者
	}
	return nil
}

// Close 关闭文件和 CSV 写入器，并确保所有缓冲数据落盘。
func (w *SafeCSVWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// 强制将缓冲数据写入文件
	w.csvW.Flush()
	if err := w.csvW.Error(); err != nil {
		w.file.Close()
		close(w.errCh)
		return fmt.Errorf("CSV Flush 错误: %w", err)
	}

	// 关闭文件句柄
	if err := w.file.Close(); err != nil {
		close(w.errCh)
		return fmt.Errorf("关闭文件失败: %w", err)
	}

	// 最终关闭错误通道
	close(w.errCh)
	return nil
}

// Error 返回错误通道，供外部监听写入错误
func (w *SafeCSVWriter) Error() <-chan error {
	return w.errCh
}

// =============================================================================
// 改造后的并发写入函数
// =============================================================================

// WriteConcurrently 并发地将数据切片中的独立数字（由空格分隔）
// 一行一行写入到 CSV 文件中。
func WriteConcurrently(data []string, concurrency int, filename string) error {
	// 1. 创建线程安全的 CSV 写入器
	writer, err := NewSafeCSVWriter(filename, concurrency*10) // 错误通道缓冲大一点
	if err != nil {
		return fmt.Errorf("创建 CSV 写入器失败: %w", err)
	}

	// 2. 监听写入错误（在后台协程运行）
	go func() {
		for csvErr := range writer.Error() {
			// 这里我们只记录错误，不中断主流程
			log.Printf("CSV 写入器捕获到错误: %v", csvErr)
		}
	}()

	// 3. 启动并发写入
	var wg sync.WaitGroup
	// 任务通道，用于发送单个数字 (写入任务)
	taskCh := make(chan string, concurrency*4)

	// 启动 worker 协程
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for numberStr := range taskCh {
				// 核心改造：将单个数字作为 CSV 的一行（一个字段）写入
				fields := []string{numberStr}

				if err := writer.Write(fields); err != nil {
					// 注意：这里的错误已经被 writer.Error() 监听了，但同时返回给 worker
					// 也可以在这里 log，确保信息不丢失。
					log.Printf("worker %d 写入任务失败 [%s]: %v", workerID, numberStr, err)
				}
			}
		}(i)
	}

	// 4. 发送所有任务 (在新协程中执行，以便主线程可以等待 worker)
	totalRowsSent := 0

	go func() {
		defer close(taskCh) // 任务发送完毕后，关闭通道

		// 遍历输入的每一行字符串
		for _, line := range data {
			// 使用 strings.Fields 按空格分割出所有非空数字
			numbers := strings.Fields(line)

			// 将每个数字作为一个独立任务发送
			for _, num := range numbers {
				// 确保数据非空（strings.Fields 已经处理了大部分情况）
				if num != "" {
					taskCh <- num
					totalRowsSent++
				}
			}
		}

		fmt.Printf("任务发送完成，总计生成 %d 个独立写入任务。\n",
			totalRowsSent)
	}()

	// 5. 等待所有 worker 协程完成工作
	wg.Wait()

	// 6. 最终关闭，确保所有数据落盘
	if err := writer.Close(); err != nil {
		return fmt.Errorf("关闭 CSV 写入器失败: %w", err)
	}

	fmt.Printf("成功并发写入 %d 行数据到 %s（使用 %d 个协程）。\n", totalRowsSent, filename, concurrency)
	return nil
}

// 启动写入到csv中
// 需要确保写入的数据是账号，有充值，改过密码
func RunWirteCsv() {
	userList := make([]int64, 0, 1000)
	// 后台登录
	if ctx, err := login.RunAdminSitLogin(); err != nil {
		logger.LogError("csv写入的后台登录失败", err)
		return
	} else {
		// 获取后台用户列表的userid
		if resp, UserInfo, err := memberlist.GetUserListApi(ctx, 5, 200, 0); err != nil {
			logger.Logger.Warn("csv写入的获取用户列表录失败", err)
			return
		} else {
			fmt.Println("csv写入", resp)
			for _, v := range UserInfo {
				userList = append(userList, v.UserId)
			}
			// 进行用户的修改密码和充值
			if len(userList) <= 0 {
				logger.Logger.Warn("会员id列表的值为空")
				return
			}
			if err := BatchUpdateUsers(ctx, userList); err != nil {
				logger.Logger.Warn("会员id多线程执行失败")
			}

			// 返回用户的账号列表
			userNameList := make([]string, 0, 1000)
			for _, v := range userList {
				if _, username, err := memberlist.GetUserAmount(ctx, int(v)); err != nil {
					logger.Logger.Warn("用户id返回用户账号转换失败", err)
					continue
				} else {
					userNameList = append(userNameList, username)
				}
			}
			//fmt.Println("修改号的账号列表", userNameList)
			//
			WriteConcurrently(userNameList, 5, config.CSVADDR)
		}
	}

}
