package main

import (
	"crypto/sha256"
	"fmt"
)

func coder(link string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(link))
	hashBytes := hash.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)
	return hashString[:6], nil 
}
