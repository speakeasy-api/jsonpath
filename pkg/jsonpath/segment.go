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

func (s Segment) Query(value *yaml.Node, root *yaml.Node) []*yaml.Node {
	if s.Child != nil {
		return s.Child.Query(value, root)
	} else if s.Descendant != nil {
		return s.Descendant.Query(value, root)
	} else {
		panic("no segment type")
	}
}

type SegmentKind int

type ChildSegment struct {
	kind      ChildSegmentSubKind
	dotName   string
	selectors []Selector
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

func (s ChildSegment) Query(value *yaml.Node, root *yaml.Node) []*yaml.Node {
	result := []*yaml.Node{}

	switch s.kind {
	case ChildSegmentDotWildcard:
		// Handle wildcard - get all children
		for _, child := range value.Content {
			result = append(result, child)
		}

	case ChildSegmentDotMemberName:
		// Handle member access
		if value.Kind == yaml.MappingNode {
			// In YAML mapping nodes, keys and values alternate

			for i := 0; i < len(value.Content); i += 2 {
				key := value.Content[i]
				val := value.Content[i+1]

				if key.Value == s.dotName {
					result = append(result, val)
					break
				}
			}
		}

	case ChildSegmentLongHand:
		// Handle long hand selectors
		for _, selector := range s.selectors {
			result = append(result, selector.Query(value, root)...)
		}
	default:
		panic("unknown child segment kind")
	}

	return result

}

type DescendantSegment struct {
	innerSegment Segment
}

func (s DescendantSegment) ToString() string {
	builder := strings.Builder{}
	builder.WriteString("..")
	builder.WriteString(s.innerSegment.ToString())
	return builder.String()
}

func (s DescendantSegment) Query(value *yaml.Node, root *yaml.Node) []*yaml.Node {
	// run the inner segment against this node
	result := s.innerSegment.Query(value, root)
	children := descend(value, root)
	for _, child := range children {
		result = append(result, s.innerSegment.Query(child, root)...)
	}
	return result
}

func descend(value *yaml.Node, root *yaml.Node) []*yaml.Node {
	result := []*yaml.Node{value}
	for _, child := range value.Content {
		result = append(result, descend(child, root)...)
	}
	return result
}
