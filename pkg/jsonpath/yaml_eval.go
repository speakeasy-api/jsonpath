package jsonpath

import (
	"gopkg.in/yaml.v3"
	"regexp"
	"strconv"
)

func (l literal) Equals(value literal) bool {
	if l.integer != nil && value.integer != nil {
		return *l.integer == *value.integer
	}
	if l.float64 != nil && value.float64 != nil {
		return *l.float64 == *value.float64
	}
	if l.integer != nil && value.float64 != nil {
		return float64(*l.integer) == *value.float64
	}
	if l.float64 != nil && value.integer != nil {
		return *l.float64 == float64(*value.integer)
	}
	if l.string != nil && value.string != nil {
		return *l.string == *value.string
	}
	if l.bool != nil && value.bool != nil {
		return *l.bool == *value.bool
	}
	if l.null != nil && value.null != nil {
		return *l.null == *value.null
	}
	if l.node != nil && value.node != nil {
		return equalsNode(l.node, value.node)
	}
	return false
}

func equalsNode(a *yaml.Node, b *yaml.Node) bool {
	// decode into interfaces, then compare
	if a.Tag != b.Tag {
		return false
	}
	switch a.Tag {
	case "!!str":
		return a.Value == b.Value
	case "!!int":
		return a.Value == b.Value
	case "!!float":
		return a.Value == b.Value
	case "!!bool":
		return a.Value == b.Value
	case "!!null":
		return a.Value == b.Value
	case "!!seq":
		for i := 0; i < len(a.Content); i++ {
			if !equalsNode(a.Content[i], b.Content[i]) {
				return false
			}
		}
	case "!!map":
		for i := 0; i < len(a.Content); i += 2 {
			if !equalsNode(a.Content[i], b.Content[i]) {
				return false
			}
			if !equalsNode(a.Content[i+1], b.Content[i+1]) {
				return false
			}
		}
	}
	return true
}

func (l literal) LessThan(value literal) bool {
	if l.integer != nil && value.integer != nil {
		return *l.integer < *value.integer
	}
	if l.float64 != nil && value.float64 != nil {
		return *l.float64 < *value.float64
	}
	if l.integer != nil && value.float64 != nil {
		return float64(*l.integer) < *value.float64
	}
	if l.float64 != nil && value.integer != nil {
		return *l.float64 < float64(*value.integer)
	}
	return false
}

func (l literal) LessThanOrEqual(value literal) bool {
	if l.integer != nil && value.integer != nil {
		return *l.integer <= *value.integer
	}
	if l.float64 != nil && value.float64 != nil {
		return *l.float64 <= *value.float64
	}
	return false
}

func (c comparable) Evaluate(node *yaml.Node, root *yaml.Node) literal {
	if c.literal != nil {
		return *c.literal
	}
	if c.singularQuery != nil {
		return c.singularQuery.Evaluate(node, root)
	}
	if c.functionExpr != nil {
		return c.functionExpr.Evaluate(node, root)
	}
	return literal{}
}

func (e functionExpr) length(node *yaml.Node, root *yaml.Node) literal {
	switch node.Kind {
	case yaml.SequenceNode:
		res := len(node.Content)
		return literal{integer: &res}
	case yaml.MappingNode:
		res := len(node.Content) / 2
		return literal{integer: &res}
	case yaml.ScalarNode:
		res := len(node.Value)
		return literal{integer: &res}
	default:
		return literal{}
	}
}

func (e functionExpr) count(node *yaml.Node, root *yaml.Node) literal {
	args := e.args[0].filterQuery.Query(node, root)
	//
	res := len(args)
	return literal{integer: &res}
}

func (e functionExpr) match(node *yaml.Node, root *yaml.Node) literal {
	if node.Kind != yaml.ScalarNode {
		return literal{bool: &[]bool{false}[0]}
	}
	arg1 := e.args[0].Evaluate(node, root)
	arg2 := e.args[1].Evaluate(node, root)
	if arg1.string == nil || arg2.string == nil {
		return literal{bool: &[]bool{false}[0]}
	}
	matched, _ := regexp.MatchString(*arg2.string, *arg1.string)
	return literal{bool: &matched}
}

func (e functionExpr) search(node *yaml.Node, root *yaml.Node) literal {
	if node.Kind != yaml.ScalarNode {
		return literal{bool: &[]bool{false}[0]}
	}
	arg1 := e.args[0].Evaluate(node, root)
	arg2 := e.args[1].Evaluate(node, root)
	if arg1.string == nil || arg2.string == nil {
		return literal{bool: &[]bool{false}[0]}
	}
	matched, _ := regexp.MatchString(*arg2.string, *arg1.string)
	return literal{bool: &matched}
}

func (e functionExpr) value(node *yaml.Node, root *yaml.Node) literal {
	args := e.args[0].filterQuery.Query(node, root)
	if len(args) == 1 {
		return nodeToLiteral(args[0])
	}
	return literal{}
}

func nodeToLiteral(node *yaml.Node) literal {
	switch node.Tag {
	case "!!str":
		return literal{string: &node.Value}
	case "!!int":
		i, _ := strconv.Atoi(node.Value)
		return literal{integer: &i}
	case "!!float":
		f, _ := strconv.ParseFloat(node.Value, 64)
		return literal{float64: &f}
	case "!!bool":
		b, _ := strconv.ParseBool(node.Value)
		return literal{bool: &b}
	case "!!null":
		b := true
		return literal{null: &b}
	default:
		return literal{node: node}
	}
}

func (e functionExpr) Evaluate(node *yaml.Node, root *yaml.Node) literal {
	switch e.funcType {
	case functionTypeLength:
		return e.length(node, root)
	case functionTypeCount:
		return e.count(node, root)
	case functionTypeMatch:
		return e.match(node, root)
	case functionTypeSearch:
		return e.search(node, root)
	case functionTypeValue:
		return e.value(node, root)
	}
	return literal{}
}

func (q singularQuery) Evaluate(node *yaml.Node, root *yaml.Node) literal {
	if q.relQuery != nil {
		return q.relQuery.Evaluate(node, root)
	}
	if q.absQuery != nil {
		return q.absQuery.Evaluate(node, root)
	}
	return literal{}
}

func (q relQuery) Evaluate(node *yaml.Node, root *yaml.Node) literal {
	result := q.Query(node, root)
	if len(result) == 1 {
		return nodeToLiteral(result[0])
	}
	return literal{}

}

func (a functionArgument) Evaluate(node *yaml.Node, root *yaml.Node) literal {
	if a.literal != nil {
		return *a.literal
	}
	if a.filterQuery != nil {
		result := a.filterQuery.Query(node, root)
		if len(result) == 1 {
			return nodeToLiteral(result[0])
		}
		return literal{}
	}
	if a.logicalExpr != nil {
		if a.logicalExpr.Matches(node, root) {
			return literal{bool: &[]bool{true}[0]}
		}
		return literal{bool: &[]bool{false}[0]}
	}
	if a.functionExpr != nil {
		return a.functionExpr.Evaluate(node, root)
	}
	return literal{}
}

func (q absQuery) Evaluate(node *yaml.Node, root *yaml.Node) literal {
	result := q.Query(root, root)
	if len(result) == 1 {
		return nodeToLiteral(result[0])
	}
	return literal{}
}
