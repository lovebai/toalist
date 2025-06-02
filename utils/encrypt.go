package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
)

var key = "toalistsecretkey1234567890123450"

func MD5Encrypt(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str + "toalist123321")) // 加盐 直接写死了
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}

func AESEncrypt(plaintext string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("invalid key length")
	}
	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nil, nonce, []byte(plaintext), nil)
	result := append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(result), nil
}

func AESDecrypt(encoded string) (string, error) {
	if len(key) != 32 {
		return "", errors.New("invalid key length")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}
	if len(data) < 12 {
		return "", errors.New("invalid ciphertext")
	}
	nonce := data[:12]
	ciphertext := data[12:]
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
