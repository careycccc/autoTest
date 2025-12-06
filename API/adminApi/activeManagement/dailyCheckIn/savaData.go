package dailycheckin

import (
	"autoTest/store/logger"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// headers 是 Excel 的表头，使用结构体的中文注释
var headers = []string{
	"会员账号",
	"会员ID",
	"会员当前VIP等级",
	"当天充值金额",
	"连续签到天数",
	"当天签到奖励金额",
	"是否自动领取奖励",
	"是否黑名单",
	"是否进入活动详情页",
	"警告信息",
}

// WriteTodayExcel 将当天的签到数据写入 Excel 文件
// - records: 要写入的 UserDailyCheckInInfo 切片
// - 返回: 生成的文件名和可能的错误
// 过程: 创建新 Excel，写入表头和数据行，确保布尔值写入为字符串 "true"/"false" 以避免兼容问题
func WriteTodayExcel(records []UserDailyCheckInInfo) (filename string, err error) {
	// 生成文件名: 当前日期如 "2006-01-02.xlsx"
	filename = time.Now().Format("2006-01-02") + ".xlsx"

	// 创建新 Excel 文件
	f := excelize.NewFile()

	// 延迟保存文件，并在出错时记录日志
	defer func() {
		if err := f.SaveAs(filename); err != nil {
			log.Printf("保存 Excel 文件失败: %v", err)
		}
	}()

	sheet := "Sheet1" // 使用默认 sheet 名

	// 写入表头（第一行）
	for col, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1) // A1, B1 等
		f.SetCellValue(sheet, cell, h)
	}

	// 写入数据行（从第二行开始）
	for rowIdx, r := range records {
		row := rowIdx + 2 // 行号从 2 开始

		// 逐字段写入，确保类型正确
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), r.UserAccount) // 字符串
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.UserId)      // int
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), r.VipLevel)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), r.RechargeAmount) // float64
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), r.CheckinNumberDay)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), r.CheckinAward)

		// 布尔值写入为小写字符串 "true"/"false"，确保读取兼容性
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), strconv.FormatBool(r.IsAutoReceiveAward))
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), strconv.FormatBool(r.IsBlacklist))
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), strconv.FormatBool(r.IsDetailsPage))

		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), r.AlarmInformation) // 字符串
	}

	// 输出日志
	fmt.Printf("今日签到数据已写入: %s （共 %d 条记录）\n", filename, len(records))
	return filename, nil
}

// RestoreYesterdayExcel 读取昨天的 Excel 文件并恢复为结构体切片
// - 自动计算昨天的文件名
// - 调用 RestoreFromExcel 进行实际读取
func RestoreYesterdayExcel() ([]UserDailyCheckInInfo, error) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02") + ".xlsx"
	return RestoreFromExcel(yesterday)
}

// RestoreFromExcel 从指定 Excel 文件恢复数据
// - filename: 要读取的文件路径
// - 返回: 恢复的 UserDailyCheckInInfo 切片和可能的错误
// 过程: 打开文件，读取所有行，跳过表头，解析每个字段，支持各种布尔表示形式
func RestoreFromExcel(filename string) ([]UserDailyCheckInInfo, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, fmt.Errorf("文件不存在: %s", filename)
	}

	// 打开 Excel 文件
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("打开 Excel 文件失败: %v", err)
	}
	defer f.Close() // 确保关闭文件

	// 获取第一个 sheet
	sheet := f.GetSheetName(0)

	// 获取所有行
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("读取 sheet 失败: %v", err)
	}

	// 检查是否有数据
	if len(rows) <= 1 {
		return nil, fmt.Errorf("文件无数据（只有表头或为空）")
	}

	var records []UserDailyCheckInInfo

	// 从第二行开始解析（i=0 是表头）
	for i := 1; i < len(rows); i++ {
		row := rows[i]

		// 如果行长度不足 10，补齐空字符串防止 panic
		for len(row) < 10 {
			row = append(row, "")
		}

		// 解析每个字段，提供默认值防止转换失败
		userId := 0
		if row[1] != "" {
			userId, _ = strconv.Atoi(row[1])
		}

		vipLevel := 0
		if row[2] != "" {
			vipLevel, _ = strconv.Atoi(row[2])
		}

		rechargeAmount := 0.0
		if row[3] != "" {
			rechargeAmount, _ = strconv.ParseFloat(row[3], 64)
		}

		checkinNumberDay := 0
		if row[4] != "" {
			checkinNumberDay, _ = strconv.Atoi(row[4])
		}

		checkinAward := 0.0
		if row[5] != "" {
			checkinAward, _ = strconv.ParseFloat(row[5], 64)
		}

		// 布尔字段使用兼容函数转换
		isAutoReceiveAward := toBool(row[6])
		isBlacklist := toBool(row[7])
		isDetailsPage := toBool(row[8])

		// 警告信息直接取字符串
		alarmInformation := row[9]

		// 追加到切片
		records = append(records, UserDailyCheckInInfo{
			UserAccount:        row[0],
			UserId:             userId,
			VipLevel:           vipLevel,
			RechargeAmount:     rechargeAmount,
			CheckinNumberDay:   checkinNumberDay,
			CheckinAward:       checkinAward,
			IsAutoReceiveAward: isAutoReceiveAward,
			IsBlacklist:        isBlacklist,
			IsDetailsPage:      isDetailsPage,
			AlarmInformation:   alarmInformation,
		})
	}

	// 输出日志
	fmt.Printf("成功从 %s 恢复 %d 条数据\n", filepath.Base(filename), len(records))
	return records, nil
}

// toBool 将各种形式的布尔值转换为 bool
// 支持: bool 类型、字符串 ("true"/"false"/"1"/"0"/"是"/"否"/"yes"/"no" 等)、数字 (非0为true)
func toBool(v interface{}) bool {
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case bool:
		return val
	case string:
		s := strings.ToLower(strings.TrimSpace(val))
		return s == "true" || s == "1" || s == "是" || s == "yes" || s == "y" || s == "t"
	case int, int8, int16, int32, int64:
		return val.(int64) != 0
	case uint, uint8, uint16, uint32, uint64:
		return val.(uint64) != 0
	case float32, float64:
		return val.(float64) != 0
	default:
		return false
	}
}

// main 函数：示例使用
// 包括写入示例数据和恢复昨天数据的演示
func SaveToDayData(todayData []UserDailyCheckInInfo) {
	// 写入今天的 Excel
	_, err := WriteTodayExcel(todayData)
	if err != nil {
		log.Fatalf("写入失败: %v", err)
	}
}

// 恢复昨天的数据
func RecoverYesterdayData() []UserDailyCheckInInfo {
	// 示例: 恢复昨天的数据（假设昨天文件存在）
	yesterdayData, err := RestoreYesterdayExcel()
	if err != nil {
		fmt.Println("恢复昨天数据失败:", err)
		return nil
	}

	return yesterdayData
}

// 从昨天的这个人的数据中恢复
func RecoverYesterdayDataByAccount(account string) UserDailyCheckInInfo {
	if len(LastDayUserDailyCheckInInfoList) == 0 {
		logger.Logger.Warn("昨天的数据为空")
		return UserDailyCheckInInfo{}
	} else {
		for _, info := range LastDayUserDailyCheckInInfoList {
			if info.UserAccount == account {
				return info
			}
		}
		logger.Logger.Warn(account, "昨天的数据没有找到这个人的数据")
		return UserDailyCheckInInfo{}
	}
}
