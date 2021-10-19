/**
 * @Author: Lee
 * @Description:
 * @File:  captcha
 * @Version: 1.0.0
 * @Date: 2021/10/19 10:40 下午
 */

package verify_code

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/mojocn/base64Captcha"
	"time"
)

type CaptchaStore struct {
	Ctx         context.Context
	RedisClient *redis.Client
	Prefix      string
}

func (s *CaptchaStore) Set(id string, value string) error {
	key := s.Prefix + id
	res := s.RedisClient.Set(s.Ctx, key, value, time.Minute*5)
	return res.Err()
}

func (s *CaptchaStore) Get(id string, clear bool) string {
	key := s.Prefix + id
	res := s.RedisClient.Get(s.Ctx, key)
	if clear {
		_ = s.RedisClient.Del(s.Ctx, key)
	}
	return res.Val()
}

func (s *CaptchaStore) Verify(id string, answer string, clear bool) bool {
	code := s.Get(id, clear)
	return code == answer
}

func (s *CaptchaStore) GetVerifyCode() (id string, content string, err error) {
	// 验证码生成器
	drive := base64Captcha.NewDriverString(35, 90, 0,
		base64Captcha.OptionShowHollowLine, 5, "123456789abcdefghijklmnopqrstuvwxyz",
		nil, nil, []string{"wqy-microhei.ttc"})
	captcha := base64Captcha.NewCaptcha(drive, s)
	id, content, err = captcha.Generate()
	return
}
