package jsonpath

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type segment struct {
	Child      *innerSegment
	Descendant *innerSegment
}

type segmentSubKind int

const (
	segmentDotWildcard   segmentSubKind = iota // .*
	segmentDotMemberName                       // .property
	segmentLongHand                            // [ selector[] ]
)

func (s segment) ToString() string {
	if s.Child != nil {
		if s.Child.kind != segmentLongHand {
			return "." + s.Child.ToString()
		} else {
			return s.Child.ToString()
		}
	} else if s.Descendant != nil {
		return ".." + s.Descendant.ToString()
	} else {
		panic("no segment")
	}
}

type innerSegment struct {
	kind      segmentSubKind
	dotName   string
	selectors []*selector
}

func (s innerSegment) ToString() string {
	builder := strings.Builder{}
	switch s.kind {
	case segmentDotWildcard:
		builder.WriteString("*")
		break
	case segmentDotMemberName:
		builder.WriteString(s.dotName)
		break
	case segmentLongHand:
		builder.WriteString("[")
		for i, selector := range s.selectors {
			builder.WriteString(selector.ToString())
			if i < len(s.selectors)-1 {
				builder.WriteString(", ")
			}
		}
		builder.WriteString("]")
		break
	default:
		panic("unknown child segment kind")
	}
	return builder.String()
}

func descend(value *yaml.Node, root *yaml.Node) []*yaml.Node {
	result := []*yaml.Node{value}
	for _, child := range value.Content {
		result = append(result, descend(child, root)...)
	}
	return result
}
