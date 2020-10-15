package parser

import (
	"io/ioutil"
	"strings"
	"testing"
)


func NewParserForTestFile(t *testing.T, filePath string) *Parser {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Could not read test file %s: %v", filePath, err)
	}

	reader := strings.NewReader(string(content))

	return NewParser(reader)
}