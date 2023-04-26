package auth

import (
	"crud-books/config"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_GenAccessToken(t *testing.T) {
	jwtSecret := "asd21531asda"
	aTTL, _ := time.ParseDuration("10m")
	rTTL, _ := time.ParseDuration("10m")
	userId := "12454161464"
	cfg := config.Config{
		JwtSecret:       jwtSecret,
		AccessTokenTTL:  aTTL,
		RefreshTokenTTL: rTTL,
	}

	jwtEng := NewJwtEngine(&cfg)
	got, err := jwtEng.GenAccessToken(userId)
	require.NoError(t, err)
	fmt.Println(got)
}
