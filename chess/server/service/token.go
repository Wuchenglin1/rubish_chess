package service

import (
	"chess/server/model"
	"chess/server/tool"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

func CreateToken(u model.User, tokenType string, duration time.Duration) (string, error) {
	claims := model.Claims{
		UserInfo:  u,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  []string{"player"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			Issuer:    "clinyu",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	cfg := tool.GetConfig()
	str, err := token.SignedString([]byte(cfg.Signature))
	if err != nil {
		return "", err
	}
	return str, nil
}

func ParseToken(tokenStr string) (model.User, error) {
	cfg := tool.GetConfig()
	token, err := jwt.ParseWithClaims(tokenStr, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Signature), nil
	})
	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims.UserInfo, nil
	} else {
		return model.User{}, err
	}
}
