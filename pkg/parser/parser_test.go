package parser

import (
	"testing"
)

func TestParserCanParseEmptyFile(t *testing.T) {
	parser := NewParserForTestFile(t, "test_data/empty.gsx")

	rootNode, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error while parsing: %v", err)
	}

	if rootNode != nil {
		t.Errorf("Root node for empty file should be nil but was %v", rootNode)
	}
}

func TestParserCanParseSingleElementFile(t *testing.T) {
	expectedNode := &Node{
		Type: "single",
		Children: []*Node{},
	}
	parser := NewParserForTestFile(t, "test_data/single_element.gsx")

	rootNode, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error while parsing: %v", err)
	}

	compareNodes(t, expectedNode, rootNode)
}

func TestParserCanParseChildrenFile(t *testing.T) {
	expectedNode := &Node{
		Type: "Parent",
		Children: []*Node{
			{
				Type: "Child",
			},
		},
	}
	parser := NewParserForTestFile(t, "test_data/children.gsx")

	rootNode, err := parser.Parse()
	if err != nil {
		t.Fatalf("Error while parsing: %v", err)
	}

	compareNodes(t, expectedNode, rootNode)
}

func compareNodes(t *testing.T, expected *Node, actual *Node) {
	if expected == nil && actual != nil {
		t.Fatal("Node is not nil but should be")
	}

	if expected != nil && actual == nil {
		t.Fatal("Node is nil but should not be")
	}

	if actual.Type != expected.Type {
		t.Errorf("Node type should be %s but was %s", expected.Type, actual.Type)
	}

	if len(expected.Children) != len(actual.Children) {
		t.Errorf("Number of children in node %s should be %d but was %d", expected.Type, len(expected.Children), len(actual.Children))
	} else {
		for i, expChild := range expected.Children {
			actualChild := actual.Children[i]

			compareNodes(t, expChild, actualChild)
		}
	}
	
}