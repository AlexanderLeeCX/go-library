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

//const (
//	keyPublic = `
//-----BEGIN PUBLIC KEY-----
//MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAw1D3Vn0SjJ5y+nSMqsH3
//IjzD3Wa96KyZpNy8zbPCAz2Sw5fG2Yl1bIo9FBlS7dwneOxUNvcKiegpn5eJTwup
//Qnnbv2Ay650OmYSK82J1xwg4YNCJcsaRr3UMSGAfHJ7D1fHg6ELOZFkbhZODrNiN
//ddriF179VY0SkIxTQAf8BTtkXsEnrtjokxJHXuv6Zwn4qhWdShpaRVYTMqsmi9pB
//2ywjYJ+NkTHHm7tjP/z8vKjrXjW0xfb2YxPRgg64MjrlUtn6b4j32df8yCIqgT4a
//iBxY7rFYG2QePPFWvJ3Ms3fVDaZrdSpyszZw3429tPPN1G5feXbKv4kooFiYKv+D
//QwIDAQAB
//-----END PUBLIC KEY-----
//`
//	keyPrivate = `
//-----BEGIN RSA PRIVATE KEY-----
//MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDDUPdWfRKMnnL6
//dIyqwfciPMPdZr3orJmk3LzNs8IDPZLDl8bZiXVsij0UGVLt3Cd47FQ29wqJ6Cmf
//l4lPC6lCedu/YDLrnQ6ZhIrzYnXHCDhg0IlyxpGvdQxIYB8cnsPV8eDoQs5kWRuF
//k4Os2I112uIXXv1VjRKQjFNAB/wFO2RewSeu2OiTEkde6/pnCfiqFZ1KGlpFVhMy
//qyaL2kHbLCNgn42RMcebu2M//Py8qOteNbTF9vZjE9GCDrgyOuVS2fpviPfZ1/zI
//IiqBPhqIHFjusVgbZB488Va8ncyzd9UNpmt1KnKzNnDfjb20883Ubl95dsq/iSig
//WJgq/4NDAgMBAAECggEBAIobhVl1xRkDrV+l7BWOY/akqwax8JVG/rmRkDuP4R8z
//ecSuXOBTj2F5emjs4zPoGU0rJv1av+v16wC7QU9QepXT3uu61Sa/fqRVEX+53ngn
//Ot5Sdu5etIMxq8a9mSI+rVFp4FO7cX+JdqmEPnaJBbYRWQ+XjmDhCQCHCRLc0nrL
//YQ1cmKXhPz9lX9e/fEUgMpn6MvxrR1qxQuF9nXQ4+Eg9MCgOAPtc5GHOSzGOzMnK
//dcO41rpsB39/0FT2cGW4U592wJ9zbX4oFVsHL5o6Zdjp7MDXq2qgiVL/jka+3tNX
//30KzVUvKc9vDwSUbSWI0vGNqy6MNnMaP9cLFdz13f4kCgYEA43+PX8FyERAc8apQ
//lPPv9Eev3ENyVeucTXUcjWVwiXoYtbdQfLJOuG82jdQ4QZ6wtAuEIrc5A0NujnFG
//X1R29UDLNwyNJ68/NOa6kwCcLG0p+5p4SgihmaWwmuoKckaS4gqrS9Lvm4hgqkCb
//olF88wAKA0oddPZoxpnlTL6EFS8CgYEA28k/nRnBsuBorh17v2dEHRDZcBeTX8K7
//qJccZTkMrBeSrA/nlcyNd7u8vT9i9ecttcRV4yhrjNmXMrhyAaMh4azq1qSjptYM
//3WhuaSmi6yR9B++pjiKai2vggrBeWEDOjPEwhWl/PQXE9eMj8bddvPswHH2ywHa6
//XrqLs2V2Vi0CgYBNeAmttOUP9Gm2zaWFI5BJogO7wOf1ZDcklUW0zJ9G4WH6t0Lc
//Q6fU3GI6Z9MEXXKUzPshCz2J4/OI4//vxIaBu5+3zjlfEyk17YAJQQLtifrq584g
//f9HvzWFXT21hPrET8kgkmN7pGsa4EyosWw1ufkvqlNl1E9fYEV3pBVNbFwKBgQCl
//NZbC0ayfeD5Xu0Pc8ZPqwVKhBqe6ENgc91HZ6NNUvPd8rQvot3UTrqRGIVKTA26B
//to7VDPojSyBzeOAByQ1b5S41oFZ/v2C2QZzVIf4cATaW85khhXNkH/gIZOjWMAjT
//Oy286ztAtIiESHQpaytkNfDJSddHAzg+or0GYdtdFQKBgQDJ5IOClC/e+avjckFK
//B//aF79XkN6ASD7DdvBswUIUtSBlJYQNIniE4LMVY72hvTb1mZEE3OmYYmS+Ix3v
//ZMRF/FHY7JSiLfhSOtTmTBEudWDX3B83kzdYuhmjbFxEitwd1v1yZaeLhrKKbP+g
//CW7AQfywMUzAuZcSC+r7/S5W+g==
//-----END RSA PRIVATE KEY-----
//`
//)

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
