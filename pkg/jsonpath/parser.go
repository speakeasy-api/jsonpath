package jsonpath

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
	segments  []Segment
	current   int
}

// NewParser creates a new Parser with the given tokens.
func NewParser(tokenizer *Tokenizer) *Parser {
	return &Parser{tokenizer: tokenizer, tokens: tokenizer.Tokenize()}
}

// Parse parses the JSONPath tokens and returns the root node of the AST.
//
//	jsonpath-query      = root-identifier segments
func (p *Parser) Parse() error {
	if len(p.tokens) == 0 {
		return fmt.Errorf("empty JSONPath expression")
	}

	if p.tokens[p.current].Token != ROOT {
		return p.parseFailure(p.tokens[p.current], "expected '$'")
	}
	p.current++

	for p.current < len(p.tokens) {
		segment, err := p.parseSegment()
		if err != nil {
			return err
		}
		p.segments = append(p.segments, segment)
	}
	return nil
}

// parseDescendantSegment parses a descendant segment (preceded by "..").
func (p *Parser) parseDescendantSegment() (*DescendantSegment, error) {
	if p.tokens[p.current].Token != RECURSIVE {
		return nil, p.parseFailure(p.tokens[p.current], "expected '..'")
	}

	// three kinds of descendant segments
	// "..*" -> recursive wildcard
	// "..STRING_LITERAL" -> recursive hunt for an identifier
	// "..[INNER_SEGMENT]" -> recursive hunt for some inner expression (e.g. a filter expression)
	nextToken := p.tokens[p.current+1]
	if nextToken.Token == WILDCARD {
		node := &DescendantSegment{SubKind: DescendantWildcardSelector}
		p.current += 2
		return node, nil
	} else if nextToken.Token == STRING_LITERAL {
		node := &DescendantSegment{SubKind: DescendantDotNameSelector}
		p.current += 2
		return node, nil
	} else if nextToken.Token == BRACKET_LEFT {
		node := &DescendantSegment{SubKind: DescendantLongSelector}
		previousCurrent := p.current
		p.current += 2
		innerSegment, err := p.parseSegment()
		node.LongFormInner = innerSegment
		if err != nil {
			p.current = previousCurrent
			return nil, err
		}
		if p.tokens[p.current].Token != BRACKET_RIGHT {
			p.current = previousCurrent
			return nil, p.parseFailure(p.tokens[p.current], "expected ']'")
		}
		return node, nil
	}

	return nil, p.parseFailure(nextToken, "unexpected descendant segment")
}

func (p *Parser) parseFailure(target TokenInfo, msg string) error {
	return errors.New(p.tokenizer.ErrorTokenString(target, msg))
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

func (p *Parser) parseSegment() (Segment, error) {
	currentToken := p.tokens[p.current]
	if currentToken.Token == RECURSIVE {
		node, err := p.parseDescendantSegment()
		if err != nil {
			return nil, err
		}
		return node, nil
	}
	return p.parseChildSegment()
}

func (p *Parser) parseChildSegment() (Segment, error) {
	// .*
	// .STRING_LITERAL
	// []
	firstToken := p.tokens[p.current]
	if firstToken.Token == CHILD && p.tokens[p.current+1].Token == WILDCARD {
		p.current += 2
		return &ChildSegment{SubKind: ChildSegmentDotWildcard, Tokens: []TokenInfo{firstToken, p.tokens[p.current+1]}}, nil
	} else if firstToken.Token == CHILD && p.tokens[p.current+1].Token == STRING_LITERAL {
		p.current += 2
		return &ChildSegment{SubKind: ChildSegmentDotMemberName, Tokens: []TokenInfo{firstToken, p.tokens[p.current+1]}}, nil
	} else if firstToken.Token == CHILD && p.tokens[p.current+1].Token == BRACKET_LEFT {
		prior := p.current
		p.current += 2
		innerSegment, err := p.parseSelector()
		if err != nil {
			p.current = prior
			return nil, err
		}
		if p.tokens[p.current].Token != BRACKET_RIGHT {
			prior = p.current
			return nil, p.parseFailure(p.tokens[p.current], "expected ']'")
		}
		return &ChildSegment{SubKind: ChildSegmentSelector, Tokens: []TokenInfo{firstToken, p.tokens[p.current]}, InnerSelector: innerSegment}, nil
	}
	return nil, p.parseFailure(firstToken, "unexpected token when parsing child segment")
}

type SelectorSubKind int

const (
	SelectorSubKindWildcard SelectorSubKind = iota
	SelectorSubKindName
	SelectorSubKindArraySlice
	SelectorSubKindArrayIndex
	SelectorSubKindFilter
)

type Selector struct {
	SubKind SelectorSubKind
}

func (p *Parser) parseSelector() (*Selector, error) {
	return nil, nil
}
