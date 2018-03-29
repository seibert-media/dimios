package assert

import (
	"fmt"
	"reflect"
)

func buildError(format string, message string, args ...interface{}) error {
	if len(message) == 0 {
		return fmt.Errorf(format, args...)
	}
	var a []interface{}
	b := append(a, message)
	c := append(b, args...)
	return fmt.Errorf("%s, "+format, c...)
}

func sameType(a interface{}, b interface{}) bool {
	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(b)
	return aType == bType
}

func less(a interface{}, b interface{}) bool {
	if sameType(a, b) {
		switch a.(type) {
		case int:
			return a.(int) < b.(int)
		case int8:
			return a.(int8) < b.(int8)
		case int16:
			return a.(int16) < b.(int16)
		case int32:
			return a.(int32) < b.(int32)
		case int64:
			return a.(int64) < b.(int64)
		case uint:
			return a.(uint) < b.(uint)
		case uint8:
			return a.(uint8) < b.(uint8)
		case uint16:
			return a.(uint16) < b.(uint16)
		case uint32:
			return a.(uint32) < b.(uint32)
		case uint64:
			return a.(uint64) < b.(uint64)
		case float32:
			return a.(float32) < b.(float32)
		case float64:
			return a.(float64) < b.(float64)
		}
		panic(fmt.Sprintf("unsupport typ: %v", reflect.TypeOf(a)))
	}
	return false
}

func isByteArray(value interface{}) bool {
	if value == nil {
		return false
	}
	typ := reflect.TypeOf(value)
	return typ.Kind() == reflect.Slice && typ.Elem().Kind() == reflect.Uint8
}
