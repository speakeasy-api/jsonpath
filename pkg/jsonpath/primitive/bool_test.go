package primitive

import "testing"

func TestParseBool(t *testing.T) {
	tests := []struct {
		input string
		want  bool
		err   bool
	}{
		{"true", true, false},
		{"false", false, false},
		{"True", true, false},
		{"FaLse", false, false},
		{"invalid", false, true},
		{"1", false, true},
		{"0", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseBool(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("ParseBool() error = %v, wantErr %v", err, tt.err)
				return
			}
			if got != tt.want {
				t.Errorf("ParseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
