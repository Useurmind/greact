package greact

import "fmt"

type StateHook struct {
	State interface{}
}

type HookManager struct {
	currentHook      int
	currentComponent Component
	hooks            map[Component][]interface{}
}

func (h *HookManager) SetComponent(component Component) {
	h.currentComponent = component
}

func (h *HookManager) GetOrCreateHook(hook interface{}) interface{} {
	if h.currentComponent == nil {
		panic("No component set when using hook")
	}

	currentHook := h.currentHook
	h.currentHook++
	compHooks := h.hooks[h.currentComponent]

	if len(compHooks) > currentHook {
		fmt.Println("Returning existing hook")
		return compHooks[currentHook]
	} else {
		fmt.Println("Creating new hook")
		compHooks = append(compHooks, hook)
		return hook
	}
}

var hookManager = &HookManager{}

func UseState(initialValue interface{}) (interface{}, func(interface{})) {
	stateHook := hookManager.GetOrCreateHook(&StateHook{State: initialValue}).(*StateHook)
	return stateHook.State, func(state interface{}) { stateHook.State = state }
}