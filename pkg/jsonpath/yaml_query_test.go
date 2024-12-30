package jsonpath

import (
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath/token"
	"gopkg.in/yaml.v3"
	"reflect"
	"strings"
	"testing"
)

func TestQuery(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		yaml     string
		expected []string
	}{
		{
			name:     "Root node",
			input:    "$",
			yaml:     "foo",
			expected: []string{"foo"},
		},
		{
			name:  "Single child",
			input: "$.store",
			yaml: `
store:
  book: 
    - title: Book 1
    - title: Book 2
`,
			expected: []string{
				"book:\n    - title: Book 1\n    - title: Book 2",
			},
		},
		{
			name:  "Multiple children",
			input: "$.store.book[*].title",
			yaml: `
store:
  book: 
    - title: Book 1
    - title: Book 2
`,
			expected: []string{"Book 1", "Book 2"},
		},
		{
			name:     "Array index",
			input:    "$[1]",
			yaml:     "[foo, bar, baz]",
			expected: []string{"bar"},
		},
		{
			name:     "Array slice",
			input:    "$[1:3]",
			yaml:     "[foo, bar, baz, qux]",
			expected: []string{"bar", "baz"},
		},
		{
			name:     "Array slice with step",
			input:    "$[0:5:2]",
			yaml:     "[foo, bar, baz, qux, quux]",
			expected: []string{"foo", "baz", "quux"},
		},
		{
			name:  "Filter expression",
			input: "$.store.book[?(@.price < 10)].title",
			yaml: `
store:
  book:
    - title: Book 1 
      price: 9.99
    - title: Book 2
      price: 12.99
`,
			expected: []string{"Book 1"},
		},
		{
			name:  "Nested filter expression",
			input: "$.store.book[?(@.price < 10 && @.category == 'fiction')].title",
			yaml: `
store:
  book:
    - title: Book 1
      price: 9.99
      category: fiction
    - title: Book 2
      price: 8.99
      category: non-fiction
`,
			expected: []string{"Book 1"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var root yaml.Node
			err := yaml.Unmarshal([]byte(test.yaml), &root)
			if err != nil {
				t.Errorf("Error parsing YAML: %v", err)
				return
			}

			tokenizer := token.NewTokenizer(test.input)
			parser := NewParser(tokenizer, tokenizer.Tokenize())
			err = parser.Parse()
			if err != nil {
				t.Errorf("Error parsing JSON path: %v", err)
				return
			}

			result := parser.path.Query(&root, &root)
			var actual []string
			for _, node := range result {
				actual = append(actual, nodeToString(node))
			}

			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("Expected:\n%v\nGot:\n%v", test.expected, actual)
			}
		})
	}
}

func nodeToString(node *yaml.Node) string {
	var builder strings.Builder
	err := yaml.NewEncoder(&builder).Encode(node)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(builder.String())
}
