package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenAes(t *testing.T) {

	// ========== 1. 生成 AES-128 密钥 (16字节，通用业务场景) ==========
	key128, err := GenerateAesKey(16)
	if err != nil {
		fmt.Printf("生成AES-128密钥失败: %v\n", err)
	} else {
		fmt.Println("✅ AES-128 密钥(二进制):", key128)
		fmt.Println("✅ AES-128 密钥(16进制):", AesKeyToHexStr(key128))
		fmt.Println("✅ AES-128 密钥(Base64):", AesKeyToBase64Str(key128))
		fmt.Println("--------------------------------------------------------")
	}

	// ========== 2. 生成 AES-256 密钥 (32字节，推荐支付/金融场景，你的订单表首选) ==========
	key256, err := GenerateAesKey(32)
	if err != nil {
		fmt.Printf("生成AES-256密钥失败: %v\n", err)
	} else {
		fmt.Println("✅ AES-256 密钥(二进制):", key256)
		fmt.Println("✅ AES-256 密钥(16进制):", AesKeyToHexStr(key256))
		fmt.Println("✅ AES-256 密钥(Base64):", AesKeyToBase64Str(key256))
		fmt.Println("--------------------------------------------------------")
	}

	// ========== 3. 反向转换示例 (解密时使用) ==========
	hexKey := AesKeyToHexStr(key256)
	originKey, _ := HexStrToAesKey(hexKey)
	assert.Equal(t, string(originKey), string(key256))

}
