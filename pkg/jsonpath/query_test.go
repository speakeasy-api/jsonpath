package jsonpath

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery_ToString(t *testing.T) {
	tests := []struct {
		name     string
		query    Query
		expected string
	}{
		{
			name: "root query",
			query: Query{
				Kind:     TokenInfo{Token: ROOT},
				Segments: []Segment{},
			},
			expected: "$",
		},
		{
			name: "current query",
			query: Query{
				Kind:     TokenInfo{Token: CURRENT},
				Segments: []Segment{},
			},
			expected: "@",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.query.ToString()
			assert.Equal(t, tt.expected, result)
		})
	}
}
