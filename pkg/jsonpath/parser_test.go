package jsonpath_test

import (
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		invalid bool
	}{
		{
			name:  "Root node",
			input: "$",
		},
		{
			name:  "Single Dot child",
			input: "$.store",
		},
		{
			name:  "Single Bracket child",
			input: "$['store']",
		},
		{
			name:  "Bracket child",
			input: "$['store']['book']",
		},
		{
			name:  "Array index",
			input: "$[0]",
		},
		{
			name:  "Array slice",
			input: "$[1:3]",
		},
		{
			name:  "Array slice with step",
			input: "$[0:5:2]",
		},
		{
			name:  "Array slice with negative step",
			input: "$[5:1:-2]",
		},
		{
			name:  "Filter expression",
			input: "$[?(@.price < 10)]",
		},
		{
			name:  "Nested filter expression",
			input: "$[?(@.price < 10 && @.category == 'fiction')]",
		},
		{
			name:  "Function call",
			input: "$.books[?(length(@) > 100)]",
		},
		{
			name:    "Invalid missing closing ]",
			input:   "$.paths.['/pet'",
			invalid: true,
		},
		{
			name:    "Invalid extra input",
			input:   "$.paths.['/pet')",
			invalid: true,
		},
		{
			name:  "Valid filter",
			input: "$.paths[?(1 == 1)]",
		},
		{
			name:    "Invalid filter",
			input:   "$.paths[?(true]",
			invalid: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path, err := jsonpath.NewPath(test.input)
			if test.invalid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.input, path.String())
		})
	}
}
