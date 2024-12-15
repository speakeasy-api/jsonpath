package primitive

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

// ParseStringLiteral parses a string literal from the input string
func ParseStringLiteral(input string) (string, error) {
	if len(input) < 2 {
		return "", fmt.Errorf("invalid string literal")
	}

	quote := input[0]
	if quote != '"' && quote != '\'' {
		return "", fmt.Errorf("invalid string literal")
	}

	if input[len(input)-1] != quote {
		return "", fmt.Errorf("expected an ending quote")
	}

	var result strings.Builder
	var escaped bool
	for _, r := range input[1 : len(input)-1] {
		if escaped {
			switch r {
			case 'b':
				result.WriteByte('\b')
			case 't':
				result.WriteByte('\t')
			case 'n':
				result.WriteByte('\n')
			case 'f':
				result.WriteByte('\f')
			case 'r':
				result.WriteByte('\r')
			case '/', '\\', '"', '\'':
				result.WriteRune(r)
			case 'u':
				if tail, err := parseUnicodeSequence(input); err == nil {
					result.WriteString(tail)
				} else {
					return "", err
				}
			default:
				return "", fmt.Errorf("invalid escape sequence")
			}
			escaped = false
		} else if r == '\\' {
			escaped = true
		} else {
			result.WriteRune(r)
		}
	}
	return result.String(), nil
}

func parseUnicodeSequence(input string) (string, error) {
	if len(input) < 6 {
		return "", fmt.Errorf("invalid unicode sequence")
	}

	hex := input[2:6]
	u, err := strconv.ParseUint(hex, 16, 16)
	if err != nil {
		return "", fmt.Errorf("invalid unicode sequence: %v", err)
	}

	r := rune(u)
	if utf16.IsSurrogate(r) {
		if len(input) < 12 {
			return "", fmt.Errorf("invalid unicode surrogate pair")
		}
		hex2 := input[8:12]
		u2, err := strconv.ParseUint(hex2, 16, 16)
		if err != nil {
			return "", fmt.Errorf("invalid unicode surrogate pair: %v", err)
		}
		r2 := rune(u2)
		if !utf16.IsSurrogate(r2) {
			return "", fmt.Errorf("invalid unicode surrogate pair")
		}
		combined := utf16.DecodeRune(r, r2)
		return string(combined), nil
	}

	if !utf8.ValidRune(r) {
		return "", fmt.Errorf("invalid unicode sequence")
	}

	return string(r), nil
}

func isValidUnescapedChar(r rune, quote rune) bool {
	if r == quote {
		return false
	}
	return r >= ' ' && r != '\\' && r <= unicode.MaxRune
}
