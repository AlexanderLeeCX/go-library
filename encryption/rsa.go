/**
 * @Author: Lee
 * @Description:
 * @File:  rsa
 * @Version: 1.0.0
 * @Date: 2021/10/19 10:30 下午
 */

package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// RsaEncryptString rsa加密
// @pubKeyByte 生成的原始公钥 -----BEGIN RSA PUBLIC KEY-----xxx-----END RSA PUBLIC KEY-----
func RsaEncryptString(plainStr string, pubKeyByte []byte) (string, error) {
	block, _ := pem.Decode(pubKeyByte)
	if block == nil {
		return "", errors.New("pem.Decode fail")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := rsa.EncryptPKCS1v15(rand.Reader, pubKey.(*rsa.PublicKey), []byte(plainStr))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}

// RsaDecryptString 解密
// @priKeyByte 生成的原始密钥 -----BEGIN RSA PRIVATE KEY-----xxx-----END RSA PRIVATE KEY-----
func RsaDecryptString(encryptedStr string, priKeyByte []byte) (string, error) {
	block, _ := pem.Decode(priKeyByte)
	if block == nil {
		return "", errors.New("pem.Decode fail")
	}
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", err
	}
	plainBytes, err := rsa.DecryptPKCS1v15(rand.Reader, priKey, encryptedBytes)
	if err != nil {
		return "", err
	}
	return string(plainBytes), nil
}

// PemKeyString 原始生成的密钥转换为字符串
func PemKeyString(keyByte []byte) string {
	block, _ := pem.Decode(keyByte)
	if block == nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(block.Bytes)
}
