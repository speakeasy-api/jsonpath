package jsonpath

import "gopkg.in/yaml.v3"

type JsonPath interface {
	Evaluate(document *yaml.Node) ([]*yaml.Node, error)
}

type jsonpath struct {
	query JsonPathQuery
}

type JsonPathError struct {
	ShortMessage string `json:"short_message"`
	LongMessage  string `json:"long_message"`
}

func Parse(selector string) (JsonPath, *JsonPathError) {
	return nil, nil
}
