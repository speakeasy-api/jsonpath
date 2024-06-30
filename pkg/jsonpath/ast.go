package jsonpath

type Kind int

const (
	ChildSegmentKind Kind = iota
	DescendantSegmentKind
)

// Node represents a node in the JSONPath AST.
type Node interface {
	Kind() Kind // The Kind of the node
}

type Segment = Node

// JSONPath expression built from segments that have
// been syntactically restricted in a certain way (Section 2.3.5.1)
type Query struct {
	RootNode TokenInfo
	Segments []Segment
}

var _ Node = DescendantSegment{}

type DescendantKind int

const (
	DescendantWildcardSelector = iota
	DescendantDotNameSelector
	DescendantLongSelector
)

var _ Segment = DescendantSegment{}

type DescendantSegment struct {
	SubKind       DescendantKind
	LongFormInner Segment
}

func (d DescendantSegment) Kind() Kind {
	return DescendantSegmentKind
}

// One of the constructs that selects children ([<selectors>])
// or descendants (..[<selectors>]) of an input value
// segment             = child-segment / descendant-segment

// child-segment       = bracketed-selection /
//
//	("."
//	 (wildcard-selector /
//	  member-name-shorthand))
type ChildSegment struct {
	Kind Kind
}

// Expr represents a JSONPath expression.
type Expr interface {
	Node
	exprNode()
}

// RootNode represents the root node of a JSONPath expression.
type RootNode struct {
	Dollar TokenInfo // position of "$"
}

func (n *RootNode) Pos() TokenInfo { return n.Dollar }
func (n *RootNode) End() TokenInfo {
	return TokenInfo{Token: n.Dollar.Token, Line: n.Dollar.Line, Column: n.Dollar.Column + 1}
}
func (n *RootNode) exprNode() {}

// CurrentNode represents the current node in a JSONPath expression.
type CurrentNode struct {
	At TokenInfo // position of "@"
}

// ComparisonNode represents a comparison expression in a JSONPath filter.
type ComparisonNode struct {
	Lhs      Expr      // left-hand side expression
	Operator TokenInfo // comparison operator
	Rhs      Expr      // right-hand side expression
}

// BooleanNode represents a boolean literal in a JSONPath expression.
type BooleanNode struct {
	Value TokenInfo // boolean value
}

func (n *BooleanNode) Pos() TokenInfo { return n.Value }
func (n *BooleanNode) End() TokenInfo {
	return TokenInfo{Token: n.Value.Token, Line: n.Value.Line, Column: n.Value.Column + len(n.Value.Literal)}
}
func (n *BooleanNode) exprNode() {}

// NumberNode represents a numeric literal in a JSONPath expression.
type NumberNode struct {
	Value TokenInfo // numeric value
}

func (n *NumberNode) Pos() TokenInfo { return n.Value }
func (n *NumberNode) End() TokenInfo {
	return TokenInfo{Token: n.Value.Token, Line: n.Value.Line, Column: n.Value.Column + len(n.Value.Literal)}
}
func (n *NumberNode) exprNode() {}

// StringNode represents a string literal in a JSONPath expression.
type StringNode struct {
	Value TokenInfo // string value
}

func (n *StringNode) Pos() TokenInfo { return n.Value }
func (n *StringNode) End() TokenInfo {
	return TokenInfo{Token: n.Value.Token, Line: n.Value.Line, Column: n.Value.Column + len(n.Value.Literal)}
}
func (n *StringNode) exprNode() {}

// NullNode represents a null literal in a JSONPath expression.
type NullNode struct {
	Null TokenInfo // position of "null"
}

func (n *NullNode) Pos() TokenInfo { return n.Null }
func (n *NullNode) End() TokenInfo {
	return TokenInfo{Token: n.Null.Token, Line: n.Null.Line, Column: n.Null.Column + 4}
}
func (n *NullNode) exprNode() {}

// FunctionCallNode represents a function call in a JSONPath expression.
type FunctionCallNode struct {
	Name   TokenInfo // function name
	Lparen TokenInfo // position of "("
	Args   []Expr    // function arguments
	Rparen TokenInfo // position of ")"
}

func (n *FunctionCallNode) Pos() TokenInfo { return n.Name }
func (n *FunctionCallNode) End() TokenInfo {
	return TokenInfo{Token: n.Rparen.Token, Line: n.Rparen.Line, Column: n.Rparen.Column + 1}
}
func (n *FunctionCallNode) exprNode() {}
