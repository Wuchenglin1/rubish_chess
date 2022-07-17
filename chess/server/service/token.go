package service

import (
	"chess/server/model"
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
	//cfg := tool.GetConfig()
	str, err := token.SignedString([]byte("clinyu"))
	if err != nil {
		return "", err
	}
	return str, nil
}

func ParseToken(tokenStr string) (model.User, error) {
	//cfg := tool.GetConfig()
	token, err := jwt.ParseWithClaims(tokenStr, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("clinyu"), nil
	})
	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		return claims.UserInfo, nil
	} else {
		return model.User{}, err
	}
}
