package utils

import (
	"crypto/sha256"
	"fmt"
)

func HashPassword(rawPassword string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(rawPassword)))
}
