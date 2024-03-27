package kit

import (
	"reflect"
)

const (
	// OK request success
	OK = "200"
	// Failed request failed
	Failed = "500"
)

// Response define a struct for http request
type Response struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// New create a new instance from src
func New(src interface{}) interface{} {
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr { //如果是指针类型
		typ = typ.Elem()               //获取源实际类型(否则为指针类型)
		dst := reflect.New(typ).Elem() //创建对象
		return dst.Addr().Interface()  //返回指针
	}
	dst := reflect.New(typ).Elem() //创建对象
	return dst.Interface()
}
