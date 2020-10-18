package greact

import (
	"testing"
)

func TestSetVNodeWorks(t *testing.T) {
	HookManagerInstance.SetVNode(NewVNode(nil))
	HookManagerInstance.SetVNode(nil)
}