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

func runEffectHook(node *VNode, effect func() func(), args ...interface{}) *EffectHook {
	HookManagerInstance.SetVNode(node)
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