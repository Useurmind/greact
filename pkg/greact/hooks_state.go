package greact

import "fmt"

type StateHook struct {
	State interface{}
}

func UseState(initialValue interface{}) (interface{}, func(interface{})) {
	hook, _ := HookManagerInstance.GetOrCreateHook(&StateHook{State: initialValue})
	stateHook := hook.(*StateHook)
	requestRerender := HookManagerInstance.GetRequestRerender()
	return stateHook.State, func(state interface{}) {
		fmt.Printf("Setting state to %v\n", state)
		stateHook.State = state
		requestRerender()
	}
}