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
	rootElement greact.Element
	rootNode *greact.VNode
}

func Render(root greact.Element) {

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
	jsBody.Call("appendChild", jsRoot)

	renderer := &Renderer{
		jsDoc: jsDoc,
		jsRoot: jsRoot,
		rootElement: root,
		rootNode: greact.NewVNode(),
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

	// children := r.jsRoot.Get("children")
	// for i := 0; i < children.Get("length").Int(); i++ {
	// 	child := children.Index(i)
	// 	r.jsRoot.Call("removeChild", child)
	// }

	// r.rootNode.JSNode = r.jsRoot

	return r.renderElement(r.jsRoot, 0, r.rootNode, r.rootElement)
}

func (r *Renderer) renderElement(jsParent js.Value, index int, node *greact.VNode, element greact.Element) error {
	if element == nil {
		return nil
	}

	fmt.Printf("Rendering node %v with index %d / element %v\n", node, index, element)

	node.NextElement = element
	
	switch e := element.(type) {
	case *greact.HTMLElement:

		reused := r.ensureJSChild(jsParent, node)
		jsElement := node.JSNode.(js.Value)

		if reused {
			currHTMLElem := node.CurrentElement.(*greact.HTMLElement)
			// remove old props from node
			for k := range currHTMLElem.Props {
				if k[0] == 'o' && k[1] == 'n' {
					listenerName := strings.ToLower(k[2:len(k)])
					listener := node.EventListeners[listenerName]

					delete(node.EventListeners, listenerName)

					jsElement.Call("removeEventListener", listenerName, listener)
				} else {
					jsElement.Set(k, js.Undefined())
				}
			}
		}

		// apply new props to node
		if e.Props != nil {
			for k, v := range e.Props {
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

		// render children
		numberChildren := 0
		if e.Children != nil {
			numberChildren = len(e.Children)
			for i, child := range e.Children {
				childNode := node.GetChild(i)

				r.renderElement(jsElement, i, childNode, child)
			}
		}


		// remove unused nodes
		fmt.Printf("Keep %d child nodes and pop %d\n", numberChildren, len(node.Children) - numberChildren)
		poppedNodes := node.PopChildren(numberChildren)
		for _, node := range poppedNodes {
			jsChild := node.JSNode.(js.Value)
			fmt.Printf("Remove js node %s\n", jsChild.Get("id"))
			jsElement.Call("removeChild", jsChild)
		}
		

		node.CurrentElement = e

		return nil
	case *greact.ComponentElement:
		// curretCElem := node.CurrentElement.(*ComponentElement)
		greact.HookManagerInstance.SetVNode(node)
		renderedElement := e.Component.Render()
		greact.HookManagerInstance.SetVNode(nil)

		// if node.CurrentElement != nil {
		// 	// keep the component the same to k
		// 	if CompareTypes(curretCElem.Component, e.Component) {

		// 	}
		// }

		renderedNode := node.GetChild(0)

		err := r.renderElement(jsParent, index, renderedNode, renderedElement)

		node.CurrentElement = e

		return err
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

// return whether the old node can be reused
func (r *Renderer) ensureJSChild(jsParent js.Value, node *greact.VNode) bool {
	nextHTMLElem := node.NextElement.(*greact.HTMLElement)

	// first check if we can reuse existing dom node
	if node.CurrentElement != nil {

		switch currHTMLElem := node.CurrentElement.(type) {
		case *greact.HTMLElement:
			oldJSElem := node.JSNode.(js.Value)
			if (currHTMLElem.Tag == nextHTMLElem.Tag) {
				// yes we can reuse it
				fmt.Printf("Reuse js node %s\n", oldJSElem.Get("id"))
				return true
			} 
	
			// no we need to replace it with a different kind of node
			fmt.Printf("Replace js node %s\n", oldJSElem.Get("id"))
			jsElement := r.jsDoc.Call("createElement", nextHTMLElem.Tag)
			jsParent.Call("replaceChild", jsElement, oldJSElem)
			node.JSNode = jsElement
			return false
		}
		
	}

	// no we need to create a new node
	fmt.Printf("Create new js node %s\n", nextHTMLElem.Props["id"])
	jsElement := r.jsDoc.Call("createElement", nextHTMLElem.Tag)
	jsParent.Call("appendChild", jsElement)
	node.JSNode = jsElement
	return false
}