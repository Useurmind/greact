package greact

import (
	"testing"
)

func TestCompareTypes(t *testing.T) {
	equals := CompareTypes(HTMLElement{}, HTMLElement{})
	if !equals {
		t.Errorf("Expected HTMLElements to have same type, but didnt")
	}

	equals = CompareTypes(&VNode{}, &VNode{})
	if !equals {
		t.Errorf("Expected *VNodes to have same type, but didnt")
	}

	equals = CompareTypes(VNode{}, &VNode{})
	if equals {
		t.Errorf("Expected VNode to have other type than &VNode, but didnt")
	}
}
