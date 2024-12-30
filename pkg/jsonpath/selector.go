package jsonpath

import (
	"fmt"
	"strconv"
)

type SelectorSubKind int

const (
	SelectorSubKindWildcard SelectorSubKind = iota
	SelectorSubKindName
	SelectorSubKindArraySlice
	SelectorSubKindArrayIndex
	SelectorSubKindFilter
)

type Slice struct {
	Start *int
	End   *int
	Step  *int
}

type Selector struct {
	Kind   SelectorSubKind
	name   string
	index  int
	slice  *Slice
	filter *FilterSelector
}

func (s Selector) ToString() string {
	switch s.Kind {
	case SelectorSubKindName:
		return "\"" + s.name + "\""
	case SelectorSubKindArrayIndex:
		// int to string
		return "[" + strconv.Itoa(s.index) + "]"
	default:
		panic(fmt.Sprintf("unimplemented selector kind: %v", s.Kind))
	}
	return ""
}
