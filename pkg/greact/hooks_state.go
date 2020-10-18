package greact

type StateHook struct {
	State interface{}
}

func UseState(initialValue interface{}) (interface{}, func(interface{})) {
	hook, _ := HookManagerInstance.GetOrCreateHook(&StateHook{State: initialValue})
	stateHook := hook.(*StateHook)
	return stateHook.State, func(state interface{}) { stateHook.State = state }
}