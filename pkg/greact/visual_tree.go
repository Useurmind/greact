package greact

import (
	"fmt"
)

type DOMNode interface{}
type EventListener interface{}

type VTree struct {
	rootNode *VNode
	mainNode *VNode
	mainElement Element
}

func NewVTree(rootDOMNode DOMNode, mainElement Element) *VTree {
	rootNode := NewVNode(nil)
	rootNode.DOMNode = rootDOMNode
	mainNode := rootNode.GetChild(0)
	return &VTree{
		rootNode: rootNode,
		mainNode: mainNode,
		mainElement: mainElement,
	}
}

type VNode struct {
	CurrentElement Element
	DOMNode         DOMNode

	Parent *VNode
	Children []*VNode
	EventListeners map[string]EventListener
	
	hookCounter     int
	hooks            []Hook
}

func NewVNode(parent *VNode) *VNode {
	return &VNode{
		Parent: parent,
		Children: make([]*VNode, 0),
		EventListeners: make(map[string]EventListener),
		hooks: make([]Hook, 0),
	}
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

func (n *VNode) GetChild(index int) *VNode {
	if len(n.Children) < index {
		panic(fmt.Errorf("Requested child %d but got two few childs yet to append %d", index, len(n.Children)))
	}

	if index < len(n.Children) {
		return n.Children[index]
	}

	newChild := NewVNode(n)
	n.Children = append(n.Children, newChild)

	return newChild
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
	for _, hook := range n.hooks {
		switch h := hook.(type) {
		case LifecycleHook:
			h.OnMounted()
		}
	}
}

func (n *VNode) OnRendering() {
	for _, hook := range n.hooks {
		switch h := hook.(type) {
		case LifecycleHook:
			h.OnRendering()
		}
	}
}

func (n *VNode) OnRendered() {
	for _, hook := range n.hooks {
		switch h := hook.(type) {
		case LifecycleHook:
			h.OnRendered()
		}
	}
}

func (n *VNode) OnUnmounting() {
	for _, hook := range n.hooks {
		switch h := hook.(type) {
		case LifecycleHook:
			h.OnUnmounting()
		}
	}
}