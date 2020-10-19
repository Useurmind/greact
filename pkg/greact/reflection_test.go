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

func TestCompareValues(t *testing.T) {
	equals := CompareValues(1, 1)
	if !equals {
		t.Errorf("Two ints should be equal")
	}
	equals = CompareValues(1, 2)
	if equals {
		t.Errorf("Two ints should not be equal")
	}
	
	equals = CompareValues("asd", "asd")
	if !equals {
		t.Errorf("Two strings should be equal")
	}
	equals = CompareValues("asd", "asdaf")
	if equals {
		t.Errorf("Two strings should not be equal")
	}
	
	equals = CompareValues(true, true)
	if !equals {
		t.Errorf("Two bools should be equal")
	}
	equals = CompareValues(true, false)
	if equals {
		t.Errorf("Two bools should not be equal")
	}

	fun := func() {}
	equals = CompareValues(fun, fun)
	if !equals {
		t.Errorf("The same function should not equal")
	}
	equals = CompareValues(func() {}, func() {})
	if equals {
		t.Errorf("Two independent functions should not be equal")
	}
}

func TestCopyInterfaceValues(t *testing.T) {
	var i interface{} = 1
	var i2 interface{} = ""
	ic := CopyInterfaceValue(i)

	if ic.(int) != 1 {
		t.Errorf("Copied interface should have same value")
	}

	ic = i2
	if i.(int) == 2 {
		t.Errorf("Original should not change through assignment")
	}
}