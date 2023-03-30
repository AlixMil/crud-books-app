package jwt_package

import (
	"crud-books/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaims struct {
	Token string `json:"token"`
	jwt.RegisteredClaims
}

type JwtTokener struct {
	signinKey []byte
	tokenTTL  time.Duration
}

func New(cfg config.Config) *JwtTokener {
	return &JwtTokener{
		signinKey: []byte(cfg.JWTSecret),
		tokenTTL:  time.Second * time.Duration(cfg.JWTTokenTTL),
	}
}

func (j JwtTokener) GenerateToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Second * j.tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   userId,
	})

	tokenString, err := token.SignedString(j.signinKey)
	if err != nil {
		return "", fmt.Errorf("failed while signing token in generate token func, error: %w", err)
	}

	return tokenString, nil
}

func (j JwtTokener) ParseToken(token string) (string, error) {
	acceptedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.signinKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed parse jwt in parsetoken func, error: %w", err)
	}

	if !acceptedToken.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := acceptedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("invalid claims - subject")
	}

	return subject, nil
}
