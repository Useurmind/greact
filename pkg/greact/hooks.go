package greact

import "fmt"

type StateHook struct {
	State interface{}
}

type HookManager struct {
	currentNode *VNode
	hookCounters     map[*VNode]int
	hooks            map[*VNode][]interface{}
}

func (h *HookManager) SetVNode(node *VNode) {
	h.currentNode = node

	// always reset this counter
	h.hookCounters[node] = 0

	_, ok := h.hooks[node]
	if !ok {
		h.hooks[node] = make([]interface{}, 0)
	}
}

func (h *HookManager) GetOrCreateHook(hook interface{}) interface{} {
	if h.currentNode == nil {
		panic("No node set when using hook")
	}

	hookCount := h.hookCounters[h.currentNode]
	hooks := h.hooks[h.currentNode]

	if len(hooks) > hookCount {
		fmt.Println("Returning existing hook")
		return hooks[hookCount]
	} 

	fmt.Println("Creating new hook")
	h.hooks[h.currentNode] = append(hooks, hook)
	h.hookCounters[h.currentNode] = h.hookCounters[h.currentNode] + 1
	return hook
}

var HookManagerInstance = &HookManager{
	hookCounters: make(map[*VNode]int),
	hooks: make(map[*VNode][]interface{}),
}

func UseState(initialValue interface{}) (interface{}, func(interface{})) {
	stateHook := HookManagerInstance.GetOrCreateHook(&StateHook{State: initialValue}).(*StateHook)
	return stateHook.State, func(state interface{}) { stateHook.State = state }
}