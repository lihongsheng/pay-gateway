package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

const (
	saltToken = "c368a753dbf1784d39359db751fcfac9"
	saltKey   = "be142eb8e109ea0fb4a3793ec11cb4ab"
)

func removeAllSpaces(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func hmacSha256Hex(key, input string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil)) // 64 hex 字符
}

func GenerateTokenAndEncodingKey(input string) (string, string) {
	input = removeAllSpaces(input)
	if len(input) < 15 {
		input = input + strings.Repeat("0", 15-len(input))
	}
	hexToken := hmacSha256Hex(saltToken, input)
	token := hexToken[:32]

	hexKey := hmacSha256Hex(saltKey, input)
	encodingKey := hexKey[:43]

	return token, encodingKey
}

// HashCrc 实现与JS版本一致的CRC32哈希算法
// str: 输入字符串
// length: 输出结果的固定长度
// return: 固定长度的纯数字字符串，不足补0，超长截断
func HashCrc(str string, length int) string {
	crc := uint32(0xFFFFFFFF)
	const polynomial = uint32(0xEDB88320)
	for i := 0; i < len(str); i++ {
		byteVal := uint32(str[i])
		crc ^= byteVal
		for j := 0; j < 8; j++ {
			if crc&1 == 1 {
				crc = (crc >> 1) ^ polynomial
			} else {
				crc = crc >> 1
			}
		}
	}
	crc ^= 0xFFFFFFFF
	crcStr := strconv.FormatUint(uint64(crc), 10)
	if len(crcStr) < length {
		crcStr = strings.Repeat("0", length-len(crcStr)) + crcStr
	} else if len(crcStr) > length {
		crcStr = crcStr[:length]
	}
	return crcStr
}

func GetRequestBody(request *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("read request body err: %v", err)
	}

	_ = request.Body.Close()
	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body, nil
}
