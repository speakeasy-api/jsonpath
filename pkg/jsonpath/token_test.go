package jsonpath

import (
	"fmt"
	"testing"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenInfo
	}{
		{
			name:  "Root",
			input: "$",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
			},
		},
		{
			name:  "Child",
			input: "$.store.book",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 2, Literal: "store", Len: 5},
				{Token: CHILD, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 8, Literal: "book", Len: 4},
			},
		},
		{
			name:  "Wildcard",
			input: "$.*",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: WILDCARD, Line: 1, Column: 2, Literal: "", Len: 1},
			},
		},
		{
			name:  "Recursive",
			input: "$..author",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: RECURSIVE, Line: 1, Column: 1, Literal: "", Len: 2},
				{Token: STRING_LITERAL, Line: 1, Column: 3, Literal: "author", Len: 6},
			},
		},
		{
			name:  "Union",
			input: "$..book[0,1]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: RECURSIVE, Line: 1, Column: 1, Literal: "", Len: 2},
				{Token: STRING_LITERAL, Line: 1, Column: 3, Literal: "book", Len: 4},
				{Token: BRACKET_LEFT, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: NUMBER, Line: 1, Column: 8, Literal: "0", Len: 1},
				{Token: UNION, Line: 1, Column: 9, Literal: "", Len: 1},
				{Token: NUMBER, Line: 1, Column: 10, Literal: "1", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 11, Literal: "", Len: 1},
			},
		},
		{
			name:  "Slice",
			input: "$..book[0:2]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: RECURSIVE, Line: 1, Column: 1, Literal: "", Len: 2},
				{Token: STRING_LITERAL, Line: 1, Column: 3, Literal: "book", Len: 4},
				{Token: BRACKET_LEFT, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: NUMBER, Line: 1, Column: 8, Literal: "0", Len: 1},
				{Token: ARRAY_SLICE, Line: 1, Column: 9, Literal: "", Len: 1},
				{Token: NUMBER, Line: 1, Column: 10, Literal: "2", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 11, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter",
			input: "$.store.book[?(@.price < 10)]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 2, Literal: "store", Len: 5},
				{Token: CHILD, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 8, Literal: "book", Len: 4},
				{Token: BRACKET_LEFT, Line: 1, Column: 12, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 13, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 14, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 15, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 16, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 17, Literal: "price", Len: 5},
				{Token: LT, Line: 1, Column: 23, Literal: "", Len: 1},
				{Token: NUMBER, Line: 1, Column: 25, Literal: "10", Len: 2},
				{Token: PAREN_RIGHT, Line: 1, Column: 27, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 28, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter equality",
			input: "$[?(@.child=='x')]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: BRACKET_LEFT, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 2, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 3, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 4, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 5, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: EQ, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: STRING, Line: 1, Column: 13, Literal: "x", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 16, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 17, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter inequality",
			input: "$[?(@.child!='x')]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: BRACKET_LEFT, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 2, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 3, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 4, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 5, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: NE, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: STRING, Line: 1, Column: 13, Literal: "x", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 16, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 17, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter greater than",
			input: "$[?(@.child>1)]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: BRACKET_LEFT, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 2, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 3, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 4, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 5, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: GT, Line: 1, Column: 11, Literal: "", Len: 1},
				{Token: NUMBER, Line: 1, Column: 12, Literal: "1", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 13, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 14, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter greater than or equal",
			input: "$[?(@.child>=1)]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: BRACKET_LEFT, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 2, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 3, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 4, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 5, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: GE, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: NUMBER, Line: 1, Column: 13, Literal: "1", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 14, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 15, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter less than",
			input: "$[?(@.child<1)]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: BRACKET_LEFT, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 2, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 3, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 4, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 5, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: LT, Line: 1, Column: 11, Literal: "", Len: 1},
				{Token: NUMBER, Line: 1, Column: 12, Literal: "1", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 13, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 14, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter less than or equal",
			input: "$[?(@.child<=1)]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: BRACKET_LEFT, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 2, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 3, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 4, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 5, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: LE, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: NUMBER, Line: 1, Column: 13, Literal: "1", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 14, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 15, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter logical AND",
			input: "$[?(@.child&&@.other)]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: BRACKET_LEFT, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 2, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 3, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 4, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 5, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: AND, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: CURRENT, Line: 1, Column: 13, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 14, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 15, Literal: "other", Len: 5},
				{Token: PAREN_RIGHT, Line: 1, Column: 20, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 21, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter logical OR",
			input: "$[?(@.child||@.other)]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: BRACKET_LEFT, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 2, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 3, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 4, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 5, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: OR, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: CURRENT, Line: 1, Column: 13, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 14, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 15, Literal: "other", Len: 5},
				{Token: PAREN_RIGHT, Line: 1, Column: 20, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 21, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter logical NOT",
			input: "$[?(!@.child)]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: BRACKET_LEFT, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 2, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 3, Literal: "", Len: 1},
				{Token: NOT, Line: 1, Column: 4, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 5, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 6, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 7, Literal: "child", Len: 5},
				{Token: PAREN_RIGHT, Line: 1, Column: 12, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 13, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter regular expression (illegal right now)",
			input: "$[?(@.child=~/.*/)]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 1, Literal: "", Len: 2},
				{Token: PAREN_LEFT, Line: 1, Column: 3, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 4, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 5, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: MATCHES, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: ILLEGAL, Line: 1, Column: 13, Literal: "", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 17, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 18, Literal: "", Len: 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenizer := NewTokenizer(test.input)
			tokens := tokenizer.Tokenize()

			if len(tokens) != len(test.expected) {
				msg := tokenizer.ErrorTokenString(tokens[0], fmt.Sprintf("Expected %d tokens, got %d", len(test.expected), len(tokens)))
				t.Error(msg)
			}

			for i, expectedToken := range test.expected {
				actualToken := tokens[i]
				if actualToken != expectedToken {
					msg := tokenizer.ErrorString(actualToken, fmt.Sprintf("Expected token %+v, got %+v", expectedToken, actualToken))
					t.Error(msg)
				}
			}
		})
	}
}
