package utils

import (
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

// 处理基础的struct,返回payloadStruct, payloadList
func BaseStructHandler() (*model.BaseStruct, []interface{}) {
	payloadStruct := &model.BaseStruct{}
	timestamp, random, language := request.GetTimeRandom()
	payloadList := []interface{}{random, language, "", timestamp}
	return payloadStruct, payloadList
}

// 需要传入需要生成多少一个id的个数，并且返回id的列表
func RandmoUserId(generateCount int) []string {
	// 模拟高并发生成100万个ID
	var wg sync.WaitGroup
	generated := sync.Map{} // 存储已生成的ID，检查重复
	collisionCount := 0
	// generateCount := 1000000
	idList := make([]string, 0, generateCount)
	for i := 0; i < generateCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id, err := RandmoUserCount()
			if err != nil {
				logger.LogError("Error generating ID", err)
				return
			}
			// 检查重复
			if _, exists := generated.LoadOrStore(id, true); exists {
				logger.Logger.Info("用户的重复检测", id)
				collisionCount++
			}
			// 生成的用户id，可以进行接下来的操作
			idList = append(idList, id)
		}()
	}

	wg.Wait()
	logger.Logger.Info("已生成的用户数", generateCount, "重复的用户数", collisionCount)
	return idList
}

// 随机生成用户以今日的日期开头的
func RandmoUserCount() (string, error) {
	// 获取当前日期
	now := time.Now()
	month := now.Month()
	day := now.Day()

	// 格式化月和日
	var prefix string
	if month < 10 {
		prefix = fmt.Sprintf("%d%02d", month, day) // 月1位+日2位=3位
	} else {
		prefix = fmt.Sprintf("%02d%02d", month, day) // 月2位+日2位=4位
	}

	// 根据前缀长度决定随机数位数
	var randomLength int
	if len(prefix) == 3 {
		randomLength = 7
	} else {
		randomLength = 6
	}

	// 生成随机数
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(randomLength)), nil) // 10^randomLength
	randNum, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// 格式化随机数，补0到指定长度
	randStr := fmt.Sprintf("%0*d", randomLength, randNum)

	// 合并前缀和随机数
	return "91" + prefix + randStr, nil
}

// WriteExcelFromSlice 封装函数：切片数据写入新Excel，返回完整Excel对象
// 参数1: dataSlice - [][]interface{} 切片，每个子切片对应一行，每列顺序对应表头
// 参数2: sourcePath - string 源Excel文件路径
// 返回: *excelize.File - 全新Excel文件对象（表头来自源文件第一行，数据从第二行开始）
func WriteExcelFromSlice(dataSlice [][]interface{}, sourcePath string) (*excelize.File, error) {
	// 1. 读取源文件第一行表头
	fSource, err := excelize.OpenFile(sourcePath)
	if err != nil {
		logger.LogError("打开源文件失败", err)
		return nil, fmt.Errorf("打开源文件失败: %v", err)
	}
	defer fSource.Close()

	sheetName := fSource.GetSheetName(0)
	rows, err := fSource.GetRows(sheetName)
	if err != nil || len(rows) == 0 {
		logger.LogError("读取源文件表头失败", err)
		return nil, fmt.Errorf("读取源文件表头失败")
	}
	headers := rows[0] // 第一行作为表头

	// 2. 创建新Excel文件
	fNew := excelize.NewFile()
	newSheet := "Sheet1"

	// 3. 写入表头（第一行）
	for col, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+col)
		fNew.SetCellValue(newSheet, cell, header)
	}

	// 4. 写入数据（从第二行开始）
	for rowIdx, rowData := range dataSlice {
		rowNum := rowIdx + 2 // 第2行、第3行...

		for colIdx, cellValue := range rowData {
			if colIdx >= len(headers) {
				break // 超出表头列数，停止写入
			}
			cell := fmt.Sprintf("%c%d", 'A'+colIdx, rowNum)
			fNew.SetCellValue(newSheet, cell, cellValue)
		}
	}

	// 5. 自动调整列宽
	lastCol := 'A' + len(headers) - 1
	fNew.SetColWidth(newSheet, "A", fmt.Sprintf("%c", lastCol), 15)
	return fNew, nil
}
