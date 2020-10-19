package greact

type RenderQueue struct {
	actions []RenderAction
	postRenderActions []ExecutableAction
}

func NewRenderQueue() *RenderQueue {
	return &RenderQueue{
		actions: make([]RenderAction, 0),
		postRenderActions: make([]ExecutableAction, 0),
	}
}

func (q *RenderQueue) Clear(action RenderAction) {
	q.actions = make([]RenderAction, 0)
}

func (q *RenderQueue) AddAction(action RenderAction) {
	q.actions = append(q.actions, action)
}

func (q *RenderQueue) AddPostRenderAction(action ExecutableAction) {
	q.postRenderActions = append(q.postRenderActions, action)
}

func (q *RenderQueue) ExecutePostRenderQueue() error {
	for _, action:= range q.postRenderActions {
		err := action.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func (q *RenderQueue) LastAction() RenderAction {
	return q.actions[len(q.actions)-1]
}

func (q *RenderQueue) GetActions() []RenderAction {
	return q.actions
}