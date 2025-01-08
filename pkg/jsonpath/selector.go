package jsonpath

import (
	"fmt"
	"strconv"
	"strings"
)

type selectorSubKind int

const (
	selectorSubKindWildcard selectorSubKind = iota
	selectorSubKindName
	selectorSubKindArraySlice
	selectorSubKindArrayIndex
	selectorSubKindFilter
)

type slice struct {
	start *int
	end   *int
	step  *int
}

type selector struct {
	kind   selectorSubKind
	name   string
	index  int
	slice  *slice
	filter *filterSelector
}

func (s selector) ToString() string {
	switch s.kind {
	case selectorSubKindName:
		return "'" + escapeString(s.name) + "'"
	case selectorSubKindArrayIndex:
		// int to string
		return strconv.Itoa(s.index)
	case selectorSubKindFilter:
		return "?" + s.filter.ToString()
	case selectorSubKindWildcard:
		return "*"
	case selectorSubKindArraySlice:
		builder := strings.Builder{}
		if s.slice.start != nil {
			builder.WriteString(strconv.Itoa(*s.slice.start))
		}
		builder.WriteString(":")
		if s.slice.end != nil {
			builder.WriteString(strconv.Itoa(*s.slice.end))
		}

		if s.slice.step != nil {
			builder.WriteString(":")
			builder.WriteString(strconv.Itoa(*s.slice.step))
		}
		return builder.String()
	default:
		panic(fmt.Sprintf("unimplemented selector kind: %v", s.kind))
	}
	return ""
}
