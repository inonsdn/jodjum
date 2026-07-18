package token

import (
	"errors"
	"server/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenParseResult struct {
	UserId    uuid.UUID
	SessionId uuid.UUID
}

type Claims struct {
	SessionId string `json:"sid"`
	jwt.RegisteredClaims
}

type TokenService struct {
	tokenRepo TokenRepo
	secret    []byte
	issuer    string
	expire    time.Duration
}

func NewService(cfg *config.TokenConfig, tokenRepo TokenRepo) *TokenService {
	return &TokenService{
		tokenRepo: tokenRepo,
		secret:    []byte(cfg.Secret),
		issuer:    cfg.Issuer,
		expire:    time.Duration(cfg.Expire) * time.Minute,
	}
}

func (t *TokenService) Generate(userId uuid.UUID, sessionId uuid.UUID) (string, error) {
	nowTime := time.Now()
	claims := Claims{
		SessionId: sessionId.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    t.issuer,
			Subject:   userId.String(),
			ExpiresAt: jwt.NewNumericDate(nowTime.Add(t.expire)),
			NotBefore: jwt.NewNumericDate(nowTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(t.secret)

	return signedToken, err
}

func (t *TokenService) Parse(tokenString string) (TokenParseResult, error) {
	parseResult := TokenParseResult{}
	claims := Claims{}
	jwtToken, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (any, error) {
			return t.secret, nil
		},
		jwt.WithIssuer(t.issuer),
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
	)

	if err != nil {
		return parseResult, err
	}
	if !jwtToken.Valid {
		return parseResult, errors.New("Invalid token, please generate new token")
	}
	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return parseResult, errors.New("invalid token subject")
	}

	sessionId, err := uuid.Parse(claims.SessionId)
	if err != nil {
		return parseResult, errors.New("invalid token session id")
	}

	if claims.SessionId == "" {
		return parseResult, errors.New("session id is missing")
	}

	parseResult.UserId = userId
	parseResult.SessionId = sessionId

	return parseResult, nil
}
