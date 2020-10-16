package greact

import (
	"fmt"
	"strings"
	"syscall/js"
	_ "time"
)

type Renderer struct {
	jsDoc js.Value
	jsRoot js.Value
	root Element
}

func Render(root Element) {

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
		root: root,
	}

	for i := 0; i < 1; i++{
		fmt.Println("Render loop")
		err := renderer.renderRoot()
		handleError(err)
	}

	fmt.Println("Waiting to exit")
	c := make(chan bool)
	<-c
}

func (r *Renderer) renderRoot() error {
	fmt.Println("Render root")

	children := r.jsRoot.Get("children")
	for i := 0; i < children.Get("length").Int(); i++ {
		child := children.Index(i)
		r.jsRoot.Call("removeChild", child)
	}

	return r.renderElement(r.jsRoot, 0, r.root)
}

func (r *Renderer) renderElement(jsParent js.Value, index int, element Element) error {
	if element == nil {
		return nil
	}
	
	switch e := element.(type) {
	case *HTMLElement:
		var jsElement js.Value
		siblings := jsParent.Get("children")
		if siblings.Get("length").Int() > index {
			jsElement = siblings.Index(index)
		} else {
			jsElement = r.jsDoc.Call("createElement", e.Tag)
			jsParent.Call("appendChild", jsElement)
		}

		if e.Props != nil {
			for k, v := range e.Props {
				// event listeners are somewhat special
				if k[0] == 'o' && k[1] == 'n' {
					jsElement.Call("addEventListener", strings.ToLower(k[2:len(k)]), js.FuncOf(r.wrapGoFunction(v.(func()))))
				} else {
					jsElement.Set(k, v)
				}
			}
		}

		if e.Children != nil {
			for i, child := range e.Children {
				r.renderElement(jsElement, i, child)
			}
		}

		return nil
	case *ComponentElement:
		hookManager.SetComponent(e.Component)
		renderedElement := e.Component.Render()
		hookManager.SetComponent(nil)
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

func (r *Renderer) wrapGoFunction(fn func()) func(js.Value, []js.Value) interface {} {
    return func(_ js.Value, _ []js.Value) interface {} {
		fn()
		r.renderRoot()
        return nil
    }
}
