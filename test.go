package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateRandomKey() string {
	key := make([]byte, 32)
	rand.Read(key)
	return base64.StdEncoding.EncodeToString(key)
}

func main() {
	key := generateRandomKey()
	fmt.Println(key)
}
