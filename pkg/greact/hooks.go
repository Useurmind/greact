package greact

import "fmt"

type Hook interface {

}

type LifecycleHook interface {
	OnMounted()
	OnRendering()
	OnRendered()
	OnUnmounting()
}

type HookManager struct {
	hookContext *HookContext
}

type HookContext struct {
	CurrentNode *VNode
	RequestRerender func(*VNode)
}

func (h *HookManager) SetContext(hookContext *HookContext) {
	h.hookContext = hookContext

	if hookContext != nil && hookContext.CurrentNode != nil {
		// always reset this counter
		hookContext.CurrentNode.hookCounter = 0
	}
}

func (h *HookManager) ReplaceHook(index int, hook Hook) {
	node := h.hookContext.CurrentNode
	node.hooks[index] = hook
}

func (h *HookManager) GetOrCreateHook(hook Hook) (Hook, int) {
	if h.hookContext == nil {
		panic("No hook context set when using hook")
	}

	node := h.hookContext.CurrentNode
	hookCount := node.hookCounter
	hooks := node.hooks
	
	node.hookCounter = node.hookCounter + 1

	if len(hooks) > hookCount {
		fmt.Println("Returning existing hook")
		return hooks[hookCount], hookCount
	} 

	fmt.Println("Creating new hook")
	node.hooks = append(hooks, hook)
	return hook, hookCount
}

func (h *HookManager) GetRequestRerender() func() {
	node := h.hookContext.CurrentNode
	requestRerender := h.hookContext.RequestRerender

	return func() {
		node.requestedRender = true
		requestRerender(node)
	}
}

var HookManagerInstance = &HookManager{
}