package pkg

import (
	"errors"
	"fmt"
)

var (
	ParseError = errors.New("parse error")
)

// Parser represents a JSONPath parser.
type Parser struct {
	tokenizer *Tokenizer
	tokens    []TokenInfo
	root      Node
	current   int
}

// NewParser creates a new Parser with the given tokens.
func NewParser(tokenizer *Tokenizer) *Parser {
	return &Parser{tokenizer: tokenizer, tokens: tokenizer.Tokenize()}
}

// Parse parses the JSONPath tokens and returns the root node of the AST.
//
//	jsonpath-query      = root-identifier segments
func (p *Parser) Parse() (Node, error) {
	if len(p.tokens) == 0 {
		return nil, fmt.Errorf("empty JSONPath expression")
	}

	if p.root != nil {
		return p.root, nil
	}

	if p.tokens[p.current].Token != ROOT {
		return nil, fmt.Errorf("%s: %w", p.tokenizer.ErrorString(p.tokens[p.current], fmt.Sprintf("unexpected token (expected '$')")), ParseError)
	}
	root, err := p.parseRoot()
	for p.current < len(p.tokens) {
		p.parseSegment()
	}

	if p.current < len(p.tokens) {
		return nil, fmt.Errorf(p.tokenizer.ErrorString(p.tokens[p.current], fmt.Sprintf("unexpected token")))
	}

	p.root = root

	return root, nil
}

// parseRoot parses the root node.
func (p *Parser) parseRoot() (*RootNode, error) {
	node := &RootNode{Dollar: p.tokens[p.current]}
	p.current++
	return node, nil
}

// parseCurrent parses the current node.
func (p *Parser) parseCurrent() (*CurrentNode, error) {
	node := &CurrentNode{At: p.tokens[p.current]}
	p.current++
	return node, nil
}

// parseIdentifier parses an identifier node.
func (p *Parser) parseIdentifier() (*IdentifierNode, error) {
	node := &IdentifierNode{Name: p.tokens[p.current]}
	p.current++
	return node, nil
}

// parseWildcard parses a wildcard node.
func (p *Parser) parseWildcard() (*WildcardNode, error) {
	node := &WildcardNode{Star: p.tokens[p.current]}
	p.current++
	return node, nil
}

// parseRecursiveDescent parses a recursive descent node.
func (p *Parser) parseRecursiveDescent() (*RecursiveDescentNode, error) {
	node := &RecursiveDescentNode{DoubleDot: p.tokens[p.current]}
	p.current++
	return node, nil
}

// parseSubscriptOrFilter parses a subscript or filter node.
func (p *Parser) parseSubscriptOrFilter() (Expr, error) {
	if p.peek(COLON) || p.peek(NUMBER) {
		return p.parseSlice()
	} else if p.peek(WILDCARD) || p.peek(STRING_LITERAL) || p.peek(STRING) || p.peek(NUMBER) || p.peek(BOOLEAN) || p.peek(NULL) {
		return p.parseSubscript()
	} else if p.peek(PAREN_LEFT) || p.peek(CURRENT) || p.peek(ROOT) || p.peek(RECURSIVE) {
		return p.parseFilter()
	} else {
		return nil, fmt.Errorf("unexpected token %s at line %d, column %d",
			p.tokens[p.current].Literal, p.tokens[p.current].Line, p.tokens[p.current].Column)
	}
}

// parseSubscript parses a subscript node.
func (p *Parser) parseSubscript() (*SubscriptNode, error) {
	node := &SubscriptNode{Lbrack: p.tokens[p.current]}
	p.current++

	index, err := p.parseRootQuery()
	if err != nil {
		return nil, err
	}
	node.Index = index

	if !p.expect(BRACKET_RIGHT) {
		return nil, fmt.Errorf("expected ']' at line %d, column %d",
			p.tokens[p.current].Line, p.tokens[p.current].Column)
	}
	node.Rbrack = p.tokens[p.current]
	p.current++

	return node, nil
}

