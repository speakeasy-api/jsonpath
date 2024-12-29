package jsonpath

// filter-selector     = "?" S logical-expr
type FilterSelector struct {
	// logical-expr        = logical-or-expr
	Expression *LogicalOrExpr
}

// logical-or-expr     = logical-and-expr *(S "||" S logical-and-expr)
type LogicalOrExpr struct {
	Expressions []*LogicalAndExpr
}

// logical-and-expr    = basic-expr *(S "&&" S basic-expr)
type LogicalAndExpr struct {
	Expressions []*BasicExpr
}

// RelQuery rel-query = current-node-identifier segments
// current-node-identifier = "@"
type RelQuery struct {
	Segments []*Segment
}

// FilterQuery filter-query        = rel-query / jsonpath-query
type FilterQuery struct {
	RelQuery      *RelQuery
	JsonPathQuery *JsonPathQuery
}

// FunctionArgument function-argument   = literal /
//
//	filter-query / ; (includes singular-query)
//	logical-expr /
//	function-expr
type FunctionArgument struct {
	Literal      *Literal
	FilterQuery  *FilterQuery
	LogicalExpr  *LogicalOrExpr
	FunctionExpr *FunctionExpr
}

//function-name       = function-name-first *function-name-char
//function-name-first = LCALPHA
//function-name-char  = function-name-first / "_" / DIGIT
//LCALPHA             = %x61-7A  ; "a".."z"
//

type FunctionType int

const (
	FunctionTypeLength FunctionType = iota
	FunctionTypeCount
	FunctionTypeMatch
	FunctionTypeSearch
	FunctionTypeValue
)

func (f FunctionType) String() string {
	switch f {
	case FunctionTypeLength:
		return "length"
	case FunctionTypeCount:
		return "count"
	case FunctionTypeMatch:
		return "match"
	case FunctionTypeSearch:
		return "search"
	case FunctionTypeValue:
		return "value"
	default:
		return "unknown"
	}
}

// FunctionExpr function-expr       = function-name "(" S [function-argument
// *(S "," S function-argument)] S ")"
type FunctionExpr struct {
	Type FunctionType
	Args []*FunctionArgument
}

// TestExpr test-expr           = [logical-not-op S]
//
//	(filter-query / ; existence/non-existence
//	 function-expr) ; LogicalType or NodesType
type TestExpr struct {
	Not          bool
	FilterQuery  *FilterQuery
	FunctionExpr *FunctionExpr
}

// BasicExpr basic-expr          =
//
//	 paren-expr /
//		comparison-expr /
//		test-expr
type BasicExpr struct {
	ParenExpr      *ParenExpr
	ComparisonExpr *ComparisonExpr
	TestExpr       *TestExpr
}

// Literal literal = number /
// . string-literal /
// . true / false / null
type Literal struct {
	Integer *int
	Float64 *float64
	String  *string
	Bool    *bool
	Null    *bool
}

type AbsQuery JsonPathQuery

// SingularQuery singular-query = rel-singular-query / abs-singular-query
type SingularQuery struct {
	RelQuery *RelQuery
	AbsQuery *AbsQuery
}

// Comparable
//
//	comparable = literal /
//	singular-query / ; singular query value
//	function-expr    ; ValueType
type Comparable struct {
	Literal       *Literal
	SingularQuery *SingularQuery
}

// ComparisonExpr represents a comparison expression
//
//	comparison-expr     = comparable S comparison-op S comparable
//	literal             = number / string-literal /
//	                      true / false / null
//	comparable          = literal /
//	                      singular-query / ; singular query value
//	                      function-expr    ; ValueType
//	comparison-op       = "==" / "!=" /
//	                      "<=" / ">=" /
//	                      "<"  / ">"
type ComparisonExpr struct {
	Left  *Comparable
	Op    ComparisonOperator
	Right *Comparable
}

// ExistExpr represents an existence expression
type ExistExpr struct {
	Query string
}

// ParenExpr represents a parenthesized expression
//
//	paren-expr          = [logical-not-op S] "(" S logical-expr S ")"
type ParenExpr struct {
	// "!"
	Not bool
	// "(" LogicalOrExpr ")"
	Expr *LogicalOrExpr
}

// ComparisonOperator represents a comparison operator
type ComparisonOperator int

const (
	EqualTo ComparisonOperator = iota
	NotEqualTo
	LessThan
	LessThanEqualTo
	GreaterThan
	GreaterThanEqualTo
)
