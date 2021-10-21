/**
 * @Author: Lee
 * @Description:
 * @File:  jwt_test
 * @Version: 1.0.0
 * @Date: 2021/10/21 12:11 下午
 */

package tests

import (
	"fmt"
	"go-library/encryption"
	"log"
	"testing"
	"time"
)

const secKey = "1122333"
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MzQ3OTgxODcsInVpZCI6MSwic2Vzc2lvbl9rZXkiOiIiLCJ1c2VyX2FnZW50IjoiIiwiY3JlYXRlX2F0IjowfQ.boAWxqIW8Ze2xuV-Rb1ATJQsUlEZ1-n6isARYaKOEts"

func TestGenerateToken(t *testing.T) {
	j := &encryption.Jwt{
		secKey,
	}
	claims := encryption.CustomClaims{
		Uid: 1,
	}
	a := time.Now().Add(time.Second * 60).Unix()
	claims.ExpiresAt = a
	token, err := j.GenerateToken(claims)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(token)
}

func TestVerifyToken(t *testing.T) {
	j := &encryption.Jwt{
		secKey,
	}
	claims, err := j.VerifyToken(token)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(claims.Uid)
}

func TestRefreshToken(t *testing.T) {
	j := &encryption.Jwt{
		secKey,
	}
	token, err := j.RefreshToken(token, 7200)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(token)

}
