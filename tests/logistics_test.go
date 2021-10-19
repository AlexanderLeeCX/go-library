/**
 * @Author: Lee
 * @Description:
 * @File:  logistics_test
 * @Version: 1.0.0
 * @Date: 2021/10/16 7:13 下午
 */

package tests

import (
	"context"
	"fmt"
	"go-library/logistics"
	"go-library/verify_code"
	"log"
	"testing"
)

const appCode = "0bba53f3bc734d3baf4a7f3260673437"

// TestGetExpressList 测试获取全球物流公司接口
func TestGetExpressList(t *testing.T) {

	l, err := logistics.NewLogistics(appCode)
	if err != nil {
		log.Fatalln(err)
	}
	expressList, resp, err := l.GetExpressList()
	fmt.Println(resp)
	if err != nil {
		log.Fatalln(err)
	}
	for _, express := range expressList {
		fmt.Println(express)
	}
}

// TestGetExpInfo 测试获取物流信息
func TestGetExpInfo(t *testing.T) {
	l, err := logistics.NewLogistics(appCode)
	if err != nil {
		log.Fatalln(err)
	}
	expInfo, resp, err := l.GetExpInfo("75804477518305", "zhongtong", "")
	fmt.Println(resp)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(expInfo)
}

// TestGetVerifyCode 测试获取验证码
func TestGetVerifyCode(t *testing.T) {
	captcha := &verify_code.CaptchaStore{RedisClient: nil, Ctx: context.Background(), Prefix: ""}
	fmt.Println(captcha)
}
