package token

import (
	"time"

	"github.com/danyouknowme/awayfromus/pkg/util"
	"github.com/dgrijalva/jwt-go"
)

func CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(util.AppConfig.SecretKey))
}
