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
	EOF     = Token{Name: "EOF" }
	NONE    = Token{Name: "NONE"}
	ILLEGAL = Token{Name: "ILLEGAL" }
	GSX_OPEN_ELEMENT  = Token{Name: "GSX_OPEN_ELEMENT", Rune: '<'}
	GSX_CLOSE_ELEMENT  = Token{Name: "GSX_CLOSE_ELEMENT", Rune: '>'}
	IDENT = Token{Name: "IDENT"}
)

func NewIdent(value string) Token {
	token := IDENT
	token.Value = value
	return token
}

type Position struct {
	line   int
	column int
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
			return l.pos, GSX_OPEN_ELEMENT, nil				
		case GSX_CLOSE_ELEMENT.Rune:
			return l.pos, GSX_CLOSE_ELEMENT, nil
		default:
			if unicode.IsSpace(r) {
				continue // nothing to do here, just move on
			} else if unicode.IsLetter(r) {
				// backup and let lexIdent rescan the beginning of the ident
				startPos := l.pos
				err = l.backup()
				if err != nil {
					return l.pos, NONE, err
				}
				token, err := l.lexIdent()
				if err != nil {
					return l.pos, NONE, err
				}
				return startPos, token, nil
			} else {
				return l.pos, ILLEGAL, fmt.Errorf("Encountered illegal token %c", r)
			}
		}
	}
}

func (l *Lexer) backup() error {
	if err := l.reader.UnreadRune(); err != nil {
		return err
	}
	
	l.pos.column--
	return nil
}

func (l *Lexer) nextLine() {
	l.pos.line++
	l.pos.column = 0
}

func (l *Lexer) lexIdent() (Token, error) {
	var value = ""
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return NewIdent(value), nil
			}

			return NONE, err
		}

		l.pos.column++
		// take every rune that is allowed in an identifier
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			value = value + string(r)
		} else {
			// scanned something not in the identifier
			l.backup()
			return NewIdent(value), nil
		}
	}
}

