/**
 * @Author: Lee
 * @Description:
 * @File:  redis
 * @Version: 1.0.0
 * @Date: 2021/10/22 3:38 下午
 */

package databases

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	host     string
	port     int
	password string
	poolSize int
	minIdle  int
	timeout  int
}

func NewRedis(host string, port int, password string, poolSize int, minIdle int, timeout int) *Redis {
	return &Redis{
		host: host, port: port, password: password, poolSize: poolSize, minIdle: minIdle, timeout: timeout,
	}
}

// NewClient 创建redis连接对象
func (r *Redis) NewClient(db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", r.host, r.port),
		Password:     r.password,
		PoolSize:     r.poolSize,
		MinIdleConns: r.minIdle,
		IdleTimeout:  time.Duration(r.timeout) * time.Second,
		DB:           db,
	})
}
