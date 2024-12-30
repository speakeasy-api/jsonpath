package jsonpath

import (
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestLiteralEquals(t *testing.T) {
	testCases := []struct {
		name     string
		literal1 Literal
		literal2 Literal
		expected bool
	}{
		{
			name:     "Equal integers",
			literal1: Literal{Integer: intPtr(10)},
			literal2: Literal{Integer: intPtr(10)},
			expected: true,
		},
		{
			name:     "Different integers",
			literal1: Literal{Integer: intPtr(10)},
			literal2: Literal{Integer: intPtr(20)},
			expected: false,
		},
		{
			name:     "Equal floats",
			literal1: Literal{Float64: float64Ptr(3.14)},
			literal2: Literal{Float64: float64Ptr(3.14)},
			expected: true,
		},
		{
			name:     "Different floats",
			literal1: Literal{Float64: float64Ptr(3.14)},
			literal2: Literal{Float64: float64Ptr(2.71)},
			expected: false,
		},
		{
			name:     "Equal strings",
			literal1: Literal{String: stringPtr("hello")},
			literal2: Literal{String: stringPtr("hello")},
			expected: true,
		},
		{
			name:     "Different strings",
			literal1: Literal{String: stringPtr("hello")},
			literal2: Literal{String: stringPtr("world")},
			expected: false,
		},
		{
			name:     "Equal bools",
			literal1: Literal{Bool: boolPtr(true)},
			literal2: Literal{Bool: boolPtr(true)},
			expected: true,
		},
		{
			name:     "Different bools",
			literal1: Literal{Bool: boolPtr(true)},
			literal2: Literal{Bool: boolPtr(false)},
			expected: false,
		},
		{
			name:     "Equal nulls",
			literal1: Literal{Null: boolPtr(true)},
			literal2: Literal{Null: boolPtr(true)},
			expected: true,
		},
		{
			name:     "Different nulls",
			literal1: Literal{Null: boolPtr(true)},
			literal2: Literal{Null: boolPtr(false)},
			expected: false,
		},
		{
			name:     "Different types",
			literal1: Literal{Integer: intPtr(10)},
			literal2: Literal{String: stringPtr("10")},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.literal1.Equals(tc.literal2)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestLiteralLessThan(t *testing.T) {
	testCases := []struct {
		name     string
		literal1 Literal
		literal2 Literal
		expected bool
	}{
		{
			name:     "Integer less than",
			literal1: Literal{Integer: intPtr(10)},
			literal2: Literal{Integer: intPtr(20)},
			expected: true,
		},
		{
			name:     "Integer not less than",
			literal1: Literal{Integer: intPtr(20)},
			literal2: Literal{Integer: intPtr(10)},
			expected: false,
		},
		{
			name:     "Float less than",
			literal1: Literal{Float64: float64Ptr(3.14)},
			literal2: Literal{Float64: float64Ptr(6.28)},
			expected: true,
		},
		{
			name:     "Float not less than",
			literal1: Literal{Float64: float64Ptr(6.28)},
			literal2: Literal{Float64: float64Ptr(3.14)},
			expected: false,
		},
		{
			name:     "Different types",
			literal1: Literal{Integer: intPtr(10)},
			literal2: Literal{String: stringPtr("10")},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.literal1.LessThan(tc.literal2)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestLiteralLessThanOrEqual(t *testing.T) {
	testCases := []struct {
		name     string
		literal1 Literal
		literal2 Literal
		expected bool
	}{
		{
			name:     "Integer less than or equal",
			literal1: Literal{Integer: intPtr(10)},
			literal2: Literal{Integer: intPtr(20)},
			expected: true,
		},
		{
			name:     "Integer equal",
			literal1: Literal{Integer: intPtr(10)},
			literal2: Literal{Integer: intPtr(10)},
			expected: true,
		},
		{
			name:     "Integer not less than or equal",
			literal1: Literal{Integer: intPtr(20)},
			literal2: Literal{Integer: intPtr(10)},
			expected: false,
		},
		{
			name:     "Float less than or equal",
			literal1: Literal{Float64: float64Ptr(3.14)},
			literal2: Literal{Float64: float64Ptr(6.28)},
			expected: true,
		},
		{
			name:     "Float equal",
			literal1: Literal{Float64: float64Ptr(3.14)},
			literal2: Literal{Float64: float64Ptr(3.14)},
			expected: true,
		},
		{
			name:     "Float not less than or equal",
			literal1: Literal{Float64: float64Ptr(6.28)},
			literal2: Literal{Float64: float64Ptr(3.14)},
			expected: false,
		},
		{
			name:     "Different types",
			literal1: Literal{Integer: intPtr(10)},
			literal2: Literal{String: stringPtr("10")},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.literal1.LessThanOrEqual(tc.literal2)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestComparableEvaluate(t *testing.T) {
	testCases := []struct {
		name       string
		comparable Comparable
		node       *yaml.Node
		root       *yaml.Node
		expected   Literal
	}{
		{
			name:       "Literal",
			comparable: Comparable{Literal: &Literal{Integer: intPtr(10)}},
			node:       yamlNodeFromString("foo"),
			root:       yamlNodeFromString("foo"),
			expected:   Literal{Integer: intPtr(10)},
		},
		{
			name:       "SingularQuery",
			comparable: Comparable{SingularQuery: &SingularQuery{AbsQuery: &AbsQuery{Segments: []*Segment{}}}},
			node:       yamlNodeFromString("10"),
			root:       yamlNodeFromString("10"),
			expected:   Literal{Integer: intPtr(10)},
		},
		{
			name:       "FunctionExpr",
			comparable: Comparable{FunctionExpr: &FunctionExpr{Type: FunctionTypeLength, Args: []*FunctionArgument{{FilterQuery: &FilterQuery{RelQuery: &RelQuery{Segments: []*Segment{}}}}}}},
			node:       yamlNodeFromString(`["a", "b", "c"]`),
			root:       yamlNodeFromString(`["a", "b", "c"]`),
			expected:   Literal{Integer: intPtr(3)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.comparable.Evaluate(tc.node, tc.root)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestFunctionExprEvaluate(t *testing.T) {
	testCases := []struct {
		name     string
		funcExpr FunctionExpr
		node     *yaml.Node
		root     *yaml.Node
		expected Literal
	}{
		{
			name:     "Length of scalar",
			funcExpr: FunctionExpr{Type: FunctionTypeLength},
			node:     yamlNodeFromString("hello"),
			root:     yamlNodeFromString("hello"),
			expected: Literal{Integer: intPtr(5)},
		},
		{
			name:     "Length of sequence",
			funcExpr: FunctionExpr{Type: FunctionTypeLength},
			node:     yamlNodeFromString(`["a", "b", "c"]`),
			root:     yamlNodeFromString(`["a", "b", "c"]`),
			expected: Literal{Integer: intPtr(3)},
		},
		{
			name:     "Length of mapping",
			funcExpr: FunctionExpr{Type: FunctionTypeLength},
			node:     yamlNodeFromString(`{"a": 1, "b": 2}`),
			root:     yamlNodeFromString(`{"a": 1, "b": 2}`),
			expected: Literal{Integer: intPtr(2)},
		},
		{
			name:     "Count of nodes",
			funcExpr: FunctionExpr{Type: FunctionTypeCount, Args: []*FunctionArgument{{FilterQuery: &FilterQuery{RelQuery: &RelQuery{Segments: []*Segment{}}}}}},
			node:     yamlNodeFromString(`["a", "b", "c"]`),
			root:     yamlNodeFromString(`["a", "b", "c"]`),
			expected: Literal{Integer: intPtr(1)}, // Count of a node list is 1 (unintuitive I know)
		},
		{
			name:     "Count of node wildcard",
			funcExpr: FunctionExpr{Type: FunctionTypeCount, Args: []*FunctionArgument{{FilterQuery: &FilterQuery{RelQuery: &RelQuery{Segments: []*Segment{{Child: &ChildSegment{kind: ChildSegmentDotWildcard}}}}}}}},
			node:     yamlNodeFromString(`["a", "b", "c"]`),
			root:     yamlNodeFromString(`["a", "b", "c"]`),
			expected: Literal{Integer: intPtr(3)}, // Count of a node list is 1 (unintuitive I know)
		},
		// Add more test cases for match, search, and value functions
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.funcExpr.Evaluate(tc.node, tc.root)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestSingularQueryEvaluate(t *testing.T) {
	testCases := []struct {
		name     string
		query    SingularQuery
		node     *yaml.Node
		root     *yaml.Node
		expected Literal
	}{
		{
			name:     "RelQuery",
			query:    SingularQuery{RelQuery: &RelQuery{Segments: []*Segment{}}},
			node:     yamlNodeFromString("10"),
			root:     yamlNodeFromString("10"),
			expected: Literal{Integer: intPtr(10)},
		},
		{
			name:     "AbsQuery",
			query:    SingularQuery{AbsQuery: &AbsQuery{Segments: []*Segment{}}},
			node:     yamlNodeFromString("10"),
			root:     yamlNodeFromString("10"),
			expected: Literal{Integer: intPtr(10)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.query.Evaluate(tc.node, tc.root)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestRelQueryEvaluate(t *testing.T) {
	testCases := []struct {
		name     string
		query    RelQuery
		node     *yaml.Node
		root     *yaml.Node
		expected Literal
	}{
		{
			name:     "Single node",
			query:    RelQuery{Segments: []*Segment{}},
			node:     yamlNodeFromString("10"),
			root:     yamlNodeFromString("10"),
			expected: Literal{Integer: intPtr(10)},
		},
		{
			name:     "Child segment",
			query:    RelQuery{Segments: []*Segment{{Child: &ChildSegment{kind: ChildSegmentDotMemberName, dotName: "foo"}}}},
			node:     yamlNodeFromString(`{"foo": "bar"}`),
			root:     yamlNodeFromString(`{"foo": "bar"}`),
			expected: Literal{String: stringPtr("bar")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.query.Evaluate(tc.node, tc.root)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestAbsQueryEvaluate(t *testing.T) {
	testCases := []struct {
		name     string
		query    AbsQuery
		node     *yaml.Node
		root     *yaml.Node
		expected Literal
	}{
		{
			name:     "Root node",
			query:    AbsQuery{Segments: []*Segment{}},
			node:     yamlNodeFromString("10"),
			root:     yamlNodeFromString("10"),
			expected: Literal{Integer: intPtr(10)},
		},
		{
			name:     "Child segment",
			query:    AbsQuery{Segments: []*Segment{{Child: &ChildSegment{kind: ChildSegmentDotMemberName, dotName: "foo"}}}},
			node:     yamlNodeFromString(`{"foo": "bar"}`),
			root:     yamlNodeFromString(`{"foo": "bar"}`),
			expected: Literal{String: stringPtr("bar")},
		},
		// Add more test cases for other segment types
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.query.Evaluate(tc.node, tc.root)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}

func yamlNodeFromString(s string) *yaml.Node {
	var node yaml.Node
	err := yaml.Unmarshal([]byte(s), &node)
	if err != nil {
		panic(err)
	}
	if len(node.Content) != 1 {
		panic("expected single node")
	}
	return node.Content[0]
}
