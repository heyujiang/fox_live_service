package jwttokenx

import (
	"errors"
	"golang.org/x/exp/slog"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/cast"
)

var ErrJwtTokenExpired = errors.New("toke expired")

func GenerateToken(secret string, userId int, username string, expireIn int64, refreshTokenIn int64) (string, error) {
	nowTime := time.Now().Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":            userId,
		"username":          username,
		"expireTime":        nowTime + expireIn - 120,
		"refreshExpireTime": nowTime + refreshTokenIn - 120,
	})
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		slog.Error("generate token failed", slog.String("err", err.Error()))
		return "", err
	}
	return token, nil
}

func parseToken(token, secret string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

func ParseValuePasswordFromToken(token, secret string) (int, string, int64, int64, error) {
	jwtToken, err := parseToken(token, secret)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return 0, "", 0, 0, ErrJwtTokenExpired
			}
		}
		return 0, "", 0, 0, err
	}
	if !jwtToken.Valid {
		return 0, "", 0, 0, errors.New("invalid jwt token")
	}
	var claims = jwtToken.Claims.(jwt.MapClaims)
	var userId int
	var username string
	var expireTime, refreshExpireTime int64

	if v, ok := claims["userId"]; ok {
		userId = cast.ToInt(v)
	}
	if v, ok := claims["username"]; ok {
		username = v.(string)
	}
	if v, ok := claims["expireTime"]; ok {
		expireTime = cast.ToInt64(v)
	}
	if v, ok := claims["refreshExpireTime"]; ok {
		refreshExpireTime = cast.ToInt64(v)
	}
	return userId, username, expireTime, refreshExpireTime, nil
}
