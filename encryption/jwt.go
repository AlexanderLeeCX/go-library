/**
 * @Author: Lee
 * @Description:
 * @File:  jwt
 * @Version: 1.0.0
 * @Date: 2021/10/19 10:57 下午
 */

package encryption

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

type TokenClaims struct {
	jwt.StandardClaims
	Uid        int64  `json:"uid"`
	SessionKey string `json:"session_key"`
	UserAgent  string `json:"user_agent"`
	CreateAt   int64  `json:"create_at"`
}

type Jwt struct {
	privateKey string
	publicKey  string
}

// GenerateRS256Token 生成RS256 token
func (j *Jwt) GenerateRS256Token(claims *TokenClaims) (string, error) {
	privateKey, err := j.GetRSAPrivateKey()
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tokenStr, err := token.SignedString(privateKey)
	return tokenStr, err
}

// VerifyRS256Token 校验RS256 token
func (j *Jwt) VerifyRS256Token(tokenStr string) (*TokenClaims, error) {
	publicKey, err := j.GetRSAPublicKey()
	if err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 基于JWT的第一部分中的alg字段值进行一次验证
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("验证Token的加密类型错误")
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, err
	}
	if err := token.Claims.Valid(); err != nil {
		return nil, err
	}
	return claims, nil
}

// RefreshRS256Token 刷新token
func (j *Jwt) RefreshRS256Token(token string, tokenExpires int) (*TokenClaims, string, error) {
	claims, err := j.VerifyRS256Token(token)
	if err != nil {
		return claims, "", err
	}
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(tokenExpires)).Unix()
	tokenStr, err := j.GenerateRS256Token(claims)
	return claims, tokenStr, err
}

// GenerateHS256Token 生成HS256 token
func (j *Jwt) GenerateHS256Token(claims *TokenClaims, secret string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := at.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return token, nil
}

// VerifyHS256Token 校验HS256 token
func (j *Jwt) VerifyHS256Token(token string, secret string) (*TokenClaims, error) {
	at, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := at.Claims.(*TokenClaims)
	if !ok {
		return nil, err
	}
	if err := at.Claims.Valid(); err != nil {
		return nil, err
	}
	return claims, nil
}

// RefreshHS256Token 刷新token
func (j *Jwt) RefreshHS256Token(token string, secret string, tokenExpires int) (*TokenClaims, string, error) {
	claims, err := j.VerifyHS256Token(token, secret)
	if err != nil {
		return claims, "", err
	}
	claims.ExpiresAt = time.Now().Add(time.Second * time.Duration(tokenExpires)).Unix()
	tokenStr, err := j.GenerateHS256Token(claims, secret)
	return claims, tokenStr, err
}

// GetRSAPublicKey 获取RSA公钥
func (j *Jwt) GetRSAPublicKey() (publicKey *rsa.PublicKey, err error) {
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(j.publicKey))
	return publicKey, err
}

// GetRSAPrivateKey 获取RSA私钥
func (j *Jwt) GetRSAPrivateKey() (privateKey *rsa.PrivateKey, err error) {
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(j.privateKey))
	return privateKey, err
}
