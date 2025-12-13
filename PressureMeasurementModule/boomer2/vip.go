package boomer2

import (
	"autoTest/API/deskApi/vip" // 修改为你的实际包路径
	"context"
	"time"
)

// 标准包装函数：查询 VIP 信息
func VipQueryUserInfo(tokenCtx *context.Context) (int64, bool, string) {
	start := time.Now()

	// 调用你的真实接口（参数根据实际情况修改）
	_, _, err := vip.GetUserVipInfo(tokenCtx) // 示例 userID

	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		return elapsed, false, "VIP查询失败: " + err.Error()
	}

	return elapsed, true, ""
}
