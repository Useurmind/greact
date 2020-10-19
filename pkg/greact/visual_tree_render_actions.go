package greact

type RenderAction interface {
}

type InsertDOMNodeAction struct {
	Element    *HTMLElement
	Node       *VNode
}

type ReuseDOMNodeAction struct {
	OldElement *HTMLElement
	NewElement *HTMLElement
	Node       *VNode
}

type ReplaceDOMNodeAction struct {
	OldElement *HTMLElement
	NewElement *HTMLElement
	Node       *VNode
}

type UnsetDOMNodePropsAction struct {
	OldElement *HTMLElement
	NewElement *HTMLElement	
	Node    *VNode
}

type SetDOMNodePropsAction struct {
	NewElement *HTMLElement
	Node    *VNode
}

type RemoveDOMNodeAction struct {
	Node *VNode
}
