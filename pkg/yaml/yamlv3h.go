package yaml

import (
	goYAML "github.com/goccy/go-yaml"
	ast "github.com/goccy/go-yaml/ast"
	"io"
)

type Unmarshaler interface{ UnmarshalYAML(value *Node) error }

type Marshaler interface{ MarshalYAML() (interface{}, error) }

type Decoder struct {
	reader  io.Reader
	decoder *goYAML.Decoder
	options decoderOptions
}

type Encoder struct {
	w io.Writer
}

type TypeError struct{ Errors []string }

type Kind uint32

const (
	DocumentNode Kind = 1 << iota
	SequenceNode
	MappingNode
	ScalarNode
	AliasNode
)

type Style uint32

const (
	TaggedStyle Style = 1 << iota
	DoubleQuotedStyle
	SingleQuotedStyle
	LiteralStyle
	FoldedStyle
	FlowStyle
)

type Node struct {
	Kind        Kind    `yaml:"kind"`
	Style       Style   `yaml:"style"`
	Tag         string  `yaml:"tag"`
	Value       string  `yaml:"value"`
	Anchor      string  `yaml:"anchor"`
	Alias       *Node   `yaml:"alias"`
	Content     []*Node `yaml:"content"`
	HeadComment string  `yaml:"head_comment"`
	LineComment string  `yaml:"line_comment"`
	FootComment string  `yaml:"foot_comment"`
	Line        int     `yaml:"line"`
	Column      int     `yaml:"column"`
	ast         ast.Node
}

type IsZeroer interface{ IsZero() bool }
