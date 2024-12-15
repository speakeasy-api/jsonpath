package primitive

import (
	"fmt"
	"strings"
)

// ParseBool parses a string into a boolean value
func ParseBool(input string) (bool, error) {
	lowered := strings.ToLower(input)
	if lowered == "true" {
		return true, nil
	} else if lowered == "false" {
		return false, nil
	}
	return false, fmt.Errorf("invalid boolean value: %s", input)
}
