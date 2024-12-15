package primitive

import (
	"testing"
)

func TestParseStringLiteral(t *testing.T) {
	tests := []struct {
		input string
		want  string
		err   bool
	}{
		{`"test"`, "test", false},
		{`"test\n"`, "test\n", false},
		{`"test\ntest"`, "test\ntest", false},
		{`"test\""`, "test\"", false},
		{`"tes't"`, "tes't", false},
		{`'test'`, "test", false},
		{`'te"st'`, "te\"st", false},
		{`'te\'st'`, "te'st", false},
		{`"invalid escape \x"`, "", true},
		{`"invalid unicode \u123"`, "", true},
		{`"invalid unicode surrogate \uD800"`, "", true},
		{`"invalid unicode surrogate pair \uD800\u1234"`, "", true},
		{`"invalid quote'`, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseStringLiteral(tt.input)
			if (err != nil) != tt.err {
				t.Errorf("ParseStringLiteral() error = %v, wantErr %v", err, tt.err)
				return
			}
			if got != tt.want {
				t.Errorf("ParseStringLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}
