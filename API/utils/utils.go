package utils

import (
	"autoTest/store/logger"
	"autoTest/store/model"
	"autoTest/store/request"
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
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

// GetTodayStartAndEnd 返回今天 00:00:00 和 23:59:59 的毫秒时间戳（基于 +08 时区）
func GetTodayStartAndEnd() (int64, int64) {
	// 使用 +08 时区（Asia/Singapore）
	location := time.FixedZone("Asia/Singapore", 8*60*60)

	// 获取当前时间并加载到 +08 时区
	now := time.Now().In(location)

	// 提取年月日
	year, month, day := now.Date()

	// 构造当天的 00:00:00
	start := time.Date(year, month, day, 0, 0, 0, 0, location)

	// 构造当天的 23:59:59
	end := time.Date(year, month, day, 23, 59, 59, 999999999, location) // 精确到纳秒，接近 23:59:59.999...

	// 转换为毫秒时间戳
	return start.UnixMilli(), end.UnixMilli()
}

// 都强制返回 Asia/Singapore（+08）昨天的 00:00:00.000 ～ 23:59:59.999
func GetYesterdayStartEndMilli() (startMs, endMs int64) {
	// 强制使用新加坡时区（关键！）
	loc, _ := time.LoadLocation("Asia/Singapore")

	// 当前时间转成新加坡时间
	nowInSG := time.Now().In(loc)

	// 构造“今天”00:00:00（新加坡时间）
	todayZero := time.Date(
		nowInSG.Year(),
		nowInSG.Month(),
		nowInSG.Day(),
		0, 0, 0, 0,
		loc,
	)

	// 昨天 00:00:00
	yesterdayZero := todayZero.AddDate(0, 0, -1)

	// 昨天 23:59:59.000（你后端真正用的结束时间）
	yesterdayEnd := time.Date(
		yesterdayZero.Year(),
		yesterdayZero.Month(),
		yesterdayZero.Day(),
		23, 59, 59, 0, // 注意：秒和纳秒都设为 0，得到 .000
		loc,
	)

	return yesterdayZero.UnixMilli(), yesterdayEnd.UnixMilli()
}

// Location 北京时间（东八区）
var BeijingLocation = time.FixedZone("Asia/Shanghai", 8*3600)

// SriLankaLocation 斯里兰卡时间 = UTC+5:30
var SriLankaLocation = time.FixedZone("Asia/Colombo", 5*3600+30*60) // +0530
//	startStr: 开始时间字符串，如 "2025-11-01 00:00:00" 或 "2025-11-01 15:04:05"
//	endStr:   结束时间字符串
//	layout:   时间格式，默认为 "2006-01-02 15:04:05"，可自定义

// startUnix      秒级开始时间戳
// startUnixMilli 毫秒级开始时间戳
// endUnix        秒级结束时间戳
// endUnixMilli   毫秒级结束时间戳
// err            错误信息
func ParseTimeRangeToTimestamp(startStr, endStr string, layout ...string) (startUnix, startUnixMilli, endUnix, endUnixMilli int64, err error) {
	// 默认布局
	timeLayout := "2006-01-02 15:04:05"
	if len(layout) > 0 && layout[0] != "" {
		timeLayout = layout[0]
	}
	// 解析时间（使用北京时间）
	startTime, err := time.ParseInLocation(timeLayout, startStr, SriLankaLocation)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("解析开始时间失败: %v", err)
	}
	endTime, err := time.ParseInLocation(timeLayout, endStr, SriLankaLocation)
	if err != nil {
		return 0, 0, 0, 0, fmt.Errorf("解析结束时间失败: %v", err)
	}
	if endTime.Before(startTime) {
		return 0, 0, 0, 0, fmt.Errorf("结束时间不能早于开始时间")
	}
	return startTime.Unix(), startTime.UnixMilli(), endTime.Unix(), endTime.UnixMilli(), nil
}

// 将浮点数保留2位小数 四舍五入
func Rounding(num float64) float64 {
	formattedNum := fmt.Sprintf("%.2f", num)
	// 将字符串转成float64
	result, _ := strconv.ParseFloat(formattedNum, 64)
	return result
}

// GetDayStartEnd 返回当天的开始时间（00:00:00）和结束时间（23:59:59）
// 格式为：2006-01-02 15:04:05
func GetDayStartEnd() (string, string) {
	// 获取当前时间
	now := time.Now()

	// 当天零点（开始时间）
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// 当天 23:59:59（结束时间）
	end := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	// 格式化输出
	const layout = "2006-01-02 15:04:05"
	return start.Format(layout), end.Format(layout)
}

// 切片去重
// SliceUnique 通用去重函数（支持 int、int64、string 等任何可比较类型）
func SliceUnique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := []T{}
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
