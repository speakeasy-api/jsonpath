package token

import (
	"fmt"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath/config"
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
			name:  "child",
			input: "$.store.book",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 2, Literal: "store", Len: 5},
				{Token: CHILD, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 8, Literal: "book", Len: 4},
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
				{Token: STRING, Line: 1, Column: 3, Literal: "author", Len: 6},
			},
		},
		{
			name:  "Union",
			input: "$..book[0,1]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: RECURSIVE, Line: 1, Column: 1, Literal: "", Len: 2},
				{Token: STRING, Line: 1, Column: 3, Literal: "book", Len: 4},
				{Token: BRACKET_LEFT, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: INTEGER, Line: 1, Column: 8, Literal: "0", Len: 1},
				{Token: COMMA, Line: 1, Column: 9, Literal: "", Len: 1},
				{Token: INTEGER, Line: 1, Column: 10, Literal: "1", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 11, Literal: "", Len: 1},
			},
		},
		{
			name:  "Slice",
			input: "$..book[0:2]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: RECURSIVE, Line: 1, Column: 1, Literal: "", Len: 2},
				{Token: STRING, Line: 1, Column: 3, Literal: "book", Len: 4},
				{Token: BRACKET_LEFT, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: INTEGER, Line: 1, Column: 8, Literal: "0", Len: 1},
				{Token: ARRAY_SLICE, Line: 1, Column: 9, Literal: "", Len: 1},
				{Token: INTEGER, Line: 1, Column: 10, Literal: "2", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 11, Literal: "", Len: 1},
			},
		},
		{
			name:  "Filter",
			input: "$.store.book[?(@.price < 10)]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 2, Literal: "store", Len: 5},
				{Token: CHILD, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 8, Literal: "book", Len: 4},
				{Token: BRACKET_LEFT, Line: 1, Column: 12, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 13, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 14, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 15, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 16, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 17, Literal: "price", Len: 5},
				{Token: LT, Line: 1, Column: 23, Literal: "", Len: 1},
				{Token: INTEGER, Line: 1, Column: 25, Literal: "10", Len: 2},
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
				{Token: STRING, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: EQ, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: STRING_LITERAL, Line: 1, Column: 13, Literal: "x", Len: 3},
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
				{Token: STRING, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: NE, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: STRING_LITERAL, Line: 1, Column: 13, Literal: "x", Len: 3},
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
				{Token: STRING, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: GT, Line: 1, Column: 11, Literal: "", Len: 1},
				{Token: INTEGER, Line: 1, Column: 12, Literal: "1", Len: 1},
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
				{Token: STRING, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: GE, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: INTEGER, Line: 1, Column: 13, Literal: "1", Len: 1},
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
				{Token: STRING, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: LT, Line: 1, Column: 11, Literal: "", Len: 1},
				{Token: INTEGER, Line: 1, Column: 12, Literal: "1", Len: 1},
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
				{Token: STRING, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: LE, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: INTEGER, Line: 1, Column: 13, Literal: "1", Len: 1},
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
				{Token: STRING, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: AND, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: CURRENT, Line: 1, Column: 13, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 14, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 15, Literal: "other", Len: 5},
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
				{Token: STRING, Line: 1, Column: 6, Literal: "child", Len: 5},
				{Token: OR, Line: 1, Column: 11, Literal: "", Len: 2},
				{Token: CURRENT, Line: 1, Column: 13, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 14, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 15, Literal: "other", Len: 5},
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
				{Token: STRING, Line: 1, Column: 7, Literal: "child", Len: 5},
				{Token: PAREN_RIGHT, Line: 1, Column: 12, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 13, Literal: "", Len: 1},
			},
		},
		{
			name:  "Underscore is string literal character",
			input: "$.pagination._.next_results_cursor",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 2, Literal: "pagination", Len: 10},
				{Token: CHILD, Line: 1, Column: 12, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 13, Literal: "_", Len: 1},
				{Token: CHILD, Line: 1, Column: 14, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 15, Literal: "next_results_cursor", Len: 19},
			},
		},
		{
			name:  "Function call with no arguments",
			input: "$.books[?(@.length())]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 2, Literal: "books", Len: 5},
				{Token: BRACKET_LEFT, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 8, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 9, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 10, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 11, Literal: "", Len: 1},
				{Token: FUNCTION, Line: 1, Column: 12, Literal: "length", Len: 6},
				{Token: PAREN_LEFT, Line: 1, Column: 18, Literal: "", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 19, Literal: "", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 20, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 21, Literal: "", Len: 1},
			},
		},
		{
			name:  "Function call with one argument",
			input: "$.books[?(@.count('fiction'))]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 2, Literal: "books", Len: 5},
				{Token: BRACKET_LEFT, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 8, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 9, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 10, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 11, Literal: "", Len: 1},
				{Token: FUNCTION, Line: 1, Column: 12, Literal: "count", Len: 5},
				{Token: PAREN_LEFT, Line: 1, Column: 17, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 18, Literal: "fiction", Len: 9},
				{Token: PAREN_RIGHT, Line: 1, Column: 27, Literal: "", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 28, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 29, Literal: "", Len: 1},
			},
		},
		{
			name:  "Function call with multiple arguments",
			input: "$.books[?(@.match('fiction', 'adventure'))]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 2, Literal: "books", Len: 5},
				{Token: BRACKET_LEFT, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 8, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 9, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 10, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 11, Literal: "", Len: 1},
				{Token: FUNCTION, Line: 1, Column: 12, Literal: "match", Len: 5},
				{Token: PAREN_LEFT, Line: 1, Column: 17, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 18, Literal: "fiction", Len: 9},
				{Token: COMMA, Line: 1, Column: 27, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 29, Literal: "adventure", Len: 11},
				{Token: PAREN_RIGHT, Line: 1, Column: 40, Literal: "", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 41, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 42, Literal: "", Len: 1},
			},
		},
		{
			name:  "Nested function calls",
			input: "$.books[?(@.count(@.search('fiction')))]",
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 2, Literal: "books", Len: 5},
				{Token: BRACKET_LEFT, Line: 1, Column: 7, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 8, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 9, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 10, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 11, Literal: "", Len: 1},
				{Token: FUNCTION, Line: 1, Column: 12, Literal: "count", Len: 5},
				{Token: PAREN_LEFT, Line: 1, Column: 17, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 18, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 19, Literal: "", Len: 1},
				{Token: FUNCTION, Line: 1, Column: 20, Literal: "search", Len: 6},
				{Token: PAREN_LEFT, Line: 1, Column: 26, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 27, Literal: "fiction", Len: 9},
				{Token: PAREN_RIGHT, Line: 1, Column: 36, Literal: "", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 37, Literal: "", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 38, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 39, Literal: "", Len: 1},
			},
		},
		//{
		//	name:  "Filter regular expression (illegal right now)",
		//	input: "$[?(@.child=~/.*/)]",
		//	expected: []TokenInfo{
		//		{Token: ROOT, Line: 1, Column: 0, literal: "", Len: 1},
		//		{Token: FILTER, Line: 1, Column: 1, literal: "", Len: 2},
		//		{Token: PAREN_LEFT, Line: 1, Column: 3, literal: "", Len: 1},
		//		{Token: CURRENT, Line: 1, Column: 4, literal: "", Len: 1},
		//		{Token: CHILD, Line: 1, Column: 5, literal: "", Len: 1},
		//		{Token: STRING, Line: 1, Column: 6, literal: "child", Len: 5},
		//		{Token: MATCHES, Line: 1, Column: 11, literal: "", Len: 2},
		//		{Token: ILLEGAL, Line: 1, Column: 13, literal: "", Len: 1},
		//		{Token: PAREN_RIGHT, Line: 1, Column: 17, literal: "", Len: 1},
		//		{Token: BRACKET_RIGHT, Line: 1, Column: 18, literal: "", Len: 1},
		//	},
		//},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenizer := NewTokenizer(test.input)
			tokens := tokenizer.Tokenize()

			if len(tokens) != len(test.expected) {
				msg := tokenizer.ErrorTokenString(&tokens[0], fmt.Sprintf("Expected %d tokens, got %d", len(test.expected), len(tokens)))
				t.Error(msg)
			}

			for i, expectedToken := range test.expected {
				actualToken := tokens[i]
				if actualToken != expectedToken {
					msg := tokenizer.ErrorString(&actualToken, fmt.Sprintf("Expected token %+v, got %+v", expectedToken, actualToken))
					t.Error(msg)
				}
			}
		})
	}
}

