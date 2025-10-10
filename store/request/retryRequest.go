package request

import (
	"autoTest/store/config"
	"autoTest/store/model"
	"fmt"
	"time"
)

// RetryableFuncWithResult defines an interface for retryable functions
type RetryableFuncWithResult interface {
	Call(params ...interface{}) (interface{}, error)
}

// func2WithResult adapter for functions with 2 string parameters, returning (string, error)
type Func2WithResult func(string, string) (*model.Response, error)

func (f Func2WithResult) Call(params ...interface{}) (interface{}, error) {
	if len(params) < 2 {
		return nil, fmt.Errorf("missing parameters: expected 2, got %d", len(params))
	}
	userName, ok := params[0].(string)
	if !ok {
		return nil, fmt.Errorf("first parameter is not a string: %v", params[0])
	}
	password, ok := params[1].(string)
	if !ok {
		return nil, fmt.Errorf("second parameter is not a string: %v", params[1])
	}
	return f(userName, password)
}

// 需要进行重试的操作
func RetryOperationWithResult(fn RetryableFuncWithResult, params ...interface{}) (interface{}, error) {
	const maxRetries = config.MAXRtryNUMBER //最大重试次数
	const retryInterval = config.FIXEDTIME  // 最大重试时间

	for attempt := 1; attempt <= maxRetries; attempt++ {
		result, err := fn.Call(params...)
		if err == nil {
			return result, nil
		}

		if err.Error() == "请求失败: Requests are too frequent, Please try again later" {
			if attempt == maxRetries {
				return nil, fmt.Errorf("max retries reached: %w", err)
			}
			fmt.Printf("Attempt %d failed with error: %v, retrying in %v...\n", attempt, err, retryInterval)
			time.Sleep(retryInterval)
			continue
		}

		return nil, err
	}

	return nil, nil
}

//调用实例

//  result, err := retryOperationWithResult(func2WithResult(UserloginY1), userName, password)
//     if err != nil {
//         fmt.Printf("UserloginY1 failed: %v\n", err)
//     } else {
//         fmt.Printf("UserloginY1 succeeded, result: %v\n", result)
//     }
