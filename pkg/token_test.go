package pkg

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
				{Token: ROOT, Line: 1, Column: 0, Literal: ""},
				{Token: EOF, Line: 1, Column: 1, Literal: ""},
			},
		},
		{
			name:  "Child",
			input: "$.store.book",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: ""},
				{Token: CHILD, Line: 1, Column: 1, Literal: ""},
				{Token: LITERAL, Line: 1, Column: 2, Literal: "store"},
				{Token: CHILD, Line: 1, Column: 7, Literal: ""},
				{Token: LITERAL, Line: 1, Column: 8, Literal: "book"},
				{Token: EOF, Line: 1, Column: 12, Literal: ""},
			},
		},
		{
			name:  "Wildcard",
			input: "$.*",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: ""},
				{Token: CHILD, Line: 1, Column: 1, Literal: ""},
				{Token: WILDCARD, Line: 1, Column: 2, Literal: ""},
				{Token: EOF, Line: 1, Column: 3, Literal: ""},
			},
		},
		{
			name:  "Recursive",
			input: "$..author",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: ""},
				{Token: RECURSIVE, Line: 1, Column: 1, Literal: ""},
				{Token: CHILD, Line: 1, Column: 2, Literal: ""},
				{Token: LITERAL, Line: 1, Column: 3, Literal: "author"},
				{Token: EOF, Line: 1, Column: 9, Literal: ""},
			},
		},
		{
			name:  "Union",
			input: "$..book[0,1]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: ""},
				{Token: RECURSIVE, Line: 1, Column: 1, Literal: ""},
				{Token: CHILD, Line: 1, Column: 2, Literal: ""},
				{Token: LITERAL, Line: 1, Column: 3, Literal: "book"},
				{Token: BRACKET_LEFT, Line: 1, Column: 7, Literal: ""},
				{Token: NUMBER, Line: 1, Column: 8, Literal: "0"},
				{Token: UNION, Line: 1, Column: 9, Literal: ""},
				{Token: NUMBER, Line: 1, Column: 10, Literal: "1"},
				{Token: BRACKET_RIGHT, Line: 1, Column: 11, Literal: ""},
				{Token: EOF, Line: 1, Column: 12, Literal: ""},
			},
		},
		{
			name:  "Slice",
			input: "$..book[0:2]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: ""},
				{Token: RECURSIVE, Line: 1, Column: 1, Literal: ""},
				{Token: CHILD, Line: 1, Column: 2, Literal: ""},
				{Token: LITERAL, Line: 1, Column: 3, Literal: "book"},
				{Token: BRACKET_LEFT, Line: 1, Column: 7, Literal: ""},
				{Token: NUMBER, Line: 1, Column: 8, Literal: "0"},
				{Token: SLICE, Line: 1, Column: 9, Literal: ""},
				{Token: NUMBER, Line: 1, Column: 10, Literal: "2"},
				{Token: BRACKET_RIGHT, Line: 1, Column: 11, Literal: ""},
				{Token: EOF, Line: 1, Column: 12, Literal: ""},
			},
		},
		{
			name:  "Filter",
			input: "$.store.book[?(@.price < 10)]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: ""},
				{Token: CHILD, Line: 1, Column: 1, Literal: ""},
				{Token: LITERAL, Line: 1, Column: 2, Literal: "store"},
				{Token: CHILD, Line: 1, Column: 7, Literal: ""},
				{Token: LITERAL, Line: 1, Column: 8, Literal: "book"},
				{Token: BRACKET_LEFT, Line: 1, Column: 12, Literal: ""},
				{Token: FILTER, Line: 1, Column: 13, Literal: ""},
				{Token: PAREN_LEFT, Line: 1, Column: 14, Literal: ""},
				{Token: CURRENT, Line: 1, Column: 15, Literal: ""},
				{Token: CHILD, Line: 1, Column: 16, Literal: ""},
				{Token: LITERAL, Line: 1, Column: 17, Literal: "price"},
				{Token: LITERAL, Line: 1, Column: 23, Literal: "<"},
				{Token: NUMBER, Line: 1, Column: 25, Literal: "10"},
				{Token: PAREN_RIGHT, Line: 1, Column: 27, Literal: ""},
				{Token: BRACKET_RIGHT, Line: 1, Column: 28, Literal: ""},
				{Token: EOF, Line: 1, Column: 29, Literal: ""},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenizer := NewTokenizer(test.input)
			tokens := tokenizer.Tokenize()

			if len(tokens) != len(test.expected) {
				t.Fatalf("Expected %d tokens, got %d", len(test.expected), len(tokens))
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
