/**
 * @Author: Lee
 * @Description:
 * @File:  time
 * @Version: 1.0.0
 * @Date: 2021/10/22 3:57 下午
 */

package utils

import (
	"time"
)

// GetTimeStamp 获取当前时间戳
func GetTimeStamp() int64 {
	return time.Now().Unix()
}
