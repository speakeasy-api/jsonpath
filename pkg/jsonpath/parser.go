package jsonpath

import (
	"errors"
	"fmt"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath/token"
	"strconv"
)

var (
	ParseError = errors.New("parse error")
)

// Parser represents a JSONPath parser.
type Parser struct {
	tokenizer *token.Tokenizer
	tokens    []token.TokenInfo
	path      Query
	current   int
}

// NewParser creates a new Parser with the given tokens.
func NewParser(tokenizer *token.Tokenizer, tokens []token.TokenInfo) *Parser {
	return &Parser{tokenizer, tokens, Query{}, 0}
}

// Parse parses the JSONPath tokens and returns the root node of the AST.
//
//	jsonpath-query      = root-identifier segments
func (p *Parser) Parse() error {
	if len(p.tokens) == 0 {
		return fmt.Errorf("empty JSONPath expression")
	}

	if p.tokens[p.current].Token != token.ROOT {
		return p.parseFailure(p.tokens[p.current], "expected '$'")
	}
	p.current++

	for p.current < len(p.tokens) {
		segment, err := p.parseSegment()
		if err != nil {
			return err
		}
		p.path.Segments = append(p.path.Segments, segment)
	}
	return nil
}

// parseDescendantSegment parses a descendant segment (preceded by "..").
func (p *Parser) parseDescendantSegment() (*DescendantSegment, error) {
	if p.tokens[p.current].Token != token.RECURSIVE {
		return nil, p.parseFailure(p.tokens[p.current], "expected '..'")
	}

	// three kinds of descendant segments
	// "..*" -> recursive wildcard
	// "..STRING" -> recursive hunt for an identifier
	// "..[INNER_SEGMENT]" -> recursive hunt for some inner expression (e.g. a filter expression)
	nextToken := p.tokens[p.current+1]
	if nextToken.Token == token.WILDCARD {
		node := &DescendantSegment{SubKind: DescendantSegmentSubKindWildcard}
		p.current += 2
		return node, nil
	} else if nextToken.Token == token.STRING {
		node := &DescendantSegment{SubKind: DescendantSegmentSubKindDotName}
		p.current += 2
		return node, nil
	} else if nextToken.Token == token.BRACKET_LEFT {
		node := &DescendantSegment{SubKind: DescendantSegmentSubKindLongHand}
		previousCurrent := p.current
		p.current += 2
		innerSegment, err := p.parseSegment()
		node.innerSegment = innerSegment
		if err != nil {
			p.current = previousCurrent
			return nil, err
		}

		if p.tokens[p.current].Token != token.BRACKET_RIGHT {
			p.current = previousCurrent
			return nil, p.parseFailure(p.tokens[p.current], "expected ']'")
		}
		p.current += 1
		return node, nil
	}

	return nil, p.parseFailure(nextToken, "unexpected descendant segment")
}

func (p *Parser) parseFailure(target token.TokenInfo, msg string) error {
	return errors.New(p.tokenizer.ErrorTokenString(target, msg))
}

// peek returns true if the upcoming token matches the given token type.
func (p *Parser) peek(token token.Token) bool {
	return p.current < len(p.tokens) && p.tokens[p.current+1].Token == token
}

// expect consumes the current token if it matches the given token type.
func (p *Parser) expect(token token.Token) bool {
	if p.peek(token) {
		p.current++
		return true
	}
	return false
}

// isComparisonOperator returns true if the given token is a comparison operator.
func (p *Parser) isComparisonOperator(tok token.Token) bool {
	return tok == token.EQ || tok == token.NE || tok == token.GT || tok == token.GE || tok == token.LT || tok == token.LE
}

func (p *Parser) parseSegment() (*Segment, error) {
	currentToken := p.tokens[p.current]
	if currentToken.Token == token.RECURSIVE {
		node, err := p.parseDescendantSegment()
		if err != nil {
			return nil, err
		}
		return &Segment{Descendant: node}, nil
	}
	return p.parseChildSegment()
}

