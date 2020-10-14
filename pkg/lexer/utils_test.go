package lexer

import (
	"io/ioutil"
	"strings"
	"testing"
)


func NewLexerForTestFile(t *testing.T, filePath string) *Lexer {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Could not read test file %s: %v", filePath, err)
	}

	reader := strings.NewReader(string(content))

	return NewLexer(reader)
}