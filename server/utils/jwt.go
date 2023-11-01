package utils

import (
	"github.com/golang-jwt/jwt"
	"redirect/model"
	"time"
)

var (
	jwtSecret = []byte(JwtConfig.JwtKey)
	jwtIssuer = "gin-redirect"
)

// GenerateToken 生成token
func GenerateToken(account string, userId int) (token string, err error) {
	jwtExpireTime := time.Duration(JwtConfig.JwtExpiredTime) * time.Second
	expireTime := time.Now().Add(jwtExpireTime).Unix()
	claims := model.Claims{
		Account: account,
		Id:      userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    jwtIssuer,
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return
	}

	return
}

// ParseToken 解析token
func ParseToken(token string) (*model.Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*model.Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
