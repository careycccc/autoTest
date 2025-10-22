package utils

import (
	"autoTest/store/config"
	"autoTest/store/logger"
	"bufio"
	"crypto/md5"
	"encoding/binary"
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

// ReadYAMLByLine 逐行读取YAML文件，返回字符串列表
func ReadYAMLByLine(filePath string) ([]string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 创建字符串切片存储每一行
	var lines []string
	scanner := bufio.NewScanner(file)

	// 逐行读取
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// 检查扫描过程中是否出现错误
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
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
func GenerateRandomInt(min, max int64) (float64, error) {
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
	f, _ := randomInt.Float64()
	return f, nil
}

// 随机生成一个数字
func RandmoNumber(number int) int64 {
	// 初始化随机数种子
	var b [8]byte // 一个 int64 需要 8 个字节
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err) // 在实际应用中，应该优雅地处理错误
	}
	return int64(binary.BigEndian.Uint64(b[:]))
}

// 常见的邮箱域名列表
var domains = []string{
	"gmail.com", "yahoo.com", "hotmail.com", "outlook.com",
	"163.com", "126.com", "qq.com", "sina.com",
	"foxmail.com", "sohu.com", "139.com", "189.com",
	"aliyun.com", "protonmail.com", "icloud.com",
	"aol.com", "zoho.com", "mail.com", "inbox.com",
}

// 随机生成邮箱地址
func GenerateRandomEmail() string {
	// 生成随机用户名长度 (6-12个字符)
	usernameLen := 6 + RandInt(0, 6)
	username := GenerateRandomString(usernameLen)

	// 随机选择域名
	domain := domains[RandInt(0, len(domains)-1)]

	return username + "@" + domain
}

// 生成随机字符串
func GenerateRandomString(length int) string {
	// 可用的字符集
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := strings.Builder{}

	for i := 0; i < length; i++ {
		idx := RandInt(0, len(chars)-1)
		result.WriteByte(chars[idx])
	}

	return result.String()
}

// 生成指定范围的随机整数
func RandInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return min + int(n.Int64())
}

// GenerateBankCard 优化版：批量读取，解决并发阻塞，生成ifsc
func GenerateBankCard(length int, prefix ...string) (string, error) {
	if length < 10 || length > 19 {
		return "", fmt.Errorf("长度必须10-19位")
	}

	cardPrefix := "4"
	if len(prefix) > 0 {
		cardPrefix = prefix[0]
	}
	if len(cardPrefix) >= length {
		return "", fmt.Errorf("前缀过长")
	}

	// ⭐ 关键优化：一次性读取所有需要的字节
	neededBytes := length - len(cardPrefix) - 1
	if neededBytes <= 0 {
		return cardPrefix, nil
	}

	buffer := make([]byte, neededBytes)
	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("随机数生成失败: %v", err)
	}

	card := cardPrefix
	for _, b := range buffer {
		card += fmt.Sprintf("%d", int(b)%10)
	}

	checkDigit := calculateLuhnCheckDigit(card)
	return card + fmt.Sprintf("%d", checkDigit), nil
}

