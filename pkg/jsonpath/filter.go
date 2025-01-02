package jsonpath

import "gopkg.in/yaml.v3"

// filter-selector     = "?" S logical-expr
type filterSelector struct {
	// logical-expr        = logical-or-expr
	expression *logicalOrExpr
}

// logical-or-expr     = logical-and-expr *(S "||" S logical-and-expr)
type logicalOrExpr struct {
	expressions []*logicalAndExpr
}

// logical-and-expr    = basic-expr *(S "&&" S basic-expr)
type logicalAndExpr struct {
	expressions []*basicExpr
}

// relQuery rel-query = current-node-identifier segments
// current-node-identifier = "@"
type relQuery struct {
	segments []*segment
}

// filterQuery filter-query        = rel-query / jsonpath-query
type filterQuery struct {
	relQuery      *relQuery
	jsonPathQuery *jsonPathAST
}

// functionArgument function-argument   = literal /
//
//	filter-query / ; (includes singular-query)
//	logical-expr /
//	function-expr
type functionArgument struct {
	literal      *literal
	filterQuery  *filterQuery
	logicalExpr  *logicalOrExpr
	functionExpr *functionExpr
}

//function-name       = function-name-first *function-name-char
//function-name-first = LCALPHA
//function-name-char  = function-name-first / "_" / DIGIT
//LCALPHA             = %x61-7A  ; "a".."z"
//

type functionType int

const (
	functionTypeLength functionType = iota
	functionTypeCount
	functionTypeMatch
	functionTypeSearch
	functionTypeValue
)

var functionTypeMap = map[string]functionType{
	"length": functionTypeLength,
	"count":  functionTypeCount,
	"match":  functionTypeMatch,
	"search": functionTypeSearch,
	"value":  functionTypeValue,
}

func (f functionType) String() string {
	for k, v := range functionTypeMap {
		if v == f {
			return k
		}
	}
	return "unknown"
}

// functionExpr function-expr       = function-name "(" S [function-argument
// *(S "," S function-argument)] S ")"
type functionExpr struct {
	funcType functionType
	args     []*functionArgument
}

// testExpr test-expr           = [logical-not-op S]
//
//	(filter-query / ; existence/non-existence
//	 function-expr) ; LogicalType or NodesType
type testExpr struct {
	not          bool
	filterQuery  *filterQuery
	functionExpr *functionExpr
}

// basicExpr basic-expr          =
//
//	 paren-expr /
//		comparison-expr /
//		test-expr
type basicExpr struct {
	parenExpr      *parenExpr
	comparisonExpr *comparisonExpr
	testExpr       *testExpr
}

// literal literal = number /
// . string-literal /
// . true / false / null
type literal struct {
	// we generally decompose these into their component parts for easier evaluation
	integer *int
	float64 *float64
	string  *string
	bool    *bool
	null    *bool
	node    *yaml.Node
}

type absQuery jsonPathAST

// singularQuery singular-query = rel-singular-query / abs-singular-query
type singularQuery struct {
	relQuery *relQuery
	absQuery *absQuery
}

// comparable
//
//	comparable = literal /
//	singular-query / ; singular query value
//	function-expr    ; ValueType
type comparable struct {
	literal       *literal
	singularQuery *singularQuery
	functionExpr  *functionExpr
}

// comparisonExpr represents a comparison expression
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
type comparisonExpr struct {
	left  *comparable
	op    comparisonOperator
	right *comparable
}

// existExpr represents an existence expression
type existExpr struct {
	query string
}

// parenExpr represents a parenthesized expression
//
//	paren-expr          = [logical-not-op S] "(" S logical-expr S ")"
type parenExpr struct {
	// "!"
	not bool
	// "(" logicalOrExpr ")"
	expr *logicalOrExpr
}

// comparisonOperator represents a comparison operator
type comparisonOperator int

const (
	equalTo comparisonOperator = iota
	notEqualTo
	lessThan
	lessThanEqualTo
	greaterThan
	greaterThanEqualTo
)
