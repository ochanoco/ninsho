package ninsho

import (
	"crypto/rand"
	"fmt"
)

func secureRandom(b int) (string, error) {
	k := make([]byte, b)
	if _, err := rand.Read(k); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", k), nil
}
