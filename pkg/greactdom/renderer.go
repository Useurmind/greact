package greactdom

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/useurmind/greact/pkg/greact"
)

type Renderer struct {
	jsDoc js.Value
	jsRoot js.Value
	vTree *greact.VTree
}

func Render(root greact.Element) {

	jsDoc := js.Global().Get("document")
	if !jsDoc.Truthy() {
		HandleError(fmt.Errorf("Could not retrieve js document"))
		return
	}

	jsBody := jsDoc.Get("body")
	if !jsBody.Truthy() {
		HandleError(fmt.Errorf("Could not retrieve js body"))
		return
	}

	jsRoot := jsDoc.Call("createElement", "div")
	if !jsRoot.Truthy() {
		HandleError(fmt.Errorf("Could not create js root"))
		return
	}
	jsRoot.Set("id", "root")
	jsBody.Call("appendChild", jsRoot)

	renderer := &Renderer{
		jsDoc: jsDoc,
		jsRoot: jsRoot,
		vTree: greact.NewVTree(jsRoot, root),
	}

	for i := 0; i < 1; i++{
		fmt.Println("Render loop")
		err := renderer.renderRoot()
		HandleError(err)
	}

	fmt.Println("Waiting to exit")
	c := make(chan bool)
	<-c
}

func (r *Renderer) renderRoot() error {
	fmt.Println("Render root")

	return r.vTree.Render(r)
}

func (r *Renderer) HandleInsertDOMNodeAction(action *greact.InsertDOMNodeAction) error {
	element := action.Element
	node := action.Node

	fmt.Printf("Create new js node %s\n", element.Props["id"])
	jsElement := r.jsDoc.Call("createElement", element.Tag)
	
	jsParent := node.FindParentDOMNode().(js.Value)
	jsParent.Call("appendChild", jsElement)
	node.DOMNode = jsElement

	return nil
}

func (r *Renderer) HandleReuseDOMNodeAction(action *greact.ReuseDOMNodeAction) error {
	return nil
}

func (r *Renderer) HandleReplaceDOMNodeAction(action *greact.ReplaceDOMNodeAction) error {
	oldElement := action.OldElement
	newElement := action.NewElement
	node := action.Node
	oldJSElem := node.DOMNode.(js.Value)

	fmt.Printf("Replace js node %s\n", oldElement.Props["id"])
	jsElement := r.jsDoc.Call("createElement", newElement)
	jsParent := node.FindParentDOMNode().(js.Value)
	jsParent.Call("replaceChild", jsElement, oldJSElem)
	node.DOMNode = jsElement

	return nil
}

func (r *Renderer) HandleUnsetDOMNodeProps(action *greact.UnsetDOMNodePropsAction) error {
	oldElement := action.OldElement
	node := action.Node
	jsElement := node.DOMNode.(js.Value)

	// remove old props from node
	for k := range oldElement.Props {
		if k[0] == 'o' && k[1] == 'n' {
			listenerName := strings.ToLower(k[2:len(k)])
			listener := node.EventListeners[listenerName]

			delete(node.EventListeners, listenerName)

			jsElement.Call("removeEventListener", listenerName, listener)
		} else {
			jsElement.Set(k, js.Undefined())
		}
	}

	return nil
}

func (r *Renderer) HandleSetDOMNodeProps(action *greact.SetDOMNodePropsAction) error {
	element := action.NewElement
	node := action.Node
	jsElement := node.DOMNode.(js.Value)

	if element.Props != nil {
		for k, v := range element.Props {
			// event listeners are somewhat special
			if k[0] == 'o' && k[1] == 'n' {
				listenerName := strings.ToLower(k[2:len(k)])
				listener := js.FuncOf(r.wrapGoFunction(v.(func())))
				node.EventListeners[listenerName] = listener

				jsElement.Call("addEventListener", listenerName, listener)
			} else {
				jsElement.Set(k, v)
			}
		}
	}

	return nil
}

func (r *Renderer) HandleRemoveDOMNode(action *greact.RemoveDOMNodeAction) error {
	node := action.Node
	parentJSElement := node.FindParentDOMNode().(js.Value)
	jsElement := node.DOMNode.(js.Value)

	fmt.Printf("Remove js node %s\n", jsElement.Get("id"))
	parentJSElement.Call("removeChild", jsElement)

	return nil
}

func HandleError(err error) {
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