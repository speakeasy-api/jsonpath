package jsonpath

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type JsonPathQuery struct {
	// "$"
	Segments []*Segment
}

func (q JsonPathQuery) ToString() string {
	b := strings.Builder{}
	b.WriteString("$")
	for _, segment := range q.Segments {
		b.WriteString(segment.ToString())
	}
	return b.String()
}

func (q JsonPathQuery) Query(current *yaml.Node, root *yaml.Node) []*yaml.Node {
	var result []*yaml.Node
	//if q.Kind.Token == token.ROOT {
	result = append(result, root)
	/*} else if q.Kind.Token == token.CURRENT {
		result = append(result, current)
	}*/

	for _, segment := range q.Segments {
		newValue := []*yaml.Node{}
		for _, value := range result {
			newValue = append(newValue, segment.Query(value, root)...)
		}
		result = newValue
	}
	return result
}
