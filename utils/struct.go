/**
 * @Author: Lee
 * @Description:
 * @File:  struct
 * @Version: 1.0.0
 * @Date: 2021/10/22 3:57 下午
 */

package utils

import (
	"reflect"
)

// StructToMap 结构体转Map
func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		data[obj1.Field(i).Name] = obj2.Field(i).Interface()
	}
	return data
}
