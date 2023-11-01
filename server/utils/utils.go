package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/btcsuite/btcd/btcutil/base58"
)

// ExitOnError 退出程序
func ExitOnError(message string, err error) {
	if err != nil {
		log.Printf("[%s] - %s", message, err)
		os.Exit(-1)
	}
}

// PrintOnError 打印错误
func PrintOnError(message string, err error) {
	if err != nil {
		log.Printf("[%s] - %s", message, err)
	}
}

// RaiseError 返回错误
func RaiseError(message string) error {
	if !EmptyString(message) {
		return fmt.Errorf(message)
	}
	return nil
}

//CompareArrays 比较两个数组值是否相同
func CompareArrays(arr1 []string, arr2 []string) bool {
	if len(arr1) == 0 || len(arr2) == 0 {
		return false
	}

	result := reflect.DeepEqual(arr1, arr2)
	if result == false {
		for _, item := range arr1 {
			if !InArray(item, arr2) {
				return false
			}
		}
	}
	return true
}

func ContainsString(src []string, dest string) bool {
	for _, item := range src {
		if item == dest {
			return true
		}
	}
	return false
}

// DifferenceStringsArr 取前者src与后者dest两个字符串列表的差集
func DifferenceStringsArr(src []string, dest []string) []string {
	res := make([]string, 0)
	for _, item := range src {
		if !ContainsString(dest, item) {
			res = append(res, item)
		}
	}
	return res
}

// EmptyString 判断字符串是否为空
func EmptyString(str string) bool {
	str = strings.TrimSpace(str)
	return strings.EqualFold(str, "")
}

// UserAgentIpHash 生成用户代理和IP的哈希值
func UserAgentIpHash(useragent string, ip string) string {
	input := fmt.Sprintf("%s-%s-%s-%d", useragent, ip, time.Now().String(), rand.Int())
	data, _ := Sha256Of(input)
	str := Base58Encode(data)
	return str[:10]
}

// Sha256Of 计算字符串的哈希值
func Sha256Of(input string) ([]byte, error) {
	algorithm := sha256.New()
	_, err := algorithm.Write([]byte(strings.TrimSpace(input)))
	if err != nil {
		return nil, err
	}
	return algorithm.Sum(nil), nil
}

// Base58Encode base58编码
func Base58Encode(data []byte) string {
	return base58.Encode(data)
}

// PasswordBase58Hash 密码加密
func PasswordBase58Hash(password string) (string, error) {
	data, err := Sha256Of(password)
	if err != nil {
		return "", err
	}
	return base58.Encode(data), nil
}

// SplitStrIdsToInt 分隔字符串ids
// strIds like "1,2,3"
// intIds []int like [1,2,3]
func SplitStrIdsToInt(strIds string, sep string) []int {
	var intIds []int
	arr := strings.Split(strIds, sep)
	if len(arr) == 0 {
		return intIds
	}
	for _, id := range arr {
		nid, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		intIds = append(intIds, nid)
	}
	return intIds
}

//InArray 查找字符是否在数组中
func InArray(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}

	return false
}

//HashPassword 计算md5
func HashPassword(password, salt string) string {
	h := md5.New()
	io.WriteString(h, password+salt)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//ValidateIP 验证IP格式，目前支持IP和IP段
func ValidateIP(ipStr string) bool {
	ipStr = strings.TrimSpace(ipStr)
	if IsIpv4(ipStr) {
		return true
	}

	//ip段 192.168.0.1/24
	if IsCIDR(ipStr) {
		return true
	}

	// 匹配起始ip,如: 192.168.0.2-192.168.0.252
	if IsRange(ipStr) {
		return true
	}

	return false
}

// IsIpv4 正则匹配ipv4
func IsIpv4(ipStr string) bool {
	pattern := "((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}"
	matched, _ := regexp.MatchString(pattern, ipStr)
	if matched {
		ipAddr := net.ParseIP(ipStr)
		if ipAddr != nil {
			return true
		}
	}
	return false
}

// IsCIDR 验证ip段
func IsCIDR(ipStr string) bool {
	isCIDR := strings.Contains(ipStr, "/")
	if isCIDR {
		_, _, err := net.ParseCIDR(ipStr)
		if err == nil {
			return true
		}
	}
	return false
}

func IsRange(ipStr string) bool {
	isRange := strings.Contains(ipStr, "-")
	if isRange {
		ss := strings.Split(ipStr, "-")
		if IsIpv4(ss[0]) && IsIpv4(ss[1]) {
			return true
		}
	}
	return false
}

//IpRangeContains IP段检查
func IpRangeContains(ip string, cidr string) bool {
	if IsIpv4(cidr) {
		if ip == cidr {
			return true
		}
	}
	//ip段 192.168.0.1/24
	if IsCIDR(cidr) {
		// Parse the IP range
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			return false
		}
		// Parse the IP to check
		ipAddr := net.ParseIP(ip)
		return ipNet.Contains(ipAddr)
	}

	// 匹配起始ip,如: 192.168.0.2-192.168.0.252
	if IsRange(cidr) {
		ss := strings.Split(cidr, "-")
		return IpBetween(ip, ss[0], ss[1])
	}

	return false
}

func TrimString(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\t", "")
	str = strings.ReplaceAll(str, "\n", "")
	return str
}

func AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(fmt.Sprintf("key 长度必须 16/24/32长度: %s", err.Error()))
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	//使用RawURLEncoding 不要使用StdEncoding
	//不要使用StdEncoding  放在url参数中回导致错误
	return base64.RawURLEncoding.EncodeToString(cryted)

}

func AesDecrypt(cryted string, key string) string {
	//使用RawURLEncoding 不要使用StdEncoding
	//不要使用StdEncoding  放在url参数中回导致错误
	crytedByte, _ := base64.RawURLEncoding.DecodeString(cryted)
	k := []byte(key)

	// 分组秘钥
	block, err := aes.NewCipher(k)
	if err != nil {
		panic(fmt.Sprintf("key 长度必须 16/24/32长度: %s", err.Error()))
	}
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//RandSimplePassword 随机生成简单密码
func RandSimplePassword() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}

//IpBetween 判断ip是否在某个区间
func IpBetween(ip, start, end string) bool {
	from := net.ParseIP(start)
	to := net.ParseIP(end)
	ipStr := net.ParseIP(ip)
	if from == nil || to == nil || ipStr == nil {
		fmt.Println("An ip input is nil") // or return an error!?
		return false
	}

	from16 := from.To16()
	to16 := to.To16()
	ipStr16 := ipStr.To16()
	if from16 == nil || to16 == nil || ipStr16 == nil {
		fmt.Println("An ip did not convert to a 16 byte") // or return an error!?
		return false
	}

	if bytes.Compare(ipStr16, from16) >= 0 && bytes.Compare(ipStr16, to16) <= 0 {
		return true
	}
	return false
}
