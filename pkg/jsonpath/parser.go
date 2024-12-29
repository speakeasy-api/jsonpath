package jsonpath

import (
	"errors"
	"fmt"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath/token"
	"strconv"
)

// Parser represents a JSONPath parser.
type Parser struct {
	tokenizer *token.Tokenizer
	tokens    []token.TokenInfo
	path      JsonPathQuery
	current   int
}

// NewParser creates a new Parser with the given tokens.
func NewParser(tokenizer *token.Tokenizer, tokens []token.TokenInfo) *Parser {
	return &Parser{tokenizer, tokens, JsonPathQuery{}, 0}
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
	} else if p.tokens[p.current].Token == token.INTEGER {
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
		return p.parseFilterSelector()
	}

	return nil, p.parseFailure(p.tokens[p.current], "unexpected token when parsing selector")
}

func (p *Parser) parseSliceSelector() (*Slice, error) {
	// slice-selector = [start S] ":" S [end S] [":" [S step]]
	var start, end, step *int

	// Parse the start index
	if p.tokens[p.current].Token == token.INTEGER {
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
	if p.tokens[p.current].Token == token.INTEGER {
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
		if p.tokens[p.current].Token == token.INTEGER {
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
	if p.tokens[p.current].Token != token.FILTER {
		return nil, p.parseFailure(p.tokens[p.current], "expected '?'")
	}
	p.current++

	if p.tokens[p.current].Token != token.PAREN_LEFT {
		return nil, p.parseFailure(p.tokens[p.current], "expected '('")
	}
	p.current++

	expr, err := p.parseLogicalOrExpr()
	if err != nil {
		return nil, err
	}

	if p.tokens[p.current].Token != token.PAREN_RIGHT {
		return nil, p.parseFailure(p.tokens[p.current], "expected ')'")
	}
	p.current++

	return &Selector{Kind: SelectorSubKindFilter, filter: &FilterSelector{expr}}, nil
}

func (p *Parser) parseLogicalOrExpr() (*LogicalOrExpr, error) {
	var expr LogicalOrExpr

	for {
		andExpr, err := p.parseLogicalAndExpr()
		if err != nil {
			return nil, err
		}
		expr.Expressions = append(expr.Expressions, andExpr)

		if p.tokens[p.current].Token != token.OR {
			break
		}
		p.current++
	}

	return &expr, nil
}

func (p *Parser) parseLogicalAndExpr() (*LogicalAndExpr, error) {
	var expr LogicalAndExpr

	for {
		basicExpr, err := p.parseBasicExpr()
		if err != nil {
			return nil, err
		}
		expr.Expressions = append(expr.Expressions, basicExpr)

		if p.tokens[p.current].Token != token.AND {
			break
		}
		p.current++
	}

	return &expr, nil
}

func (p *Parser) parseBasicExpr() (*BasicExpr, error) {
	//basic-expr          = paren-expr /
	//	                    comparison-expr /
	//                      test-expr

	switch p.tokens[p.current].Token {
	case token.NOT:
		p.current++
		if p.tokens[p.current].Token == token.PAREN_LEFT {
			p.current++
			expr, err := p.parseLogicalOrExpr()
			if err != nil {
				return nil, err
			}
			if p.tokens[p.current].Token != token.PAREN_RIGHT {
				return nil, p.parseFailure(p.tokens[p.current], "expected ')'")
			}
			p.current++
			return &BasicExpr{ParenExpr: &ParenExpr{Not: true, Expr: expr}}, nil
		}
		return nil, p.parseFailure(p.tokens[p.current], "expected '(' after '!'")
	case token.PAREN_LEFT:
		p.current++
		expr, err := p.parseLogicalOrExpr()
		if err != nil {
			return nil, err
		}
		if p.tokens[p.current].Token != token.PAREN_RIGHT {
			return nil, p.parseFailure(p.tokens[p.current], "expected ')'")
		}
		p.current++
		return &BasicExpr{ParenExpr: &ParenExpr{Not: false, Expr: expr}}, nil
	}
	prevCurrent := p.current
	comparisonExpr, comparisonErr := p.parseComparisonExpr()
	if comparisonErr == nil {
		return &BasicExpr{ComparisonExpr: comparisonExpr}, nil
	}
	p.current = prevCurrent
	testExpr, testErr := p.parseTestExpr()
	if testErr == nil {
		return &BasicExpr{TestExpr: testExpr}, nil
	}
	p.current = prevCurrent
	return nil, p.parseFailure(p.tokens[p.current], fmt.Sprintf("could not parse query: expected either TestExpr [err: %s] or ComparisonExpr: [err: %s]", testErr.Error(), comparisonErr.Error()))
}

func (p *Parser) parseComparisonExpr() (*ComparisonExpr, error) {
	left, err := p.parseComparable()
	if err != nil {
		return nil, err
	}

	if !p.isComparisonOperator(p.tokens[p.current].Token) {
		return nil, p.parseFailure(p.tokens[p.current], "expected comparison operator")
	}
	operator := p.tokens[p.current].Token
	var op ComparisonOperator
	switch operator {
	case token.EQ:
		op = EqualTo
	case token.NE:
		op = NotEqualTo
	case token.LT:
		op = LessThan
	case token.LE:
		op = LessThanEqualTo
	case token.GT:
		op = GreaterThan
	case token.GE:
		op = GreaterThanEqualTo
	default:
		return nil, p.parseFailure(p.tokens[p.current], "expected comparison operator")
	}
	p.current++

	right, err := p.parseComparable()
	if err != nil {
		return nil, err
	}

	return &ComparisonExpr{Left: left, Op: op, Right: right}, nil
}

func (p *Parser) parseComparable() (*Comparable, error) {
	switch p.tokens[p.current].Token {
	case token.STRING_LITERAL:
		literal := p.tokens[p.current].Literal
		p.current++
		return &Comparable{&Literal{String: &literal}, nil}, nil
	case token.INTEGER:
		literal := p.tokens[p.current].Literal
		p.current++
		i, err := strconv.Atoi(literal)
		if err != nil {
			return nil, p.parseFailure(p.tokens[p.current], "expected integer")
		}
		return &Comparable{Literal: &Literal{Integer: &i}}, nil
	case token.FLOAT:
		literal := p.tokens[p.current].Literal
		p.current++
		f, err := strconv.ParseFloat(literal, 64)
		if err != nil {
			return nil, p.parseFailure(p.tokens[p.current], "expected float")
		}
		return &Comparable{Literal: &Literal{Float64: &f}}, nil
	case token.TRUE:
		p.current++
		res := true
		return &Comparable{Literal: &Literal{Bool: &res}}, nil
	case token.FALSE:
		p.current++
		res := false
		return &Comparable{Literal: &Literal{Bool: &res}}, nil
	case token.NULL:
		p.current++
		res := true
		return &Comparable{Literal: &Literal{Null: &res}}, nil
	case token.ROOT:
		p.current++
		query, err := p.parseSingleQuery()
		if err != nil {
			return nil, err
		}
		return &Comparable{SingularQuery: &SingularQuery{AbsQuery: &AbsQuery{Segments: query.Segments}}}, nil
	case token.CURRENT:
		p.current++
		query, err := p.parseSingleQuery()
		if err != nil {
			return nil, err
		}
		return &Comparable{SingularQuery: &SingularQuery{RelQuery: &RelQuery{Segments: query.Segments}}}, nil
	default:
		return nil, p.parseFailure(p.tokens[p.current], "expected literal or query")
	}
}

func (p *Parser) parseQuery() (*JsonPathQuery, error) {
	var query JsonPathQuery
	for p.current < len(p.tokens) {
		segment, err := p.parseSegment()
		if err != nil {
			return nil, err
		}
		query.Segments = append(query.Segments, segment)
	}
	return &query, nil
}

func (p *Parser) parseTestExpr() (*TestExpr, error) {
	//test-expr           = [logical-not-op S]
	//                  (filter-query / ; existence/non-existence
	//                   function-expr) ; LogicalType or NodesType
	//filter-query        = rel-query / jsonpath-query
	//rel-query           = current-node-identifier segments
	//current-node-identifier = "@"
	not := false
	if p.tokens[p.current].Token == token.NOT {
		not = true
		p.current++
	}
	switch p.tokens[p.current].Token {
	case token.CURRENT:
		query, err := p.parseQuery()
		if err != nil {
			return nil, err
		}
		return &TestExpr{FilterQuery: &FilterQuery{RelQuery: &RelQuery{Segments: query.Segments}}, Not: not}, nil
	case token.ROOT:
		query, err := p.parseQuery()
		if err != nil {
			return nil, err
		}
		return &TestExpr{FilterQuery: &FilterQuery{JsonPathQuery: &JsonPathQuery{Segments: query.Segments}}, Not: not}, nil
	default:
		funcExpr, err := p.parseFunctionExpr()
		if err != nil {
			return nil, err
		}
		return &TestExpr{FunctionExpr: funcExpr, Not: not}, nil
	}

	return nil, p.parseFailure(p.tokens[p.current], "unexpected token when parsing test expression")
}

func (p *Parser) parseFunctionExpr() (*FunctionExpr, error) {
	return nil, p.parseFailure(p.tokens[p.current], "unimplemented function expr")
}

func (p *Parser) parseSingleQuery() (*JsonPathQuery, error) {
	var query JsonPathQuery
	for p.current < len(p.tokens) {
		try := p.current
		segment, err := p.parseSegment()
		if err != nil {
			// rollback
			p.current = try
			break
		}
		query.Segments = append(query.Segments, segment)
	}
	return &query, nil
}
