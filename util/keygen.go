package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func GenerateHMACKey() ([]byte, error) {
	key := make([]byte, 32) // 32字节即256位
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func GenerateHMACKeyString() (string, error) {
	key, err := GenerateHMACKey()
	if err != nil {
		return "", err
	}
	keyStr := base64.StdEncoding.EncodeToString(key)
	return keyStr, nil
}

func main() {
	keyStr, err := GenerateHMACKeyString()
	if err != nil {
		fmt.Printf("Failed to generate HMAC key: %s\n", err)
		return
	}
	fmt.Printf("Generated HMAC key: %s\n", keyStr)
}
