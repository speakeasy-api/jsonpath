package jsonpath_test

import (
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath/config"
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

func TestParserPropertyNameExtension(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		enabled bool
		valid   bool
	}{
		{
			name:    "Simple property name disabled",
			input:   "$.store~",
			enabled: false,
			valid:   false,
		},
		{
			name:    "Simple property name enabled",
			input:   "$.store~",
			enabled: true,
			valid:   true,
		},
		{
			name:    "Property name in filter disabled",
			input:   "$[?(@~)]",
			enabled: false,
			valid:   false,
		},
		{
			name:    "Property name in filter enabled",
			input:   "$[?(@~)]",
			enabled: true,
			valid:   true,
		},
		{
			name:    "Property name with bracket notation enabled",
			input:   "$['store']~",
			enabled: true,
			valid:   true,
		},
		{
			name:    "Property name with bracket notation disabled",
			input:   "$['store']~",
			enabled: false,
			valid:   false,
		},
		{
			name:    "Chained property names enabled",
			input:   "$.store~.name~",
			enabled: true,
			valid:   true,
		},
		{
			name:    "Property name in complex filter enabled",
			input:   "$[?(@~ && @.price < 10)]",
			enabled: true,
			valid:   true,
		},
		{
			name:    "Property name in complex filter disabled",
			input:   "$[?(@~ && @.price < 10)]",
			enabled: false,
			valid:   false,
		},
		{
			name:    "Missing closing a filter expression shouldn't crash",
			input:   "$.paths.*.*[?(!@.servers)",
			enabled: false,
			valid:   false,
		},
		{
			name:    "Missing closing a filter expression shouldn't crash",
			input:   "$.paths.*.*[?(!@.servers)",
			enabled: false,
			valid:   false,
		},
		{
			name:    "Missing closing a array crash",
			input:   "$.paths.*[?@[\"x-my-ignore\"]",
			enabled: false,
			valid:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var opts []config.Option
			if test.enabled {
				opts = append(opts, config.WithPropertyNameExtension())
			}

			path, err := jsonpath.NewPath(test.input, opts...)
			if !test.valid {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, test.input, path.String())
		})
	}
}
