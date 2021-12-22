package util

import "reflect"

func Zero(v interface{}) interface{} {
	vType := reflect.TypeOf(v)
	return reflect.Zero(vType).Interface()
}
