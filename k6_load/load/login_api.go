// load/api_login.go
package load

// import (
// 	"context"
// 	"fmt"

// 	// ← 改成你的 model 包实际路径
// 	yourclient "your-project-path/client" // ← 改成你放 LoginY1 的包路径

// 	"go.k6.io/k6/js/common"
// 	"go.k6.io/k6/js/modules"
// )

// // 我们用一个全局背景 context，后面每个 VU 会基于它创建子 context
// var bgCtx = context.Background()

// // LoginY1 暴露给 k6 的函数
// // 返回值设计成 map[string]any，方便 k6 断言
// func (*Load) LoginY1(username, password string) map[string]any {
// 	// 每个虚拟用户一个独立的 context（避免污染）
// 	ctx := bgCtx

// 	resp, newCtx, err := yourclient.LoginY1(ctx, username, password)
// 	if err != nil {
// 		common.Throw(modules.VUFromContext().Runtime(), fmt.Errorf("login failed: %w", err))
// 	}

// 	// 把可能有用的新 ctx 存到 VU 的状态里，供后续接口使用（可选）
// 	// 目前你只压登录的话可以先不用，后面需要链式调用再打开下面这行
// 	// modules.VUFromContext().State().Context = newCtx

// 	// 把 Response 转成 map，方便 JS 直接断言
// 	result := map[string]any{
// 		"code":    resp.Code,
// 		"msg":     resp.Msg,
// 		"msgCode": resp.MsgCode,
// 		"data":    resp.Data,
// 	}

// 	// 如果你想把原始 []byte 也带出来调试用，可以加：
// 	// "raw": string(resp.RawBody),

// 	return result
// }
