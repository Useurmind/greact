package greact

import (
	"testing"
)

func TestInitialStateCorrect(t *testing.T) {
	HookManagerInstance.SetContext(hookContextWithNode(NewVNode(nil)))

	state, _ := UseState(0)
	if state.(int) != 0 {
		t.Errorf("Initial state should be 0 but was %d", state.(int))
	}
}

func TestCanSetState(t *testing.T) {
	node := NewVNode(nil)

	HookManagerInstance.SetContext(hookContextWithNode(node))
	_, setState := UseState(0)

	setState(1)

	HookManagerInstance.SetContext(hookContextWithNode(node))
	state, _ := UseState(0)
	if state.(int) != 1 {
		t.Errorf("Changed state should be 1 but was %d", state.(int))
	}
}
