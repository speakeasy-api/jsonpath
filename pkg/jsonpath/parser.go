package jsonpath

import (
	"errors"
	"fmt"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath/token"
	"strconv"
	"strings"
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
		p.path.segments = append(p.path.segments, segment)
	}
	return nil
}

// parseDescendantSegment parses a descendant segment (preceded by "..").
func (p *Parser) parseDescendantSegment() (*descendantSegment, error) {
	if p.tokens[p.current].Token != token.RECURSIVE {
		return nil, p.parseFailure(p.tokens[p.current], "expected '..'")
	}

	// three kinds of descendant segments
	// "..*" -> recursive wildcard
	// "..STRING" -> recursive hunt for an identifier
	// "..[INNER_SEGMENT]" -> recursive hunt for some inner expression (e.g. a filter expression)
	if p.peek(token.WILDCARD) {
		node := &descendantSegment{subKind: descendantSegmentSubKindWildcard}
		p.current += 2
		return node, nil
	} else if p.peek(token.STRING) {
		node := &descendantSegment{subKind: descendantSegmentSubKindDotName}
		p.current += 2
		return node, nil
	} else if p.peek(token.BRACKET_LEFT) {
		node := &descendantSegment{subKind: descendantSegmentSubKindLongHand}
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

	return nil, p.parseFailure(p.tokens[p.current], "unexpected descendant segment")
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

func (p *Parser) parseSegment() (*segment, error) {
	currentToken := p.tokens[p.current]
	if currentToken.Token == token.RECURSIVE {
		node, err := p.parseDescendantSegment()
		if err != nil {
			return nil, err
		}
		return &segment{Descendant: node}, nil
	}
	return p.parseChildSegment()
}

func (p *Parser) parseChildSegment() (*segment, error) {
	// .*
	// .STRING
	// []
	firstToken := p.tokens[p.current]
	if firstToken.Token == token.CHILD && p.peek(token.WILDCARD) {
		p.current += 2
		return &segment{&childSegment{childSegmentDotWildcard, "", nil}, nil}, nil
	} else if firstToken.Token == token.CHILD && p.peek(token.STRING) {
		dotName := p.tokens[p.current+1].Literal
		p.current += 2
		return &segment{&childSegment{childSegmentDotMemberName, dotName, nil}, nil}, nil
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
		return &segment{&childSegment{kind: childSegmentLongHand, dotName: "", selectors: []*Selector{innerSegment}}, nil}, nil
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

	return &Selector{Kind: SelectorSubKindFilter, filter: &filterSelector{expr}}, nil
}

func (p *Parser) parseLogicalOrExpr() (*logicalOrExpr, error) {
	var expr logicalOrExpr

	for {
		andExpr, err := p.parseLogicalAndExpr()
		if err != nil {
			return nil, err
		}
		expr.expressions = append(expr.expressions, andExpr)

		if p.tokens[p.current].Token != token.OR {
			break
		}
		p.current++
	}

	return &expr, nil
}

func (p *Parser) parseLogicalAndExpr() (*logicalAndExpr, error) {
	var expr logicalAndExpr

	for {
		basicExpr, err := p.parseBasicExpr()
		if err != nil {
			return nil, err
		}
		expr.expressions = append(expr.expressions, basicExpr)

		if p.tokens[p.current].Token != token.AND {
			break
		}
		p.current++
	}

	return &expr, nil
}

func (p *Parser) parseBasicExpr() (*basicExpr, error) {
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
			return &basicExpr{parenExpr: &parenExpr{not: true, expr: expr}}, nil
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
		return &basicExpr{parenExpr: &parenExpr{not: false, expr: expr}}, nil
	}
	prevCurrent := p.current
	comparisonExpr, comparisonErr := p.parseComparisonExpr()
	if comparisonErr == nil {
		return &basicExpr{comparisonExpr: comparisonExpr}, nil
	}
	p.current = prevCurrent
	testExpr, testErr := p.parseTestExpr()
	if testErr == nil {
		return &basicExpr{testExpr: testExpr}, nil
	}
	p.current = prevCurrent
	return nil, p.parseFailure(p.tokens[p.current], fmt.Sprintf("could not parse query: expected either testExpr [err: %s] or comparisonExpr: [err: %s]", testErr.Error(), comparisonErr.Error()))
}

func (p *Parser) parseComparisonExpr() (*comparisonExpr, error) {
	left, err := p.parseComparable()
	if err != nil {
		return nil, err
	}

	if !p.isComparisonOperator(p.tokens[p.current].Token) {
		return nil, p.parseFailure(p.tokens[p.current], "expected comparison operator")
	}
	operator := p.tokens[p.current].Token
	var op comparisonOperator
	switch operator {
	case token.EQ:
		op = equalTo
	case token.NE:
		op = notEqualTo
	case token.LT:
		op = lessThan
	case token.LE:
		op = lessThanEqualTo
	case token.GT:
		op = greaterThan
	case token.GE:
		op = greaterThanEqualTo
	default:
		return nil, p.parseFailure(p.tokens[p.current], "expected comparison operator")
	}
	p.current++

	right, err := p.parseComparable()
	if err != nil {
		return nil, err
	}

	return &comparisonExpr{left: left, op: op, right: right}, nil
}

func (p *Parser) parseComparable() (*comparable, error) {
	//	comparable = literal /
	//	singular-query / ; singular query value
	//	function-expr    ; ValueType
	if literal, err := p.parseLiteral(); err == nil {
		return &comparable{literal: literal}, nil
	}
	if funcExpr, err := p.parseFunctionExpr(); err == nil {
		return &comparable{functionExpr: funcExpr}, nil
	}
	switch p.tokens[p.current].Token {
	case token.ROOT:
		p.current++
		query, err := p.parseSingleQuery()
		if err != nil {
			return nil, err
		}
		return &comparable{singularQuery: &singularQuery{absQuery: &absQuery{segments: query.segments}}}, nil
	case token.CURRENT:
		p.current++
		query, err := p.parseSingleQuery()
		if err != nil {
			return nil, err
		}
		return &comparable{singularQuery: &singularQuery{relQuery: &relQuery{segments: query.segments}}}, nil
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
		query.segments = append(query.segments, segment)
	}
	return &query, nil
}

func (p *Parser) parseTestExpr() (*testExpr, error) {
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
		p.current++
		query, err := p.parseSingleQuery()
		if err != nil {
			return nil, err
		}
		return &testExpr{filterQuery: &filterQuery{relQuery: &relQuery{segments: query.segments}}, not: not}, nil
	case token.ROOT:
		p.current++
		query, err := p.parseSingleQuery()
		if err != nil {
			return nil, err
		}
		return &testExpr{filterQuery: &filterQuery{jsonPathQuery: &JsonPathQuery{segments: query.segments}}, not: not}, nil
	default:
		funcExpr, err := p.parseFunctionExpr()
		if err != nil {
			return nil, err
		}
		return &testExpr{functionExpr: funcExpr, not: not}, nil
	}

	return nil, p.parseFailure(p.tokens[p.current], "unexpected token when parsing test expression")
}

func (p *Parser) parseFunctionExpr() (*functionExpr, error) {
	functionName := p.tokens[p.current].Literal
	if p.tokens[p.current+1].Token != token.PAREN_LEFT {
		return nil, p.parseFailure(p.tokens[p.current+1], "expected '('")
	}
	p.current += 2
	args := []*functionArgument{}
	for p.current < len(p.tokens) {
		arg, err := p.parseFunctionArgument()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if p.tokens[p.current].Token != token.COMMA {
			break
		}
		p.current++
	}
	if p.tokens[p.current].Token != token.PAREN_RIGHT {
		return nil, p.parseFailure(p.tokens[p.current], "expected ')'")
	}
	p.current++
	return &functionExpr{funcType: functionTypeMap[functionName], args: args}, nil
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
		query.segments = append(query.segments, segment)
	}
	//if len(query.segments) == 0 {
	//	return nil, p.parseFailure(p.tokens[p.current], "expected at least one segment")
	//}
	return &query, nil
}

