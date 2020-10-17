package greact

import (
	"fmt"
)


type VNode struct {
	CurrentElement Element
	NextElement    Element
	JSNode         interface{}

	Children []*VNode
	EventListeners map[string]interface{}
}

func NewVNode() *VNode {
	return &VNode{
		Children: make([]*VNode, 0),
		EventListeners: make(map[string]interface{}),
	}
}

func (n *VNode) GetChild(index int) *VNode {
	if len(n.Children) < index {
		panic(fmt.Errorf("Requested child %d but got two few childs yet to append %d", index, len(n.Children)))
	}

	if index < len(n.Children) {
		return n.Children[index]
	}

	newChild := NewVNode()
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