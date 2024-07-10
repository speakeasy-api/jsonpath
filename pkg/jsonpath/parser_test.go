package jsonpath

//
//func TestParser(t *testing.T) {
//	tests := []struct {
//		name     string
//		input    string
//		expected string
//	}{
//		{
//			name:     "Root node",
//			input:    "$",
//			expected: "$\n",
//		},
//		{
//			name:     "Current node",
//			input:    "@",
//			expected: "@\n",
//		},
//		{
//			name:     "Wildcard",
//			input:    "*",
//			expected: "*\n",
//		},
//		{
//			name:     "Recursive descent",
//			input:    "..",
//			expected: "..\n",
//		},
//		{
//			name:     "Dot child",
//			input:    "$.store.book",
//			expected: "$\n├── store\n└── book\n",
//		},
//		{
//			name:     "Bracket child",
//			input:    "$['store']['book']",
//			expected: "$\n├── ['store']\n└── ['book']\n",
//		},
//		{
//			name:     "Array index",
//			input:    "$[0]",
//			expected: "$\n└── [0]\n",
//		},
//		{
//			name:     "Array slice",
//			input:    "$[1:3]",
//			expected: "$\n└── [1:3]\n   ├── 1\n   └── 3\n",
//		},
//		{
//			name:     "Array slice with step",
//			input:    "$[0:5:2]",
//			expected: "$\n└── [0:5:2]\n   ├── 0\n   ├── 5\n   └── 2\n",
//		},
//		{
//			name:     "Array slice with negative step",
//			input:    "$[5:1:-2]",
//			expected: "$\n└── [5:1:-2]\n   ├── 5\n   ├── 1\n   └── -2\n",
//		},
//		{
//			name:     "Filter expression",
//			input:    "$[?(@.price < 10)]",
//			expected: "$\n└── [?@ < 10]\n   ├── @\n   └── 10\n",
//		},
//		{
//			name:     "Nested filter expression",
//			input:    "$[?(@.price < 10 && @.category == 'fiction')]",
//			expected: "$\n└── [?@ < 10 && @ == 'fiction']\n   ├── @ < 10\n   │   ├── @\n   │   └── 10\n   └── @ == 'fiction'\n       ├── @\n       └── 'fiction'\n",
//		},
//		{
//			name:     "Function call",
//			input:    "$.books[?(@.length() > 100)]",
//			expected: "$\n├── books\n└── [?length() > 100]\n   └── length() > 100\n       ├── length()\n       └── 100\n",
//		},
//	}
//
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			tokenizer := NewTokenizer(test.input)
//
//			parser := NewParser(tokenizer)
//			result, err := parser.Parse()
//
//			if err != nil {
//				t.Errorf("Unexpected error: %v", err)
//				return
//			}
//
//			actual := PrintNode(result)
//			if actual != test.expected {
//				t.Errorf("Expected:\n%s\nGot:\n%s", test.expected, actual)
//			}
//		})
//	}
//}
