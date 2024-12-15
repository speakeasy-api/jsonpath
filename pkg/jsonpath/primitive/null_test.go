package primitive

import "testing"

func TestParseNull(t *testing.T) {
	tests := []struct {
		input string
		err   bool
	}{
		{"null", false},
		{"Null", false},
		{"NULL", false},
		{"invalid", true},
		{"nil", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			err := ParseNull(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("ParseNull() error = %v, wantErr %v", err, tt.err)
			}
		})
	}
}
