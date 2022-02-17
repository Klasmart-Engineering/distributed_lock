package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func RandNum() string {
	d := make([]byte, 8)
	rand.Read(d)
	return hex.EncodeToString(d)
}