func TestTokenizer_categorize(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		illegal bool
		simple  bool
	}{
		{name: "identity", path: "", simple: true},
		{name: "root", path: "$", simple: true},
		{name: "unmatched closing parenthesis", path: ")", illegal: true},
		{name: "unmatched closing square bracket", path: "]", illegal: true},
		{name: "dot child", path: "$.child", simple: true},
		{name: "dot child with implicit root", path: ".child"},
		{name: "undotted child with implicit root", path: "child"},
		{name: "dot child with no name", path: "$.", simple: true},
		{name: "dot child with missing dot", path: "$a", simple: true},
		{name: "dot child with trailing dot", path: "$.child.", simple: true},
		{name: "dot child of dot child", path: "$.child1.child2", simple: true},
		{name: "dot child with array subscript", path: "$.child[*]"},
		{name: "dot child with malformed array subscript", path: "$.child[1:2:3:4]"},
		{name: "dot child with array subscript with zero step", path: "$.child[1:2:0]"},
		{name: "dot child with non-integer array subscript", path: "$.child[1:2:a]"},
		{name: "dot child with unclosed array subscript", path: "$.child[*", illegal: true},
		{name: "dot child with missing array subscript", path: "$.child[]", simple: true},
		{name: "dot child with embedded space", path: "$.child more", simple: true},
		{name: "bracket child", path: "$['child']", simple: true},
		{name: "bracket child with double quotes", path: `$["child"]`, simple: true},
		{name: "bracket child with unmatched quotes", path: `$["child']`, illegal: true},
		{name: "bracket child with empty name", path: "$['']", simple: true},
		{name: "bracket child of bracket child", path: "$['child1']['child2']", simple: true},
		{name: "double quoted bracket child of bracket child", path: `$['child1']["child2"]`, simple: true},
		{name: "bracket child union", path: "$['child','child2']"},
		{name: "bracket child union with whitespace", path: "$[ 'child' , 'child2' ]"},
		{name: "bracket child union with mixed quotes", path: `$[ 'child' , "child2" ]`},
		{name: "bracket child quoted union literal", path: "$[',']", simple: true},
		{name: "bracket child with array subscript", path: "$['child'][*]"},
		{name: "bracket child with malformed array subscript", path: "$['child'][1:2:3:4]"},
		{name: "bracket child with non-integer array subscript", path: "$['child'][1:2:a]"},
		{name: "bracket child with unclosed array subscript", path: "$['child'][*", illegal: true},
		{name: "bracket child with missing array subscript", path: "$['child'][]", simple: true},
		{name: "bracket child followed by space", path: "$['child'] ", illegal: true},
		{name: "bracket dotted child", path: "$['child1.child2']", simple: true},
		{name: "bracket child with array subscript", path: "$['child'][*]"},
		{name: "property name dot child", path: "$.child~", illegal: true},
		{name: "property name dot child with implicit root", path: ".child~", illegal: true},
		{name: "property name undotted child with implicit root", path: "child~", illegal: true},
		{name: "property name dot child with no name", path: "$.~", illegal: true},
		{name: "property name dot child with missing dot", path: "$a~", illegal: true},
		{name: "property name dot child with trailing chars", path: "$.child~.test", illegal: true},
		{name: "property name undotted child with trailing chars", path: "child~.test", illegal: true},
		{name: "property name dot child with trailing dot", path: "$.child.~", illegal: true},
		{name: "property name dot child of dot child", path: "$.child1.child2~", illegal: true},
		{name: "property name dot child with wildcard array subscript", path: "$.child[*]~", illegal: true},
		{name: "property name dot child with an array subscript", path: "$.child[0]~", illegal: true},
		{name: "property name dot child with array subscript with zero step", path: "$.child[1:2:0]~", illegal: true},
		{name: "property name dot child with non-integer array subscript", path: "$.child[1:2:a]~", illegal: true},
		{name: "property name dot child with unclosed array subscript", path: "$.child[*~", illegal: true},
		{name: "property name dot child with missing array subscript", path: "$.child[]~", illegal: true},
		{name: "property name dot child with embedded space", path: "$.child more~", illegal: true},
		{name: "property name bracket child", path: "$['child']~", illegal: true},
		{name: "property name bracket child with double quotes", path: `$["child"]~`, illegal: true},
		{name: "property name bracket child with unmatched quotes", path: `$["child']~`, illegal: true},
		{name: "property name bracket child with empty name", path: "$['']~", illegal: true},
		{name: "property name bracket child of bracket child", path: "$['child1']['child2']~", illegal: true},
		{name: "property name double quoted bracket child of bracket child", path: `$['child1']["child2"]~`, illegal: true},
		{name: "property name bracket child union", path: "$['child','child2']~", illegal: true},
		{name: "property name bracket child union with whitespace", path: "$[ 'child' , 'child2' ]~", illegal: true},
		{name: "property name bracket child union with mixed quotes", path: `$[ 'child' , "child2" ]~`, illegal: true},
		{name: "property name bracket child quoted union literal", path: "$[',']~", illegal: true},
		{name: "property name bracket child with wildcard array subscript", path: "$['child'][*]~", illegal: true},
		{name: "property name bracket child with wildcard array subscript and trailing chars", path: "$['child'][*]~.child", illegal: true},
		{name: "property name bracket child with ~ in name", path: "$['child~']~", illegal: true},
		{name: "bracket child with array subscript", path: "$['child'][1]~", illegal: true},
		{name: "property name bracket child with non-integer array subscript", path: "$['child'][1:2:a]~", illegal: true},
		{name: "property name bracket child with unclosed array subscript", path: "$['child'][*~", illegal: true},
		{name: "property name bracket child with missing array subscript", path: "$['child'][]~", illegal: true},
		{name: "property name bracket child separated a  by space", path: "$['child'] ~", illegal: true},
		{name: "property name bracket child followed by space", path: "$['child']~ ", illegal: true},
		{name: "property name bracket dotted child", path: "$['child1.child2']~", illegal: true},
		{name: "array union", path: "$[0,1]"},
		{name: "array union with whitespace", path: "$[ 0 , 1 ]"},
		{name: "bracket child with malformed array subscript", path: "$['child'][1:2:3:4]"},
		{name: "bracket child with malformed array subscript in union", path: "$['child'][0,1:2:3:4]"},
		{name: "bracket child with non-integer array subscript", path: "$['child'][1:2:a]"},
		{name: "bracket child of dot child", path: "$.child1['child2']", simple: true},
		{name: "array slice of root", path: "$[1:3]"},
		{name: "dot child of bracket child", path: "$['child1'].child2", simple: true},
		{name: "recursive descent", path: "$..child"},
		{name: "recursive descent of dot child", path: "$.child1..child2"},
		{name: "recursive descent of bracket child", path: "$['child1']..child2"},
		{name: "repeated recursive descent", path: "$..child1..child2"},
		{name: "recursive descent with dot child", path: "$..child1.child2"},
		{name: "recursive descent with bracket child", path: "$..child1['child2']"},
		{name: "recursive descent with missing name", path: "$.."},
		{name: "recursive descent with array access", path: "$..[0]"},
		{name: "recursive descent with filter", path: "$..[?(@.child)]"},
		{name: "recursive descent with bracket child", path: "$..['child']"},
		{name: "recursive descent with double quoted bracket child", path: `$..["child"]`},
		{name: "wildcarded children", path: "$.*"},
		{name: "simple filter", path: "$[?(@.child)]"},
		{name: "simple filter with leading whitespace", path: "$[?( @.child)]"},
		{name: "simple filter with trailing whitespace", path: "$[?( @.child )]"},
		{name: "simple filter with bracket", path: "$[?((@.child))]"},
		{name: "simple filter with bracket with extra whitespace", path: "$[?( ( @.child ) )]"},
		{name: "simple filter with more complex subpath", path: "$[?((@.child[0]))]"},
		{name: "missing filter ", path: "$[?()]"},
		{name: "unclosed filter", path: "$[?(", illegal: true},
		{name: "filter with missing operator", path: "$[?(@.child @.other)]"},
		{name: "filter with malformed term", path: "$[?([)]", illegal: true},
		{name: "filter with misplaced open bracket", path: "$[?(@.child ()]", illegal: true},
		{name: "simple negative filter", path: "$[?(!@.child)]"},
		{name: "misplaced filter negation", path: "$[?(@.child !@.other)]"},
		{name: "simple negative filter with extra whitespace", path: "$[?( ! @.child)]"},
		{name: "simple filter with root expression", path: "$[?($.child)]"},
		{name: "filter integer equality, literal on the right", path: "$[?(@.child==1)]"},
		{name: "filter string equality, literal on the right", path: "$[?(@.child=='x')]"},
		{name: "filter string equality with apparent boolean", path: `$[?(@.child=="true")]`},
		{name: "filter string equality with apparent null", path: `$[?(@.child=="null")]`},
		{name: "filter string equality, double-quoted literal on the right", path: `$[?(@.child=="x")]`},
		{name: "filter integer equality with invalid literal", path: "$[?(@.child==-)]", illegal: true},
		{name: "filter integer equality with integer literal which is too large", path: "$[?(@.child==9223372036854775808)]"},
		{name: "filter integer equality with invalid float literal", path: "$[?(@.child==1.2.3)]", illegal: true},
		{name: "filter integer equality with invalid string literal", path: "$[?(@.child=='x)]", illegal: true},
		{name: "filter integer equality, literal on the left", path: "$[?(1==@.child)]"},
		{name: "filter float equality, literal on the left", path: "$[?(1.5==@.child)]"},
		{name: "filter fractional float equality, literal on the left", path: "$[?(-1.5e-1==@.child)]"},
		{name: "filter fractional float equality, literal on the right", path: "$[?(@.child== -1.5e-1 )]"},
		{name: "filter boolean true equality, literal on the right", path: "$[?(@.child== true )]"},
		{name: "filter boolean false equality, literal on the right", path: "$[?(@.child==false)]"},
		{name: "filter boolean true equality, literal on the left", path: "$[?(true==@.child)]"},
		{name: "filter boolean false equality, literal on the left", path: "$[?( false ==@.child)]"},
		{name: "filter null equality, literal on the right", path: "$[?(@.child==null)]"},
		{name: "filter null true equality, literal on the left", path: "$[?(null==@.child)]"},
		{name: "filter equality with missing left hand value", path: "$[?(==@.child)]"},
		{name: "filter equality with missing left hand value inside bracket", path: "$[?((==@.child))]"},
		{name: "filter equality with missing right hand value", path: "$[?(@.child==)]"},
		{name: "filter integer equality, root path on the right", path: "$[?(@.child==$.x)]"},
		{name: "filter integer equality, root path on the left", path: "$[?($.x==@.child)]"},
		{name: "filter string equality, literal on the right", path: "$[?(@.child=='x')]"},
		{name: "filter string equality, literal on the left", path: "$[?('x'==@.child)]"},
		{name: "filter string equality, literal on the left with unmatched string delimiter", path: "$[?('x==@.child)]", illegal: true},
		{name: "filter string equality with unmatched string delimiter", path: "$[?(@.child=='x)]", illegal: true},
		{name: "filter integer inequality, literal on the right", path: "$[?(@.child!=1)]"},
		{name: "filter inequality with missing left hand operator", path: "$[?(!=1)]"},
		{name: "filter equality with missing right hand value", path: "$[?(@.child!=)]"},
		{name: "filter greater than, integer literal on the right", path: "$[?(@.child>1)]"},
		{name: "filter greater than, decimal literal on the right", path: "$[?(@.child> 1.5)]"},
		{name: "filter greater than, path to path", path: "$[?(@.child1>@.child2)]"},
		{name: "filter greater than with left hand operand missing", path: "$[?(>1)]"},
		{name: "filter greater than with missing right hand value", path: "$[?(@.child>)]"},
		{name: "filter greater than, string on the right", path: "$[?(@.child>'x')]"},
		{name: "filter greater than, string on the left", path: "$[?('x'>@.child)]"},
		{name: "filter greater than or equal, integer literal on the right", path: "$[?(@.child>=1)]"},
		{name: "filter greater than or equal, decimal literal on the right", path: "$[?(@.child>=1.5)]"},
		{name: "filter greater than or equal with left hand operand missing", path: "$[?(>=1)]"},
		{name: "filter greater than or equal with missing right hand value", path: "$[?(@.child>=)]"},
		{name: "filter greater than or equal, string on the right", path: "$[?(@.child>='x')]"},
		{name: "filter greater than or equal, string on the left", path: "$[?('x'>=@.child)]"},
		{name: "filter less than, integer literal on the right", path: "$[?(@.child<1)]"},
		{name: "filter less than, decimal literal on the right", path: "$[?(@.child< 1.5)]"},
		{name: "filter less than with left hand operand missing", path: "$[?(<1)]"},
		{name: "filter less than with missing right hand value", path: "$[?(@.child<)]"},
		{name: "filter less than, string on the right", path: "$[?(@.child<'x')]"},
		{name: "filter less than, string on the left", path: "$[?('x'<@.child)]"},
		{name: "filter less than or equal, integer literal on the right", path: "$[?(@.child<=1)]"},
		{name: "filter less than or equal, decimal literal on the right", path: "$[?(@.child<=1.5)]"},
		{name: "filter less than or equal with left hand operand missing", path: "$[?(<=1)]"},
		{name: "filter less than or equal with missing right hand value", path: "$[?(@.child<=)]"},
		{name: "filter less than or equal, string on the right", path: "$[?(@.child<='x')]"},
		{name: "filter less than or equal, string on the left", path: "$[?('x'<=@.child)]"},
		{name: "filter conjunction", path: "$[?(@.child&&@.other)]"},
		{name: "filter conjunction with literals and whitespace", path: "$[?(@.child == 'x' && -9 == @.other)]"},
		{name: "filter conjunction with bracket children", path: "$[?(@['child'][*]&&@['other'])]"},
		{name: "filter invalid leading conjunction", path: "$[?(&&", illegal: true},
		{name: "filter conjunction with extra whitespace", path: "$[?(@.child && @.other)]"},
		{name: "filter disjunction", path: "$[?(@.child||@.other)]"},
		{name: "filter invalid leading disjunction", path: "$[?(||", illegal: true},
		{name: "filter disjunction with extra whitespace", path: "$[?(@.child || @.other)]"},
		{name: "simple filter of child", path: "$.child[?(@.child)]"},
		{name: "filter with missing end", path: "$[?(@.child", illegal: true},
		{name: "nested filter (edge case)", path: "$[?(@.y[?(@.z)])]"},
		{name: "filter negation", path: "$[?(!@.child)]"},
		{name: "filter negation of comparison (edge case)", path: "$[?(!@.child>1)]"},
		{name: "filter negation of bracket", path: "$[?(!(@.child))]"},
		{name: "filter regular expression", path: "$[?(@.child=~/.*/)]", illegal: true},
		{name: "filter regular expression with escaped /", path: `$[?(@.child=~/\/.*/)]`, illegal: true},
		{name: "filter regular expression with escaped \\", path: `$[?(@.child=~/\\/)]`, illegal: true},
		{name: "filter regular expression with missing leading /", path: `$[?(@.child=~.*/)]`, illegal: true},
		{name: "filter regular expression with missing trailing /", path: `$[?(@.child=~/.*)]`, illegal: true},
		{name: "filter regular expression to match string literal", path: `$[?('x'=~/.*/)]`, illegal: true},
		{name: "filter regular expression to match integer literal", path: `$[?(0=~/.*/)]`, illegal: true},
		{name: "filter regular expression to match float literal", path: `$[?(.1=~/.*/)]`, illegal: true},
		{name: "filter invalid regular expression", path: `$[?(@.child=~/(.*/)]`, illegal: true},
		{name: "unescaped single quote in bracket child name", path: `$['single'quote']`, illegal: true},
		{name: "escaped single quote in bracket child name", path: `$['single\']quote']`, simple: true},
		{name: "escaped backslash in bracket child name", path: `$['\\']`, simple: true},
		{name: "unescaped single quote after escaped backslash in bracket child name", path: `$['single\\'quote']`, illegal: true},
		{name: "unsupported escape sequence in bracket child name", path: `$['\n']`, simple: true},
		{name: "unclosed and empty bracket child name with space", path: `$[ '`, illegal: true},
		{name: "unclosed and empty bracket child name with formfeed", path: "[\f'", illegal: true},
		{name: "filter involving value of current node on left hand side", path: "$[?(@==1)]"},
		{name: "filter involving value of current node on right hand side", path: "$[?(1==@ || 2== @ )]"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//defer func() {
			//	if r := recover(); r != nil {
			//		t.Errorf("Tokenizer panicked for path: %s\nPanic: %v", tc.path, r)
			//	}
			//}()

			tokenizer := NewTokenizer(tc.path)
			tokenizedJsonPath := tokenizer.Tokenize()
			foundIllegal := false
			for _, token := range tokenizedJsonPath {
				if token.Token == ILLEGAL {
					foundIllegal = true
					if !tc.illegal {
						t.Errorf("%s", tokenizer.ErrorString(&token, "Illegal Token"))
					}
				}
			}
			if tc.illegal && !foundIllegal {
				t.Errorf("%s", tokenizer.ErrorTokenString(&tokenizedJsonPath[0], "Expected an illegal token"))
			}

			if tc.simple && foundIllegal {
				t.Errorf("Expected a simple path, but found an illegal token")
			}

			if tc.simple && !tokenizedJsonPath.IsSimple() {
				for _, token := range tokenizedJsonPath {

					simple := false
					for _, subToken := range SimpleTokens {
						if token.Token == subToken {
							simple = true
						}
					}
					if !simple {
						t.Errorf("%s", tokenizer.ErrorString(&token, "Expected a simple path, but found a non-simple token"))
					}
				}
			}
			if !tc.simple && tokenizedJsonPath.IsSimple() {
				t.Errorf("%s", tokenizer.ErrorTokenString(&tokenizedJsonPath[0], "Expected a non-simple path, but found it was simple"))
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		err      bool
	}{
		{
			name:     "Valid double quoted string",
			input:    `"test"`,
			expected: "test",
		},
		{
			name:     "Valid double quoted string with newline",
			input:    `"test\n"`,
			expected: "test\n",
		},
		{
			name:     "Valid double quoted string with escaped quote",
			input:    `"test\""`,
			expected: `test"`,
		},
		{
			name:     "Valid single quoted string",
			input:    `'test'`,
			expected: "test",
		},
		{
			name:     "Valid single quoted string with double quote",
			input:    `'te"st'`,
			expected: `te"st`,
		},
		{
			name:     "Valid single quoted string with escaped single quote",
			input:    `'te\'st'`,
			expected: `te'st`,
		},
		{
			name:  "Invalid Unicode control character",
			input: "\u0000",
			err:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tokenizer := NewTokenizer(test.input)
			tokens := tokenizer.Tokenize()

			if test.err {
				if len(tokens) != 1 || tokens[0].Token != ILLEGAL {
					t.Errorf("Expected an illegal token, but got: %v", tokens)
				}
			} else {
				if len(tokens) != 1 {
					t.Errorf("Expected a single token, but got: %v", tokens)
				} else if tokens[0].Token != STRING_LITERAL {
					t.Errorf("Expected a STRING_LITERAL token, but got: %v", tokens[0])
				} else if tokens[0].Literal != test.expected {
					t.Errorf("Expected literal '%s', but got '%s'", test.expected, tokens[0].Literal)
				} else if tokens[0].Len != len(test.input) {
					t.Errorf("Expected length %d, but got %d", len(test.input), tokens[0].Len)
				}
			}
		})
	}
}

