package greact

import "fmt"

func (t *VTree) ComputeRenderQueue() (*RenderQueue, error) {
	renderQueue := &RenderQueue{}
	err := t.WalkTree(renderQueue, 0, t.mainNode, t.mainElement)
	if err != nil {
		return renderQueue, err
	}
	return renderQueue, nil
}

func (t *VTree) WalkTree(queue *RenderQueue, index int, node *VNode, element Element) error {
	if element == nil {
		return nil
	}

	switch e := element.(type) {
	case *HTMLElement:
		return t.WalkHTMLElement(queue, index, node, e)
	case *ComponentElement:
		return t.WalkComponentElement(queue, index, node, e)
	}

	return fmt.Errorf("Could not render element missing tag/component")
}

func (t *VTree) WalkComponentElement(queue *RenderQueue, index int, node *VNode, element *ComponentElement) error {
	HookManagerInstance.SetVNode(node)
	renderedElement := element.Component.Render()
	HookManagerInstance.SetVNode(nil)

	renderedNode := node.GetChild(0)

	err := t.WalkTree(queue, 0, renderedNode, renderedElement)

	node.CurrentElement = element

	return err
}

func (t *VTree) WalkHTMLElement(queue *RenderQueue, index int, node *VNode, element *HTMLElement) error {
	err := t.EnsureDOMNodeAction(queue, node, element)
	if err != nil {
		return err
	}

	switch queue.LastAction().(type) {
	case *ReuseDOMNodeAction:
		queue.AddAction(&UnsetDOMNodePropsAction{
			OldElement: node.CurrentElement.(*HTMLElement),
			NewElement: element,
			Node: node,
		})
	}

	queue.AddAction(&SetDOMNodePropsAction{
		NewElement: element,
		Node: node,
	})

	// render children
	numberChildren := 0
	if element.Children != nil {
		numberChildren = len(element.Children)
		for i, child := range element.Children {
			childNode := node.GetChild(i)

			t.WalkTree(queue, i, childNode, child)
		}
	}


	// remove unused nodes
	poppedNodes := node.PopChildren(numberChildren)
	for _, childNode := range poppedNodes {
		queue.AddAction(&RemoveDOMNodeAction{
			Node: childNode,
		})
	}
	

	node.CurrentElement = element

	return nil
}

func (t *VTree) EnsureDOMNodeAction(queue *RenderQueue, node *VNode, element *HTMLElement) error {
	// first check if we can reuse existing dom node
	if node.CurrentElement != nil {

		switch currHTMLElem := node.CurrentElement.(type) {
		case *HTMLElement:
			if (currHTMLElem.Tag == element.Tag) {
				// yes we can reuse it
				queue.AddAction(&ReuseDOMNodeAction{
					OldElement: currHTMLElem,
					NewElement: element,
					Node: node,
				})
				return nil
			} 
	
			queue.AddAction(&ReplaceDOMNodeAction{
				OldElement: currHTMLElem,
				NewElement: element,
				Node: node,
			})
			return nil
		}
		
	}

	// no we need to create a new node
	queue.AddAction(&InsertDOMNodeAction{
		Element: element,
		Node: node,
	})
	return nil
}