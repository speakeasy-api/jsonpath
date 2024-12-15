package primitive

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseNumber parses a string into a numeric value
func ParseNumber(input string) (float64, error) {
	// Check if the input string represents an integer
	if _, err := strconv.ParseInt(input, 10, 64); err == nil {
		return strconv.ParseFloat(input, 64)
	}

	// Check if the input string represents a floating-point number
	if strings.ContainsAny(input, ".eE") {
		return strconv.ParseFloat(input, 64)
	}

	return 0, fmt.Errorf("invalid number: %s", input)
}
