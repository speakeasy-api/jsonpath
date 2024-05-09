package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

// Token represents a lexical token in a JSONPath expression.
type Token int

// The list of tokens.
const (
	ILLEGAL Token = iota
	EOF
	LITERAL
	NUMBER
	STRING
	BOOLEAN
	NULL
	ROOT
	CURRENT
	WILDCARD
	RECURSIVE
	UNION
	CHILD
	SUBSCRIPT
	SLICE
	FILTER
	PAREN_LEFT
	PAREN_RIGHT
	BRACKET_LEFT
	BRACKET_RIGHT
	BRACE_LEFT
	BRACE_RIGHT
	COLON
	COMMA
	DOT
	PIPE
	QUESTION
)

var tokens = [...]string{
	ILLEGAL:       "ILLEGAL",
	EOF:           "EOF",
	LITERAL:       "LITERAL",
	NUMBER:        "NUMBER",
	STRING:        "STRING",
	BOOLEAN:       "BOOLEAN",
	NULL:          "NULL",
	ROOT:          "$",
	CURRENT:       "@",
	WILDCARD:      "*",
	RECURSIVE:     "..",
	UNION:         ",",
	CHILD:         ".",
	SUBSCRIPT:     "[]",
	SLICE:         ":",
	FILTER:        "?",
	PAREN_LEFT:    "(",
	PAREN_RIGHT:   ")",
	BRACKET_LEFT:  "[",
	BRACKET_RIGHT: "]",
	BRACE_LEFT:    "{",
	BRACE_RIGHT:   "}",
	COLON:         ":",
	COMMA:         ",",
	DOT:           ".",
	PIPE:          "|",
	QUESTION:      "?",
}

// String returns the string representation of the token.
func (tok Token) String() string {
	if tok >= 0 && tok < Token(len(tokens)) {
		return tokens[tok]
	}
	return "token(" + strconv.Itoa(int(tok)) + ")"
}

func (t Tokenizer) ErrorString(target TokenInfo, msg string) string {
	var errorBuilder strings.Builder

	// Write the error message with line and column information
	errorBuilder.WriteString(fmt.Sprintf("Error at line %d, column %d: %s\n", target.Line, target.Column, msg))

	// Find the start and end positions of the line containing the target token
	lineStart := 0
	lineEnd := len(t.input)
	for i := target.Line - 1; i > 0; i-- {
		if pos := strings.LastIndexByte(t.input[:lineStart], '\n'); pos != -1 {
			lineStart = pos + 1
			break
		}
	}
	if pos := strings.IndexByte(t.input[lineStart:], '\n'); pos != -1 {
		lineEnd = lineStart + pos
	}

	// Extract the line containing the target token
	line := t.input[lineStart:lineEnd]
	errorBuilder.WriteString(line)
	errorBuilder.WriteString("\n")

	// Calculate the number of spaces before the target token
	spaces := strings.Repeat(" ", target.Column)

	// Write the caret symbol pointing to the target token
	errorBuilder.WriteString(spaces)
	errorBuilder.WriteString("^\n")

	return errorBuilder.String()
}

// TokenInfo represents a token and its associated information.
type TokenInfo struct {
	Token   Token
	Line    int
	Column  int
	Literal string
}

// Tokenizer represents a JSONPath tokenizer.
type Tokenizer struct {
	input  string
	pos    int
	line   int
	column int
	tokens []TokenInfo
}

// NewTokenizer creates a new JSONPath tokenizer for the given input string.
func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		input: input,
		line:  1,
	}
}

// Tokenize tokenizes the input string and returns a slice of TokenInfo.
func (t *Tokenizer) Tokenize() []TokenInfo {
	for t.pos < len(t.input) {
		t.skipWhitespace()
		if t.pos >= len(t.input) {
			break
		}

		switch ch := t.input[t.pos]; {
		case ch == '$':
			t.addToken(ROOT, "")
		case ch == '@':
			t.addToken(CURRENT, "")
		case ch == '*':
			t.addToken(WILDCARD, "")
		case ch == '.':
			if t.peek() == '.' {
				t.addToken(RECURSIVE, "")
			} else {
				t.addToken(CHILD, "")
			}
		case ch == ',':
			t.addToken(UNION, "")
		case ch == ':':
			t.addToken(SLICE, "")
		case ch == '?':
			t.addToken(FILTER, "")
		case ch == '(':
			t.addToken(PAREN_LEFT, "")
		case ch == ')':
			t.addToken(PAREN_RIGHT, "")
		case ch == '[':
			t.addToken(BRACKET_LEFT, "")
		case ch == ']':
			t.addToken(BRACKET_RIGHT, "")
		case ch == '{':
			t.addToken(BRACE_LEFT, "")
		case ch == '}':
			t.addToken(BRACE_RIGHT, "")
		case ch == '|':
			t.addToken(PIPE, "")
		case ch == '"' || ch == '\'':
			t.scanString(rune(ch))

		case isDigit(ch):
			t.scanNumber()
		case isLetter(ch):
			t.scanLiteral()
		default:
			t.addToken(ILLEGAL, string(ch))
		}
		t.pos++
		t.column++
	}

	t.addToken(EOF, "")
	return t.tokens
}

func (t *Tokenizer) addToken(token Token, literal string) {
	t.tokens = append(t.tokens, TokenInfo{
		Token:   token,
		Line:    t.line,
		Column:  t.column,
		Literal: literal,
	})
}

func (t *Tokenizer) scanString(quote rune) {
	start := t.pos + 1
	for i := start; i < len(t.input); i++ {
		if t.input[i] == byte(quote) {
			t.addToken(STRING, t.input[start:i])
			t.pos = i
			t.column += i - start + 1
			return
		}
	}
	t.addToken(ILLEGAL, t.input[start:])
	t.pos = len(t.input) - 1
	t.column = len(t.input) - 1
}

func (t *Tokenizer) scanNumber() {
	start := t.pos
	for i := start; i < len(t.input); i++ {
		if !isDigit(t.input[i]) {
			t.addToken(NUMBER, t.input[start:i])
			t.pos = i - 1
			t.column += i - start - 1
			return
		}
	}
	t.addToken(NUMBER, t.input[start:])
	t.pos = len(t.input) - 1
	t.column = len(t.input) - 1
}

func (t *Tokenizer) scanLiteral() {
	start := t.pos
	for i := start; i < len(t.input); i++ {
		if !isLetter(t.input[i]) {
			literal := t.input[start:i]
			switch literal {
			case "true", "false":
				t.addToken(BOOLEAN, literal)
			case "null":
				t.addToken(NULL, literal)
			default:
				t.addToken(LITERAL, literal)
			}
			t.pos = i - 1
			t.column += i - start - 1
			return
		}
	}
	literal := t.input[start:]
	switch literal {
	case "true", "false":
		t.addToken(BOOLEAN, literal)
	case "null":
		t.addToken(NULL, literal)
	default:
		t.addToken(LITERAL, literal)
	}
	t.pos = len(t.input) - 1
	t.column = len(t.input) - 1

}

func (t *Tokenizer) skipWhitespace() {
	for t.pos < len(t.input) {
		ch := t.input[t.pos]
		if ch == '\n' {
			t.line++
			t.pos++
			t.column = 0
		} else if !isSpace(ch) {
			break
		} else {
			t.pos++
			t.column++
		}
	}
}

func (t *Tokenizer) peek() byte {
	if t.pos+1 < len(t.input) {
		return t.input[t.pos+1]
	}
	return 0
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r'
}
