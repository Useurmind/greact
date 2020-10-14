package lexer

import (
	"testing"
)
	

func TestLexerCanParseEmptyFile(t *testing.T) {
	lexer := NewLexerForTestFile(t, "test_data/empty.gsx")
	expPos := Position{ column: 0, line: 1 }
	expToken := EOF

	pos, token, err := lexer.Next()
	if err != nil {
		t.Fatalf("Could not get next lexer token %v", err)
	}

	if pos.column != expPos.column || pos.line != expPos.line {
		t.Errorf("Position should be %v but was %v", expPos, pos)
	}

	if token.Name != expToken.Name {
		t.Errorf("Token should be %v but was %v", expToken, token)
	}
}

func TestLexerCanParseSimpleFile(t *testing.T) {
	lexer := NewLexerForTestFile(t, "test_data/simple.gsx")
	expPos := Position{ column: 1, line: 2 }  // first line is a space line
	expTokens := []Token { 
		GSX_OPEN_ELEMENT,
		GSX_CLOSE_ELEMENT, 
		NewIdent("some-identifier_123"),
	}

	for _, expToken := range expTokens {
		pos, token, err := lexer.Next()
		if err != nil {
			t.Fatalf("Could not get next lexer token %v", err)
		}

		if pos.column != expPos.column || pos.line != expPos.line {
			t.Errorf("Position should be %v but was %v", expPos, pos)
		}

		if token.Name != expToken.Name {
			t.Errorf("Token should be %v but was %v", expToken, token)
		}

		expPos.line++
	}
}