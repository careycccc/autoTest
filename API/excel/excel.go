package excel

import (
	"autoTest/API/adminApi/login"
	memberlist "autoTest/API/adminApi/memberList/memberList"
	"autoTest/API/utils"
	"autoTest/store/logger"
	"context"
	"fmt"
)

func RunExcel() {
	// 示例使用
	userIdlist := getUserId(3000, 0)
	if userIdlist == nil {
		logger.LogError("获取用户ID失败", nil)
		return
	}
	dataSlice := generateUserDataSlice(userIdlist)
	sourcePath := "./1111.xlsx" // 你的源文件路径  需要放在项目的根目录中

	// 调用函数
	newExcel, err := utils.WriteExcelFromSlice(dataSlice, sourcePath)
	if err != nil {
		logger.LogError("写入Excel失败", err)
		return
	}

	// 保存新文件
	if err := newExcel.SaveAs("output.xlsx"); err != nil {
		logger.LogError("保存Excel失败", err)

		return
	}

	logger.Logger.Info("✅ 成功！新文件已保存: output.xlsx")
	logger.Logger.Info(fmt.Sprintf("📊 写入 %d 行数据", len(dataSlice)))

}

// 获取userid userNumber 获取多少用户  userType 用户类型 0 正式账号 1 测试账号 2 游客账号
func getUserId(userNumber int, userType int8) []int {
	ctx := context.Background()
	_, ctxToken, err := login.AdminSitLogin(&ctx)
	if err != nil {
		logger.LogError("登录失败", err)
		return nil
	}
	_, list, err := memberlist.GetUserListApi(ctxToken, userNumber, 1, userType)
	if err != nil {
		logger.LogError("获取用户列表失败", err)
		return nil
	}
	userList := make([]int, 0, userNumber)
	for _, user := range list {
		logger.Logger.Info("用户ID为", user.UserId)
		userList = append(userList, int(user.UserId))
	}
	return userList
}

// GenerateUserDataSlice 根据传入的多个userid生成对应格式的dataSlice
func generateUserDataSlice(userIDs []int) [][]interface{} {
	var dataSlice [][]interface{}

	// 固定字段：-1 和 "测试112"
	level := -1
	userType := "测试114"

	// 遍历所有userid，逐个append
	for _, userID := range userIDs {
		dataSlice = append(dataSlice, []interface{}{userID, level, userType})
	}

	return dataSlice
}
