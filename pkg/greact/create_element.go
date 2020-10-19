package greact

import (
	"reflect"
	"github.com/google/go-cmp/cmp"
)


type Props map[string]interface{}

type Element interface {
	GetKey() string

	GetChildren() []Element

	Equal(other Element) bool
}

type HTMLElement struct {
	Tag      string
	Props    Props
	Children []Element
}

func (e *HTMLElement) GetKey() string {
	key, ok := e.Props["key"]
	if !ok {
		return ""
	}
	return key.(string)
}

func (e *HTMLElement) GetChildren() []Element {
	return e.Children
}

func (e *HTMLElement) Equal(other Element) bool {
	switch o := other.(type) {
	case *HTMLElement:
		if e.Tag != o.Tag {
			return false
		}

		for k,v := range e.Props {
			vOther, ok := o.Props[k]
			if !ok {
				return false
			}

			if vOther != v {
				return false
			}
		}

		return true
	default:
		return false
	}
}

type ComponentElement struct {
	Component Component
	Props     interface{}
	Children  []Element
}

func NewComponentElement(component Component, props interface{}, children ...Element) *ComponentElement {
	element := &ComponentElement{
		Component: component,
		Props:     props,
		Children:  filterNilChildren(children...),
	}

	if props != nil {
		ApplyProps(element.Component, props)
	}

	return element
}

func (e *ComponentElement) GetKey() string {
	if e.Props == nil {
		return ""
	}

	return reflect.ValueOf(e.Props).Elem().FieldByName("Key").String()
}

func (e *ComponentElement) GetChildren() []Element {
	return e.Children
}

func (e *ComponentElement) Equal(other Element) bool {
	switch o := other.(type) {
	case *ComponentElement:
		cTypeEqual := CompareTypes(e.Component, o.Component)

		if !cTypeEqual {
			return false
		}

		cmp.Equal(e.Props, o.Props)

		return true
	default:
		return false
	}
}

func CreateElement(elementType interface{}, props interface{}, children ...Element) Element {
	var element Element
	switch e := elementType.(type) {
	case string:
		var appliedProps Props
		if props != nil {
			appliedProps = props.(Props)
		}
		element = &HTMLElement{
			Tag:      e,
			Props:    appliedProps,
			Children: filterNilChildren(children...),
		}
	case Component:
		element = NewComponentElement(e, props, children...)
	default:
		panic("Unkown element type")
	}

	return element
}

func filterNilChildren(children ...Element) []Element {
	nonNil := make([]Element, 0)

	for _, child := range children {
		if child != nil {
			nonNil = append(nonNil, child)
		}
	}

	return nonNil
}

