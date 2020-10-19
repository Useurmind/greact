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
	var err error
	
	if node.CurrentElement == nil {
		// if the previous element was nil this is a mount
		queue.AddPostRenderAction(&ComponentMountedAction{
			Node: node,
		})
	} else if !CompareTypes(node.CurrentElement, element) {
		// if the component type changed this is an unmount and remount
		node.OnUnmounting()
		queue.AddPostRenderAction(&ComponentMountedAction{
			Node: node,
		})
	}

	node.CurrentElement = element

	HookManagerInstance.SetVNode(node)
	node.OnRendering()
	renderedElement := element.Component.Render()
	HookManagerInstance.SetVNode(nil)

	// if the previous element was nil this is a mount
	queue.AddPostRenderAction(&ComponentRenderedAction{
		Node: node,
	})

	keptChildren := 0	
	if renderedElement != nil {
		keptChildren = 1
		renderedNode, _ := node.GetChild(0)
		err = t.WalkTree(queue, 0, renderedNode, renderedElement)
	}
	
	err = t.removeUnusedChildNodes(queue, node, keptChildren)
	if err != nil {
		return err
	}

	return err
}

func (t *VTree) WalkHTMLElement(queue *RenderQueue, index int, node *VNode, element *HTMLElement) error {
	// unmounting happens in this function
	err := t.EnsureDOMNodeAction(queue, node, element)
	if err != nil {
		return err
	}

	node.CurrentElement = element

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
			childNode, _ := node.GetChild(i)

			t.WalkTree(queue, i, childNode, child)
		}
	}

	err = t.removeUnusedChildNodes(queue, node, numberChildren)
	if err != nil {
		return err
	}
	

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
		case *ComponentElement:
			// if the previous element is a component we unmount it
			node.OnUnmounting()
		}		
	}


	// no we need to create a new node
	queue.AddAction(&InsertDOMNodeAction{
		Element: element,
		Node: node,
	})
	return nil
}

func (t *VTree) removeUnusedChildNodes(queue *RenderQueue, node *VNode, keepChildren int) error {
	poppedNodes := node.PopChildren(keepChildren)
	for _, childNode := range poppedNodes {
		// deleted nodes will have their elements removed
		childNode.OnUnmounting()
		childNode.UnmountChildrenRecurse()

		queue.AddAction(&RemoveDOMNodeAction{
			Node: childNode,
		})
	}

	return nil
}