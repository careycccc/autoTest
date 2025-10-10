package request

import (
	"autoTest/store/config"
	"autoTest/store/logger"
	"autoTest/store/model"
	"errors"
	"reflect"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries     int
	RetryDelay     time.Duration
	TooFrequentMsg string
}

// 默认重试配置
var defaultRetryConfig = RetryConfig{
	MaxRetries:     config.MAXRtryNUMBER,                          // 最大重试次数
	RetryDelay:     config.MAXWaitTIME,                            // 最大重试的间隔时间
	TooFrequentMsg: "Too frequent access, please try again later", // 重试的理由
}

// RetryWrapper 通用的重试包装函数
func RetryWrapper(fn interface{}, args ...interface{}) ([]interface{}, error) {
	return retryWrapperWithConfig(fn, defaultRetryConfig, args...)
}

// retryWrapperWithConfig 带配置的重试包装函数
func retryWrapperWithConfig(fn interface{}, config RetryConfig, args ...interface{}) ([]interface{}, error) {
	// 获取函数的反射值
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		err := errors.New("provided fn is not a function")
		logger.LogError("传入的参数不是一个函数", err)
		return nil, err
	}

	// 准备参数
	fnArgs := make([]reflect.Value, len(args))
	for i, arg := range args {
		fnArgs[i] = reflect.ValueOf(arg)
	}

	for attempt := 0; attempt <= config.MaxRetries; attempt++ {
		// 调用函数
		results := fnValue.Call(fnArgs)

		// 检查返回值数量（假设最后一个返回值是error）
		if len(results) == 0 {
			return nil, errors.New("function has no return values")
		}

		// 获取最后一个返回值（error）
		errVal := results[len(results)-1]
		if !errVal.IsNil() {
			if attempt < config.MaxRetries {
				time.Sleep(config.RetryDelay)
				continue
			}
			return nil, errVal.Interface().(error)
		}

		// 查找 *Response 类型的返回值
		var resp *model.Response
		for i := 0; i < len(results)-1; i++ {
			if r, ok := results[i].Interface().(*model.Response); ok && r != nil {
				resp = r
				break
			}
		}

		// 如果找到 Response，检查是否需要重试
		if resp != nil && resp.Msg == config.TooFrequentMsg && attempt < config.MaxRetries {
			time.Sleep(config.RetryDelay)
			continue
		}

		// 准备返回值
		output := make([]interface{}, len(results))
		for i := range results {
			output[i] = results[i].Interface()
		}
		return output, nil
	}

	return nil, errors.New("max retries reached with too frequent access")
}

// 使用实例
// func main(){
// ctx := context.Background()
// exampleFunc需要被重试的函数
// ctx, 42, "test"是重试函数的参数，依次填写再后面就可以了，
// results 是重试函数的返回值，是一个[]interface{}，可以依次用下标取出
// results, err := RetryWrapper(exampleFunc, ctx, 42, "test")
// if err != nil {
//     // 处理错误
//     println("Error:", err.Error())
//     return
// }
// }
