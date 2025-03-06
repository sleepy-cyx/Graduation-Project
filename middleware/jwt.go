package middleware

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// 密钥，用于签名和解签名

const (
	TokenExpireTime = time.Hour * 10000
)

var jwtKey = []byte("key")

type Claims struct {
	Username string `json:"username"`
	Id       uint32 `json:"id"`
	jwt.StandardClaims
}

// GenerateToken 生成jwt
func GenerateToken(username string, id uint32) (string, error) {
	expireTime := time.Now().Add(TokenExpireTime)
	claims := Claims{
		Username: username,
		Id:       id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(jwtKey)
}

// ParseToken 解析jwt
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	claim, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(jwtKey), err
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := claim.Claims.(*Claims); ok && claim.Valid { // 校验token

		return claims, nil
	}
	return nil, errors.New("invalid token")
}
