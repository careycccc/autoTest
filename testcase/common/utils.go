package common

import (
	"log"
	"os"
)

// SetupAllureResultsDir 初始化 Allure 结果目录
func SetupAllureResultsDir() string {
	resultsDir := os.Getenv("ALLURE_RESULTS_DIR")
	if resultsDir == "" {
		resultsDir = "./allure-results"
		log.Println("ALLURE_RESULTS_DIR 未设置，默认使用:", resultsDir)
	}
	if err := os.MkdirAll(resultsDir, 0755); err != nil {
		log.Fatalf("无法创建 allure-results 目录: %v", err)
	}
	return resultsDir
}

// CheckAllureResultsDir 检查 Allure 结果目录内容
func CheckAllureResultsDir(resultsDir string) {
	files, err := os.ReadDir(resultsDir)
	if err != nil {
		log.Printf("无法读取 allure-results 目录: %v", err)
		return
	}
	if len(files) == 0 {
		log.Println("allure-results 目录为空")
	}
	for _, file := range files {
		log.Println("allure-results 文件:", file.Name())
	}
}
