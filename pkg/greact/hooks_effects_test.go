package greact

import (
	"testing"
)

func TestUseEffectExecutesWithoutCleanup(t *testing.T) {
	node := NewVNode(nil)
	executed := false
	runEffectHook(node, func() func() {
		executed = true
		return nil
	}, 1)

	if !executed {
		t.Errorf("Effect did not execute")
	}
}

func TestUseEffectExecutesWithCleanup(t *testing.T) {
	node := NewVNode(nil)

	executed := false
	cleanedUp := false
	runEffectHook(node, func() func() {
		executed = true
		return func() {
			cleanedUp = true
		}
	}, 1)

	if !executed {
		t.Errorf("Effect did not execute")
	}

	if !cleanedUp {
		t.Errorf("Effect did not cleanup")
	}
}

func TestUseEffectDoesNotExecuteWhenArgsDidNotChange(t *testing.T) {
	node := NewVNode(nil)
	executed := false
	cleanedUp := false
	runEffectHook(node, func() func() {
		executed = true
		return func() {
			cleanedUp = true
		}
	}, 1)

	executed = false
	cleanedUp = false
	runEffectHook(node, func() func() {
		executed = true
		return func() {
			cleanedUp = true
		}
	}, 1)

	if executed {
		t.Errorf("Effect should not have execute because args did not change")
	}

	if cleanedUp {
		t.Errorf("Effect should not have cleaned up")
	}
}

func TestUseEffectDoesExecuteWhenArgsChange(t *testing.T) {
	node := NewVNode(nil)
	executed := false
	cleanedUp := false

	runEffectHook(node, func() func() {
		executed = true
		return func() {
			cleanedUp = true
		}
	}, 1)

	executed = false
	cleanedUp = false
	runEffectHook(node, func() func() {
		executed = true
		return func() {
			cleanedUp = true
		}
	}, 2)

	if !executed {
		t.Errorf("Effect should have executed because args did change")
	}

	if !cleanedUp {
		t.Errorf("Effect should have cleaned up")
	}
}

func TestEffectPerformedOncePerCycle(t *testing.T) {
	node := NewVNode(nil)
	executed := 0
	cleanedUp := 0

	HookManagerInstance.SetContext(hookContextWithNode(node))
	UseEffect(func() func() {
		executed = executed + 1
		return func() {
			cleanedUp = cleanedUp + 1
		}
	}, 1)

	// executed during 1 cycle in tree walk
	node.OnRendering()
	// executed during 1 cycle in post render
	node.OnMounted()
	node.OnRendered()

	if cleanedUp != 0 {
		t.Errorf("After first cycle cleanup should not have happened but was %d", cleanedUp)
	}

	// executed during 2 cycle in tree walk
	node.OnUnmounting()

	if executed != 1 {
		t.Errorf("Expected effect to be performed once but was %d", executed)
	}
	if cleanedUp != 1 {
		t.Errorf("Expected cleanup to be performed once but was %d", cleanedUp)
	}
}

func TestEffectCleanupBeforeNextCycle(t *testing.T) {
	node := NewVNode(nil)
	executed := 0
	cleanedUp := 0

	// imitate render
	HookManagerInstance.SetContext(hookContextWithNode(node))
	UseEffect(func() func() {
		executed = executed + 1
		return func() {
			cleanedUp = cleanedUp + 1
		}
	}, 1)

	// executed during 1 cycle in tree walk
	node.OnRendering()
	// executed during 1 cycle in post render
	node.OnMounted()
	node.OnRendered()

	if cleanedUp != 0 {
		t.Errorf("After first cycle cleanup should not have happened but was %d", cleanedUp)
	}

	// executed during 2 cycle in tree walk
	node.OnRendering()

	// imitate render
	HookManagerInstance.SetContext(hookContextWithNode(node))
	UseEffect(func() func() {
		executed = executed + 1
		return func() {
			cleanedUp = cleanedUp + 1
		}
	}, 2)

	// executed during 2 cycle in post render
	node.OnRendered()

	if executed != 2 {
		t.Errorf("Expected effect to be performed twice but was %d", executed)
	}
	if cleanedUp != 1 {
		t.Errorf("Expected cleanup to be performed once but was %d", cleanedUp)
	}
}

