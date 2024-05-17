package jsonpath

type Kind int

const (
	ChildSegmentKind Kind = iota
	DescendantSegmentKind
)

// Node represents a node in the JSONPath AST.
type Node interface {
	Pos() TokenInfo // position of the first token belonging to the node
	End() TokenInfo // position of the first token immediately after the node
	String() string // print pretty.
}

// JSONPath expression built from segments that have
// been syntactically restricted in a certain way (Section 2.3.5.1)
type Query struct {
	RootNode TokenInfo
	Segments []Segment
}

// One of the constructs that selects children ([<selectors>])
// or descendants (..[<selectors>]) of an input value
// segment             = child-segment / descendant-segment
type Segment struct {
	Kind              Kind
	ChildSegment      ChildSegment
	DescendantSegment DescendantSegment
}

// child-segment       = bracketed-selection /
//
//	("."
//	 (wildcard-selector /
//	  member-name-shorthand))
type ChildSegment struct {
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

func (n *CurrentNode) Pos() TokenInfo { return n.At }
func (n *CurrentNode) End() TokenInfo {
	return TokenInfo{Token: n.At.Token, Line: n.At.Line, Column: n.At.Column + 1}
}
func (n *CurrentNode) exprNode() {}

// IdentifierNode represents an identifier in a JSONPath expression.
type IdentifierNode struct {
	Name TokenInfo // identifier name
}

func (n *IdentifierNode) Pos() TokenInfo { return n.Name }
func (n *IdentifierNode) End() TokenInfo {
	return TokenInfo{Token: n.Name.Token, Line: n.Name.Line, Column: n.Name.Column + len(n.Name.Literal)}
}
func (n *IdentifierNode) exprNode() {}

// WildcardNode represents a wildcard in a JSONPath expression.
type WildcardNode struct {
	Star TokenInfo // position of "*"
}

func (n *WildcardNode) Pos() TokenInfo { return n.Star }
func (n *WildcardNode) End() TokenInfo {
	return TokenInfo{Token: n.Star.Token, Line: n.Star.Line, Column: n.Star.Column + 1}
}
func (n *WildcardNode) exprNode() {}

// RecursiveDescentNode represents a recursive descent operator in a JSONPath expression.
type RecursiveDescentNode struct {
	DoubleDot TokenInfo // position of ".."
}

func (n *RecursiveDescentNode) Pos() TokenInfo { return n.DoubleDot }
func (n *RecursiveDescentNode) End() TokenInfo {
	return TokenInfo{Token: n.DoubleDot.Token, Line: n.DoubleDot.Line, Column: n.DoubleDot.Column + 2}
}
func (n *RecursiveDescentNode) exprNode() {}

// SubscriptNode represents a subscript operator in a JSONPath expression.
type SubscriptNode struct {
	Lbrack TokenInfo // position of "["
	Index  Expr      // subscript index expression
	Rbrack TokenInfo // position of "]"
}

func (n *SubscriptNode) Pos() TokenInfo { return n.Lbrack }
func (n *SubscriptNode) End() TokenInfo {
	return TokenInfo{Token: n.Rbrack.Token, Line: n.Rbrack.Line, Column: n.Rbrack.Column + 1}
}
func (n *SubscriptNode) exprNode() {}

// SliceNode represents a slice operator in a JSONPath expression.
type SliceNode struct {
	Lbrack TokenInfo // position of "["
	Start  Expr      // start index expression
	Colon1 TokenInfo // position of first ":"
	Finish Expr      // end index expression
	Colon2 TokenInfo // position of second ":", if any
	Step   Expr      // step expression
	Rbrack TokenInfo // position of "]"
}

func (n *SliceNode) Pos() TokenInfo { return n.Lbrack }
func (n *SliceNode) End() TokenInfo {
	return TokenInfo{Token: n.Rbrack.Token, Line: n.Rbrack.Line, Column: n.Rbrack.Column + 1}
}
func (n *SliceNode) exprNode() {}

// UnionNode represents a union operator in a JSONPath expression.
type UnionNode struct {
	Lhs   Expr      // left-hand side expression
	Comma TokenInfo // position of ","
	Rhs   Expr      // right-hand side expression
}

func (n *UnionNode) Pos() TokenInfo { return n.Lhs.Pos() }
func (n *UnionNode) End() TokenInfo { return n.Rhs.End() }
func (n *UnionNode) exprNode()      {}

// FilterNode represents a filter expression in a JSONPath expression.
type FilterNode struct {
	Lbrack TokenInfo // position of "["
	Expr   Expr      // filter expression
	Rbrack TokenInfo // position of "]"
}

func (n *FilterNode) Pos() TokenInfo { return n.Lbrack }
func (n *FilterNode) End() TokenInfo {
	return TokenInfo{Token: n.Rbrack.Token, Line: n.Rbrack.Line, Column: n.Rbrack.Column + 1}
}
func (n *FilterNode) exprNode() {}

// ComparisonNode represents a comparison expression in a JSONPath filter.
type ComparisonNode struct {
	Lhs      Expr      // left-hand side expression
	Operator TokenInfo // comparison operator
	Rhs      Expr      // right-hand side expression
}

func (n *ComparisonNode) Pos() TokenInfo { return n.Lhs.Pos() }
func (n *ComparisonNode) End() TokenInfo { return n.Rhs.End() }
func (n *ComparisonNode) exprNode()      {}

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
