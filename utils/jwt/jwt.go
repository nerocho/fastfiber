package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 生成最终的 JWT token
func GenerateToken(userId int64, expire time.Duration, secret string) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Hour * expire)
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expireTime), // 过期时间
		ID:        strconv.FormatInt(userId, 10),  // 用户id
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte(secret))
}

// 解析和校验 token
func ParseToken(tokenString string, secret string) (int64, error) {

	var claims jwt.RegisteredClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("SigningMethod 不正确: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, jwt.ErrSignatureInvalid
	}

	return strconv.ParseInt(claims.ID, 10, 64)
}