func TestEffectDoesNotExecuteIfNothingChanged(t *testing.T) {
	node := NewVNode(nil)
	executed := 0
	cleanedUp := 0

	// imitate render
	HookManagerInstance.SetContext(hookContextWithNode(node))
	UseEffect(func() func() {
		executed = executed + 1
		return func() {
			cleanedUp = cleanedUp + 1
		}
	}, 1)

	// executed during 1 cycle in tree walk
	node.OnRendering()
	// executed during 1 cycle in post render
	node.OnMounted()
	node.OnRendered()

	if cleanedUp != 0 {
		t.Errorf("After first cycle cleanup should not have happened but was %d", cleanedUp)
	}

	// executed during 2 cycle in tree walk
	node.OnRendering()

	if cleanedUp != 0 {
		t.Errorf("After second cycle before rendering cleanup should not have happened but was %d", cleanedUp)
	}

	// imitate render with no changes
	HookManagerInstance.SetContext(hookContextWithNode(node))
	UseEffect(func() func() {
		executed = executed + 1
		return func() {
			cleanedUp = cleanedUp + 1
		}
	}, 1)

	// executed during 2 cycle in post render
	node.OnRendered()

	// same effect should still be active
	if executed != 1 {
		t.Errorf("Expected effect to be performed once but was %d", executed)
	}
	if cleanedUp != 0 {
		t.Errorf("Expected cleanup to be performed zero times but was %d", cleanedUp)
	}
}

func TestEffectClosureWorks(t *testing.T) {
	node := NewVNode(nil)
	closure1Executed := false
	closure2Executed := false

	// executed during 1 cycle in tree walk
	node.OnRendering()

	// imitate render
	HookManagerInstance.SetContext(hookContextWithNode(node))
	state, setState := UseState(1)
	stateInt1 := state.(int)
	UseEffect(func() func() {
		closure1Executed = true
		if stateInt1 != 1 {
			t.Errorf("Expected state to be 1 in effect but was %d", stateInt1)
		}
		return func() {
			if stateInt1 != 1 {
				t.Errorf("Expected state to be 1 in cleanup but was %d", stateInt1)
			}
		}
	}, 1)

	setState(2)

	// executed during 1 cycle in post render
	node.OnMounted()
	node.OnRendered()

	// executed during 2 cycle in tree walk
	node.OnRendering()

	// imitate render with no changes
	HookManagerInstance.SetContext(hookContextWithNode(node))
	state, setState = UseState(1)
	stateInt2 := state.(int)
	UseEffect(func() func() {
		closure2Executed = true
		if stateInt2 != 2 {
			t.Errorf("Expected state to be 2 in effect but was %d", stateInt2)
		}
		return func() {
			if stateInt2 != 2 {
				t.Errorf("Expected state to be 2 in cleanup but was %d", stateInt2)
			}
		}
	}, 2)

	node.OnRendered()
	node.OnUnmounting()

	if !closure1Executed || !closure2Executed {
		t.Errorf("Expected both closures to be executed")
	}
}

func runEffectHook(node *VNode, effect func() func(), args ...interface{}) *EffectHook {
	HookManagerInstance.SetContext(hookContextWithNode(node))
	UseEffect(effect, args...)

	effectHook := node.hooks[0].(*EffectHook)
	effectHook.executeEffect()
	effectHook.cleanupEffect()

	return effectHook
}

func TestUseEffectHook(t *testing.T) {
	node := NewVNode(nil)
	executed := false
	cleanedUp := false

	runEffectHook(node, func() func() {
		executed = true
		return func() {
			cleanedUp = true
		}
	}, 1)

	executed = false
	cleanedUp = false
	runEffectHook(node, func() func() {
		executed = true
		return func() {
			cleanedUp = true
		}
	}, 2)

	if !executed {
		t.Errorf("Effect should have executed because args did change")
	}

	if !cleanedUp {
		t.Errorf("Effect should have cleaned up")
	}
}
