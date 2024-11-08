package jsonpath

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"regexp"
	"sync"
)

// JsonPathType represents the basic types in JSONPath
type JsonPathType int

const (
	JsonPathTypeNodes JsonPathType = iota
	JsonPathTypeValue
	JsonPathTypeLogical
)

func (t JsonPathType) String() string {
	switch t {
	case JsonPathTypeNodes:
		return "NodesType"
	case JsonPathTypeLogical:
		return "LogicalType"
	case JsonPathTypeValue:
		return "ValueType"
	default:
		return "Unknown"
	}
}

// LogicalType represents true/false values in JSONPath
type LogicalType bool

const (
	LogicalTrue  LogicalType = true
	LogicalFalse LogicalType = false
)

// ValueType represents a JSON value or Nothing
type ValueType struct {
	value *yaml.Node
	isNil bool
}

func NewValueTypeFromNode(node *yaml.Node) ValueType {
	if node == nil {
		return ValueType{isNil: true}
	}
	return ValueType{value: node}
}

func (v ValueType) IsNothing() bool {
	return v.isNil
}

// FunctionArgType represents different types of function arguments
type FunctionArgType int

const (
	FunctionArgTypeLiteral FunctionArgType = iota
	FunctionArgTypeSingularQuery
	FunctionArgTypeValue
	FunctionArgTypeNodelist
	FunctionArgTypeLogical
)

func (t FunctionArgType) String() string {
	switch t {
	case FunctionArgTypeLiteral:
		return "literal"
	case FunctionArgTypeSingularQuery:
		return "singular query"
	case FunctionArgTypeValue:
		return "value type"
	case FunctionArgTypeNodelist:
		return "nodes type"
	case FunctionArgTypeLogical:
		return "logical type"
	default:
		return "unknown"
	}
}

// FunctionValidationError represents errors that can occur during function validation
type FunctionValidationError struct {
	msg string
}

func (e FunctionValidationError) Error() string {
	return e.msg
}

// Function represents a JSONPath function
type Function struct {
	Name       string
	ResultType FunctionArgType
	Validator  func([]FunctionExprArg) error
	Evaluator  func([]*yaml.Node) ([]*yaml.Node, error)
}

// FunctionExpr represents a function expression
type FunctionExpr struct {
	name       string
	args       []FunctionExprArg
	returnType FunctionArgType
	evaluator  func([]*yaml.Node) ([]*yaml.Node, error)
}

// FunctionExprArg represents different types of function arguments
type FunctionExprArg interface {
	TypeKind() FunctionArgType
	Evaluate(current, root *yaml.Node) ([]*yaml.Node, error)
}

// Registry manages function registration
type Registry struct {
	functions map[string]*Function
	mu        sync.RWMutex
}

var defaultRegistry = &Registry{
	functions: make(map[string]*Function),
}

func (r *Registry) Register(f *Function) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.functions[f.Name] = f
}

func (r *Registry) Get(name string) (*Function, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.functions[name]
	return f, ok
}

// ConvertibleTo checks if one type can convert to another
func (t FunctionArgType) ConvertibleTo(target JsonPathType) bool {
	switch {
	case t == FunctionArgTypeLiteral && target == JsonPathTypeValue:
		return true
	case t == FunctionArgTypeValue && target == JsonPathTypeValue:
		return true
	case t == FunctionArgTypeSingularQuery && (target == JsonPathTypeValue || target == JsonPathTypeNodes || target == JsonPathTypeLogical):
		return true
	case t == FunctionArgTypeNodelist && (target == JsonPathTypeNodes || target == JsonPathTypeLogical):
		return true
	case t == FunctionArgTypeLogical && target == JsonPathTypeLogical:
		return true
	default:
		return false
	}
}

func lengthFunction() *Function {
	return &Function{
		Name:       "length",
		ResultType: FunctionArgTypeValue,
		Validator: func(args []FunctionExprArg) error {
			if len(args) != 1 {
				return &FunctionValidationError{msg: "length function requires exactly 1 argument"}
			}
			argType := args[0].TypeKind()
			if !argType.ConvertibleTo(JsonPathTypeValue) {
				return &FunctionValidationError{msg: "length function argument must be convertible to value type"}
			}
			return nil
		},
		Evaluator: func(args []*yaml.Node) ([]*yaml.Node, error) {
			if len(args) != 1 || args[0] == nil {
				return nil, nil
			}

			node := args[0]
			var length int

			switch node.Kind {
			case yaml.SequenceNode:
				length = len(node.Content)
			case yaml.MappingNode:
				length = len(node.Content) / 2 // Mapping nodes have key-value pairs
			case yaml.ScalarNode:
				length = len(node.Value)
			default:
				return nil, nil
			}

			result := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!int",
				Value: fmt.Sprintf("%d", length),
			}
			return []*yaml.Node{result}, nil
		},
	}
}

