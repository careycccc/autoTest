package desklogintest

// import (
// 	"autoTest/API/deskApi/login"
// 	"autoTest/testcase/common"
// 	"context"
// 	"log"
// 	"testing"

// 	"github.com/ozontech/allure-go/pkg/framework/provider"
// 	"github.com/ozontech/allure-go/pkg/framework/runner"
// )

// func TestLoginY1(t *testing.T) {
// 	resultsDir := common.SetupAllureResultsDir()
// 	runner.Run(t, "登录 /api/Home/Login 测试", func(t provider.T) {
// 		t.Title("测试登录 /api/Home/Login")
// 		t.Description("验证登录 /api/Home/Login 返回正确的响应和token")
// 		t.Labels(common.CommonLabels("/login", "critical")...)
// 		ctx := context.Background()
// 		log.Println("调用 LoginY1")
// 		resp, _, err := login.LoginY1(ctx, "919111997678", "qwer1234")
// 		if err != nil {
// 			log.Printf("LoginY1 失败: %v", err)
// 			t.Errorf("LoginY1 failed: %v", err)
// 			t.WithNewStep("发送登录请求", func(s provider.StepCtx) {
// 				s.Assert().False(true, "LoginY1 调用失败: %v", err.Error())
// 			})
// 			t.Fail()
// 			return
// 		}
// 		t.WithNewStep("发送登录请求", func(s provider.StepCtx) {
// 			common.VerifyLoginResponse(s, resp, "TestLogin")
// 		})
// 	})
// 	common.FlushLogs()
// 	common.CheckAllureResultsDir(resultsDir)
// }
