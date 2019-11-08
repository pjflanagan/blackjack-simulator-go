package utils

import (
	"reflect"
)

// Contains check whether slice or array s has element e
func Contains(s interface{}, e interface{}) bool {
	switch reflect.TypeOf(s).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(s)
		for i := 0; i < s.Len(); i++ {
			if s.Index(i).Interface() == e {
				return true
			}
		}
	}
	return false
}
