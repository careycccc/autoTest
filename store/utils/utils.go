package utils

import (
	"autoTest/store/config"
	"autoTest/store/logger"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"

	// "math/rand"
	"crypto/rand"
	"os"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Md5Info 计算 MD5 哈希值
func Md5Info(data string, uppercase bool) string {
	hash := md5.New()
	hash.Write([]byte(data))
	result := hex.EncodeToString(hash.Sum(nil))
	if uppercase {
		return strings.ToUpper(result)
	}
	return strings.ToLower(result)
}

// GetSignature 传入一个map 生成签名
func GetSignature(body map[string]interface{}, verifyPwd *string) string {
	// 过滤字段
	filteredObj := make(map[string]interface{})
	keys := make([]string, 0, len(body))
	for key := range body {
		keys = append(keys, key)
	}
	sort.Strings(keys) // 按键排序

	for _, key := range keys {
		value := body[key]
		// 检查 value 不为 nil 且不为空字符串，且 key 不在排除列表中，且 value 不是数组
		if value != nil && value != "" && key != "signature" && key != "timestamp" && key != "track" {
			// 确保 value 不是切片（相当于 Python 的 list）
			if _, ok := value.([]interface{}); !ok {
				filteredObj[key] = value
			}
		}
	}

	// 转换为 JSON 字符串
	jsonData, err := json.Marshal(filteredObj)
	if err != nil {
		return "" // 错误处理，可根据需求调整
	}

	encoder := string(jsonData)
	if verifyPwd != nil {
		encoder += *verifyPwd
	}
	// 计算 MD5
	return Md5Info(encoder, true)
}

// 传入一个结构体获取返回签名
func GetSignature2(body any, verifyPwd *string) string {
	// 过滤字段并转换为 map 以便排序
	filteredObj := make(map[string]interface{})
	excludeKeys := map[string]bool{"signature": true, "timestamp": true, "track": true}

	// 使用反射获取结构体字段
	val := reflect.ValueOf(body)
	typ := reflect.TypeOf(body)
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i).Interface()
		jsonTag := field.Tag.Get("json")
		// 提取 JSON 字段名（忽略 ",omitempty" 部分）
		key := jsonTag
		if idx := len(jsonTag); idx > 9 && jsonTag[idx-9:] == ",omitempty" {
			key = jsonTag[:idx-9]
		}

		// 过滤条件：非空值、不在排除列表中
		if !excludeKeys[key] && !isEmpty(value) {
			filteredObj[key] = value
		}
	}

	// 按键排序
	keys := make([]string, 0, len(filteredObj))
	for key := range filteredObj {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// 转换为 JSON 字符串
	jsonData, err := json.Marshal(filteredObj)
	if err != nil {
		return "" // 错误处理
	}

	encoder := string(jsonData)
	if verifyPwd != nil {
		encoder += *verifyPwd
	}

	// 计算 MD5
	return Md5Info(encoder, true)
}

// isEmpty 检查值是否为空（nil、零值或空字符串）
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	}
	return false
}

// 用来解析请求
func Unmarshal(strResbody string) (result map[string]interface{}) {
	//是一个字符串
	error := json.Unmarshal([]byte(strResbody), &result)
	if error != nil {
		log.Fatalf("解析响应失败~~:%v", error)
	}
	return

}

// 读取yaml
// ReadYAML 从指定路径读取 YAML 文件并解析到结构体
func ReadYAML(filePath string, result interface{}) error {
	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取 YAML 文件失败: %w", err)
	}

	// 解析 YAML 数据到结构体
	if err := yaml.Unmarshal(data, result); err != nil {
		return fmt.Errorf("解析 YAML 数据失败: %w", err)
	}

	return err
}

// WriteYAML 将数据追加写入到 ./yaml/subUser.yaml 文件，每行一个数据
func WriteYAML(data ...interface{}) error {
	file, err := os.OpenFile(config.SUBUSERYAML, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.LogError("文件创建或者打开失败", err)
		return err
	}
	defer file.Close()

	for _, v := range data {
		// 将数据转换为字符串并写入，带换行符
		_, err := file.WriteString(fmt.Sprintf("%v\n", v))
		if err != nil {
			logger.LogError("文件写入失败", err)
			return err
		}
	}
	return nil
}

// 生成随机浏览器指纹
func GenerateCryptoRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	// fmt.Println("本次指纹", string(bytes))
	return string(bytes)
}

// 随机生成 min max的数据
func GenerateRandomInt(min, max int64) (int64, error) {
	if min > max {
		return 0, fmt.Errorf("min must be less than or equal to max")
	}

	// 生成一个大于等于0且小于max-min的随机数
	randomInt, err := rand.Int(rand.Reader, big.NewInt(max-min))
	if err != nil {
		return 0, err
	}

	// 将随机数加上最小值，得到最终的随机数范围在[min, max]之间
	randomInt.Add(randomInt, big.NewInt(min))
	return randomInt.Int64(), nil
}
