package utils

import (
	"fmt"
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

// Round rounds a number up or down
func Round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}

func ToPercent(num int, denom int) string {
	if denom == 0 {
		return "NA"
	}
	return fmt.Sprintf("%.2f%%", float32(num)/float32(denom)*100)
}
