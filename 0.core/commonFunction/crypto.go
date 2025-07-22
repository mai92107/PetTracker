package common

import (
	"batchLog/0.core/global"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
)

func BcryptHash(text string)(string,error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash %s, error: %+v", text, err)
    }
	return string(hashedPassword), nil
}

func BcryptCompare(hashedData, data string)bool{
	err := bcrypt.CompareHashAndPassword([]byte(hashedData),[]byte(data))
	return err == nil
}

func Encryption(text string) (string, error) {
	block, err := aes.NewCipher([]byte(global.ConfigSetting.CryptoSecretKey))
	if err != nil {
		return "", err
	}

	// CBC 模式需要 IV（初始向量），AES block 大小是 16 bytes
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// 補齊 PKCS7 padding
	paddedText := pad([]byte(text), aes.BlockSize)

	ciphertext := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedText)

	// 最後輸出為：IV + 加密內容，再 base64 encode
	result := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(result), nil
}

func Decryption(encrypted string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	if len(raw) < aes.BlockSize {
		return "", fmt.Errorf("cipher too short")
	}

	iv := raw[:aes.BlockSize]
	ciphertext := raw[aes.BlockSize:]

	block, err := aes.NewCipher([]byte(global.ConfigSetting.CryptoSecretKey))
	if err != nil {
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	// 去除 padding
	plaintext, err := unpad(ciphertext)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func unpad(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding size")
	}
	padding := int(src[length-1])
	if padding > length {
		return nil, fmt.Errorf("invalid padding")
	}
	return src[:length-padding], nil
}