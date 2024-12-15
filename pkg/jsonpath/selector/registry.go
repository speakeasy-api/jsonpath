package selector

import (
	"fmt"
	"strings"
)

// ParseFunctionExpr parses a function expression from the input string
func ParseFunctionExpr(input string) (*FunctionExpr, error) {
	input = strings.TrimSpace(input)

	name, args, err := parseFunctionNameAndArgs(input)
	if err != nil {
		return nil, err
	}

	// Retrieve the function from the registry
	function, ok := defaultRegistry.Get(name)
	if !ok {
		return nil, fmt.Errorf("undefined function: %s", name)
	}

	// Validate the function arguments
	if err := function.Validator(args); err != nil {
		return nil, fmt.Errorf("invalid function arguments: %v", err)
	}

	return &FunctionExpr{
		name:       name,
		args:       args,
		returnType: function.ResultType,
		evaluator:  function.Evaluator,
	}, nil
}

// parseFunctionNameAndArgs parses the function name and arguments from the input string
func parseFunctionNameAndArgs(input string) (string, []FunctionExprArg, error) {
	// Find the function name and arguments
	openParenIndex := strings.IndexByte(input, '(')
	if openParenIndex == -1 {
		return "", nil, fmt.Errorf("invalid function expression: missing opening parenthesis")
	}

	closeParenIndex := strings.LastIndexByte(input, ')')
	if closeParenIndex == -1 {
		return "", nil, fmt.Errorf("invalid function expression: missing closing parenthesis")
	}

	functionName := strings.TrimSpace(input[:openParenIndex])
	argumentsStr := strings.TrimSpace(input[openParenIndex+1 : closeParenIndex])

	// Parse the function arguments
	var arguments []FunctionExprArg
	if argumentsStr != "" {
		argumentStrs := strings.Split(argumentsStr, ",")
		arguments = make([]FunctionExprArg, len(argumentStrs))

		for i, argStr := range argumentStrs {
			argStr = strings.TrimSpace(argStr)
			arg, err := parseFunctionArgument(argStr)
			if err != nil {
				return "", nil, fmt.Errorf("invalid function argument: %v", err)
			}
			arguments[i] = arg
		}
	}

	return functionName, arguments, nil
}

// parseFunctionArgument parses a single function argument from the input string
func parseFunctionArgument(input string) (FunctionExprArg, error) {
	// Try parsing the argument as different types
	if arg, err := filter.ParseLiteral(input); err == nil {
		return &LiteralArgument{Value: arg}, nil
	}

	if arg, err := filter.ParseSingularPath(input); err == nil {
		return &SingularQueryArgument{Query: arg}, nil
	}

	if arg, err := ParseQuery(input); err == nil {
		return &FilterQueryArgument{Query: arg}, nil
	}

	if arg, err := ParseFunctionExpr(input); err == nil {
		return &FunctionExprArgument{Expr: arg}, nil
	}

	if arg, err := filter.ParseLogicalOrExpr(input); err == nil {
		return &LogicalExprArgument{Expr: arg}, nil
	}

	return nil, fmt.Errorf("invalid function argument: %s", input)
}

// LiteralArgument represents a literal value argument
type LiteralArgument struct {
	Value interface{}
}

// TypeKind returns the argument type for LiteralArgument
func (a *LiteralArgument) TypeKind() FunctionArgType {
	return FunctionArgTypeLiteral
}

// Evaluate evaluates the literal argument
func (a *LiteralArgument) Evaluate(current, root interface{}) (interface{}, error) {
	return a.Value, nil
}

// SingularQueryArgument represents a singular query argument
type SingularQueryArgument struct {
	Query filter.SingularPath
}

// TypeKind returns the argument type for SingularQueryArgument
func (a *SingularQueryArgument) TypeKind() FunctionArgType {
	return FunctionArgTypeSingularQuery
}

// Evaluate evaluates the singular query argument
func (a *SingularQueryArgument) Evaluate(current, root interface{}) (interface{}, error) {
	return a.Query.Evaluate(current, root)
}

// FilterQueryArgument represents a filter query argument
type FilterQueryArgument struct {
	Query filter.Query
}

// TypeKind returns the argument type for FilterQueryArgument
func (a *FilterQueryArgument) TypeKind() FunctionArgType {
	return FunctionArgTypeNodelist
}

// Evaluate evaluates the filter query argument
func (a *FilterQueryArgument) Evaluate(current, root interface{}) (interface{}, error) {
	return a.Query.Evaluate(current, root)
}

// FunctionExprArgument represents a function expression argument
type FunctionExprArgument struct {
	Expr *FunctionExpr
}

// TypeKind returns the argument type for FunctionExprArgument
func (a *FunctionExprArgument) TypeKind() FunctionArgType {
	return FunctionArgTypeValue
}

// Evaluate evaluates the function expression argument
func (a *FunctionExprArgument) Evaluate(current, root interface{}) (interface{}, error) {
	return a.Expr.Evaluate(current, root)
}

// LogicalExprArgument represents a logical expression argument
type LogicalExprArgument struct {
	Expr filter.LogicalOrExpr
}

// TypeKind returns the argument type for LogicalExprArgument
func (a *LogicalExprArgument) TypeKind() FunctionArgType {
	return FunctionArgTypeLogical
}

// Evaluate evaluates the logical expression argument
func (a *LogicalExprArgument) Evaluate(current, root interface{}) (interface{}, error) {
	return a.Expr.Evaluate(current, root)
}
