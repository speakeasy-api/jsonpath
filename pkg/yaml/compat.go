package yaml

import (
	"bytes"
	"github.com/goccy/go-yaml/lexer"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"
	"strconv"

	"io"
	"reflect"
	"unsafe"
)
import goYAML "github.com/goccy/go-yaml"
import ast "github.com/goccy/go-yaml/ast"

func Unmarshal(in []byte, out interface{}) (err error) {
	// if out of type Node, then we'll pass in a goYAML ast node instead
	// then convert it back.
	if outNode, ok := out.(*Node); ok {
		n, err := parser.ParseBytes(in, parser.ParseComments)
		if err != nil {
			return err
		}
		if n != nil && len(n.Docs) == 1 {
			result := astToNode(n.Docs[0])
			if result != nil {
				*outNode = *result
				return nil
			}
		}
	}
	return goYAML.Unmarshal(in, out)

}

func getDocumentNode(dec *goYAML.Decoder) *ast.File {
	// access struct field "parsedFile" inside the struct
	v := reflect.ValueOf(dec)
	v = v.Elem()
	v = v.FieldByName("parsedFile")
	v = reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem()
	// cast it back
	vInt := v.Interface()
	return vInt.(*ast.File)
}
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		reader:  r,
		decoder: nil,
		options: decoderOptions{},
	}
}

type decoderOptions struct {
	knownFields bool
}

func newDecoder(t *Decoder) *goYAML.Decoder {
	opts := []goYAML.DecodeOption{}
	if t.options.knownFields {
		opts = append(opts, goYAML.Strict())
	}
	return goYAML.NewDecoder(t.reader, opts...)
}

func (dec *Decoder) KnownFields(enable bool) {
	// KnownFields ensures that the keys in decoded mappings to
	// exist as fields in the struct being decoded into.
	dec.options.knownFields = enable
}
func (dec *Decoder) Decode(v interface{}) (err error) {
	dec.decoder = newDecoder(dec)
	return dec.decoder.Decode(v)
}
func (n *Node) Decode(v interface{}) (err error) {
	newValue, err := goYAML.ValueToNode(v)
	if err != nil {
		return err
	}
	newNode := astToNode(newValue)
	*n = *newNode
	return nil
}

func astToNode(value ast.Node) *Node {
	switch value := value.(type) {
	case *ast.StringNode:
		return &Node{
			ast:    value,
			Kind:   ScalarNode,
			Style:  getStringStyle(value),
			Tag:    string(token.StringTag),
			Value:  value.Value,
			Line:   value.Token.Position.Line,
			Column: value.Token.Position.Column,
		}
	case *ast.NullNode:
		return &Node{
			ast:    value,
			Kind:   ScalarNode,
			Tag:    string(token.NullTag),
			Value:  value.Token.Value,
			Line:   value.Token.Position.Line,
			Column: value.Token.Position.Column,
		}
	case *ast.IntegerNode:
		return &Node{
			ast:    value,
			Kind:   ScalarNode,
			Tag:    string(token.IntegerTag),
			Value:  value.Token.Value,
			Line:   value.Token.Position.Line,
			Column: value.Token.Position.Column,
		}
	case *ast.FloatNode:
		return &Node{
			ast:    value,
			Kind:   ScalarNode,
			Tag:    string(token.FloatTag),
			Value:  value.Token.Value,
			Line:   value.Token.Position.Line,
			Column: value.Token.Position.Column,
		}
	case *ast.MergeKeyNode:
		return nil
	case *ast.BoolNode:
		return &Node{
			ast:    value,
			Kind:   ScalarNode,
			Tag:    string(token.BooleanTag),
			Value:  value.Token.Value,
			Line:   value.Token.Position.Line,
			Column: value.Token.Position.Column,
		}
	case *ast.InfinityNode:
		return &Node{
			ast:    value,
			Kind:   ScalarNode,
			Tag:    string(token.FloatTag),
			Value:  value.Token.Value,
			Line:   value.Token.Position.Line,
			Column: value.Token.Position.Column,
		}
	case *ast.NanNode:
		return &Node{
			ast:    value,
			Kind:   ScalarNode,
			Tag:    string(token.FloatTag),
			Value:  value.Token.Value,
			Line:   value.Token.Position.Line,
			Column: value.Token.Position.Column,
		}
	case *ast.LiteralNode:
		subNode := astToNode(value.Value)
		if value.Start.Indicator == token.BlockScalarIndicator {
			subNode.Style |= LiteralStyle
			subNode.Line = value.Start.Position.Line
			subNode.Column = value.Start.Position.Column
		}
		return subNode
	case *ast.DirectiveNode:
		return nil
	case *ast.TagNode:
		subNode := astToNode(value.Value)
		subNode.Style |= TaggedStyle
		subNode.Tag = value.Start.Value
		subNode.Line = value.Start.Position.Line
		subNode.Column = value.Start.Position.Column

		return subNode
	case *ast.DocumentNode:
		return &Node{
			ast:     value,
			Kind:    DocumentNode,
			Content: []*Node{astToNode(value.Body)},
			Line:    value.GetToken().Position.Line,
			Column:  value.GetToken().Position.Column,
		}
	case *ast.MappingNode:
		var children []*Node
		for _, childValue := range value.Values {
			children = append(children, astToNode(childValue.Key))
			children = append(children, astToNode(childValue.Value))
		}
		style := Style(0)
		if value.IsFlowStyle {
			style |= FlowStyle
		}
		return &Node{
			ast:     value,
			Style:   style,
			Kind:    MappingNode,
			Content: children,
			Line:    value.GetToken().Position.Line,
			Column:  value.GetToken().Position.Column,
		}
	case *ast.MappingKeyNode:
		return astToNode(value.Value)
	case *ast.MappingValueNode:
		return astToNode(value.Value)
	case *ast.SequenceNode:
		var children []*Node
		for _, value := range value.Values {
			children = append(children, astToNode(value))
		}
		style := Style(0)
		if value.IsFlowStyle {
			style |= FlowStyle
		}

		return &Node{
			ast:     value,
			Style:   style,
			Kind:    SequenceNode,
			Content: children,
			Line:    value.GetToken().Position.Line,
			Column:  value.GetToken().Position.Column,
		}
	case *ast.AnchorNode:
		return astToNode(value.Value)
	case *ast.AliasNode:
		return &Node{
			ast:    value,
			Kind:   AliasNode,
			Anchor: "",
			Alias:  astToNode(value.Value),
			Line:   value.GetToken().Position.Line,
			Column: value.GetToken().Position.Column,
		}
	}
	return nil
}

