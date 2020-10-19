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
	currentNode *VNode
}

func (h *HookManager) SetVNode(node *VNode) {
	h.currentNode = node

	if node != nil {
		// always reset this counter
		node.hookCounter = 0
	}
}

func (h *HookManager) ReplaceHook(index int, hook Hook) {
	h.currentNode.hooks[index] = hook
}

func (h *HookManager) GetOrCreateHook(hook Hook) (Hook, int) {
	if h.currentNode == nil {
		panic("No node set when using hook")
	}

	hookCount := h.currentNode.hookCounter
	hooks := h.currentNode.hooks
	
	h.currentNode.hookCounter = h.currentNode.hookCounter + 1

	if len(hooks) > hookCount {
		fmt.Println("Returning existing hook")
		return hooks[hookCount], hookCount
	} 

	fmt.Println("Creating new hook")
	h.currentNode.hooks = append(hooks, hook)
	return hook, hookCount
}

var HookManagerInstance = &HookManager{
}