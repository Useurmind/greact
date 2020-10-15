package lexer

import (
	"testing"
)

type TokenWithPos struct {
	token Token
	pos   Position
}

func TestLexerCanParseEmptyFile(t *testing.T) {
	lexer := NewLexerForTestFile(t, "test_data/empty.gsx")
	expPos := Position{column: 0, line: 1}
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
	expPos := Position{column: 1, line: 2} // first line is a space line
	expTokens := []Token{
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

func TestLexerCanParseElementsFile(t *testing.T) {
	lexer := NewLexerForTestFile(t, "test_data/elements.gsx")
	expTokens := []TokenWithPos{
		// line 1
		{
			token: GSX_OPEN_ELEMENT,
			pos:   NewPosition(1, 1),
		},
		{
			token: NewIdent("div"),
			pos:   NewPosition(1, 2),
		},
		{
			token: GSX_CLOSE_ELEMENT,
			pos:   NewPosition(1, 5),
		},
		{
			token: GSX_OPEN_CLOSING_ELEMENT,
			pos:   NewPosition(1, 6),
		},
		{
			token: NewIdent("div"),
			pos:   NewPosition(1, 8),
		},
		{
			token: GSX_CLOSE_ELEMENT,
			pos:   NewPosition(1, 11),
		},
		// line 2
		{
			token: GSX_OPEN_ELEMENT,
			pos:   NewPosition(2, 1),
		},
		{
			token: NewIdent("span"),
			pos:   NewPosition(2, 2),
		},
		{
			token: GSX_CLOSE_ELEMENT,
			pos:   NewPosition(2, 6),
		},
		{
			token: GSX_OPEN_CLOSING_ELEMENT,
			pos:   NewPosition(2, 7),
		},
		{
			token: NewIdent("span"),
			pos:   NewPosition(2, 9),
		},
		{
			token: GSX_CLOSE_ELEMENT,
			pos:   NewPosition(2, 13),
		},
		// line 3
		{
			token: GSX_OPEN_ELEMENT,
			pos:   NewPosition(3, 1),
		},
		{
			token: NewIdent("SelfClosing"),
			pos:   NewPosition(3, 2),
		},
		{
			token: GSX_CLOSE_SELFCLOSE_ELEMENT,
			pos:   NewPosition(3, 14),
		},
	}

	checkExpectedTokens(t, lexer, expTokens)
}

func TestLexerCanParseChildrenFile(t *testing.T) {
	lexer := NewLexerForTestFile(t, "test_data/children.gsx")
	expTokens := []TokenWithPos{
		// line 1
		{
			token: GSX_OPEN_ELEMENT,
			pos:   NewPosition(1, 1),
		},
		{
			token: NewIdent("Parent"),
			pos:   NewPosition(1, 2),
		},
		{
			token: GSX_CLOSE_ELEMENT,
			pos:   NewPosition(1, 8),
		},
		// line 2
		{
			token: GSX_OPEN_ELEMENT,
			pos:   NewPosition(2, 5),
		},
		{
			token: NewIdent("Child"),
			pos:   NewPosition(2, 6),
		},
		{
			token: GSX_CLOSE_SELFCLOSE_ELEMENT,
			pos:   NewPosition(2, 12),
		},
		// line 3
		{
			token: GSX_OPEN_CLOSING_ELEMENT,
			pos:   NewPosition(3, 1),
		},
		{
			token: NewIdent("Parent"),
			pos:   NewPosition(3, 3),
		},
		{
			token: GSX_CLOSE_ELEMENT,
			pos:   NewPosition(3, 9),
		},
	}

	checkExpectedTokens(t, lexer, expTokens)
}

func checkExpectedTokens(t *testing.T, lexer *Lexer, expectedTokens []TokenWithPos) {
	for _, expTokenWithPos := range expectedTokens {
		expToken := expTokenWithPos.token
		expPos := expTokenWithPos.pos
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