func nodeToAST(t *Node) ast.Node {
	if t.ast != nil {
		return t.ast
	}
	if t.Style&TaggedStyle != 0 || unknownTag(t.Tag) {
		newT := *t
		newT.Style &= ^TaggedStyle
		newT.Tag = ""
		tagged := &ast.TagNode{
			BaseNode: &ast.BaseNode{},
			Start: &token.Token{
				Type:   token.TagType,
				Origin: t.Tag,
				Value:  t.Tag,
			},
			Value: nodeToAST(&newT),
		}
		if tagged.Value.GetToken() != nil && tagged.Value.GetToken().Type == token.DoubleQuoteType && t.Tag == string(token.StringTag) {
			tagged.Value.GetToken().Type = token.StringType
		}
		return tagged
	}

	switch t.Kind {
	case DocumentNode:
		return &ast.DocumentNode{
			BaseNode: &ast.BaseNode{},
			Body:     nodeToAST(t.Content[0]),
		}
	case MappingNode:
		mappingNode := &ast.MappingNode{
			BaseNode: &ast.BaseNode{},
			Values:   []*ast.MappingValueNode{},
		}
		for i := 0; i < len(t.Content); i += 2 {
			mappingNode.Values = append(mappingNode.Values, &ast.MappingValueNode{
				Key:   &ast.MappingKeyNode{Value: nodeToAST(t.Content[i])},
				Value: nodeToAST(t.Content[i+1]),
			})
		}
		return mappingNode
	case SequenceNode:
		sequenceNode := &ast.SequenceNode{
			BaseNode: &ast.BaseNode{},
			Values:   []ast.Node{},
		}
		for i := 0; i < len(t.Content); i++ {
			sequenceNode.Values = append(sequenceNode.Values, nodeToAST(t.Content[i]))
		}
		return sequenceNode
	case AliasNode:
		return &ast.AliasNode{
			BaseNode: &ast.BaseNode{},
			Value:    nodeToAST(t.Alias),
		}
	case ScalarNode:
		switch t.Tag {
		case string(token.BinaryTag):
			newT := *t
			newT.Tag = "!!str"
			return &ast.TagNode{
				BaseNode: &ast.BaseNode{},
				Start: &token.Token{
					Type:   token.TagType,
					Origin: t.Tag,
					Value:  t.Tag,
				},
				Value: nodeToAST(&newT),
			}
		case string(token.StringTag):
			return &ast.StringNode{
				BaseNode: &ast.BaseNode{},
				Token: &token.Token{
					Type:  getTokenType(t.Style, t.Value),
					Value: t.Value,
					Position: &token.Position{
						Line:   t.Line,
						Column: t.Column,
					},
				},
				Value: t.Value,
			}
		case string(token.IntegerTag):
			return &ast.IntegerNode{
				BaseNode: &ast.BaseNode{},
				Token: &token.Token{
					Type:  token.IntegerType,
					Value: t.Value,
					Position: &token.Position{
						Line:   t.Line,
						Column: t.Column,
					},
				},
				Value: t.Value,
			}
		case string(token.FloatTag):
			f, _ := strconv.ParseFloat(t.Value, 64)
			return &ast.FloatNode{
				BaseNode: &ast.BaseNode{},
				Token: &token.Token{
					Type:  token.FloatType,
					Value: t.Value,
					Position: &token.Position{
						Line:   t.Line,
						Column: t.Column,
					},
				},
				Value: f,
			}
		case string(token.BooleanTag):
			bool, _ := strconv.ParseBool(t.Value)
			return &ast.BoolNode{
				BaseNode: &ast.BaseNode{},
				Token: &token.Token{
					Value: t.Value,
					Position: &token.Position{
						Line:   t.Line,
						Column: t.Column,
					},
				},
				Value: bool,
			}
		case string(token.NullTag):
			return &ast.NullNode{
				BaseNode: &ast.BaseNode{},
			}
		case string(token.TimestampTag):
			return &ast.LiteralNode{
				BaseNode: &ast.BaseNode{},
				Value: &ast.StringNode{
					BaseNode: &ast.BaseNode{},
					Token: &token.Token{
						Type:   token.TagType,
						Origin: t.Value,
						Value:  t.Value,
					},
					Value: t.Value,
				},
			}

		case "":
			// parse the value
			parsed, err := parser.ParseBytes([]byte(t.Value), parser.ParseComments)
			if err != nil || len(parsed.Docs) != 1 {
				return &ast.NullNode{
					BaseNode: &ast.BaseNode{},
				}
			}

			return parsed.Docs[0].Body
		default:
			newT := *t
			newT.Style = 0
			newT.Tag = ""
			return &ast.TagNode{
				BaseNode: &ast.BaseNode{},
				Start: &token.Token{
					Type:   token.TagType,
					Origin: t.Tag,
					Value:  t.Tag,
				},
				Value: nodeToAST(&newT),
			}

		}

	}
	return &ast.NullNode{
		BaseNode: &ast.BaseNode{},
	}
}

