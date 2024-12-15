package selector

import (
	"fmt"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath/primitive"
	"strings"
)

// Filter represents a filter selector in JSONPath
type Filter struct {
	Expression LogicalOrExpr
}

// LogicalOrExpr represents a logical OR expression
type LogicalOrExpr struct {
	Expressions []LogicalAndExpr
}

// LogicalAndExpr represents a logical AND expression
type LogicalAndExpr struct {
	Expressions []BasicExpr
}

// BasicExpr represents a basic expression in a filter
type BasicExpr interface {
	isBasicExpr()
}

// ComparisonExpr represents a comparison expression
type ComparisonExpr struct {
	Left  Comparable
	Op    ComparisonOperator
	Right Comparable
}

func (ComparisonExpr) isBasicExpr() {}

// ExistExpr represents an existence expression
type ExistExpr struct {
	Query string
}

func (ExistExpr) isBasicExpr() {}

// NotExistExpr represents a negated existence expression
type NotExistExpr struct {
	Query string
}

func (NotExistExpr) isBasicExpr() {}

// FuncExpr represents a function expression
type FuncExpr struct {
	Expr FunctionExpr
}

func (FuncExpr) isBasicExpr() {}

// NotFuncExpr represents a negated function expression
type NotFuncExpr struct {
	Expr FunctionExpr
}

func (NotFuncExpr) isBasicExpr() {}

// ParenExpr represents a parenthesized expression
type ParenExpr struct {
	Expr LogicalOrExpr
}

func (ParenExpr) isBasicExpr() {}

// NotParenExpr represents a negated parenthesized expression
type NotParenExpr struct {
	Expr LogicalOrExpr
}

func (NotParenExpr) isBasicExpr() {}

// ComparisonOperator represents a comparison operator
type ComparisonOperator int

const (
	EqualTo ComparisonOperator = iota
	NotEqualTo
	LessThan
	LessThanEqualTo
	GreaterThan
	GreaterThanEqualTo
)

// Comparable represents a comparable value in a comparison expression
type Comparable interface {
	isComparable()
}

// LiteralComparable represents a literal value as a comparable
type LiteralComparable struct {
	Value primitive.Literal
}

func (LiteralComparable) isComparable() {}

// SingularQueryComparable represents a singular query as a comparable
type SingularQueryComparable struct {
	Query string
}

func (SingularQueryComparable) isComparable() {}

// FunctionExprComparable represents a function expression as a comparable
type FunctionExprComparable struct {
	Expr FunctionExpr
}

func (FunctionExprComparable) isComparable() {}

// ParseFilter parses a filter selector from the input string
func ParseFilter(input string) (Filter, error) {
	input = strings.TrimSpace(input)
	if !strings.HasPrefix(input, "?") {
		return Filter{}, fmt.Errorf("invalid filter selector: %s", input)
	}
	input = strings.TrimPrefix(input, "?")
	input = strings.TrimSpace(input)

	expr, err := parseLogicalOrExpr(input)
	if err != nil {
		return Filter{}, fmt.Errorf("invalid filter selector: %v", err)
	}

	return Filter{Expression: expr}, nil
}

func parseLogicalOrExpr(input string) (LogicalOrExpr, error) {
	exprs := []LogicalAndExpr{}
	parts := strings.Split(input, "||")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		expr, err := parseLogicalAndExpr(part)
		if err != nil {
			return LogicalOrExpr{}, err
		}
		exprs = append(exprs, expr)
	}
	return LogicalOrExpr{Expressions: exprs}, nil
}

func parseLogicalAndExpr(input string) (LogicalAndExpr, error) {
	exprs := []BasicExpr{}
	parts := strings.Split(input, "&&")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		expr, err := parseBasicExpr(part)
		if err != nil {
			return LogicalAndExpr{}, err
		}
		exprs = append(exprs, expr)
	}
	return LogicalAndExpr{Expressions: exprs}, nil
}

func parseBasicExpr(input string) (BasicExpr, error) {
	input = strings.TrimSpace(input)
	switch {
	case strings.HasPrefix(input, "!("):
		expr, err := parseParenExpr(strings.TrimPrefix(input, "!("))
		if err != nil {
			return nil, err
		}
		return NotParenExpr{Expr: expr}, nil
	case strings.HasPrefix(input, "("):
		return parseParenExpr(input)
	case strings.HasPrefix(input, "!"):
		input = strings.TrimPrefix(input, "!")
		input = strings.TrimSpace(input)
		if expr, err := parseExistExpr(input); err == nil {
			return NotExistExpr{Query: expr.Query}, nil
		}
		if expr, err := parseFuncExpr(input); err == nil {
			return NotFuncExpr{Expr: expr}, nil
		}
	default:
		if expr, err := parseCompExpr(input); err == nil {
			return expr, nil
		}
		if expr, err := parseExistExpr(input); err == nil {
			return expr, nil
		}
		if expr, err := parseFuncExpr(input); err == nil {
			return expr, nil
		}
	}
	return nil, fmt.Errorf("invalid basic expression: %s", input)
}

func parseCompExpr(input string) (ComparisonExpr, error) {
	parts := strings.Split(input, " ")
	if len(parts) != 3 {
		return ComparisonExpr{}, fmt.Errorf("invalid comparison expression: %s", input)
	}
	left, err := parseComparable(parts[0])
	if err != nil {
		return ComparisonExpr{}, err
	}
	op, err := parseComparisonOperator(parts[1])
	if err != nil {
		return ComparisonExpr{}, err
	}
	right, err := parseComparable(parts[2])
	if err != nil {
		return ComparisonExpr{}, err
	}
	return ComparisonExpr{Left: left, Op: op, Right: right}, nil
}

func parseComparisonOperator(input string) (ComparisonOperator, error) {
	switch input {
	case "==":
		return EqualTo, nil
	case "!=":
		return NotEqualTo, nil
	case "<":
		return LessThan, nil
	case "<=":
		return LessThanEqualTo, nil
	case ">":
		return GreaterThan, nil
	case ">=":
		return GreaterThanEqualTo, nil
	default:
		return 0, fmt.Errorf("invalid comparison operator: %s", input)
	}
}

func parseExistExpr(input string) (ExistExpr, error) {
	return ExistExpr{Query: input}, nil
}

func parseFuncExpr(input string) (FuncExpr, error) {
	expr, err := ParseFunctionExpr(input)
	if err != nil {
		return FuncExpr{}, err
	}
	return FuncExpr{Expr: *expr}, nil
}

func parseParenExpr(input string) (ParenExpr, error) {
	if !strings.HasPrefix(input, "(") || !strings.HasSuffix(input, ")") {
		return ParenExpr{}, fmt.Errorf("invalid parenthesized expression: %s", input)
	}
	input = strings.TrimPrefix(input, "(")
	input = strings.TrimSuffix(input, ")")
	input = strings.TrimSpace(input)
	expr, err := parseLogicalOrExpr(input)
	if err != nil {
		return ParenExpr{}, err
	}
	return ParenExpr{Expr: expr}, nil
}

func parseLiteral(input string) (primitive.Literal, error) {
	return primitive.ParseLiteral(input)
}

func parseComparable(input string) (Comparable, error) {
	if literal, err := parseLiteral(input); err == nil {
		return LiteralComparable{Value: literal}, nil
	}
	if expr, err := ParseFunctionExpr(input); err == nil {
		return FunctionExprComparable{Expr: *expr}, nil
	}
	return SingularQueryComparable{Query: input}, nil
}
