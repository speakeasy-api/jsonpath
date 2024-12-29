package jsonpath

import (
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath/token"
	"gopkg.in/yaml.v3"
	"strings"
)

type Query struct {
	// "@" or "$"
	Kind     token.TokenInfo
	Segments []*Segment
}

func (q Query) ToString() string {
	b := strings.Builder{}
	if q.Kind.Token == token.ROOT {
		b.WriteString("$")
	} else if q.Kind.Token == token.CURRENT {
		b.WriteString("@")
	}
	for _, segment := range q.Segments {
		b.WriteString(segment.ToString())
	}
	return b.String()
}

func (q Query) Query(current *yaml.Node, root *yaml.Node) []*yaml.Node {
	var result []*yaml.Node
	if q.Kind.Token == token.ROOT {
		result = append(result, root)
	} else if q.Kind.Token == token.CURRENT {
		result = append(result, current)
	}

	for _, segment := range q.Segments {
		newValue := []*yaml.Node{}
		for _, value := range result {
			newValue = append(newValue, segment.Query(value, root)...)
		}
		result = newValue
	}
	return result
}