func calculateLuhnCheckDigit(partialCard string) int {
	sum := 0
	length := len(partialCard)
	for i := 0; i < length; i++ {
		digit := int(partialCard[length-1-i] - '0')
		if (length%2 == 0 && i%2 == 1) || (length%2 == 1 && i%2 == 0) {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return (10 - sum%10) % 10
}

// IFSCBankCode 常见银行代码（确保均为4位大写字母）
var IFSCBankCode = []string{
	"SBIN",  // State Bank of India
	"HDFC",  // HDFC Bank
	"ICIC",  // ICICI Bank
	"AXIS",  // Axis Bank
	"KKBK",  // Kotak Mahindra Bank
	"PNB",   // Punjab National Bank
	"BARB",  // Bank of Baroda
	"CANB",  // Canara Bank
	"UNION", // Union Bank of India
	"INDB",  // IndusInd Bank
	"KARB",  // Karnataka Bank
	"CNRB",  // Canara Bank (alternate)
	"SCBL",  // Standard Chartered Bank
	"MAHB",  // Bank of Maharashtra
	"VIJB",  // Vijaya Bank
	"IOBA",  // Indian Overseas Bank
	"FDRL",  // Federal Bank
	"IBKL",  // IDBI Bank
	"UCOB",  // UCO Bank
}

// RandomIFSC 生成随机IFSC代码
func RandomIFSC() (string, error) {
	// 随机选择银行代码（前4位大写字母）
	bankCode := IFSCBankCode[randomIndex(len(IFSCBankCode))]
	if len(bankCode) != 4 {
		return "", fmt.Errorf("invalid bank code length: %s (length=%d)", bankCode, len(bankCode))
	}

	// 中间固定数字（第5位，必须是0）
	middleNum := "0"

	// 生成分行代码（后6位，大写字母+数字）
	branchCode, err := generateBranchCode()
	if err != nil {
		return "", fmt.Errorf("failed to generate branch code: %w", err)
	}
	if len(branchCode) != 6 {
		return "", fmt.Errorf("invalid branch code length: %s (length=%d)", branchCode, len(branchCode))
	}

	// 拼接IFSC代码
	ifsc := bankCode + middleNum + branchCode
	if len(ifsc) != 11 {
		return "", fmt.Errorf("invalid IFSC length: %s (length=%d)", ifsc, len(ifsc))
	}

	// 验证IFSC格式
	if !validateIFSC(ifsc) {
		return "", fmt.Errorf("generated IFSC code is invalid: %s", ifsc)
	}

	return ifsc, nil
}

// validateIFSC 验证IFSC格式
func validateIFSC(ifsc string) bool {
	if len(ifsc) != 11 {
		return false
	}
	// 检查前4位是否都是大写字母
	for _, c := range ifsc[:4] {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	// 检查第5位是否为0
	if ifsc[4] != '0' {
		return false
	}
	// 检查后6位是否为字母或数字
	for _, c := range ifsc[5:] {
		if (c < 'A' || c > 'Z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}

// randomIndex 使用crypto/rand生成随机索引
func randomIndex(max int) int {
	if max <= 0 {
		return 0
	}
	b := make([]byte, 1)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // 在生产环境中应妥善处理
	}
	return int(b[0]) % max
}

// generateBranchCode 生成6位分行代码（大写字母+数字）
func generateBranchCode() (string, error) {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)

	for i := 0; i < 6; i++ {
		idx := randomIndex(len(chars))
		result[i] = chars[idx]
	}
	return string(result), nil
}

// 生成usdt的地址的函数
// base58Alphabet 定义了Base58字符集
const base58Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// GenerateTRONAddress 生成一个符合TRON USDT地址格式的随机地址
func GenerateTRONAddress() (string, error) {
	// TRON地址固定长度34，包含前缀'T'
	length := 34
	// 预分配字节切片，减去前缀'T'
	bytes := make([]byte, length-1)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// 构建地址
	var address strings.Builder
	address.WriteString("T") // 前缀固定为'T'

	// 将随机字节映射到Base58字符集
	for i := 0; i < length-1; i++ {
		// 取模确保索引在base58Alphabet范围内
		idx := int(bytes[i]) % len(base58Alphabet)
		address.WriteByte(base58Alphabet[idx])
	}

	return address.String(), nil
}

// GenerateUPIFormat 生成 UPI 格式字符串，电话号码和银行名称均为随机
// 返回格式: 10位随机电话号码@随机4-8位银行名称
func GenerateUPIFormat() (string, error) {
	// 生成10位随机电话号码
	phoneNumber := ""
	for i := 0; i < 10; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("生成随机电话号码失败: %v", err)
		}
		phoneNumber += fmt.Sprintf("%d", num)
	}

	// 生成4-8位随机银行名称
	length, err := rand.Int(rand.Reader, big.NewInt(5))
	if err != nil {
		return "", fmt.Errorf("生成随机银行名称长度失败: %v", err)
	}
	bankNameLength := 4 + int(length.Int64()) // 随机长度4到8

	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bankName := ""
	for i := 0; i < bankNameLength; i++ {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		if err != nil {
			return "", fmt.Errorf("生成随机银行名称字符失败: %v", err)
		}
		bankName += string(chars[idx.Int64()])
	}

	// 构造 UPI 格式
	upi := fmt.Sprintf("%s@%s", phoneNumber, strings.ToLower(bankName))

	return upi, nil
}

// 生成一个长度为n的随机数字字符串
func GenerateNumberString(n int) (string, error) {
	// 如果n <= 0，返回空字符串和错误
	if n <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	// 创建一个字节切片来存储结果
	result := make([]byte, n)

	// 填充随机数字（0-9）
	for i := 0; i < n; i++ {
		// 生成0-9的随机数
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %v", err)
		}
		// 将数字转换为字符（'0' 到 '9'）
		result[i] = byte('0' + num.Int64())
	}

	// 将字节切片转换为字符串
	return string(result), nil
}
