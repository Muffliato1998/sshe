package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
)

// MD5 哈希生成工具函数（用于生成固定长度的密钥）
func md5Hash(secretKey string) []byte {
	hash := md5.Sum([]byte(secretKey))
	return hash[:]
}

// EncryptAES AES 加密函数
func EncryptAES(plainText []byte, secretKey string) (string, error) {
	// 使用 MD5 生成 16 字节的密钥
	block, err := aes.NewCipher(md5Hash(secretKey))
	if err != nil {
		return "", fmt.Errorf("error: failed to create AES cipher: %v", err)
	}

	// 填充明文至 AES 块大小的倍数
	paddingText := pkcs7Padding(plainText, block.BlockSize())

	// 创建加密模式
	cipherText := make([]byte, len(paddingText))

	iv := md5Hash(secretKey)[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText, paddingText)

	return hex.EncodeToString(cipherText), nil
}

// DecryptAES AES 解密函数
func DecryptAES(cipherHex, key string) (string, error) {
	// 使用 MD5 生成 16 字节的密钥
	block, err := aes.NewCipher(md5Hash(key))
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}

	// 将密文从十六进制转换为字节
	cipherText, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode cipher text: %v", err)
	}

	// 验证密文长度是否合法
	if len(cipherText)%block.BlockSize() != 0 {
		return "", errors.New("invalid cipher text length")
	}

	// 创建解密模式
	plainText := make([]byte, len(cipherText))
	iv := md5Hash(key)[:aes.BlockSize] // 使用 key 的前 16 字节作为 IV
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainText, cipherText)

	// 去除填充
	plainText, err = pkcs7UnPadding(plainText)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

// PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// PKCS7 去除填充
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("data is empty")
	}

	padding := int(data[length-1])
	if padding > length || padding > aes.BlockSize {
		return nil, errors.New("invalid padding")
	}

	return data[:length-padding], nil
}