// parseSlice parses a slice node.
func (p *Parser) parseSlice() (*SliceNode, error) {
	node := &SliceNode{Lbrack: p.tokens[p.current]}
	p.current++

	if p.peek(COLON) {
		node.Start = nil
	} else {
		start, err := p.parseRootQuery()
		if err != nil {
			return nil, err
		}
		node.Start = start
	}

	if !p.expect(COLON) {
		return nil, fmt.Errorf("expected ':' at line %d, column %d",
			p.tokens[p.current].Line, p.tokens[p.current].Column)
	}
	node.Colon1 = p.tokens[p.current]
	p.current++

	if p.peek(COLON) || p.peek(BRACKET_RIGHT) {
		node.Finish = nil
	} else {
		end, err := p.parseRootQuery()
		if err != nil {
			return nil, err
		}
		node.Finish = end
	}

	if p.peek(COLON) {
		node.Colon2 = p.tokens[p.current]
		p.current++

		step, err := p.parseRootQuery()
		if err != nil {
			return nil, err
		}
		node.Step = step
	}

	if !p.expect(BRACKET_RIGHT) {
		return nil, fmt.Errorf("expected ']' at line %d, column %d",
			p.tokens[p.current].Line, p.tokens[p.current].Column)
	}
	node.Rbrack = p.tokens[p.current]
	p.current++

	return node, nil
}

// parseFilter parses a filter node.
func (p *Parser) parseFilter() (*FilterNode, error) {
	node := &FilterNode{Lbrack: p.tokens[p.current]}
	p.current++

	expr, err := p.parseRootQuery()
	if err != nil {
		return nil, err
	}
	node.Expr = expr

	if !p.expect(BRACKET_RIGHT) {
		return nil, fmt.Errorf("expected ']' at line %d, column %d",
			p.tokens[p.current].Line, p.tokens[p.current].Column)
	}
	node.Rbrack = p.tokens[p.current]
	p.current++

	return node, nil
}

// parseFunctionCall parses a function call node.
func (p *Parser) parseFunctionCall() (*FunctionCallNode, error) {
	node := &FunctionCallNode{Name: p.tokens[p.current]}
	p.current++

	if !p.expect(PAREN_LEFT) {
		return nil, fmt.Errorf("expected '(' at line %d, column %d",
			p.tokens[p.current].Line, p.tokens[p.current].Column)
	}
	node.Lparen = p.tokens[p.current]
	p.current++

	for !p.peek(PAREN_RIGHT) {
		arg, err := p.parseRootQuery()
		if err != nil {
			return nil, err
		}
		node.Args = append(node.Args, arg)

		if p.peek(COMMA) {
			p.current++
		} else if !p.peek(PAREN_RIGHT) {
			return nil, fmt.Errorf("expected ',' or ')' at line %d, column %d",
				p.tokens[p.current].Line, p.tokens[p.current].Column)
		}
	}

	if !p.expect(PAREN_RIGHT) {
		return nil, fmt.Errorf("expected ')' at line %d, column %d",
			p.tokens[p.current].Line, p.tokens[p.current].Column)
	}
	node.Rparen = p.tokens[p.current]
	p.current++

	return node, nil
}

// parseComparison parses a comparison node.
func (p *Parser) parseComparison() (*ComparisonNode, error) {
	lhs, err := p.parseRootQuery()
	if err != nil {
		return nil, err
	}

	if !p.isComparisonOperator(p.tokens[p.current].Token) {
		return nil, fmt.Errorf("expected comparison operator at line %d, column %d",
			p.tokens[p.current].Line, p.tokens[p.current].Column)
	}
	operator := p.tokens[p.current]
	p.current++

	rhs, err := p.parseRootQuery()
	if err != nil {
		return nil, err
	}

	return &ComparisonNode{Lhs: lhs, Operator: operator, Rhs: rhs}, nil
}

// parseLiteral parses a literal node (boolean, number, string, or null).
func (p *Parser) parseLiteral() (Expr, error) {
	switch p.tokens[p.current].Token {
	case BOOLEAN:
		return &BooleanNode{Value: p.tokens[p.current]}, nil
	case NUMBER:
		return &NumberNode{Value: p.tokens[p.current]}, nil
	case STRING:
		return &StringNode{Value: p.tokens[p.current]}, nil
	case NULL:
		return &NullNode{Null: p.tokens[p.current]}, nil
	default:
		return nil, fmt.Errorf("unexpected token %s at line %d, column %d",
			p.tokens[p.current].Literal, p.tokens[p.current].Line, p.tokens[p.current].Column)
	}
}

// peek returns true if the current token matches the given token type.
func (p *Parser) peek(token Token) bool {
	return p.current < len(p.tokens) && p.tokens[p.current].Token == token
}

// expect consumes the current token if it matches the given token type.
func (p *Parser) expect(token Token) bool {
	if p.peek(token) {
		p.current++
		return true
	}
	return false
}

// isComparisonOperator returns true if the given token is a comparison operator.
func (p *Parser) isComparisonOperator(token Token) bool {
	return token == EQ || token == NE || token == GT || token == GE || token == LT || token == LE
}
