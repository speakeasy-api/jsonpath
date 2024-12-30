package jsonpath

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type segment struct {
	Child      *childSegment
	Descendant *descendantSegment
}

type childSegmentSubKind int

const (
	childSegmentDotWildcard   childSegmentSubKind = iota // .*
	childSegmentDotMemberName                            // .property
	childSegmentLongHand                                 // [ Selector[] ]
)

func (s segment) ToString() string {
	if s.Child != nil {
		return s.Child.ToString()
	} else if s.Descendant != nil {
		return s.Descendant.ToString()
	} else {
		panic("no segment")
	}

}

type segmentKind int

type childSegment struct {
	kind      childSegmentSubKind
	dotName   string
	selectors []*Selector
}

func (s childSegment) ToString() string {
	builder := strings.Builder{}
	switch s.kind {
	case childSegmentDotWildcard:
		builder.WriteString(".*")
		break
	case childSegmentDotMemberName:
		builder.WriteString(".")
		builder.WriteString(s.dotName)
		break
	case childSegmentLongHand:
		builder.WriteString("[")
		for i, selector := range s.selectors {
			builder.WriteString(selector.ToString())
			if i < len(s.selectors)-1 {
				builder.WriteString(",")
			}
		}
		builder.WriteString("]")
		break
	default:
		panic("unknown child segment kind")
	}
	return builder.String()
}

type descendantSegmentSubKind int

const (
	descendantSegmentSubKindWildcard descendantSegmentSubKind = iota
	descendantSegmentSubKindDotName
	descendantSegmentSubKindLongHand
)

type descendantSegment struct {
	subKind      descendantSegmentSubKind
	innerSegment *segment
}

func (s descendantSegment) ToString() string {
	builder := strings.Builder{}
	builder.WriteString("..")
	builder.WriteString(s.innerSegment.ToString())
	return builder.String()
}

func descend(value *yaml.Node, root *yaml.Node) []*yaml.Node {
	result := []*yaml.Node{value}
	for _, child := range value.Content {
		result = append(result, descend(child, root)...)
	}
	return result
}
