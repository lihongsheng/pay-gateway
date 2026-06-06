package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base32"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// 入参：keySize 密钥长度，只能传 16(AES-128)、24(AES-192)、32(AES-256)
// 返回：[]byte 密钥，error 错误信息
func GenerateAesKey(keySize int) ([]byte, error) {
	// 校验密钥长度是否符合AES规范，非规范长度直接返回错误
	if keySize != 16 && keySize != 24 && keySize != 32 {
		return nil, fmt.Errorf("AES密钥长度非法，仅支持：16(AES-128)、24(AES-192)、32(AES-256)，你传入的是：%d", keySize)
	}

	// 初始化指定长度的字节数组
	key := make([]byte, keySize)
	// 用加密安全的随机数填充数组，crypto/rand.Reader 是真随机数据源
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, fmt.Errorf("生成AES密钥失败: %w", err)
	}
	return key, nil
}

func AesKeyToHexStr(key []byte) string {
	return hex.EncodeToString(key)
}

func GenAesHexStr(keySize int) (string, error) {
	key, err := GenerateAesKey(keySize)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func GenAesBase64Str(keySize int) (string, error) {
	key, err := GenerateAesKey(keySize)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

func AesKeyToBase64Str(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}

func HexStrToAesKey(hexStr string) ([]byte, error) {
	return hex.DecodeString(hexStr)
}

func Base64StrToAesKey(b64Str string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(b64Str)
}

var (
	ErrGenBlock  = errors.New("gen block fail")
	ErrGenGcm    = errors.New("gcm: invalid nonce size")
	ErrGenNonce  = errors.New("gen nonce fail")
	ErrBase64Dec = errors.New("base64 decode fail")
	ErrDataLen   = errors.New("data len error")
	ErrDecrypt   = errors.New("decrypt fail")
)

// AesGCMEncryptHex 模式加密（带认证，更安全）
func AesGCMEncryptHex(plainText, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", ErrGenBlock
	}
	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrGenGcm
	}
	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", ErrGenNonce
	}
	// 加密并生成认证标签
	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	return base32.StdEncoding.EncodeToString(cipherText), nil
}

// AesGCMDecryptHex 模式解密
func AesGCMDecryptHex(cipherTextHex string, key []byte) (string, error) {
	cipherText, err := base32.StdEncoding.DecodeString(cipherTextHex)
	if err != nil {
		return "", ErrBase64Dec
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", ErrGenBlock
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", ErrGenGcm
	}
	// 分离nonce和密文
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", ErrDataLen
	}
	nonce := cipherText[:nonceSize]
	cipherText = cipherText[nonceSize:]
	// 解密（同时验证认证标签）
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", ErrDecrypt
	}
	return string(plainText), nil
}
