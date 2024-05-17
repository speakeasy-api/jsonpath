package jsonpath

import (
	"fmt"
	"strings"
)

const (
	EdgeItem = "├── "
	EdgeLast = "└── "
)

func (n *RootNode) String() string {
	return fmt.Sprintf("$")
}

func (n *CurrentNode) String() string {
	return fmt.Sprintf("@")
}

func (n *IdentifierNode) String() string {
	return fmt.Sprintf("%s", n.Name.Literal)
}

func (n *WildcardNode) String() string {
	return fmt.Sprintf("*")
}

func (n *RecursiveDescentNode) String() string {
	return fmt.Sprintf("..")
}

func (n *SubscriptNode) String() string {
	return fmt.Sprintf("[%s]", n.Index.String())
}

func (n *SliceNode) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	if n.Start != nil {
		sb.WriteString(n.Start.String())
	}
	sb.WriteString(":")
	if n.Finish != nil {
		sb.WriteString(n.Finish.String())
	}
	if n.Step != nil {
		sb.WriteString(":")
		sb.WriteString(n.Step.String())
	}
	sb.WriteString("]")
	return sb.String()
}

func (n *UnionNode) String() string {
	return fmt.Sprintf("%s, %s", n.Lhs.String(), n.Rhs.String())
}

func (n *FilterNode) String() string {
	return fmt.Sprintf("[?%s]", n.Expr.String())
}

func (n *ComparisonNode) String() string {
	return fmt.Sprintf("%s %s %s", n.Lhs.String(), n.Operator.Literal, n.Rhs.String())
}

func (n *BooleanNode) String() string {
	return fmt.Sprintf("%s", n.Value.Literal)
}

func (n *NumberNode) String() string {
	return fmt.Sprintf("%s", n.Value.Literal)
}

func (n *StringNode) String() string {
	return fmt.Sprintf("'%s'", n.Value.Literal)
}

func (n *NullNode) String() string {
	return fmt.Sprintf("null")
}

func (n *FunctionCallNode) String() string {
	var args []string
	for _, arg := range n.Args {
		args = append(args, arg.String())
	}
	return fmt.Sprintf("%s(%s)", n.Name.Literal, strings.Join(args, ", "))
}

func PrintNode(node Node) string {
	var sb strings.Builder
	printNode(&sb, "", node)
	return sb.String()
}

func printNode(sb *strings.Builder, prefix string, node Node) {
	sb.WriteString(prefix)
	sb.WriteString(node.String())
	sb.WriteString("\n")

	switch n := node.(type) {
	case *SubscriptNode:
		printNode(sb, prefix+EdgeItem, n.Index)
	case *SliceNode:
		if n.Start != nil {
			printNode(sb, prefix+EdgeItem, n.Start)
		}
		if n.Finish != nil {
			printNode(sb, prefix+EdgeItem, n.Finish)
		}
		if n.Step != nil {
			printNode(sb, prefix+EdgeLast, n.Step)
		}
	case *UnionNode:
		printNode(sb, prefix+EdgeItem, n.Lhs)
		printNode(sb, prefix+EdgeLast, n.Rhs)
	case *FilterNode:
		printNode(sb, prefix+EdgeItem, n.Expr)
	case *ComparisonNode:
		printNode(sb, prefix+EdgeItem, n.Lhs)
		printNode(sb, prefix+EdgeLast, n.Rhs)
	case *FunctionCallNode:
		for i, arg := range n.Args {
			if i == len(n.Args)-1 {
				printNode(sb, prefix+EdgeLast, arg)
			} else {
				printNode(sb, prefix+EdgeItem, arg)
			}
		}
	}
}
