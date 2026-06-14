package auth

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const oauthStateTokenType = "oauth_state"

type OAuthMode string

const (
	OAuthModeSignIn OAuthMode = "signin"
	OAuthModeLink   OAuthMode = "link"
)

type OAuthState struct {
	Provider     string
	Mode         OAuthMode
	UserID       uuid.UUID
	ConfirmMerge bool
}

func (s *Signer) SignOAuthState(state OAuthState, ttl time.Duration) (string, error) {
	if strings.TrimSpace(state.Provider) == "" {
		return "", errors.New("auth: oauth provider is required")
	}
	if state.Mode != OAuthModeSignIn && state.Mode != OAuthModeLink {
		return "", errors.New("auth: invalid oauth mode")
	}
	if state.Mode == OAuthModeLink && state.UserID == uuid.Nil {
		return "", errors.New("auth: link mode requires user id")
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"iss": LobbyIssuer(),
		"aud": sessionTokenAudience(),
		"typ": oauthStateTokenType,
		"jti": uuid.NewString(),
		"iat": now.Unix(),
		"nbf": now.Unix(),
		"exp": now.Add(ttl).Unix(),
		"prv": strings.ToLower(strings.TrimSpace(state.Provider)),
		"mod": string(state.Mode),
	}
	if state.UserID != uuid.Nil {
		claims["uid"] = state.UserID.String()
	}
	if state.ConfirmMerge {
		claims["cm"] = true
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = s.kid
	return token.SignedString(s.privateKey)
}

func (s *Signer) VerifyOAuthState(tokenString string) (OAuthState, error) {
	issuer := LobbyIssuer()
	audience := sessionTokenAudience()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodEdDSA {
			return nil, fmt.Errorf("auth: unexpected signing method %v", token.Header["alg"])
		}
		return s.publicKey, nil
	},
		jwt.WithIssuer(issuer),
		jwt.WithAudience(audience),
		jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
	)
	if err != nil {
		return OAuthState{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return OAuthState{}, errors.New("auth: invalid oauth state claims")
	}
	if typ, _ := claims["typ"].(string); typ != oauthStateTokenType {
		return OAuthState{}, errors.New("auth: token is not an oauth state token")
	}

	provider, _ := claims["prv"].(string)
	mode, _ := claims["mod"].(string)
	state := OAuthState{
		Provider: strings.ToLower(strings.TrimSpace(provider)),
		Mode:     OAuthMode(mode),
	}
	if uid, _ := claims["uid"].(string); uid != "" {
		parsed, err := uuid.Parse(uid)
		if err != nil {
			return OAuthState{}, fmt.Errorf("auth: invalid oauth state user id: %w", err)
		}
		state.UserID = parsed
	}
	if cm, ok := claims["cm"].(bool); ok {
		state.ConfirmMerge = cm
	}
	return state, nil
}