func unknownTag(tag string) bool {
	return tag != string(token.IntegerTag) &&
		tag != string(token.FloatTag) &&
		tag != string(token.NullTag) &&
		tag != string(token.SequenceTag) &&
		tag != string(token.MappingTag) &&
		tag != string(token.StringTag) &&
		tag != string(token.BinaryTag) &&
		tag != string(token.OrderedMapTag) &&
		tag != string(token.SetTag) &&
		tag != string(token.TimestampTag) &&
		tag != string(token.BooleanTag) &&
		tag != string(token.MergeTag) &&
		tag != ""
}

func getStringStyle(n *ast.StringNode) Style {
	switch n.Token.Type {
	case token.SingleQuoteType:
		return SingleQuotedStyle
	case token.DoubleQuoteType:
		return DoubleQuotedStyle
	}
	return 0
}
func getTokenType(n Style, value string) token.Type {
	switch n {
	case SingleQuotedStyle:
		return token.SingleQuoteType
	case DoubleQuotedStyle:
		return token.DoubleQuoteType
	}
	tokenized := lexer.Tokenize(value)
	if len(tokenized) == 1 && tokenized[0].Type != token.StringType {
		return token.DoubleQuoteType
	}

	return 0
}

func Marshal(in interface{}) (out []byte, err error) { return goYAML.Marshal(in) }
func NewEncoder(w io.Writer) *Encoder                { return &Encoder{w} }
func (e *Encoder) Encode(v interface{}) (err error) {
	return goYAML.NewEncoder(e.w, goYAML.CustomMarshaler[Node](func(t Node) ([]byte, error) {
		astNode := nodeToAST(&t)
		if astNode == nil {
			return []byte{}, nil
		}
		return []byte(astNode.String()), nil
	})).Encode(v)
}

func (n *Node) Encode(v interface{}) (err error) {
	in := reflect.ValueOf(v)
	if in.IsValid() {
		if node, ok := in.Interface().(*Node); ok {
			*n = *node
			return nil
		}

	}

	buf := bytes.NewBuffer([]byte{})
	node, err := goYAML.NewEncoder(buf).EncodeToNode(v)
	if err != nil {
		return err
	}
	*n = *astToNode(node)
	return nil
}

func (e *Encoder) SetIndent(spaces int) {}
func (e *Encoder) Close() (err error)   { return nil }
func (e *TypeError) Error() string      { return "" }
func (n *Node) IsZero() bool            { return false }
func (n *Node) LongTag() string         { return "" }
func (n *Node) ShortTag() string        { return "" }
func (n *Node) SetString(s string)      {}
