package primitive

import (
	"fmt"
	"strings"
)

// ParseNull parses a string into a null value
func ParseNull(input string) error {
	lowered := strings.ToLower(input)
	if lowered == "null" {
		return nil
	}
	return fmt.Errorf("invalid null value: %s", input)
}
