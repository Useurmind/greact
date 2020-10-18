package greact

type RenderQueue struct {
	actions []RenderAction
}

func NewRenderQueue() *RenderQueue {
	return &RenderQueue{
		actions: make([]RenderAction, 0),
	}
}

func (q *RenderQueue) Clear(action RenderAction) {
	q.actions = make([]RenderAction, 0)
}

func (q *RenderQueue) AddAction(action RenderAction) {
	q.actions = append(q.actions, action)
}

func (q *RenderQueue) LastAction() RenderAction {
	return q.actions[len(q.actions)-1]
}

func (q *RenderQueue) GetActions() []RenderAction {
	return q.actions
}

type Renderer interface {
	HandleInsertDOMNodeAction(action *InsertDOMNodeAction) error
	HandleReuseDOMNodeAction(action *ReuseDOMNodeAction) error
	HandleReplaceDOMNodeAction(action *ReplaceDOMNodeAction) error
	HandleUnsetDOMNodeProps(action *UnsetDOMNodePropsAction) error
	HandleSetDOMNodeProps(action *SetDOMNodePropsAction) error
	HandleRemoveDOMNode(action *RemoveDOMNodeAction) error
}

func RenderVTree(tree *VTree, renderer Renderer) error {
	renderQueue, err := tree.ComputeRenderQueue()
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

	return nil
}

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
