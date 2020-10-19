package greact

import (
	"testing"
)

func TestSetVNodeWorks(t *testing.T) {
	HookManagerInstance.SetVNode(NewVNode(nil))
	HookManagerInstance.SetVNode(nil)
}

func TestCreateMultipleHooksWorks(t *testing.T) {
	node := NewVNode(nil)

	// render 1
	HookManagerInstance.SetVNode(node)
	state1, _ := UseState(0)
	state2, _ := UseState(1)
	state3, _ := UseState(2)
	HookManagerInstance.SetVNode(nil)

	if state1.(int) != 0 {
		t.Errorf("Hook1 should return state 0 but was %d", state1.(int))
	}
	if state2.(int) != 1 {
		t.Errorf("Hook2 should return state 1 but was %d", state2.(int))
	}
	if state3.(int) != 2 {
		t.Errorf("Hook3 should return state 2 but was %d", state3.(int))
	}

	// render 2
	HookManagerInstance.SetVNode(node)
	state1, _ = UseState(0)
	state2, _ = UseState(1)
	state3, _ = UseState(2)
	HookManagerInstance.SetVNode(nil)

	if state1.(int) != 0 {
		t.Errorf("Hook1 should return state 0 but was %d", state1.(int))
	}
	if state2.(int) != 1 {
		t.Errorf("Hook2 should return state 1 but was %d", state2.(int))
	}
	if state3.(int) != 2 {
		t.Errorf("Hook3 should return state 2 but was %d", state3.(int))
	}
}