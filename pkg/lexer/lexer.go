package lexer

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)


type Token struct {
	Name  string
	Rune rune
	Value string
}

var (
	// Control tokens
	NOTFOUND = Token{Name: "NOTFOUND"}
	CONTINUELOOP = Token{Name: "CONTINUELOOP"}

	// language tokens
	EOF     = Token{Name: "EOF" }
	NONE    = Token{Name: "NONE"}
	ILLEGAL = Token{Name: "ILLEGAL" }
	GSX_OPEN_ELEMENT  = Token{Name: "GSX_OPEN_ELEMENT", Rune: '<'}
	GSX_CLOSE_ELEMENT  = Token{Name: "GSX_CLOSE_ELEMENT", Rune: '>'}
	GSX_OPEN_CLOSING_ELEMENT  = Token{Name: "GSX_CLOSE_ELEMENT", Rune: '<', Value: "</"}
	GSX_CLOSE_SELFCLOSE_ELEMENT  = Token{Name: "GSX_CLOSE_SELFCLOSE_ELEMENT", Rune: '/', Value: "/>"}
	GSX_IDENT = Token{Name: "IDENT"}
)

func NewIdent(value string) Token {
	token := GSX_IDENT
	token.Value = value
	return token
}

type Position struct {
	line   int
	column int
}

func NewPosition(line int, column int) Position {
	return Position{line: line, column: column}
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		pos:    Position{line: 1, column: 0},
		reader: bufio.NewReader(reader),
	}
}

func (l *Lexer) Next() (Position, Token, error) {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return l.pos, EOF, nil
			}

			// at this point there isn't much we can do, and the compiler
			// should just return the raw error to the user
			return l.pos, NONE, err
		}

		// update the column to the position of the newly read in rune
		l.pos.column++

		switch r {
		case '\n':
			l.nextLine()
		case GSX_OPEN_ELEMENT.Rune:
			return chain(l).
				backup().
				then(func() (Position, Token, error) { return l.lexGsxOpen() }).
				do()
		case GSX_CLOSE_SELFCLOSE_ELEMENT.Rune:
			return chain(l).
				backup().
				then(func() (Position, Token, error) { return l.lexGsxCloseSelfClose() }).
				do()
		case GSX_CLOSE_ELEMENT.Rune:
			return l.pos, GSX_CLOSE_ELEMENT, nil
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsLetter(r) {
				return chain(l).
					backup().
					then(func() (Position, Token, error) { return l.lexIdent() }).
					do()
			} else {
				return l.pos, ILLEGAL, fmt.Errorf("Encountered illegal token %c", r)
			}
		}
	}
}

func (l *Lexer) backup() (Position, Token, error) {
	if err := l.reader.UnreadRune(); err != nil {
		return l.pos, NONE, fmt.Errorf("Could not backup during lexing: %v", err)
	}
	
	l.pos.column--
	return l.pos, NONE, nil
}

func (l *Lexer) nextLine() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) lexGsxOpen() (Position, Token, error) {

	// read the <
	r, pos, token, err := l.readRune()
	if err != nil {
		return pos, token, err
	}
	if token == EOF {
		return pos, token, err
	}
	startPos := pos

	if r != GSX_OPEN_ELEMENT.Rune {
		return l.backup()
	}

	// read possible /
	r, pos, token, err = l.readRune()
	if err != nil {
		return pos, token, err
	}
	if token == EOF {
		return pos, token, err
	}

	if r != '/' {
		pos, token, err := l.backup()
		if err != nil {
			return pos, token, err
		}

		return l.pos, GSX_OPEN_ELEMENT, nil
	}

	return startPos, GSX_OPEN_CLOSING_ELEMENT, nil
}

func (l *Lexer) lexGsxCloseSelfClose() (Position, Token, error) {

	// read the /
	r, pos, token, err := l.readRune()
	if err != nil {
		return pos, token, err
	}
	if token == EOF {
		return pos, token, err
	}
	startPos := pos

	if r != GSX_CLOSE_SELFCLOSE_ELEMENT.Rune {
		return l.backup()
	}

	// read possible >
	r, pos, token, err = l.readRune()
	if err != nil {
		return pos, token, err
	}
	if token == EOF {
		return pos, token, err
	}

	if r == '>' {
		return startPos, GSX_CLOSE_SELFCLOSE_ELEMENT, nil
	}

	return startPos, ILLEGAL, nil
}

func (l *Lexer) lexIdent() (Position, Token, error) {
	value := ""
	var startPos Position

	return l.forEachRune(
		func(r rune) (Position, Token, error) {
			if startPos.line == 0 {
				startPos = NewPosition(l.pos.line, l.pos.column)
			}

			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
				value = value + string(r)
			} else {
				// scanned something not in the identifier
				l.backup()
				return startPos, NewIdent(value), nil
			}

			return l.pos, CONTINUELOOP, nil
		},
		func() (Position, Token, error) {
			return startPos, NewIdent(value), nil
		},
	)
}

func (l *Lexer) readRune() (rune, Position, Token, error) {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			return 'a', l.pos, EOF, nil
		}

		// at this point there isn't much we can do, and the compiler
		// should just return the raw error to the user
		return 'a', l.pos, NONE, err
	}

	// update the column to the position of the newly read in rune
	l.pos.column++

	return r, l.pos, NONE, nil
}

func (l *Lexer) forEachRune(
	handleRune func(r rune) (Position, Token, error),
	handleEOF func() (Position, Token, error)) (Position, Token, error){

	for {
		r, pos, token, err := l.readRune()
		if err != nil {
			return pos, token, err
		}

		if token == EOF {
			return handleEOF()
		}

		pos, token, err = handleRune(r)
		if err != nil {
			return pos, token, err
		}
		if token.Value != CONTINUELOOP.Value {
			return pos, token, err
		}
	}
}


type LexerChain struct {
	lexer *Lexer
	actions []func() (Position, Token, error)
}

func chain(lexer *Lexer) *LexerChain {
	return &LexerChain{
		lexer: lexer,
		actions: make([]func() (Position, Token, error), 0),
	}
}

func (c *LexerChain) then(f func() (Position, Token, error)) *LexerChain {
	c.actions = append(c.actions, f)

	return c
}

func (c *LexerChain) do() (Position, Token, error) {
	pos := Position{}
	token := NONE
	err := error(nil)
	for _, action := range c.actions {
		pos, token, err = action()
		if err != nil {
			return pos, token, err
		}
	}

	return pos, token, err
}

func (c *LexerChain) backup() *LexerChain {
	c.then(func() (Position, Token, error) { return c.lexer.backup() })

	return c
}

func (c *LexerChain) forEachRune(
	handleRune func(r rune) (Position, Token, error),
	handleEOF func() (Position, Token, error)) *LexerChain {

	c.then(func() (Position, Token, error) { return c.lexer.forEachRune(handleRune, handleEOF) })

	return c
}

