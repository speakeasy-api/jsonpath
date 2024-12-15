package selector

import (
	"fmt"
	"strconv"
	"strings"
)

// Slice represents a slice selector in JSONPath
type Slice struct {
	Start int64
	End   int64
	Step  int64
}

// NewSlice creates a new Slice instance
func NewSlice() Slice {
	return Slice{
		Start: 0,
		End:   0,
		Step:  1,
	}
}

// ParseArraySlice parses an array slice selector from the input string
func ParseArraySlice(input string) (Slice, error) {
	var slice Slice

	parts := strings.Split(input, ":")
	if len(parts) > 3 {
		return slice, fmt.Errorf("invalid array slice selector: %s", input)
	}

	if len(parts) >= 1 && parts[0] != "" {
		start, err := strconv.ParseInt(strings.TrimSpace(parts[0]), 10, 64)
		if err != nil {
			return slice, fmt.Errorf("invalid start value: %s", parts[0])
		}
		slice.Start = start
	}

	if len(parts) >= 2 && parts[1] != "" {
		end, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
		if err != nil {
			return slice, fmt.Errorf("invalid end value: %s", parts[1])
		}
		slice.End = end
	}

	if len(parts) == 3 && parts[2] != "" {
		step, err := strconv.ParseInt(strings.TrimSpace(parts[2]), 10, 64)
		if err != nil {
			return slice, fmt.Errorf("invalid step value: %s", parts[2])
		}
		slice.Step = step
	}

	return slice, nil
}
