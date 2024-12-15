package primitive

import (
	"testing"
)

func TestParseNumber(t *testing.T) {
	tests := []struct {
		input string
		want  float64
		err   bool
	}{
		{"123", 123, false},
		{"-1", -1, false},
		{"1e10", 1e10, false},
		{"1.0001", 1.0001, false},
		{"-0", -0.0, false},
		{"invalid", 0, true},
		{"1.2.3", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseNumber(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("ParseNumber() error = %v, wantErr %v", err, tt.err)
				return
			}
			if got != tt.want {
				t.Errorf("ParseNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
