/**
 * @Author: Lee
 * @Description:
 * @File:  jwt
 * @Version: 1.0.0
 * @Date: 2021/10/19 10:57 下午
 */

package encryption

import (
	"github.com/golang-jwt/jwt"
	"log"
	"time"
)

type Jwt struct {
	SecKey string
}

type CustomClaims struct {
	jwt.StandardClaims
	Uid        int64  `json:"uid"`
	SessionKey string `json:"session_key"`
	UserAgent  string `json:"user_agent"`
	CreateAt   int64  `json:"create_at"`
}

// GenerateToken 生成 token
func (j *Jwt) GenerateToken(claims CustomClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(j.SecKey))
}

// VerifyToken 解析token
func (j *Jwt) VerifyToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecKey), nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	claims := token.Claims.(*CustomClaims)
	err = claims.Valid()
	return claims, err
}

// RefreshToken 刷新token时效
func (j *Jwt) RefreshToken(tokenString string, tokenExpires int) (string, error) {
	claims, err := j.VerifyToken(tokenString)
	if err != nil {
		return "", err
	}
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(tokenExpires)).Unix()
	token, err := j.GenerateToken(*claims)
	return token, err
}
