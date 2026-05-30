package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const defaultJWTKID = "lobby-dev"

// Signer creates and verifies session JWTs.
type Signer struct {
	kid        string
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
	publicX    string
}

// LoadSignerFromEnv loads signing keys from environment or generates dev keys.
func LoadSignerFromEnv() (*Signer, error) {
	kid := os.Getenv("JWKS_KID")
	if kid == "" {
		kid = defaultJWTKID
	}

	if pemData := strings.TrimSpace(os.Getenv("JWT_PRIVATE_KEY_PEM")); pemData != "" {
		return newSignerFromPEM(kid, pemData)
	}

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("auth: generate dev signing key: %w", err)
	}

	signer := &Signer{
		kid:        kid,
		privateKey: priv,
		publicKey:  pub,
		publicX:    base64.RawURLEncoding.EncodeToString(pub),
	}

	fmt.Fprintf(os.Stderr, "auth: generated development JWT signing key (kid=%s)\n", kid)
	return signer, nil
}

func newSignerFromPEM(kid, pemData string) (*Signer, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, errors.New("auth: invalid JWT private key PEM")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("auth: parse JWT private key: %w", err)
	}

	privateKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("auth: JWT private key must be Ed25519")
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)
	return &Signer{
		kid:        kid,
		privateKey: privateKey,
		publicKey:  publicKey,
		publicX:    base64.RawURLEncoding.EncodeToString(publicKey),
	}, nil
}

// PublicJWK returns the public JWKS payload for this signer.
func (s *Signer) PublicJWK() map[string]string {
	return map[string]string{
		"kty": "OKP",
		"crv": "Ed25519",
		"alg": "EdDSA",
		"kid": s.kid,
		"x":   s.publicX,
	}
}

// JWKSHandler serves the public JWKS document.
func (s *Signer) JWKSHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{s.PublicJWK()}})
	}
}

// SignUserToken issues a session JWT for the given user ID.
func (s *Signer) SignUserToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"iat": now.Unix(),
		"exp": now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = s.kid
	return token.SignedString(s.privateKey)
}

// VerifyUserToken validates a JWT and returns the user ID.
func (s *Signer) VerifyUserToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodEdDSA {
			return nil, fmt.Errorf("auth: unexpected signing method %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("auth: invalid token claims")
	}

	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		return uuid.Nil, errors.New("auth: token missing subject")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, fmt.Errorf("auth: invalid subject: %w", err)
	}

	return userID, nil
}
