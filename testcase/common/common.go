package common

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 是全局 Zap 日志实例
var Logger *zap.Logger
var loggerOnce sync.Once
var logFilePath string // 存储日志文件路径

// logEntry 定义日志结构
type logEntry struct {
	Level     string      `json:"level"`
	Timestamp string      `json:"timestamp"`
	Caller    string      `json:"caller"`
	Msg       string      `json:"msg"`
	TestName  string      `json:"testName"`
	Response  interface{} `json:"response"`
}

// logStore 存储每个测试用例的最后日志条目
var logStore = struct {
	sync.Mutex
	entries map[string]logEntry
}{entries: make(map[string]logEntry)}

// initLogger 初始化 Zap 日志
func initLogger() *zap.Logger {
	loggerOnce.Do(func() {
		if err := os.MkdirAll("zaplogger", 0755); err != nil {
			log.Fatalf("无法创建 zaplogger 目录: %v", err)
		}
		timestamp := time.Now().Format("20060102_150405")
		logFilePath = filepath.Join("zaplogger", timestamp+".yaml")
		cfg := zap.NewProductionConfig()
		cfg.OutputPaths = []string{logFilePath}
		cfg.EncoderConfig = zapcore.EncoderConfig{
			TimeKey:       "timestamp",
			LevelKey:      "level",
			NameKey:       "logger",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "", // 禁用 stacktrace
			LineEnding:    zapcore.DefaultLineEnding,
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format("2006-01-02 15:04:05"))
			},
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		cfg.Encoding = "json"
		var err error
		Logger, err = cfg.Build(zap.AddCaller())
		if err != nil {
			log.Fatalf("无法初始化 zap 日志器: %v", err)
		}
	})
	return Logger
}

// CommonLabels 返回通用的 Allure 标签
func CommonLabels(api string, severity string) []*allure.Label {
	return []*allure.Label{
		allure.NewLabel("suite", "sit3003接口测试"),
		allure.NewLabel("environment", "staging"),
		allure.NewLabel("severity", severity),
		allure.NewLabel("api", api),
	}
}

// LogAssertion 记录断言结果
func LogAssertion(t provider.StepCtx, testName, message string, success bool, response interface{}) {
	logger := initLogger()
	level := "info"
	msg := "断言成功"
	if !success {
		level = "error"
		msg = "断言失败"
	}
	_, file, line, ok := runtime.Caller(2)
	caller := "unknown"
	if ok {
		caller = fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	entry := logEntry{
		Level:     level,
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Caller:    caller,
		Msg:       msg,
		TestName:  testName,
		Response:  response,
	}
	logStore.Lock()
	logStore.entries[testName] = entry
	logStore.Unlock()
	if success {
		logger.Info(msg, zap.String("testName", testName), zap.Any("response", response))
	} else {
		logger.Error(msg, zap.String("testName", testName), zap.Any("response", response))
	}
}

// FlushLogs 将缓存的日志写入文件
func FlushLogs() {
	logStore.Lock()
	defer logStore.Unlock()
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("无法打开日志文件进行追加: %v", err)
		return
	}
	defer f.Close()
	for _, entry := range logStore.entries {
		jsonBytes, err := json.Marshal(entry)
		if err != nil {
			log.Printf("无法将日志条目序列化为 JSON: %v", err)
			continue
		}
		if _, err := f.WriteString(string(jsonBytes) + "\n"); err != nil {
			log.Printf("无法写入日志文件: %v", err)
		}
	}
	logStore.entries = make(map[string]logEntry)
}

/*
VerifyLoginResponse 公共函数：验证登录响应（泛型，支持任意具有 Code, Msg, MsgCode, Data 字段的结构体指针）
funcName 要测试的函数的名字
这个是响应data有数据的情况
*
*/
func VerifyLoginResponse(s provider.StepCtx, resp any, funcName string) {
	v := reflect.ValueOf(resp).Elem()
	code := int(v.FieldByName("Code").Int())
	msgField := v.FieldByName("Msg")
	msg := msgField.String()
	msgCode := int(v.FieldByName("MsgCode").Int())
	data := v.FieldByName("Data").Interface()

	s.Assert().Equal(0, code, "响应代码应为 0")
	s.Assert().Equal("Succeed", msg, "响应消息应为 'Succeed'")
	s.Assert().Equal(0, msgCode, "响应消息代码应为 0")
	s.Assert().NotNil(data, "响应数据不应为空")

	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Printf("无法序列化响应: %v", err)
		s.Assert().False(true, "无法序列化响应: %v", err.Error())
		return
	}
	responseContent := string(jsonBytes)
	// log.Println("响应内容:", responseContent)
	s.WithNewStep("检查响应", func(subStep provider.StepCtx) {
		subStep.WithParameters(allure.NewParameter("响应内容", responseContent))
		subStep.Assert().True(true, "响应内容已记录")
	})

	// 假设 common.LogAssertion 接受 any resp，并内部处理；如果需要调整，可传入 code, msg 等
	LogAssertion(s, funcName, "登录检查", code == 0 && msg == "Succeed" && msgCode == 0, resp)
}

/*
VerifyLoginResponse 公共函数：验证登录响应（泛型，支持任意具有 Code, Msg, MsgCode, Data 字段的结构体指针）
funcName 要测试的函数的名字
这个是响应data为nil的情况
*
*/
func VerifyLoginResponse2(s provider.StepCtx, resp any, funcName string) {
	v := reflect.ValueOf(resp).Elem()
	code := int(v.FieldByName("Code").Int())
	msgField := v.FieldByName("Msg")
	msg := msgField.String()
	msgCode := int(v.FieldByName("MsgCode").Int())
	serviceTime := v.FieldByName("ServiceTime").Interface()

	// 核心验证逻辑，移除对data非空的强制要求
	s.Assert().Equal(0, code, "响应代码应为 0")
	s.Assert().Equal("Succeed", msg, "响应消息应为 'Succeed'")
	s.Assert().Equal(0, msgCode, "响应消息代码应为 0")
	s.Assert().NotNil(serviceTime, "时间数据不应为空")

	// 序列化响应内容
	jsonBytes, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Printf("无法序列化响应: %v", err)
		s.Assert().False(true, "无法序列化响应: %v", err.Error())
		return
	}
	responseContent := string(jsonBytes)
	//log.Println("响应内容:", responseContent)

	// 记录响应内容的子步骤
	s.WithNewStep("检查响应", func(subStep provider.StepCtx) {
		subStep.WithParameters(allure.NewParameter("响应内容", responseContent))
		subStep.Assert().True(true, "响应内容已记录")
	})

	// 修改LogAssertion调用，只验证code、msg和msgCode
	LogAssertion(s, funcName, "登录检查", code == 0 && msg == "Succeed" && msgCode == 0, resp)
}
