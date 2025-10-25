package topup

import (
	"strings"
)

// 进行充值操作
// func RunRechargeGoods() {
// 	userName := "911024518516"
// 	// 进行用户登录
// 	if _, ctx, err := registerapi.GeneralAgentRegister(userName); err != nil {
// 		logger.LogError("充值的前台用户登录失败", err)
// 	} else {
// 		if _, getRechargeGoodsList, err := GetRechargeGoodsListApi(ctx); err != nil {
// 			logger.LogError("获取充值的充值金额键盘和配置", err)
// 		} else {
// 			localList := make([]string, 0, 10)
// 			// 逻辑只有当所有的三方都充值失败了才会去充值本地,所以先把本地的保存起来
// 			for _, v := range *getRechargeGoodsList {
// 				rechargeGoodsId := v.RechargeGoodsId
// 				for _, i := range v.SupportCategorys {
// 					if checkLocalPrefix(i.RechargeType) {
// 						localList = append(localList, i.RechargeType)
// 						continue
// 					}
// 					rechargeCategoryId := i.ID

// 				}
// 			}
// 		}
// 	}
// }

// 判断Local开头，说明是本地充值
func checkLocalPrefix(s string) bool {
	return strings.HasPrefix(s, "Local")
}