func TestPropertyNameExtension(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		enabled  bool
		expected []TokenInfo
	}{
		{
			name:    "Property name extension enabled",
			input:   "$.child~",
			enabled: true,
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 2, Literal: "child", Len: 5},
				{Token: PROPERTY_NAME, Line: 1, Column: 7, Literal: "", Len: 1},
			},
		},
		{
			name:    "Property name extension disabled",
			input:   "$.child~",
			enabled: false,
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: CHILD, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING, Line: 1, Column: 2, Literal: "child", Len: 5},
				{Token: ILLEGAL, Line: 1, Column: 7, Literal: "invalid property name token without config.PropertyNameExtension set to true", Len: 1},
			},
		},
		{
			name:    "Property name extension with bracket notation enabled",
			input:   "$['child']~",
			enabled: true,
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: BRACKET_LEFT, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: STRING_LITERAL, Line: 1, Column: 2, Literal: "child", Len: 7},
				{Token: BRACKET_RIGHT, Line: 1, Column: 9, Literal: "", Len: 1},
				{Token: PROPERTY_NAME, Line: 1, Column: 10, Literal: "", Len: 1},
			},
		},
		{
			name:    "Property name extension in filter with current node",
			input:   "$[?(@~)]",
			enabled: true,
			expected: []TokenInfo{
				{Token: ROOT, Line: 1, Column: 0, Literal: "", Len: 1},
				{Token: BRACKET_LEFT, Line: 1, Column: 1, Literal: "", Len: 1},
				{Token: FILTER, Line: 1, Column: 2, Literal: "", Len: 1},
				{Token: PAREN_LEFT, Line: 1, Column: 3, Literal: "", Len: 1},
				{Token: CURRENT, Line: 1, Column: 4, Literal: "", Len: 1},
				{Token: PROPERTY_NAME, Line: 1, Column: 5, Literal: "", Len: 1},
				{Token: PAREN_RIGHT, Line: 1, Column: 6, Literal: "", Len: 1},
				{Token: BRACKET_RIGHT, Line: 1, Column: 7, Literal: "", Len: 1},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var tokenizer *Tokenizer
			if test.enabled {
				tokenizer = NewTokenizer(test.input, config.WithPropertyNameExtension())
			} else {
				tokenizer = NewTokenizer(test.input)
			}

			tokens := tokenizer.Tokenize()

			if len(tokens) != len(test.expected) {
				t.Errorf("Expected %d tokens, got %d\n%s",
					len(test.expected),
					len(tokens),
					tokenizer.ErrorTokenString(&tokens[0], "Unexpected number of tokens"))
			}

			for i, expectedToken := range test.expected {
				if i >= len(tokens) {
					break
				}
				actualToken := tokens[i]
				if actualToken != expectedToken {
					t.Error(tokenizer.ErrorString(&actualToken,
						fmt.Sprintf("Expected token %+v, got %+v", expectedToken, actualToken)))
				}
			}
		})
	}
}
