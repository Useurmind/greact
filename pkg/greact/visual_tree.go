package greact

import (
	"fmt"
)

type DOMNode interface{}
type EventListener interface{}

type VTree struct {
	rootNode    *VNode
	mainNode    *VNode
	mainElement Element
	requestRender func(node *VNode)
}

func NewVTree(rootDOMNode DOMNode, mainElement Element, requestRender func(node *VNode)) *VTree {
	rootNode := NewVNode(nil)
	rootNode.DOMNode = rootDOMNode
	mainNode, _ := rootNode.GetChild(0)
	return &VTree{
		rootNode:    rootNode,
		mainNode:    mainNode,
		mainElement: mainElement,
		requestRender: requestRender,
	}
}

type Renderer interface {
	HandleInsertDOMNodeAction(action *InsertDOMNodeAction) error
	HandleReuseDOMNodeAction(action *ReuseDOMNodeAction) error
	HandleReplaceDOMNodeAction(action *ReplaceDOMNodeAction) error
	HandleUnsetDOMNodeProps(action *UnsetDOMNodePropsAction) error
	HandleSetDOMNodeProps(action *SetDOMNodePropsAction) error
	HandleRemoveDOMNode(action *RemoveDOMNodeAction) error
}

func (tree *VTree) Render(renderer Renderer) error {
	tree.mainNode.requestedRender = true
	return tree.RenderNodeWithElement(renderer, tree.mainNode, tree.mainElement)
}

func (tree *VTree) RenderNode(renderer Renderer, node *VNode) error {
	return tree.RenderNodeWithElement(renderer, node, node.CurrentElement)
}

func (tree *VTree) RenderNodeWithElement(renderer Renderer, node *VNode, element Element) error {
	node.requestedRender = true
	renderQueue, err := tree.ComputeRenderQueue(node, element)
	if err != nil {
		return err
	}

	for _, action := range renderQueue.GetActions() {
		var err error
		switch a := action.(type) {
		case *InsertDOMNodeAction:
			err = renderer.HandleInsertDOMNodeAction(a)
		case *ReuseDOMNodeAction:
			err = renderer.HandleReuseDOMNodeAction(a)
		case *ReplaceDOMNodeAction:
			err = renderer.HandleReplaceDOMNodeAction(a)
		case *UnsetDOMNodePropsAction:
			err = renderer.HandleUnsetDOMNodeProps(a)
		case *SetDOMNodePropsAction:
			err = renderer.HandleSetDOMNodeProps(a)
		case *RemoveDOMNodeAction:
			err = renderer.HandleRemoveDOMNode(a)
		}
		if err != nil {
			return err
		}
	}

	err = renderQueue.ExecutePostRenderQueue()
	if err != nil {
		return err
	}

	return nil
}

type VNode struct {
	CurrentElement Element
	DOMNode        DOMNode

	Parent         *VNode
	Children       []*VNode
	EventListeners map[string]EventListener

	hookCounter int
	hooks       []Hook
	requestedRender bool
}

func NewVNode(parent *VNode) *VNode {
	return &VNode{
		Parent:         parent,
		Children:       make([]*VNode, 0),
		EventListeners: make(map[string]EventListener),
		hooks:          make([]Hook, 0),
	}
}

func (n *VNode) Key() string {
	if n.CurrentElement == nil {
		return ""
	}

	return n.CurrentElement.GetKey()
}

func (n *VNode) FindParentDOMNode() DOMNode {
	if n.Parent == nil {
		return nil
	}

	if n.Parent.DOMNode != nil {
		return n.Parent.DOMNode
	}

	return n.Parent.FindParentDOMNode()
}

func (n *VNode) GetChild(index int) (*VNode, bool) {
	if len(n.Children) < index {
		panic(fmt.Errorf("Requested child %d but got two few childs yet to append %d", index, len(n.Children)))
	}

	if index < len(n.Children) {
		return n.Children[index], false
	}

	newChild := NewVNode(n)
	n.Children = append(n.Children, newChild)

	return newChild, true
}

func (n *VNode) PopChildren(keepChildren int) []*VNode {
	if keepChildren >= len(n.Children) {
		return make([]*VNode, 0)
	}

	poppedChildren := n.Children[keepChildren:len(n.Children)]
	n.Children = n.Children[0:keepChildren]

	return poppedChildren
}

func (n *VNode) OnMounted() {
	fmt.Printf("Node mounted: %s\n", n.Key())
	for _, hook := range n.hooks {
		switch h := hook.(type) {
		case LifecycleHook:
			h.OnMounted()
		}
	}
}

func (n *VNode) OnRendering() {
	fmt.Printf("Node rendering: %s\n", n.Key())
	for _, hook := range n.hooks {
		switch h := hook.(type) {
		case LifecycleHook:
			h.OnRendering()
		}
	}
}

func (n *VNode) OnRendered() {
	fmt.Printf("Node rendered: %s\n", n.Key())
	for _, hook := range n.hooks {
		switch h := hook.(type) {
		case LifecycleHook:
			h.OnRendered()
		}
	}
	n.requestedRender = false
}

func (n *VNode) OnUnmounting() {
	fmt.Printf("Node unmounting: %s\n", n.Key())
	for _, hook := range n.hooks {
		switch h := hook.(type) {
		case LifecycleHook:
			h.OnUnmounting()
		}
	}

	// make sure that we do not by accident use this again
	n.CurrentElement = nil

	if len(n.hooks) > 0 {
		// reset hooks when the node is unmounting a component
		n.hooks = make([]Hook, 0)
	}
}

func (n *VNode) UnmountChildrenRecurse() {
	for _, child := range n.Children {
		// unmount all childs recursively
		child.OnUnmounting()
		child.UnmountChildrenRecurse()
	}
}
