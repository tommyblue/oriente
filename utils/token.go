package utils

import (
	"crypto/rand"
	"fmt"
)

func IDGenerator() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
