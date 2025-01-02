package jsonpath

import (
	"gopkg.in/yaml.v3"
)

type Evaluator interface {
	Query(current *yaml.Node, root *yaml.Node) []*yaml.Node
}

// jsonPathAST can be Evaluated
var _ Evaluator = jsonPathAST{}

func (q jsonPathAST) Query(current *yaml.Node, root *yaml.Node) []*yaml.Node {
	result := make([]*yaml.Node, 0)
	// If the top level node is a documentnode, unwrap it
	if root.Kind == yaml.DocumentNode && len(root.Content) == 1 {
		root = root.Content[0]
	}
	result = append(result, root)

	for _, segment := range q.segments {
		newValue := []*yaml.Node{}
		for _, value := range result {
			newValue = append(newValue, segment.Query(value, root)...)
		}
		result = newValue
	}
	return result
}

func (s segment) Query(value *yaml.Node, root *yaml.Node) []*yaml.Node {
	if s.Child != nil {
		return s.Child.Query(value, root)
	} else if s.Descendant != nil {
		// run the inner segment against this node
		var result = []*yaml.Node{}
		children := descend(value, root)
		for _, child := range children {
			result = append(result, s.Descendant.Query(child, root)...)
		}
		// make children unique by pointer value
		result = unique(result)
		return result
	} else {
		panic("no segment type")
	}
}

func unique(nodes []*yaml.Node) []*yaml.Node {
	// stably returns a new slice containing only the unique elements from nodes
	res := make([]*yaml.Node, 0)
	seen := make(map[*yaml.Node]bool)
	for _, node := range nodes {
		if _, ok := seen[node]; !ok {
			res = append(res, node)
			seen[node] = true
		}
	}
	return res
}

func (s innerSegment) Query(value *yaml.Node, root *yaml.Node) []*yaml.Node {
	result := []*yaml.Node{}

	switch s.kind {
	case segmentDotWildcard:
		// Handle wildcard - get all children
		switch value.Kind {
		case yaml.MappingNode:
			// in a mapping node, keys and values alternate
			// we just want to return the values
			for i, child := range value.Content {
				if i%2 == 1 {
					result = append(result, child)
				}
			}
		case yaml.SequenceNode:
			for _, child := range value.Content {
				result = append(result, child)
			}
		}
		return result
	case segmentDotMemberName:
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

	case segmentLongHand:
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
		if s.index >= len(value.Content) || s.index < 0 {
			return nil
		}
		return []*yaml.Node{value.Content[s.index]}
	case SelectorSubKindWildcard:
		if value.Kind == yaml.SequenceNode {
			return value.Content
		} else if value.Kind == yaml.MappingNode {
			var result []*yaml.Node
			for i, child := range value.Content {
				if i%2 == 1 {
					result = append(result, child)
				}
			}
			return result
		}
		return nil
	case SelectorSubKindArraySlice:
		if value.Kind != yaml.SequenceNode {
			return nil
		}
		start, end, step := 0, len(value.Content), 1
		if s.slice.Start != nil {
			start = *s.slice.Start
		} else {
			start = 0
		}
		if s.slice.End != nil {
			end = *s.slice.End
		} else {
			end = len(value.Content)
		}
		if s.slice.Step != nil {
			step = *s.slice.Step
		} else {
			step = 1
		}
		if step == 0 {
			return nil
		}
		var result []*yaml.Node
		if step < 0 {
			for i := end - 1; i >= start && i >= 0 && i < len(value.Content); i += step {
				result = append(result, value.Content[i])
			}
		} else {
			for i := start; i < end && i >= 0 && i < len(value.Content); i += step {
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

func (s filterSelector) Matches(node *yaml.Node, root *yaml.Node) bool {
	return s.expression.Matches(node, root)
}

func (e logicalOrExpr) Matches(node *yaml.Node, root *yaml.Node) bool {
	for _, expr := range e.expressions {
		if expr.Matches(node, root) {
			return true
		}
	}
	return false
}

func (e logicalAndExpr) Matches(node *yaml.Node, root *yaml.Node) bool {
	for _, expr := range e.expressions {
		if !expr.Matches(node, root) {
			return false
		}
	}
	return true
}

func (e basicExpr) Matches(node *yaml.Node, root *yaml.Node) bool {
	if e.parenExpr != nil {
		result := e.parenExpr.expr.Matches(node, root)
		if e.parenExpr.not {
			return !result
		}
		return result
	} else if e.comparisonExpr != nil {
		return e.comparisonExpr.Matches(node, root)
	} else if e.testExpr != nil {
		return e.testExpr.Matches(node, root)
	}
	return false
}

func (e comparisonExpr) Matches(node *yaml.Node, root *yaml.Node) bool {
	leftValue := e.left.Evaluate(node, root)
	rightValue := e.right.Evaluate(node, root)

	switch e.op {
	case equalTo:
		return leftValue.Equals(rightValue)
	case notEqualTo:
		return !leftValue.Equals(rightValue)
	case lessThan:
		return leftValue.LessThan(rightValue)
	case lessThanEqualTo:
		return leftValue.LessThanOrEqual(rightValue)
	case greaterThan:
		return rightValue.LessThan(leftValue)
	case greaterThanEqualTo:
		return rightValue.LessThanOrEqual(leftValue)
	default:
		return false
	}
}

func (e testExpr) Matches(node *yaml.Node, root *yaml.Node) bool {
	var result bool
	if e.filterQuery != nil {
		result = len(e.filterQuery.Query(node, root)) > 0
	} else if e.functionExpr != nil {
		funcResult := e.functionExpr.Evaluate(node, root)
		if funcResult.bool != nil {
			result = *funcResult.bool
		} else if funcResult.null == nil {
			result = true
		}
	}
	if e.not {
		return !result
	}
	return result
}

func (q filterQuery) Query(node *yaml.Node, root *yaml.Node) []*yaml.Node {
	if q.relQuery != nil {
		return q.relQuery.Query(node, root)
	}
	if q.jsonPathQuery != nil {
		return q.jsonPathQuery.Query(node, root)
	}
	return nil
}

func (q relQuery) Query(node *yaml.Node, root *yaml.Node) []*yaml.Node {
	result := []*yaml.Node{node}
	for _, seg := range q.segments {
		var newResult []*yaml.Node
		for _, value := range result {
			newResult = append(newResult, seg.Query(value, root)...)
		}
		result = newResult
	}
	return result
}

func (q absQuery) Query(node *yaml.Node, root *yaml.Node) []*yaml.Node {
	result := []*yaml.Node{root}
	for _, seg := range q.segments {
		var newResult []*yaml.Node
		for _, value := range result {
			newResult = append(newResult, seg.Query(value, root)...)
		}
		result = newResult
	}
	return result
}
