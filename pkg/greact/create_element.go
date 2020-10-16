package greact

type Props map[string]interface{}

type Element interface {
	GetChildren() []Element
}

type HTMLElement struct {
	Tag      string
	Props    Props
	Children []Element
}

func (e *HTMLElement) GetChildren() []Element {
	return e.Children
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
		Children:  children,
	}

	if props != nil {
		ApplyProps(element.Component, props)
	}

	return element
}

func (e *ComponentElement) GetChildren() []Element {
	return e.Children
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
			Children: children,
		}
	case Component:
		element = NewComponentElement(e, props, children...)
	default:
		panic("Unkown element type")
	}

	return element
}