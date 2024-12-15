package selector

import (
	"gopkg.in/yaml.v3"
	"testing"
)

// Helper function to create YAML nodes for testing
func createScalarNode(value string, tag string) *yaml.Node {
	return &yaml.Node{
		Kind:  yaml.ScalarNode,
		Tag:   tag,
		Value: value,
	}
}

func createSequenceNode(values ...*yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.SequenceNode,
		Tag:     "!!seq",
		Content: values,
	}
}

func createMappingNode(pairs ...*yaml.Node) *yaml.Node {
	return &yaml.Node{
		Kind:    yaml.MappingNode,
		Tag:     "!!map",
		Content: pairs,
	}
}

func TestLengthFunction(t *testing.T) {
	f := lengthFunction()

	tests := []struct {
		name     string
		input    *yaml.Node
		expected string
		wantErr  bool
	}{
		{
			name:     "string length",
			input:    createScalarNode("hello", "!!str"),
			expected: "5",
		},
		{
			name: "array length",
			input: createSequenceNode(
				createScalarNode("1", "!!int"),
				createScalarNode("2", "!!int"),
			),
			expected: "2",
		},
		{
			name: "object length",
			input: createMappingNode(
				createScalarNode("key1", "!!str"),
				createScalarNode("value1", "!!str"),
				createScalarNode("key2", "!!str"),
				createScalarNode("value2", "!!str"),
			),
			expected: "2",
		},
		{
			name:     "empty string",
			input:    createScalarNode("", "!!str"),
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.Evaluator([]*yaml.Node{tt.input})

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != 1 {
				t.Errorf("expected 1 result, got %d", len(result))
				return
			}

			if result[0].Value != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result[0].Value)
			}
		})
	}
}

func TestMatchFunction(t *testing.T) {
	f := matchFunction()

	tests := []struct {
		name     string
		text     string
		pattern  string
		expected bool
		wantErr  bool
	}{
		{
			name:     "exact match",
			text:     "hello",
			pattern:  "hello",
			expected: true,
		},
		{
			name:     "pattern match",
			text:     "hello123",
			pattern:  "hello\\d+",
			expected: true,
		},
		{
			name:     "no match",
			text:     "hello",
			pattern:  "world",
			expected: false,
		},
		{
			name:    "invalid pattern",
			text:    "hello",
			pattern: "[",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := createScalarNode(tt.text, "!!str")
			pattern := createScalarNode(tt.pattern, "!!str")

			result, err := f.Evaluator([]*yaml.Node{text, pattern})

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != 1 {
				t.Errorf("expected 1 result, got %d", len(result))
				return
			}

			got := result[0].Value == "true"
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestSearchFunction(t *testing.T) {
	f := searchFunction()

	tests := []struct {
		name     string
		text     string
		pattern  string
		expected bool
		wantErr  bool
	}{
		{
			name:     "substring match",
			text:     "hello world",
			pattern:  "world",
			expected: true,
		},
		{
			name:     "pattern match",
			text:     "abc123def",
			pattern:  "\\d+",
			expected: true,
		},
		{
			name:     "no match",
			text:     "hello",
			pattern:  "xyz",
			expected: false,
		},
		{
			name:    "invalid pattern",
			text:    "hello",
			pattern: "[",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text := createScalarNode(tt.text, "!!str")
			pattern := createScalarNode(tt.pattern, "!!str")

			result, err := f.Evaluator([]*yaml.Node{text, pattern})

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			got := result[0].Value == "true"
			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestCountFunction(t *testing.T) {
	f := countFunction()

	tests := []struct {
		name     string
		input    []*yaml.Node
		expected string
	}{
		{
			name: "multiple nodes",
			input: []*yaml.Node{
				createSequenceNode(
					createScalarNode("1", "!!str"),
					createScalarNode("2", "!!str"),
					createScalarNode("3", "!!str"),
				),
			},
			expected: "3",
		},
		{
			name: "single node",
			input: []*yaml.Node{
				createSequenceNode(
					createScalarNode("test", "!!str"),
				),
			},
			expected: "1",
		},
		{
			name: "empty sequence",
			input: []*yaml.Node{
				createSequenceNode(),
			},
			expected: "0",
		},
		{
			name: "mapping node",
			input: []*yaml.Node{
				createMappingNode(
					createScalarNode("key1", "!!str"),
					createScalarNode("value1", "!!str"),
				),
			},
			expected: "2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.Evaluator(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != 1 {
				t.Errorf("expected 1 result, got %d", len(result))
				return
			}

			if result[0].Value != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result[0].Value)
			}
		})
	}
}

func TestValueFunction(t *testing.T) {
	f := valueFunction()

	tests := []struct {
		name     string
		input    []*yaml.Node
		expected *yaml.Node
		wantNil  bool
	}{
		{
			name: "single node",
			input: []*yaml.Node{
				createScalarNode("test", "!!str"),
			},
			expected: createScalarNode("test", "!!str"),
		},
		{
			name: "multiple nodes",
			input: []*yaml.Node{
				createScalarNode("1", "!!str"),
				createScalarNode("2", "!!str"),
			},
			wantNil: true,
		},
		{
			name:    "empty input",
			input:   []*yaml.Node{},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := f.Evaluator(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if tt.wantNil {
				if result != nil {
					t.Error("expected nil result")
				}
				return
			}

			if len(result) != 1 {
				t.Errorf("expected 1 result, got %d", len(result))
				return
			}

			if result[0].Value != tt.expected.Value || result[0].Tag != tt.expected.Tag {
				t.Errorf("expected %v, got %v", tt.expected, result[0])
			}
		})
	}
}
