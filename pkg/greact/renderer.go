package greact

import (
	"fmt"
	_ "time"
	"syscall/js"
)

type Renderer struct {
	jsDoc js.Value
	jsRoot js.Value
}

func Render(root *Element) {
	jsDoc := js.Global().Get("document")
	if !jsDoc.Truthy() {
		handleError(fmt.Errorf("Could not retrieve js document"))
		return
	}

	jsBody := jsDoc.Get("body")
	if !jsBody.Truthy() {
		handleError(fmt.Errorf("Could not retrieve js body"))
		return
	}

	jsRoot := jsDoc.Call("createElement", "div")
	if !jsRoot.Truthy() {
		handleError(fmt.Errorf("Could not create js root"))
		return
	}
	jsRoot.Set("id", "root")
	jsRoot.Set("innerHTML", "Hello WASM")

	jsBody.Call("appendChild", jsRoot)

	renderer := &Renderer{
		jsDoc: jsDoc,
		jsRoot: jsRoot,
	}

	for i := 0; i < 3; i++{
		fmt.Println("Render loop")
		err := renderer.renderElement(jsRoot, 0, root)
		handleError(err)
	}
}

func (r *Renderer) renderElement(jsParent js.Value, index int, element *Element) error {
	if element == nil {
		return nil
	}

	if element.Tag != "" {
		var jsElement js.Value
		siblings := jsParent.Get("children")
		if siblings.Get("length").Int() > index {
			jsElement = siblings.Index(index)
		} else {
			jsElement = r.jsDoc.Call("createElement", element.Tag)
			jsParent.Call("appendChild", jsElement)
		}

		if element.Props != nil {
			for k, v := range element.Props {
				jsElement.Set(k, v)
			}
		}

		if element.Children != nil {
			for i, child := range element.Children {
				r.renderElement(jsElement, i, child)
			}
		}

		return nil
	} else if element.Component != nil {
		renderedElement := element.Component.Render()
		return r.renderElement(jsParent, index, renderedElement)
	}

	return fmt.Errorf("Could not render element missing tag/component")
}

func handleError(err error) {
	if err != nil {
		fmt.Errorf("ERROR: %v", err)
		panic(err)
	}
}