func (p *Parser) parseChildSegment() (*Segment, error) {
	// .*
	// .STRING
	// []
	firstToken := p.tokens[p.current]
	if firstToken.Token == token.CHILD && p.tokens[p.current+1].Token == token.WILDCARD {
		p.current += 2
		return &Segment{&ChildSegment{ChildSegmentDotWildcard, "", nil}, nil}, nil
	} else if firstToken.Token == token.CHILD && p.tokens[p.current+1].Token == token.STRING {
		dotName := p.tokens[p.current+1].Literal
		p.current += 2
		return &Segment{&ChildSegment{ChildSegmentDotMemberName, dotName, nil}, nil}, nil
	} else if firstToken.Token == token.BRACKET_LEFT {
		prior := p.current
		p.current += 1
		innerSegment, err := p.parseSelector()
		if err != nil {
			p.current = prior
			return nil, err
		}
		if p.tokens[p.current].Token != token.BRACKET_RIGHT {
			prior = p.current
			return nil, p.parseFailure(p.tokens[p.current], "expected ']'")
		}
		p.current += 1
		return &Segment{&ChildSegment{kind: ChildSegmentLongHand, dotName: "", selectors: []*Selector{innerSegment}}, nil}, nil
	}
	return nil, p.parseFailure(firstToken, "unexpected token when parsing child segment")
}

func (p *Parser) parseSelector() (*Selector, error) {
	//selector            = name-selector /
	//                      wildcard-selector /
	//                      slice-selector /
	//                      index-selector /
	//                      filter-selector

	//    name-selector       = string-literal
	if p.tokens[p.current].Token == token.STRING_LITERAL {
		name := p.tokens[p.current].Literal
		p.current++
		return &Selector{Kind: SelectorSubKindName, name: name}, nil
		//    wildcard-selector   = "*"
	} else if p.tokens[p.current].Token == token.WILDCARD {
		p.current++
		return &Selector{Kind: SelectorSubKindWildcard}, nil
	} else if p.tokens[p.current].Token == token.NUMBER {
		// peek ahead to see if it's a slice
		if p.peek(token.ARRAY_SLICE) {
			slice, err := p.parseSliceSelector()
			if err != nil {
				return nil, err
			}
			return &Selector{Kind: SelectorSubKindArraySlice, slice: slice}, nil
		}
		// else it's an index
		literal := p.tokens[p.current].Literal
		// make sure literal is an integer
		i, err := strconv.ParseInt(literal, 10, 64)
		if err != nil {
			return nil, p.parseFailure(p.tokens[p.current], "expected an integer")
		}
		p.current++

		return &Selector{Kind: SelectorSubKindArrayIndex, index: int(i)}, nil
	} else if p.tokens[p.current].Token == token.ARRAY_SLICE {
		slice, err := p.parseSliceSelector()
		if err != nil {
			return nil, err
		}
		return &Selector{Kind: SelectorSubKindArraySlice, slice: slice}, nil
	} else if p.tokens[p.current].Token == token.FILTER {
		p.current++
		return p.parseFilterSelector()
	}

	return nil, p.parseFailure(p.tokens[p.current], "unexpected token when parsing selector")
}

func (p *Parser) parseSliceSelector() (*Slice, error) {
	// slice-selector = [start S] ":" S [end S] [":" [S step]]
	var start, end, step *int

	// Parse the start index
	if p.tokens[p.current].Token == token.NUMBER {
		literal := p.tokens[p.current].Literal
		i, err := strconv.Atoi(literal)
		if err != nil {
			return nil, p.parseFailure(p.tokens[p.current], "expected an integer")
		}
		start = &i
		p.current += 1
	}

	// Expect a colon
	if p.tokens[p.current].Token != token.ARRAY_SLICE {
		return nil, p.parseFailure(p.tokens[p.current], "expected ':'")
	}
	p.current++

	// Parse the end index
	if p.tokens[p.current].Token == token.NUMBER {
		literal := p.tokens[p.current].Literal
		i, err := strconv.Atoi(literal)
		if err != nil {
			return nil, p.parseFailure(p.tokens[p.current], "expected an integer")
		}
		end = &i
		p.current++
	}

	// Check for an optional second colon and step value
	if p.tokens[p.current].Token == token.ARRAY_SLICE {
		p.current++
		if p.tokens[p.current].Token == token.NUMBER {
			literal := p.tokens[p.current].Literal
			i, err := strconv.Atoi(literal)
			if err != nil {
				return nil, p.parseFailure(p.tokens[p.current], "expected an integer")
			}
			step = &i
			p.current++
		}
	}

	return &Slice{Start: start, End: end, Step: step}, nil
}
func (p *Parser) parseFilterSelector() (*Selector, error) {
	return nil, p.parseFailure(p.tokens[p.current], "unimplemented filter selector")
}
