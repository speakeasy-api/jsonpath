package selector

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"strconv"
)

type SelectorSubKind int

const (
	SelectorSubKindWildcard SelectorSubKind = iota
	SelectorSubKindName
	SelectorSubKindArraySlice
	SelectorSubKindArrayIndex
)

type Selector struct {
	Kind  SelectorSubKind
	name  string
	index int
	slice *Slice
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

func (s Selector) Query(value *yaml.Node, root *yaml.Node) []*yaml.Node {
	switch s.Kind {
	case SelectorSubKindName:
		if value.Kind != yaml.MappingNode {
			return nil
		}
		// MappingNode children is a list of alternating keys and values
		var key string
		for i, child := range value.Content {
			if i%2 == 0 {
				key = child.Value
				continue
			}
			if key == s.name {
				return []*yaml.Node{child}
			}
		}
	case SelectorSubKindArrayIndex:
		if value.Kind != yaml.SequenceNode {
			return nil
		}
		if s.index >= len(value.Content) {
			return nil
		}
		return []*yaml.Node{value.Content[s.index]}
	}
	return nil
}
