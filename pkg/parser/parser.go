package parser

import (
	"fmt"
	"io"

	"github.com/useurmind/greact/pkg/lexer"
)

type Node struct {
	Type     string
	Parent   *Node
	Children []*Node
}

type Parser struct {
	lexer       *lexer.Lexer
	currentNode *Node
	tokenBuffer struct {
		token lexer.Token
		pos   lexer.Position
		set   bool
	}
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		lexer:       lexer.NewLexer(reader),
		currentNode: nil,
	}
}

func (p *Parser) Parse() (*Node, error) {

	node, err := p.parseNode()

	// rest of file should be "empty"

	pos, token, err := p.scanNextToken()
	if err != nil {
		return nil, err
	}

	if token.Name != lexer.EOF.Name {
		return nil, fmt.Errorf("Expected EOF after parsing root node at %v but got %v", pos, token)
	}

	return node, nil
}

func (p *Parser) parseNode() (*Node, error) {
	// parse start of element
	pos, token, err := p.scanNextToken()
	if err != nil {
		return nil, err
	}

	if token.Name == lexer.GSX_OPEN_ELEMENT.Name {
		newNode := &Node{}

		if p.currentNode != nil {
			currentNode := p.currentNode
			defer func() { p.currentNode = currentNode }()
		}

		p.currentNode = newNode
	} else {
		return nil, fmt.Errorf("Expected %s at %v", lexer.GSX_OPEN_ELEMENT.Value, pos)
	}

	pos, token, err = p.scanNextToken()
	if err != nil {
		return nil, err
	}

	if token.Name == lexer.GSX_IDENT.Name {
		p.currentNode.Type = token.Value
	} else {
		return nil, fmt.Errorf("Expected identifier at %v", pos)
	}

	// TODO parse attributes

	pos, token, err = p.scanNextToken()
	if err != nil {
		return nil, err
	}

	if token.Name == lexer.GSX_CLOSE_ELEMENT.Name {
		// continue below
	} else if token.Name == lexer.GSX_CLOSE_SELFCLOSE_ELEMENT.Name {
		// element finished because of self close
		return p.currentNode, nil
	} else {
		return nil, fmt.Errorf("Expected %s or %s at %v", lexer.GSX_CLOSE_ELEMENT.Value, lexer.GSX_CLOSE_SELFCLOSE_ELEMENT.Value, pos)
	}

	for {
		// parse children
		pos, token, err = p.scanNextToken()
		if err != nil {
			return nil, err
		}

		if token.Name == lexer.GSX_OPEN_ELEMENT.Name {
			err = p.unscanToken()
			if err != nil {
				return nil, err
			}

			child, err := p.parseNode()
			if err != nil {
				return nil, err
			}

			p.currentNode.Children = append(p.currentNode.Children, child)
			child.Parent = p.currentNode

		} else {
			err = p.unscanToken()
			if err != nil {
				return nil, err
			}
			break
		}

	}

	// parse element closing
	pos, token, err = p.scanNextToken()
	if err != nil {
		return nil, err
	}

	if token.Name == lexer.GSX_OPEN_CLOSING_ELEMENT.Name {
		// continue below
	} else {
		return nil, fmt.Errorf("Expected %s or child element at %v", lexer.GSX_OPEN_CLOSING_ELEMENT.Name, pos)
	}

	pos, token, err = p.scanNextToken()
	if err != nil {
		return nil, err
	}

	if token.Name == lexer.GSX_IDENT.Name && token.Value == p.currentNode.Type {
		// continue below
	} else {
		return nil, fmt.Errorf("Expected identifier %s at %v", p.currentNode.Type, pos)
	}

	pos, token, err = p.scanNextToken()
	if err != nil {
		return nil, err
	}

	if token.Name == lexer.GSX_CLOSE_ELEMENT.Name {
		return p.currentNode, nil
	}

	return nil, fmt.Errorf("Expected %s at %v", lexer.GSX_CLOSE_ELEMENT.Value, pos)
}

func (p *Parser) scanNextToken() (lexer.Position, lexer.Token, error) {
	if p.tokenBuffer.set {
		p.tokenBuffer.set = false
		return p.tokenBuffer.pos, p.tokenBuffer.token, nil
	}

	pos, token, err := p.lexer.Next()
	if err != nil {
		return pos, token, err
	}

	p.tokenBuffer.pos = pos
	p.tokenBuffer.token = token
	p.tokenBuffer.set = false

	return pos, token, err
}

func (p *Parser) unscanToken() error {
	if p.tokenBuffer.set {
		return fmt.Errorf("Doing unscanToken twice will lead to loss of last unscanned token and is currently not supported")
	}

	p.tokenBuffer.set = true

	return nil
}
