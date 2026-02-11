package sixearn

// import (
// 	"autoTest/API/adminApi/login"
// 	"autoTest/store/logger"
// )

// var masterSelf int = -2
// var subUserId = []subUserInfo
// var subUserInfo = {
// 	userId:0,
// 	hierarchy:0,
// 	betMoneny:0.0,
// 	depMoneny:0.0,
// }

// // 六级返佣代理
// func RunSixearn() {
// 	master := 5944630
// 	// 第一个查询出总代的自身的层级
// 	if ctx, err := login.RunAdminSitLogin(); err != nil {
// 		logger.LogError("六级返佣代理后台登录失败", err)
// 	} else {
// 		if _, masterInfo, err := MasterAgentQuery(ctx, master, false, false); err != nil {
// 			logger.LogError("查询总代信息失败", err)
// 			return
// 		} else {
// 			masterSelf = masterInfo[0].Hierarchy
// 			// 查询出自身及所有的下级
// 			if _, masterInfo, err := MasterAgentQuery(ctx, master, true, false); err != nil {
// 				logger.LogError("查询总代信息失败", err)
// 				return
// 			} else {

// 			}
// 		}
// 	}
// }
