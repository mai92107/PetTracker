package mqttUtils

import (
	mqttUtil "batchLog/0.core/mqtt"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

func Encrypt(data, secretKey string) (string, error) {
	key := []byte(secretKey)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		fmt.Print(key)
		return "", errors.New("invalid AES key size (must be 16, 24, or 32 bytes)")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	cipherText := gcm.Seal(nonce, nonce, []byte(data), nil)
	res := base64.StdEncoding.EncodeToString(cipherText)
	err = mqttUtil.PubMsgToTopic("encrypted", res)
	return res, err
}

func Decrypt(data, secretKey string) (string, error) {
	key := []byte(secretKey)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("invalid AES key size (must be 16, 24, or 32 bytes)")
	}

	dataString, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	if len(dataString) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherText := dataString[:gcm.NonceSize()], dataString[gcm.NonceSize():]

	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	res := string(plainText)
	err = mqttUtil.PubMsgToTopic("decrypted", res)
	return res, err
}
