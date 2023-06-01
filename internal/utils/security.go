package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
)

const (
	chars     string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	apiKeyLen int    = 32
)

func GetHash(payload, key string) string {
	hash := hmac.New(sha256.New, []byte(key))
	hash.Write([]byte(payload))
	dst := hash.Sum(nil)
	return hex.EncodeToString(dst)
}

func CheckPasswordHash(password, hash, key string) bool {
	newHash := GetHash(password, key)
	return hash == newHash
}

func GenerateAPIKey() string {
	var result []byte
	for num := 0; num < apiKeyLen; num++ {
		index := rand.Intn(len(chars) - 1)
		newChar := chars[index]
		result = append(result, newChar)
	}

	return string(result)
}