func (p *Parser) parseFunctionArgument() (*functionArgument, error) {
	//function-argument   = literal /
	//	filter-query / ; (includes singular-query)
	//  logical-expr /
	//	function-expr

	if lit, err := p.parseLiteral(); err == nil {
		return &functionArgument{literal: lit}, nil
	}
	if funcExpr, err := p.parseFunctionExpr(); err == nil {
		return &functionArgument{functionExpr: funcExpr}, nil
	}
	if expr, err := p.parseLogicalOrExpr(); err == nil {
		return &functionArgument{logicalExpr: expr}, nil
	}

	switch p.tokens[p.current].Token {
	case token.CURRENT:
		p.current++
		query, err := p.parseSingleQuery()
		if err != nil {
			return nil, err
		}
		return &functionArgument{filterQuery: &filterQuery{relQuery: &relQuery{segments: query.segments}}}, nil
	case token.ROOT:
		p.current++
		query, err := p.parseSingleQuery()
		if err != nil {
			return nil, err
		}
		return &functionArgument{filterQuery: &filterQuery{jsonPathQuery: &JsonPathQuery{segments: query.segments}}}, nil
	}

	return nil, p.parseFailure(p.tokens[p.current], "unexpected token for function argument")
}

func (p *Parser) parseLiteral() (*literal, error) {
	switch p.tokens[p.current].Token {
	case token.STRING_LITERAL:
		lit := p.tokens[p.current].Literal
		p.current++
		return &literal{string: &lit}, nil
	case token.INTEGER:
		lit := p.tokens[p.current].Literal
		p.current++
		i, err := strconv.Atoi(lit)
		if err != nil {
			return nil, p.parseFailure(p.tokens[p.current], "expected integer")
		}
		return &literal{integer: &i}, nil
	case token.FLOAT:
		lit := p.tokens[p.current].Literal
		p.current++
		f, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			return nil, p.parseFailure(p.tokens[p.current], "expected float")
		}
		return &literal{float64: &f}, nil
	case token.TRUE:
		p.current++
		res := true
		return &literal{bool: &res}, nil
	case token.FALSE:
		p.current++
		res := false
		return &literal{bool: &res}, nil
	case token.NULL:
		p.current++
		res := true
		return &literal{null: &res}, nil
	}
	return nil, p.parseFailure(p.tokens[p.current], "expected literal")
}

type JsonPathQuery struct {
	// "$"
	segments []*segment
}

func (q JsonPathQuery) ToString() string {
	b := strings.Builder{}
	b.WriteString("$")
	for _, segment := range q.segments {
		b.WriteString(segment.ToString())
	}
	return b.String()
}
