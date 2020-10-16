package greact

import "reflect"

func ApplyProps(component Component, props interface{}) {
	cPointer := reflect.ValueOf(component)
	cVal := cPointer.Elem()
	if cVal.Kind() != reflect.Struct {
		panic("Expected struct to set props on")
	}

	propsField := cVal.FieldByName("Props")
	if !propsField.IsValid() || !propsField.CanSet() {
		panic("Props field not valid or not settable")
	}

	pPointer := reflect.ValueOf(props)
	propsField.Set(pPointer)
}