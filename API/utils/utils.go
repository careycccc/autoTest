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
