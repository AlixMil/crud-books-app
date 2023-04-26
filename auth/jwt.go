package auth

import (
	"crud-books/config"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type jwtTokenEngine struct {
	signingKey      []byte
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
}

func NewJwtEngine(cfg *config.Config) *jwtTokenEngine {
	return &jwtTokenEngine{
		signingKey:      []byte(cfg.JwtSecret),
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
	}
}

func (j jwtTokenEngine) GenAccessToken(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("getting access token claims failed")
	}

	t := time.Time.Add(time.Now(), j.AccessTokenTTL)
	claims["id"] = userId
	claims["exp"] = jwt.NewNumericDate(t)

	signedToken, err := token.SignedString(j.signingKey)
	if err != nil {
		return "", fmt.Errorf("singing token failed, error: %w", err)
	}

	return signedToken, nil
}

func (j jwtTokenEngine) GenRefreshToken() (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("getting refresh token claims falied")
	}

	t := time.Time.Add(time.Now(), j.RefreshTokenTTL)
	claims["exp"] = jwt.NewNumericDate(t)

	signedToken, err := refreshToken.SignedString(j.signingKey)
	if err != nil {
		return "", fmt.Errorf("refresh token signing failed, error: %w", err)
	}

	return signedToken, nil
}

func (j jwtTokenEngine) GetMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			rPath := c.Request().URL.Path

			log.Println(rPath)

			if rPath == "/login" || rPath == "/register" {
				return next(c)
			}
			if rPath == "/books" && authHeader == "" {
				return next(c)
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

			tokenParts := strings.Split(tokenString, ".")
			if len(tokenParts) != 3 {
				return c.String(http.StatusUnauthorized, "malformed token")
			}

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				_, ok := t.Method.(*jwt.SigningMethodHMAC)
				if !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(j.signingKey), nil
			})
			if err != nil {
				return c.String(http.StatusUnauthorized, err.Error())
			}

			_, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return c.String(http.StatusUnauthorized, "JWT claims invalid")
			}

			c.Set("user", token)
			return next(c)
		}
	}
}
