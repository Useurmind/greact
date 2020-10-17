package greact

import (
	"testing"
)

func TestGetChildWorks(t *testing.T) {
	parent := NewVNode()

	child1 := parent.GetChild(0)
	child2 := parent.GetChild(1)

	if child1 == nil || child2 == nil || child1 == child2 {
		t.Errorf("Childs should not be nil and not be equal")
	}

	if len(parent.Children) != 2 ||
		parent.Children[0] != child1 ||
		parent.Children[1] != child2 {
		t.Errorf("Childs should be in the correct order in child collection")
	}
}

func TestPopChildrenWorks(t *testing.T) {
	parent := NewVNode()

	child1 := parent.GetChild(0)
	child2 := parent.GetChild(1)
	child3 := parent.GetChild(2)
	child4 := parent.GetChild(3)
	child5 := parent.GetChild(4)

	popped := parent.PopChildren(3)

	if len(parent.Children) != 3 ||
		parent.Children[0] != child1 ||
		parent.Children[1] != child2 ||
		parent.Children[2] != child3 {
		t.Errorf("Node should still contain first three children in order")
	}

	if len(popped) != 2 ||
		popped[0] != child4 ||
		popped[1] != child5 {
		t.Errorf("Last two children should be returned")
	}
}

func TestPopChildrenDoesNotKeepMoreThanAvailable(t *testing.T) {
	parent := NewVNode()

	child1 := parent.GetChild(0)
	child2 := parent.GetChild(1)

	popped := parent.PopChildren(3)

	if len(parent.Children) != 2 ||
		parent.Children[0] != child1 ||
		parent.Children[1] != child2 {
		t.Errorf("Node should still contain first two children in order")
	}

	if len(popped) != 0  {
		t.Errorf("Nothing should have been popped")
	}
}


func TestPopChildrenEmptysNodeChildrenOnZero(t *testing.T) {
	parent := NewVNode()

	child1 := parent.GetChild(0)
	child2 := parent.GetChild(1)

	popped := parent.PopChildren(0)

	if len(parent.Children) != 0  {
		t.Errorf("Node should not contain any children")
	}

	if len(popped) != 2 ||
		popped[0] != child1 ||
		popped[1] != child2 {
		t.Errorf("All children should have been popped")
	}
}
