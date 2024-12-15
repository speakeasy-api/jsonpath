package primitive

import (
	"testing"
)

func TestNewInteger(t *testing.T) {
	tests := []struct {
		name    string
		value   int64
		wantErr bool
	}{
		{"zero", 0, false},
		{"positive valid", MAX, false},
		{"negative valid", MIN, false},
		{"above max", MAX + 1, true},
		{"below min", MIN - 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, err := NewInteger(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewInteger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && i.Int64() != tt.value {
				t.Errorf("NewInteger() = %v, want %v", i.Int64(), tt.value)
			}
		})
	}
}

func TestInteger_CheckedOperations(t *testing.T) {
	i1, _ := NewInteger(5)
	i2, _ := NewInteger(3)
	ineg, _ := NewInteger(-3)

	tests := []struct {
		name    string
		op      func(Integer, Integer) *Integer
		a, b    Integer
		want    *Integer
		wantNil bool
	}{
		{"add positive", Integer.CheckedAdd, i1, i2, &Integer{8}, false},
		{"add negative", Integer.CheckedAdd, i1, ineg, &Integer{2}, false},
		{"sub positive", Integer.CheckedSub, i1, i2, &Integer{2}, false},
		{"mul positive", Integer.CheckedMul, i1, i2, &Integer{15}, false},

		// Test overflow cases
		{"add overflow",
			Integer.CheckedAdd,
			Integer{MAX},
			Integer{1},
			nil,
			true},
		{"sub overflow",
			Integer.CheckedSub,
			Integer{MIN},
			Integer{1},
			nil,
			true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.op(tt.a, tt.b)
			if tt.wantNil {
				if got != nil {
					t.Errorf("Expected nil result, got %v", got)
				}
				return
			}
			if got == nil {
				t.Fatal("Unexpected nil result")
			}
			if got.Int64() != tt.want.Int64() {
				t.Errorf("Got %v, want %v", got.Int64(), tt.want.Int64())
			}
		})
	}
}

func TestInteger_Abs(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		want  int64
	}{
		{"positive", 5, 5},
		{"negative", -5, 5},
		{"zero", 0, 0},
		{"min valid", MIN, -MIN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, _ := NewInteger(tt.value)
			got := i.Abs()
			if got.Int64() != tt.want {
				t.Errorf("Abs() = %v, want %v", got.Int64(), tt.want)
			}
		})
	}
}

func TestParseInteger(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"zero", "0", 0, false},
		{"positive", "42", 42, false},
		{"negative", "-42", -42, false},
		{"invalid", "abc", 0, true},
		{"too large", "9007199254740993", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseInteger(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInteger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.Int64() != tt.want {
				t.Errorf("ParseInteger() = %v, want %v", got.Int64(), tt.want)
			}
		})
	}
}
