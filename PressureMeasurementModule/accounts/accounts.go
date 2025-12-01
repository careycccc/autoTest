package accounts

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var accounts []string
var idx uint64

// 读取csv
func LoadFromCSV(path string) {
	f, _ := os.Open(path)
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Scan() // 跳标题
	for s.Scan() {
		if u := s.Text(); u != "" {
			accounts = append(accounts, u)
		}
	}
	// 会自动跳过第一个账号
	log.Printf("加载 %d 个账号", len(accounts))
}

func Next() string {
	i := atomic.AddUint64(&idx, 1) % uint64(len(accounts))
	return accounts[i]
}

// 写入csv
// SafeCSVWriter 线程安全的 CSV 写入器
type SafeCSVWriter struct {
	file   *os.File
	writer *csv.Writer
	mu     sync.Mutex // 保护 writer.Flush() 和 writer.Write()
	wg     sync.WaitGroup
	errCh  chan error    // 收集写入过程中的错误
	done   chan struct{} // 通知后台 flusher 退出
}

// NewSafeCSVWriter 创建一个安全的 CSV 写入器
// filename: CSV 文件路径（不存在会创建，已存在会追加）
// bufferSize: 内部错误通道缓冲大小，建议 100~10000
func NewSafeCSVWriter(filename string, bufferSize int) (*SafeCSVWriter, error) {
	// O_APPEND 保证多进程下也能安全追加（配合文件锁更佳，这里先用 mutex 保证单进程多协程安全）
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}

	w := csv.NewWriter(file)
	// 根据实际需求调整分隔符，例如：w.Comma = ';'

	scw := &SafeCSVWriter{
		file:   file,
		writer: w,
		errCh:  make(chan error, bufferSize),
		done:   make(chan struct{}),
	}

	// 启动后台自动 Flush（每 500ms 或者缓冲区满时刷新）
	go scw.autoFlush()

	return scw, nil
}

// autoFlush 后台定时 Flush，防止数据长时间滞留在内存
func (s *SafeCSVWriter) autoFlush() {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			s.writer.Flush()
			if err := s.writer.Error(); err != nil {
				s.errCh <- err
			}
			s.mu.Unlock()
		case <-s.done:
			return
		}
	}
}

// Write 单行写入（线程安全）
// record: []string 类型的字段切片
func (s *SafeCSVWriter) Write(record []string) error {
	s.wg.Add(1)
	// 这里使用非阻塞发送，如果缓冲满直接同步写入，防止 goroutine 泄漏
	select {
	case s.errCh <- nil: // 占位，后面统一处理错误
	default:
	}

	s.mu.Lock()
	err := s.writer.Write(record)
	// 立即 Flush 可以保证顺序，但性能低；这里选择写完后交给后台定时 Flush
	if err != nil {
		s.errCh <- err
	}
	s.mu.Unlock()

	s.wg.Done()
	return err // 返回错误让调用方决定是否继续
}

// WriteAll 批量写入（线程安全）
func (s *SafeCSVWriter) WriteAll(records [][]string) error {
	s.wg.Add(1)
	defer s.wg.Done()

	s.mu.Lock()
	defer s.mu.Unlock()

	return s.writer.WriteAll(records) // WriteAll 内部已经做了 Flush
}

// Error 返回错误通道，调用方可选择监听
func (s *SafeCSVWriter) Error() <-chan error {
	return s.errCh
}

// Close 关闭写入器，保证所有数据落盘
func (s *SafeCSVWriter) Close() error {
	// 停止后台 flusher
	close(s.done)

	// 等待所有 Write 完成
	s.wg.Wait()

	s.mu.Lock()
	s.writer.Flush()
	err := s.writer.Error()
	if err != nil {
		s.errCh <- err
	}
	s.mu.Unlock()

	close(s.errCh)

	fileErr := s.file.Close()
	if err != nil {
		return err
	}
	return fileErr
}

// WriteConcurrently 并发安全地将数据写入 CSV 文件
// data:        需要写入的字符串切片，每条字符串代表一行（会自动按逗号分割，或你传入 [][]string）
// concurrency: 并发协程数量（建议 10~200，根据机器和磁盘性能调整）
// filename:    目标 CSV 文件路径
// 返回错误信息（如果有）
func WriteConcurrently(data []string, concurrency int, filename string) error {
	// 1. 创建线程安全的 CSV 写入器
	writer, err := NewSafeCSVWriter(filename, concurrency*10) // 错误通道缓冲大一点
	if err != nil {
		return fmt.Errorf("创建 CSV 写入器失败: %w", err)
	}

	// 2. 可选：监听写入错误
	go func() {
		for err := range writer.Error() {
			if err != nil {
				log.Printf("CSV 写入错误: %v", err)
			}
		}
	}()

	// 3. 启动并发写入
	var wg sync.WaitGroup
	taskCh := make(chan string, concurrency*2) // 任务通道带点缓冲

	// 启动 worker 协程
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for line := range taskCh {
				// 默认按逗号分割字段，你也可以改成其他逻辑
				fields := strings.Split(line, ",")
				// 可选：清理空字段或 trim 空格
				for i := range fields {
					fields[i] = strings.TrimSpace(fields[i])
				}

				if err := writer.Write(fields); err != nil {
					log.Printf("worker %d 写入失败 [%s]: %v", workerID, line, err)
				}
			}
		}(i)
	}

	// 发送所有任务
	for _, line := range data {
		taskCh <- line
	}
	close(taskCh) // 通知 worker 退出

	// 等待所有写入完成
	wg.Wait()

	// 最终关闭，确保所有数据落盘
	if err := writer.Close(); err != nil {
		return fmt.Errorf("关闭 CSV 写入器失败: %w", err)
	}

	fmt.Printf("成功并发写入 %d 行数据到 %s（使用 %d 个协程）\n", len(data), filename, concurrency)
	return nil
}

// 写入示例
// func main() {
//     // 模拟 10 万条数据
//     data := make([]string, 0, 100000)
//     for i := 0; i < 100000; i++ {
//         data = append(data, fmt.Sprintf("%d,用户%d,Beijing,2025-12-01 %02d:%02d:%02d",
//             i, i, i/3600, (i/60)%60, i%60)
//     }

//     // 一键并发写入！协程数量自己定
//     if err := WriteConcurrently(data, 50, "users.csv"); err != nil {
//         log.Fatal("写入失败:", err)
//     }

//     fmt.Println("全部完成！")
// }
