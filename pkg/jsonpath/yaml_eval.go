package jsonpath

import (
	"gopkg.in/yaml.v3"
	"regexp"
	"strconv"
)

func (l Literal) Equals(value Literal) bool {
	if l.Integer != nil && value.Integer != nil {
		return *l.Integer == *value.Integer
	}
	if l.Float64 != nil && value.Float64 != nil {
		return *l.Float64 == *value.Float64
	}
	if l.Integer != nil && value.Float64 != nil {
		return float64(*l.Integer) == *value.Float64
	}
	if l.Float64 != nil && value.Integer != nil {
		return *l.Float64 == float64(*value.Integer)
	}
	if l.String != nil && value.String != nil {
		return *l.String == *value.String
	}
	if l.Bool != nil && value.Bool != nil {
		return *l.Bool == *value.Bool
	}
	if l.Null != nil && value.Null != nil {
		return *l.Null == *value.Null
	}
	return false
}

func (l Literal) LessThan(value Literal) bool {
	if l.Integer != nil && value.Integer != nil {
		return *l.Integer < *value.Integer
	}
	if l.Float64 != nil && value.Float64 != nil {
		return *l.Float64 < *value.Float64
	}
	if l.Integer != nil && value.Float64 != nil {
		return float64(*l.Integer) < *value.Float64
	}
	if l.Float64 != nil && value.Integer != nil {
		return *l.Float64 < float64(*value.Integer)
	}
	return false
}

func (l Literal) LessThanOrEqual(value Literal) bool {
	if l.Integer != nil && value.Integer != nil {
		return *l.Integer <= *value.Integer
	}
	if l.Float64 != nil && value.Float64 != nil {
		return *l.Float64 <= *value.Float64
	}
	return false
}

func (c Comparable) Evaluate(node *yaml.Node, root *yaml.Node) Literal {
	if c.Literal != nil {
		return *c.Literal
	}
	if c.SingularQuery != nil {
		return c.SingularQuery.Evaluate(node, root)
	}
	if c.FunctionExpr != nil {
		return c.FunctionExpr.Evaluate(node, root)
	}
	return Literal{Null: &[]bool{true}[0]}
}

func (e FunctionExpr) length(node *yaml.Node, root *yaml.Node) Literal {
	switch node.Kind {
	case yaml.SequenceNode:
		res := len(node.Content)
		return Literal{Integer: &res}
	case yaml.MappingNode:
		res := len(node.Content) / 2
		return Literal{Integer: &res}
	case yaml.ScalarNode:
		res := len(node.Value)
		return Literal{Integer: &res}
	default:
		return Literal{Null: &[]bool{true}[0]}
	}
}

func (e FunctionExpr) count(node *yaml.Node, root *yaml.Node) Literal {
	args := e.Args[0].FilterQuery.Query(node, root)
	//
	res := len(args)
	return Literal{Integer: &res}
}

func (e FunctionExpr) match(node *yaml.Node, root *yaml.Node) Literal {
	if node.Kind != yaml.ScalarNode {
		return Literal{Bool: &[]bool{false}[0]}
	}
	arg1 := e.Args[0].Evaluate(node, root)
	arg2 := e.Args[1].Evaluate(node, root)
	if arg1.String == nil || arg2.String == nil {
		return Literal{Bool: &[]bool{false}[0]}
	}
	matched, _ := regexp.MatchString(*arg2.String, *arg1.String)
	return Literal{Bool: &matched}
}

func (e FunctionExpr) search(node *yaml.Node, root *yaml.Node) Literal {
	if node.Kind != yaml.ScalarNode {
		return Literal{Bool: &[]bool{false}[0]}
	}
	arg1 := e.Args[0].Evaluate(node, root)
	arg2 := e.Args[1].Evaluate(node, root)
	if arg1.String == nil || arg2.String == nil {
		return Literal{Bool: &[]bool{false}[0]}
	}
	matched, _ := regexp.MatchString(*arg2.String, *arg1.String)
	return Literal{Bool: &matched}
}

func (e FunctionExpr) value(node *yaml.Node, root *yaml.Node) Literal {
	args := e.Args[0].FilterQuery.Query(node, root)
	if len(args) == 1 {
		return nodeToLiteral(args[0])
	}
	return Literal{Null: &[]bool{true}[0]}
}

func nodeToLiteral(node *yaml.Node) Literal {
	switch node.Tag {
	case "!!str":
		return Literal{String: &node.Value}
	case "!!int":
		i, _ := strconv.Atoi(node.Value)
		return Literal{Integer: &i}
	case "!!float":
		f, _ := strconv.ParseFloat(node.Value, 64)
		return Literal{Float64: &f}
	case "!!bool":
		b, _ := strconv.ParseBool(node.Value)
		return Literal{Bool: &b}
	case "!!null":
		b := true
		return Literal{Null: &b}
	default:
		return Literal{}
	}
}

func (e FunctionExpr) Evaluate(node *yaml.Node, root *yaml.Node) Literal {
	switch e.Type {
	case FunctionTypeLength:
		return e.length(node, root)
	case FunctionTypeCount:
		return e.count(node, root)
	case FunctionTypeMatch:
		return e.match(node, root)
	case FunctionTypeSearch:
		return e.search(node, root)
	case FunctionTypeValue:
		return e.value(node, root)
	}
	return Literal{Null: &[]bool{true}[0]}
}

func (q SingularQuery) Evaluate(node *yaml.Node, root *yaml.Node) Literal {
	if q.RelQuery != nil {
		return q.RelQuery.Evaluate(node, root)
	}
	if q.AbsQuery != nil {
		return q.AbsQuery.Evaluate(node, root)
	}
	return Literal{Null: &[]bool{true}[0]}
}

func (q RelQuery) Evaluate(node *yaml.Node, root *yaml.Node) Literal {
	result := q.Query(node, root)
	if len(result) == 1 {
		return nodeToLiteral(result[0])
	}
	return Literal{Null: &[]bool{true}[0]}

}

func (a FunctionArgument) Evaluate(node *yaml.Node, root *yaml.Node) Literal {
	if a.Literal != nil {
		return *a.Literal
	}
	if a.FilterQuery != nil {
		result := a.FilterQuery.Query(node, root)
		if len(result) == 1 {
			return nodeToLiteral(result[0])
		}
		return Literal{Null: &[]bool{true}[0]}
	}
	if a.LogicalExpr != nil {
		if a.LogicalExpr.Matches(node, root) {
			return Literal{Bool: &[]bool{true}[0]}
		}
		return Literal{Bool: &[]bool{false}[0]}
	}
	if a.FunctionExpr != nil {
		return a.FunctionExpr.Evaluate(node, root)
	}
	return Literal{Null: &[]bool{true}[0]}
}

func (q AbsQuery) Evaluate(node *yaml.Node, root *yaml.Node) Literal {
	result := q.Query(root, root)
	if len(result) == 1 {
		return nodeToLiteral(result[0])
	}
	return Literal{Null: &[]bool{true}[0]}
}
