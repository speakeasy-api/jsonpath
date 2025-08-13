package jsonpath

import (
	"reflect"
	"testing"

	"go.yaml.in/yaml/v4"
)

func TestLiteralEquals(t *testing.T) {
	testCases := []struct {
		name     string
		literal1 literal
		literal2 literal
		expected bool
	}{
		{
			name:     "Equal integers",
			literal1: literal{integer: intPtr(10)},
			literal2: literal{integer: intPtr(10)},
			expected: true,
		},
		{
			name:     "Different integers",
			literal1: literal{integer: intPtr(10)},
			literal2: literal{integer: intPtr(20)},
			expected: false,
		},
		{
			name:     "Equal floats",
			literal1: literal{float64: float64Ptr(3.14)},
			literal2: literal{float64: float64Ptr(3.14)},
			expected: true,
		},
		{
			name:     "Different floats",
			literal1: literal{float64: float64Ptr(3.14)},
			literal2: literal{float64: float64Ptr(2.71)},
			expected: false,
		},
		{
			name:     "Equal strings",
			literal1: literal{string: stringPtr("hello")},
			literal2: literal{string: stringPtr("hello")},
			expected: true,
		},
		{
			name:     "Different strings",
			literal1: literal{string: stringPtr("hello")},
			literal2: literal{string: stringPtr("world")},
			expected: false,
		},
		{
			name:     "Equal bools",
			literal1: literal{bool: boolPtr(true)},
			literal2: literal{bool: boolPtr(true)},
			expected: true,
		},
		{
			name:     "Different bools",
			literal1: literal{bool: boolPtr(true)},
			literal2: literal{bool: boolPtr(false)},
			expected: false,
		},
		{
			name:     "Equal nulls",
			literal1: literal{null: boolPtr(true)},
			literal2: literal{null: boolPtr(true)},
			expected: true,
		},
		{
			name:     "Different nulls",
			literal1: literal{null: boolPtr(true)},
			literal2: literal{null: boolPtr(false)},
			expected: false,
		},
		{
			name:     "Different types",
			literal1: literal{integer: intPtr(10)},
			literal2: literal{string: stringPtr("10")},
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
		literal1 literal
		literal2 literal
		expected bool
	}{
		{
			name:     "integer less than",
			literal1: literal{integer: intPtr(10)},
			literal2: literal{integer: intPtr(20)},
			expected: true,
		},
		{
			name:     "integer not less than",
			literal1: literal{integer: intPtr(20)},
			literal2: literal{integer: intPtr(10)},
			expected: false,
		},
		{
			name:     "Float less than",
			literal1: literal{float64: float64Ptr(3.14)},
			literal2: literal{float64: float64Ptr(6.28)},
			expected: true,
		},
		{
			name:     "Float not less than",
			literal1: literal{float64: float64Ptr(6.28)},
			literal2: literal{float64: float64Ptr(3.14)},
			expected: false,
		},
		{
			name:     "Different types",
			literal1: literal{integer: intPtr(10)},
			literal2: literal{string: stringPtr("10")},
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
		literal1 literal
		literal2 literal
		expected bool
	}{
		{
			name:     "integer less than or equal",
			literal1: literal{integer: intPtr(10)},
			literal2: literal{integer: intPtr(20)},
			expected: true,
		},
		{
			name:     "integer equal",
			literal1: literal{integer: intPtr(10)},
			literal2: literal{integer: intPtr(10)},
			expected: true,
		},
		{
			name:     "integer not less than or equal",
			literal1: literal{integer: intPtr(20)},
			literal2: literal{integer: intPtr(10)},
			expected: false,
		},
		{
			name:     "Float less than or equal",
			literal1: literal{float64: float64Ptr(3.14)},
			literal2: literal{float64: float64Ptr(6.28)},
			expected: true,
		},
		{
			name:     "Float equal",
			literal1: literal{float64: float64Ptr(3.14)},
			literal2: literal{float64: float64Ptr(3.14)},
			expected: true,
		},
		{
			name:     "Float not less than or equal",
			literal1: literal{float64: float64Ptr(6.28)},
			literal2: literal{float64: float64Ptr(3.14)},
			expected: false,
		},
		{
			name:     "Different types",
			literal1: literal{integer: intPtr(10)},
			literal2: literal{string: stringPtr("10")},
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
		comparable comparable
		node       *yaml.Node
		root       *yaml.Node
		expected   literal
	}{
		{
			name:       "literal",
			comparable: comparable{literal: &literal{integer: intPtr(10)}},
			node:       yamlNodeFromString("foo"),
			root:       yamlNodeFromString("foo"),
			expected:   literal{integer: intPtr(10)},
		},
		{
			name:       "singularQuery",
			comparable: comparable{singularQuery: &singularQuery{absQuery: &absQuery{segments: []*segment{}}}},
			node:       yamlNodeFromString("10"),
			root:       yamlNodeFromString("10"),
			expected:   literal{integer: intPtr(10)},
		},
		{
			name:       "functionExpr",
			comparable: comparable{functionExpr: &functionExpr{funcType: functionTypeLength, args: []*functionArgument{{filterQuery: &filterQuery{relQuery: &relQuery{segments: []*segment{}}}}}}},
			node:       yamlNodeFromString(`["a", "b", "c"]`),
			root:       yamlNodeFromString(`["a", "b", "c"]`),
			expected:   literal{integer: intPtr(3)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.comparable.Evaluate(&_index{}, tc.node, tc.root)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestSingularQueryEvaluate(t *testing.T) {
	testCases := []struct {
		name     string
		query    singularQuery
		node     *yaml.Node
		root     *yaml.Node
		expected literal
	}{
		{
			name:     "relQuery",
			query:    singularQuery{relQuery: &relQuery{segments: []*segment{}}},
			node:     yamlNodeFromString("10"),
			root:     yamlNodeFromString("10"),
			expected: literal{integer: intPtr(10)},
		},
		{
			name:     "absQuery",
			query:    singularQuery{absQuery: &absQuery{segments: []*segment{}}},
			node:     yamlNodeFromString("10"),
			root:     yamlNodeFromString("10"),
			expected: literal{integer: intPtr(10)},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.query.Evaluate(&_index{}, tc.node, tc.root)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestRelQueryEvaluate(t *testing.T) {
	testCases := []struct {
		name     string
		query    relQuery
		node     *yaml.Node
		root     *yaml.Node
		expected literal
	}{
		{
			name:     "Single node",
			query:    relQuery{segments: []*segment{}},
			node:     yamlNodeFromString("10"),
			root:     yamlNodeFromString("10"),
			expected: literal{integer: intPtr(10)},
		},
		{
			name:     "child segment",
			query:    relQuery{segments: []*segment{{kind: segmentKindChild, child: &innerSegment{kind: segmentDotMemberName, dotName: "foo"}}}},
			node:     yamlNodeFromString(`{"foo": "bar"}`),
			root:     yamlNodeFromString(`{"foo": "bar"}`),
			expected: literal{string: stringPtr("bar")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.query.Evaluate(&_index{}, tc.node, tc.root)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestAbsQueryEvaluate(t *testing.T) {
	testCases := []struct {
		name     string
		query    absQuery
		node     *yaml.Node
		root     *yaml.Node
		expected literal
	}{
		{
			name:     "Root node",
			query:    absQuery{segments: []*segment{}},
			node:     yamlNodeFromString("10"),
			root:     yamlNodeFromString("10"),
			expected: literal{integer: intPtr(10)},
		},
		{
			name:     "child segment",
			query:    absQuery{segments: []*segment{{kind: segmentKindChild, child: &innerSegment{kind: segmentDotMemberName, dotName: "foo"}}}},
			node:     yamlNodeFromString(`{"foo": "bar"}`),
			root:     yamlNodeFromString(`{"foo": "bar"}`),
			expected: literal{string: stringPtr("bar")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.query.Evaluate(&_index{}, tc.node, tc.root)
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
