package main

import (
	_ "autoTest/API/deskApi/loginApi"

	_ "autoTest/store/config"
	"autoTest/store/logger"
)

func init() {
	// 初始化日志 如果需要把日志写入到yaml文件中，就调用logger.InitLogger2()
	logger.InitLogger()
	// logger.Init(config.LogLevel)
	logger.Logger.Info("这是一个信息日志",
		"key", "value",
	)
}

func main() {
	// 模拟一个错误
	// err := someFunction()
	// if err != nil {
	// 	logger.LogError("报错消息", err)
	// }

	//ctx := context.Background()
	// resp, _, err := registerapi.RegisterMobileLoginFunc("911061997780", "P258F5N")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(resp)
	// resp, err := membermanagement.SendVerificationCode("911061997111")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(resp)
	// ctx := context.Background()
	// _, ctxT, err := login.AdminSitLogin(&ctx)
	// if err != nil {
	// 	fmt.Println("Login error:", err)
	// 	return
	// }
	// token := (*ctxT).Value(login.AuthTokenKey)
	// // financialmanagement.ArtificialRechargeFunc(ctxT, 2440315, 1250, 1)
	// // 要上传的文件路径
	// resp, err := messagemanagement.AddCarousel(ctxT, 1, 4, 5, 1, 1)
	// if err != nil {
	// 	fmt.Println("AddCarousel error:", err)
	// 	return
	// }
	// fmt.Println(resp)
	// messagemanagement.AllCustomizedCarousel()
	// uploadfile.FileUploadFunc("./assert/workerOder/1.png", token.(string))
	// if err := activeinformation.RunAddActiveInformation(); err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

}

// // someFunction 模拟一个返回错误的函数
// func someFunction() error {
// 	return fmt.Errorf("出现错误")
// }
