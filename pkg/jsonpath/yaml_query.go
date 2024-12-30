package jsonpath

import (
	"gopkg.in/yaml.v3"
)

type Evaluator interface {
	Query(current *yaml.Node, root *yaml.Node) []*yaml.Node
}

// JsonPathQuery can be Evaluated
var _ Evaluator = JsonPathQuery{}

func (q JsonPathQuery) Query(current *yaml.Node, root *yaml.Node) []*yaml.Node {
	var result []*yaml.Node
	// If the top level node is a documentnode, unwrap it
	if root.Kind == yaml.DocumentNode && len(root.Content) == 1 {
		root = root.Content[0]
	}
	result = append(result, root)

	for _, segment := range q.Segments {
		newValue := []*yaml.Node{}
		for _, value := range result {
			newValue = append(newValue, segment.Query(value, root)...)
		}
		result = newValue
	}
	return result
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

func (s DescendantSegment) Query(value *yaml.Node, root *yaml.Node) []*yaml.Node {
	// run the inner segment against this node
	result := s.innerSegment.Query(value, root)
	children := descend(value, root)
	for _, child := range children {
		result = append(result, s.innerSegment.Query(child, root)...)
	}
	return result
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
	case SelectorSubKindWildcard:
		if value.Kind == yaml.MappingNode || value.Kind == yaml.SequenceNode {
			return value.Content
		}
		return nil
	case SelectorSubKindArraySlice:
		if value.Kind != yaml.SequenceNode {
			return nil
		}
		start, end, step := 0, len(value.Content), 1
		if s.slice.Start != nil {
			start = *s.slice.Start
		}
		if s.slice.End != nil {
			end = *s.slice.End
		}
		if s.slice.Step != nil {
			step = *s.slice.Step
		}
		var result []*yaml.Node
		for i := start; i < end; i += step {
			if i >= 0 && i < len(value.Content) {
				result = append(result, value.Content[i])
			}
		}
		return result
	case SelectorSubKindFilter:
		var result []*yaml.Node
		switch value.Kind {
		case yaml.MappingNode:
			for i := 1; i < len(value.Content); i += 2 {
				if s.filter.Matches(value.Content[i], root) {
					result = append(result, value.Content[i])
				}
			}
		case yaml.SequenceNode:
			for _, child := range value.Content {
				if s.filter.Matches(child, root) {
					result = append(result, child)
				}
			}
		}
		return result
	}
	return nil
}

func (s FilterSelector) Matches(node *yaml.Node, root *yaml.Node) bool {
	return s.Expression.Matches(node, root)
}

func (e LogicalOrExpr) Matches(node *yaml.Node, root *yaml.Node) bool {
	for _, expr := range e.Expressions {
		if expr.Matches(node, root) {
			return true
		}
	}
	return false
}

func (e LogicalAndExpr) Matches(node *yaml.Node, root *yaml.Node) bool {
	for _, expr := range e.Expressions {
		if !expr.Matches(node, root) {
			return false
		}
	}
	return true
}

func (e BasicExpr) Matches(node *yaml.Node, root *yaml.Node) bool {
	if e.ParenExpr != nil {
		result := e.ParenExpr.Expr.Matches(node, root)
		if e.ParenExpr.Not {
			return !result
		}
		return result
	} else if e.ComparisonExpr != nil {
		return e.ComparisonExpr.Matches(node, root)
	} else if e.TestExpr != nil {
		return e.TestExpr.Matches(node, root)
	}
	return false
}

func (e ComparisonExpr) Matches(node *yaml.Node, root *yaml.Node) bool {
	leftValue := e.Left.Evaluate(node, root)
	rightValue := e.Right.Evaluate(node, root)

	switch e.Op {
	case EqualTo:
		return leftValue.Equals(rightValue)
	case NotEqualTo:
		return !leftValue.Equals(rightValue)
	case LessThan:
		return leftValue.LessThan(rightValue)
	case LessThanEqualTo:
		return leftValue.LessThanOrEqual(rightValue)
	case GreaterThan:
		return rightValue.LessThan(leftValue)
	case GreaterThanEqualTo:
		return rightValue.LessThanOrEqual(leftValue)
	default:
		return false
	}
}

func (e TestExpr) Matches(node *yaml.Node, root *yaml.Node) bool {
	var result bool
	if e.FilterQuery != nil {
		result = len(e.FilterQuery.Query(node, root)) > 0
	} else if e.FunctionExpr != nil {
		funcResult := e.FunctionExpr.Evaluate(node, root)
		if funcResult.Bool != nil {
			result = *funcResult.Bool
		} else if funcResult.Null == nil {
			result = true
		}
	}
	if e.Not {
		return !result
	}
	return result
}

func (q FilterQuery) Query(node *yaml.Node, root *yaml.Node) []*yaml.Node {
	if q.RelQuery != nil {
		return q.RelQuery.Query(node, root)
	}
	if q.JsonPathQuery != nil {
		return q.JsonPathQuery.Query(node, root)
	}
	return nil
}

func (q RelQuery) Query(node *yaml.Node, root *yaml.Node) []*yaml.Node {
	result := []*yaml.Node{node}
	for _, segment := range q.Segments {
		var newResult []*yaml.Node
		for _, value := range result {
			newResult = append(newResult, segment.Query(value, root)...)
		}
		result = newResult
	}
	return result
}

func (q AbsQuery) Query(node *yaml.Node, root *yaml.Node) []*yaml.Node {
	result := []*yaml.Node{root}
	for _, segment := range q.Segments {
		var newResult []*yaml.Node
		for _, value := range result {
			newResult = append(newResult, segment.Query(value, root)...)
		}
		result = newResult
	}
	return result
}
