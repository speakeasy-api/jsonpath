package selector

import (
	"testing"
)

func TestParseArraySlice(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Slice
		wantErr bool
	}{
		{"valid forward 1", "1:5:1", Slice{Start: 1, End: 5, Step: 1}, false},
		{"valid forward 2", "1:10:3", Slice{Start: 1, End: 10, Step: 3}, false},
		{"valid forward 3", "1:5:-1", Slice{Start: 1, End: 5, Step: -1}, false},
		{"valid forward 4", ":5:1", Slice{Start: 0, End: 5, Step: 1}, false},
		{"valid forward 5", "1::1", Slice{Start: 1, End: 0, Step: 1}, false},
		{"valid forward 6", "1:5", Slice{Start: 1, End: 5, Step: 1}, false},
		{"valid forward 7", "::", Slice{Start: 0, End: 0, Step: 1}, false},
		{"optional whitespace 1", "1 :5:1", Slice{Start: 1, End: 5, Step: 1}, false},
		{"optional whitespace 2", "1: 5 :1", Slice{Start: 1, End: 5, Step: 1}, false},
		{"optional whitespace 3", "1: :1", Slice{Start: 1, End: 0, Step: 1}, false},
		{"optional whitespace 4", "1:5\n:1", Slice{Start: 1, End: 5, Step: 1}, false},
		{"optional whitespace 5", "1 : 5 :1", Slice{Start: 1, End: 5, Step: 1}, false},
		{"optional whitespace 6", "1:5: 1", Slice{Start: 1, End: 5, Step: 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseArraySlice(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseArraySlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseArraySlice() = %v, want %v", got, tt.want)
			}
		})
	}
}
