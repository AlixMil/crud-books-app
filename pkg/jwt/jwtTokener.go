package jwt_package

import (
	"crud-books/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

type JwtTokener struct {
	signinKey []byte
	tokenTTL  int
}

func (j JwtTokener) GenerateTokens(userId string) (string, string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(j.tokenTTL)).Unix()

	t, err := token.SignedString(j.signinKey)
	if err != nil {
		return "", "", fmt.Errorf("singing token failed, error: %w", err)
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	trClaims := refreshToken.Claims.(jwt.MapClaims)
	trClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString(j.signinKey)
	if err != nil {
		return "", "", fmt.Errorf("refresh token signing failed, error: %w", err)
	}

	return t, rt, nil
}

func New(cfg config.Config) *JwtTokener {
	return &JwtTokener{
		signinKey: []byte(cfg.JWTSecret),
		tokenTTL:  cfg.JWTTokenTTL,
	}
}
