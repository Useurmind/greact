package greact

type Props map[string]interface{}

type Element struct {
	Tag string
	Component Component
	Props Props
	Children []*Element
}

func CreateElement(elementType interface{}, props Props, children ...*Element) *Element {
	element := &Element{
		Props: props,
		Children: children,
	}

	switch t := elementType.(type) {
	case string:
		element.Tag = t
	case Component:
		element.Component = t
	default:
		panic("Unkown element type")
	}

	return element
}