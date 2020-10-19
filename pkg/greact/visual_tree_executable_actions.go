package greact

type ExecutableAction interface {
	Execute() error
}

type ComponentMountedAction struct {
	Node *VNode
}

func (a *ComponentMountedAction) Execute() error {
	a.Node.OnMounted()
	return nil
}

type ComponentUnmountingAction struct {
	Node *VNode
}

func (a *ComponentUnmountingAction) Execute() error {
	a.Node.OnUnmounting()
	return nil
}

type ComponentRenderingAction struct {
	Node *VNode
}

func (a *ComponentRenderingAction) Execute() error {
	a.Node.OnRendering()
	return nil
}

type ComponentRenderedAction struct {
	Node *VNode
}

func (a *ComponentRenderedAction) Execute() error {
	a.Node.OnRendered()
	return nil
}