package model

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	UserInfo  User   `json:"userInfo"`
	TokenType string `json:"tokenType"` //err_token,refresh_token,access_token
	jwt.RegisteredClaims
}
