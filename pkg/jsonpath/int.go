package jsonpath

import (
	"fmt"
	"strconv"
)

// Integer represents an integer for internet JSON (RFC7493)
// The value must be within the range [-(2^53)+1, (2^53)-1]
type Integer struct {
	value int64
}

const (
	// MAX is the maximum allowed value, 2^53 - 1
	MAX int64 = 9_007_199_254_740_992 - 1
	// MIN is the minimum allowed value (-2^53) + 1
	MIN int64 = -9_007_199_254_740_992 + 1
)

// IntegerError represents errors that can occur when working with Integer
type IntegerError struct {
	msg string
}

func (e IntegerError) Error() string {
	return e.msg
}

var (
	ErrOutOfBounds = IntegerError{"integer outside valid range"}
)

// checkIsValid checks if an int64 is within the valid range
func checkIsValid(v int64) bool {
	return v >= MIN && v <= MAX
}

// ZERO represents an Integer with value 0
var ZERO = Integer{value: 0}

// NewInteger creates a new Integer from an int64
func NewInteger(value int64) (Integer, error) {
	if !checkIsValid(value) {
		return Integer{}, ErrOutOfBounds
	}
	return Integer{value: value}, nil
}

// MustNewInteger creates a new Integer from an int64, panicking if invalid
func MustNewInteger(value int64) Integer {
	i, err := NewInteger(value)
	if err != nil {
		panic("value is out of the valid range")
	}
	return i
}

// Abs returns the absolute value of the Integer
func (i Integer) Abs() Integer {
	if i.value < 0 {
		return Integer{value: -i.value}
	}
	return i
}

// CheckedAdd adds two Integers, returning nil if the result would be invalid
func (i Integer) CheckedAdd(rhs Integer) *Integer {
	sum := i.value + rhs.value
	if !checkIsValid(sum) {
		return nil
	}
	return &Integer{value: sum}
}

// CheckedSub subtracts rhs from i, returning nil if the result would be invalid
func (i Integer) CheckedSub(rhs Integer) *Integer {
	diff := i.value - rhs.value
	if !checkIsValid(diff) {
		return nil
	}
	return &Integer{value: diff}
}

// CheckedMul multiplies two Integers, returning nil if the result would be invalid
func (i Integer) CheckedMul(rhs Integer) *Integer {
	prod := i.value * rhs.value
	if !checkIsValid(prod) {
		return nil
	}
	return &Integer{value: prod}
}

// ParseInteger parses a string into an Integer
func ParseInteger(s string) (Integer, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return Integer{}, err
	}
	return NewInteger(v)
}

func (i Integer) String() string {
	return fmt.Sprintf("%d", i.value)
}

// Int64 returns the underlying int64 value
func (i Integer) Int64() int64 {
	return i.value
}
