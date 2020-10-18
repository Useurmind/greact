package greact

type EffectHook struct {
	onlyOnce  bool
	execute   bool
	effect    func() func()
	cleanup   func()
	lastProps []interface{}
}

func (h *EffectHook) OnMounted() {
	h.executeEffect()
}

func (h *EffectHook) OnRendering() {
	if !h.onlyOnce {
		h.cleanup()
	}
}

func (h *EffectHook) OnRendered() {
	if !h.onlyOnce {
		h.executeEffect()
	}
}

func (h *EffectHook) OnUnmounting() {
	h.cleanup()
}

func (h *EffectHook) executeEffect() {
	if h.execute {
		h.cleanup = h.effect()
	}
}

func (h *EffectHook) cleanupEffect() {
	if h.execute {
		if h.cleanup != nil {
			h.cleanup()
			h.cleanup = nil
		}
		h.execute = false
	}
}

func UseEffect(performEffect func() func(), dependentProps ...interface{}) {
	useEffect(performEffect, false, dependentProps...)
}

func UseEffectOnce(performEffect func() func(), dependentProps ...interface{}) {
	useEffect(performEffect, true, dependentProps...)
}

func useEffect(performEffect func() func(), onlyOnce bool, dependentProps ...interface{}) {
	newHook := &EffectHook{
		onlyOnce:  onlyOnce,
		effect:    performEffect,
		lastProps: dependentProps,
	}
	hook, hookIndex := HookManagerInstance.GetOrCreateHook(newHook)
	HookManagerInstance.ReplaceHook(hookIndex, newHook)
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

	newHook.execute = shouldExecute
}
