package jsonpath

import (
	"fmt"
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
