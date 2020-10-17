package greact

import "reflect"

func CompareTypes(v1 interface{}, v2 interface{}) bool {
	t1 := reflect.TypeOf(v1)
	t2 := reflect.TypeOf(v2)

	if t1 != t2 {
		return false
	}

	return true
}