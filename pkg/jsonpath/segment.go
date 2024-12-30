package jsonpath

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type Segment struct {
	Child      *ChildSegment
	Descendant *DescendantSegment
}

type ChildSegmentSubKind int

const (
	ChildSegmentDotWildcard   ChildSegmentSubKind = iota // .*
	ChildSegmentDotMemberName                            // .property
	ChildSegmentLongHand                                 // [ Selector[] ]
)

func (s Segment) ToString() string {
	if s.Child != nil {
		return s.Child.ToString()
	} else if s.Descendant != nil {
		return s.Descendant.ToString()
	} else {
		panic("no segment")
	}

}

type SegmentKind int

type ChildSegment struct {
	kind      ChildSegmentSubKind
	dotName   string
	selectors []*Selector
}

func (s ChildSegment) ToString() string {
	builder := strings.Builder{}
	switch s.kind {
	case ChildSegmentDotWildcard:
		builder.WriteString(".*")
		break
	case ChildSegmentDotMemberName:
		builder.WriteString(".")
		builder.WriteString(s.dotName)
		break
	case ChildSegmentLongHand:
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

type DescendantSegmentSubKind int

const (
	DescendantSegmentSubKindWildcard DescendantSegmentSubKind = iota
	DescendantSegmentSubKindDotName
	DescendantSegmentSubKindLongHand
)

type DescendantSegment struct {
	SubKind      DescendantSegmentSubKind
	innerSegment *Segment
}

func (s DescendantSegment) ToString() string {
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