// countFunction implements the count() function
func countFunction() *Function {
	return &Function{
		Name:       "count",
		ResultType: FunctionArgTypeValue,
		Validator: func(args []FunctionExprArg) error {
			if len(args) != 1 {
				return &FunctionValidationError{msg: "count function requires exactly 1 argument"}
			}
			argType := args[0].TypeKind()
			if !argType.ConvertibleTo(JsonPathTypeNodes) {
				return &FunctionValidationError{msg: "count function argument must be convertible to nodes type"}
			}
			return nil
		},
		Evaluator: func(args []*yaml.Node) ([]*yaml.Node, error) {
			if len(args) != 1 {
				return nil, nil
			}

			// Handle different node kinds
			var count int
			switch args[0].Kind {
			case yaml.SequenceNode, yaml.MappingNode:
				count = len(args[0].Content)
			case yaml.ScalarNode:
				count = 1 // A scalar node counts as 1
			default:
				count = 0
			}

			result := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!int",
				Value: fmt.Sprintf("%d", count),
			}
			return []*yaml.Node{result}, nil
		},
	}
}

// matchFunction implements the match() function
func matchFunction() *Function {
	return &Function{
		Name:       "match",
		ResultType: FunctionArgTypeLogical,
		Validator: func(args []FunctionExprArg) error {
			if len(args) != 2 {
				return &FunctionValidationError{msg: "match function requires exactly 2 arguments"}
			}
			for _, arg := range args {
				if !arg.TypeKind().ConvertibleTo(JsonPathTypeValue) {
					return &FunctionValidationError{msg: "match function arguments must be convertible to value type"}
				}
			}
			return nil
		},
		Evaluator: func(args []*yaml.Node) ([]*yaml.Node, error) {
			if len(args) != 2 || args[0] == nil || args[1] == nil {
				return nil, nil
			}

			if args[0].Kind != yaml.ScalarNode || args[1].Kind != yaml.ScalarNode {
				return createLogicalResult(false), nil
			}

			pattern := args[1].Value
			text := args[0].Value

			matched, err := regexp.MatchString("^"+pattern+"$", text)
			if err != nil {
				return nil, err
			}

			return createLogicalResult(matched), nil
		},
	}
}

// searchFunction implements the search() function
func searchFunction() *Function {
	return &Function{
		Name:       "search",
		ResultType: FunctionArgTypeLogical,
		Validator: func(args []FunctionExprArg) error {
			if len(args) != 2 {
				return &FunctionValidationError{msg: "search function requires exactly 2 arguments"}
			}
			for _, arg := range args {
				if !arg.TypeKind().ConvertibleTo(JsonPathTypeValue) {
					return &FunctionValidationError{msg: "search function arguments must be convertible to value type"}
				}
			}
			return nil
		},
		Evaluator: func(args []*yaml.Node) ([]*yaml.Node, error) {
			if len(args) != 2 || args[0] == nil || args[1] == nil {
				return nil, nil
			}

			if args[0].Kind != yaml.ScalarNode || args[1].Kind != yaml.ScalarNode {
				return createLogicalResult(false), nil
			}

			pattern := args[1].Value
			text := args[0].Value

			matched, err := regexp.MatchString(pattern, text)
			if err != nil {
				return nil, err
			}

			return createLogicalResult(matched), nil
		},
	}
}

// valueFunction implements the value() function
func valueFunction() *Function {
	return &Function{
		Name:       "value",
		ResultType: FunctionArgTypeValue,
		Validator: func(args []FunctionExprArg) error {
			if len(args) != 1 {
				return &FunctionValidationError{msg: "value function requires exactly 1 argument"}
			}
			argType := args[0].TypeKind()
			if !argType.ConvertibleTo(JsonPathTypeNodes) {
				return &FunctionValidationError{msg: "value function argument must be convertible to nodes type"}
			}
			return nil
		},
		Evaluator: func(args []*yaml.Node) ([]*yaml.Node, error) {
			if len(args) != 1 || args[0] == nil {
				return nil, nil
			}

			// Return the single node if there's exactly one
			if len(args) == 1 {
				return []*yaml.Node{args[0]}, nil
			}

			// Return nothing if there's not exactly one node
			return nil, nil
		},
	}
}

// Helper function to create logical result nodes
func createLogicalResult(value bool) []*yaml.Node {
	result := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   "!!bool",
		Value: fmt.Sprintf("%t", value),
	}
	return []*yaml.Node{result}
}

// Initialize standard functions
func init() {
	defaultRegistry.Register(lengthFunction())
	defaultRegistry.Register(countFunction())
	defaultRegistry.Register(matchFunction())
	defaultRegistry.Register(searchFunction())
	defaultRegistry.Register(valueFunction())
}
