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

func CompareValues(v1 interface{}, v2 interface{}) bool {
	t1 := reflect.TypeOf(v1)
	t2 := reflect.TypeOf(v2)

	if t1.Kind() == reflect.Func || t2.Kind() == reflect.Func {
		// its not possible to compare if the values in the closure are the same
		// therefore we need to rerender elements if they have event handlers attached
		return false
		//return reflect.ValueOf(v1).Pointer() == reflect.ValueOf(v2).Pointer()
	}

	return v1 == v2
}

func CopyInterfaceValue(v interface{}) interface{} {
	return v
}