package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"
)

func hashSecret(value string) string {
	value = strings.TrimSpace(value)
	pepper := strings.TrimSpace(os.Getenv("MAGIC_LINK_PEPPER"))
	sum := sha256.Sum256([]byte(pepper + value))
	return hex.EncodeToString(sum[:])
}

func hashMagicLinkToken(token string) string {
	return hashSecret(token)
}
