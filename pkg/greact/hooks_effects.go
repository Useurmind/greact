package greact

type EffectHook struct {
	effectActive bool
	onlyOnce  bool
	shouldExecute   bool
	shouldCleanup bool
	effect    func() func()
	cleanup   func()
	lastProps []interface{}
}

func (h *EffectHook) OnMounted() {
	h.executeEffect()
}

func (h *EffectHook) OnRendering() {
	// if !h.onlyOnce {
	// 	h.cleanupEffect()
	// }
}

func (h *EffectHook) OnRendered() {
	if !h.onlyOnce {
		h.executeEffect()
	}
}

func (h *EffectHook) OnUnmounting() {
	h.cleanupEffect()
}

func (h *EffectHook) executeEffect() {
	if h.shouldExecute {
		h.cleanupEffect()
		h.cleanup = h.effect()
		h.shouldExecute = false
	}
}

func (h *EffectHook) cleanupEffect() {
	if h.shouldCleanup {
		if h.cleanup != nil {
			h.cleanup()
			h.cleanup = nil
			h.shouldCleanup = false
		}
	}
}

func UseEffect(performEffect func() func(), dependentProps ...interface{}) {
	useEffect(performEffect, false, dependentProps...)
}

func UseEffectOnce(performEffect func() func()) {
	useEffect(performEffect, true,)
}

func useEffect(performEffect func() func(), onlyOnce bool, dependentProps ...interface{}) {
	newHook := &EffectHook{
		onlyOnce:  onlyOnce,
		effect:    performEffect,
		lastProps: dependentProps,
	}
	hook, _ := HookManagerInstance.GetOrCreateHook(newHook)
	// HookManagerInstance.ReplaceHook(hookIndex, newHook)
	oldHook := hook.(*EffectHook)

	// always execute when
	// - no props are given
	// - hook was just created
	// - any prop value changed
	// - only once will apply
	shouldExecute := len(dependentProps) == 0 || oldHook == newHook || onlyOnce
	if !shouldExecute {
		for i, val := range dependentProps {
			oldVal := oldHook.lastProps[i]

			if oldVal != val {
				shouldExecute = true
				break
			}
		}
	}

	oldHook.effect = newHook.effect
	oldHook.shouldExecute = shouldExecute
	oldHook.shouldCleanup = shouldExecute
	oldHook.lastProps = newHook.lastProps
}
