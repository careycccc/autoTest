package request

import (
	"autoTest/store/config"
	"autoTest/store/utils"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"time"
)

// get请求
/**
url 请求地址
api 接口
args[0] 添加额外的请求头设置
args[1] 添加参数
*/
func GetRequest(base_url, api string, args ...map[string]interface{}) ([]byte, *http.Response, error) {
	urlapi := base_url + api
	// 携带了参数需要设置参数
	if len(args[1]) > 0 {
		params := url.Values{}
		par := setGetParms(&params, args[1])
		urlapi = urlapi + "?" + par.Encode()
	}
	// fmt.Println("添加了参数的url---》", urlapi)
	// 创建 GET 请求
	req, err := http.NewRequest("GET", urlapi, nil)
	if err != nil {
		log.Fatalf("创建请求失败: %v", err)
	}
	// 设置请求头
	if len(args) > 0 {
		setHeaders(req, args[0])
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36")
	//req.Header.Set("Referer", "https://h5.wmgametransit.com/WinGo/WinGo_5M?Lang=en&Skin=Classic&SkinColor=Default&Token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJUb2tlblR5cGUiOiJBY2Nlc3NfVG9rZW4iLCJUZW5hbnRJZCI6IjMwMDMiLCJVc2VySWQiOiIzMDAzMDAwMTczNjMxNyIsIkFnZW50Q29kZSI6IjMwMDMwMSIsIlRlbmFudEFjY291bnQiOiIxNzM2MzE3IiwiTG9naW5JUCI6IjE3NS4xNTcuODYuMjAiLCJMb2dpblRpbWUiOiIxNzU2NDgxMzQxMzUxIiwiU3lzQ3VycmVuY3kiOiJJTlIiLCJTeXNMYW5ndWFnZSI6ImVuIiwiRGV2aWNlVHlwZSI6IlBDIiwiTG90dGVyeUxpbWl0R3JvdXBOdW0iOiIwIiwiVXNlclR5cGUiOiIwIiwibmJmIjoxNzU2NDgxMzQxLCJleHAiOjE3NTY1Njc3NDEsImlzcyI6Imp3dElzc3VlciIsImF1ZCI6ImxvdHRlcnlUaWNrZXQifQ.mtB4BS7ZpIp0xPItV-he2tISkDKC0wzMp2mWAvrfoys&RedirectUrl=https%3A%2F%2Fsit-plath5-y1.mggametransit.com%2Fgame%2Fcategory%3FcategoryCode%3DC202505280608510046&Beck=0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Connection", "keep-alive")

	// 获取请求句柄
	client := checkHttp2()
	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println("发送get请求失败")
		err := errors.New("发送get请求失败")
		return nil, nil, err
	}
	defer resp.Body.Close()
	return handlerCode(resp)
}

/*
*
paylaod 请求参数 map[string]interface{}
base_url 请求地址
api 接口地址
args[0] 添加自定义请求头 map[string]interface{}
*/
func PostRequestCofig(payload map[string]interface{}, base_url, api string, args ...map[string]interface{}) ([]byte, *http.Response, error) {
	url := base_url + api
	// fmt.Printf("本次请求的地址%v\n", url)
	// 判断传进来的paylaod是否有签名，没有就添加上
	_, exists := payload["signature"]
	if !exists {
		payload["signature"] = ""
	}
	verfiyp := ""
	signature := utils.GetSignature(payload, &verfiyp)
	if signature == "" {
		// fmt.Println("utils的签名是空的", signature)
		return nil, nil, errors.New("utils的签名是空的")
	}

	payload["signature"] = signature
	//将请求数据转换成json
	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf(" json 编码失败:%v", err)
	}
	// 发送post请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("请求失败：%v", err)
	}
	// 设置请求头
	if len(args) > 0 {
		setHeaders(req, args[0])
	}
	req.Header.Set(
		"User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36",
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Connection", "keep-alive")
	client := checkHttp2()

	//发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("发送请求失败~~~~~~~: %v", err)
		return nil, nil, err
	}

	defer resp.Body.Close()
	// fmt.Println("响应状态码", resp.StatusCode)
	respBody, resp, err := handlerCode(resp)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}

	return respBody, resp, nil
}

// 设置请求头
func setHeaders(req *http.Request, headers map[string]interface{}) {
	for key, value := range headers {
		// 将 interface{} 转换为 string
		var headerValue string
		switch v := value.(type) {
		case string:
			headerValue = v
		case fmt.Stringer:
			headerValue = v.String()
		default:
			headerValue = fmt.Sprintf("%v", v)
		}
		req.Header.Set(key, headerValue)
	}
}

// 设置get参数
func setGetParms(params *url.Values, paramsMap map[string]interface{}) url.Values {
	for key, value := range paramsMap {
		var paramsValue string
		switch v := value.(type) {
		case string:
			paramsValue = v
		case fmt.Stringer:
			paramsValue = v.String()
		default:
			paramsValue = fmt.Sprintf("%v", v)
		}
		params.Add(key, paramsValue)
	}
	return *params
}

// 响应码的处理
func handlerCode(resp *http.Response) ([]byte, *http.Response, error) {
	if resp.StatusCode == 200 {
		//获取相应的内容
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("读取响应失败：%v", err)
		}
		return respBody, resp, nil

	} else if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		//
		errString := errors.New("状态码300-400")
		return nil, resp, errString
	} else if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		//需要身份验证,或者需要
		errString := errors.New("需要身份验证,或者需要,状态码400-500")
		return nil, resp, errString
	} else {
		err := errors.New("状态码不是200~~或者是服务器错误~~~")
		return nil, nil, err
	}
}

// 验证请求的协议是不是http/2
func checkHttp2() *http.Client {
	// client := &http.Client{
	// 	// 检查 确保使用http/2
	// 	Transport: &http.Transport{
	// 		ForceAttemptHTTP2: true,
	// 	},
	// }
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 仅限测试，跳过证书验证
			},
		},
	}
	return client
}

// 获取微妙级的时间戳并返回
func GetNowTime() int64 {
	now := time.Now()
	timestampMicro := now.Unix()
	// fmt.Println("当前时间戳", timestampMicro)
	return timestampMicro
}

// 生成9位的随机数
func RandomNine() int64 {
	max_number := big.NewInt(900000000000)
	n, err := rand.Int(rand.Reader, max_number)
	if err != nil {
		fmt.Printf("生成随机数失败：%v", err)
		return -1
	}
	random_number := n.Int64() + 100000000000
	// fmt.Println("生成的随机数", random_number)
	return random_number
}

// 返回时间戳，随机数，语言，
func GetTimeRandom() (timestamp, random int64, en string) {
	random = RandomNine()
	timestamp = GetNowTime()
	return timestamp, random, config.LANGUAGE
}
