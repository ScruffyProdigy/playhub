package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
)

var loginCodePattern = regexp.MustCompile(`^\d{6}$`)

func generateLoginCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", fmt.Errorf("auth: generate login code: %w", err)
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func hashLoginCode(code string) string {
	sum := sha256.Sum256([]byte(code))
	return hex.EncodeToString(sum[:])
}

func verifyLoginCode(code, codeHash string) bool {
	if !loginCodePattern.MatchString(code) || codeHash == "" {
		return false
	}
	expected := hashLoginCode(code)
	if len(expected) != len(codeHash) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected), []byte(codeHash)) == 1
}
