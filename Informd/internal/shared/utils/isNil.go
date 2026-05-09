package utils

import "reflect"

func IsNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map:
		return rv.IsNil()
	}

	return false
}
