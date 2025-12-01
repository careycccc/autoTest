package checksign

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// Md5Info 完全和你原版一模一样
func Md5Info(data string, uppercase bool) string {
	hash := md5.New()
	hash.Write([]byte(data))
	result := hex.EncodeToString(hash.Sum(nil))
	if uppercase {
		return strings.ToUpper(result)
	}
	return strings.ToLower(result)
}

// GenerateSignatureFromJSON
// 输入：任意 JSON（string 或 []byte）
// verifyPwd：你的密钥字符串，传 nil 表示不拼接
// 输出：和大写 MD5 签名，完全匹配你们系统
func GenerateSignatureFromJSON(jsonInput interface{}, verifyPwd *string) (string, error) {
	// 1. 解析成 map
	var body map[string]interface{}
	switch v := jsonInput.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &body); err != nil {
			return "", err
		}
	case []byte:
		if err := json.Unmarshal(v, &body); err != nil {
			return "", err
		}
	case map[string]interface{}:
		body = v
	default:
		return "", fmt.Errorf("input must be json string/bytes or map")
	}

	// 2. 完全复制你的 GetSignature 逻辑
	filteredObj := make(map[string]interface{})
	keys := make([]string, 0, len(body))
	for key := range body {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := body[key]
		if value != nil && value != "" &&
			key != "signature" && key != "timestamp" && key != "track" {
			if _, ok := value.([]interface{}); !ok {
				filteredObj[key] = value
			}
		}
	}

	// jsonData, err := json.Marshal(filteredObj)
	// if err != nil {
	// 	return "", err
	// }

	// encoder := string(jsonData)
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false) // 关键！关闭 HTML 转义
	if err := enc.Encode(filteredObj); err != nil {
		return "", err
	}
	jsonData := buf.Bytes()
	if jsonData[len(jsonData)-1] == '\n' { // Encode 会多加一个换行
		jsonData = jsonData[:len(jsonData)-1]
	}
	encoder := string(jsonData)
	fmt.Println("【调试】最终参与签名的字符串是：")
	fmt.Println(encoder)
	if verifyPwd != nil {
		encoder += *verifyPwd
	}

	// 关键：第二个参数 true → 大写！
	return Md5Info(encoder, true), nil
}

// 运行检测签名
func RunCheckSign() {
	jsonStr := `{
  "browserId": "6aed08b5558eacb5d97f7b8441ac99f4",
  "language": "en",
  "loginType": "Mobile",
  "password": "qwer1234",
  "random": 203305427707,
  "track": "",
  "userName": "911201199711"
}`

	sig, _ := GenerateSignatureFromJSON(jsonStr, nil)
	fmt.Println(sig)
	// 输出：EC77D73D8B868D0A7650CD0B32D16895   ← 完全一致！一模一样！
}
