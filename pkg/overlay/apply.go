package overlay

import (
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath/config"
	"go.yaml.in/yaml/v4"
)

// ApplyTo will take an overlay and apply its changes to the given YAML
// document.
func (o *Overlay) ApplyTo(root *yaml.Node) error {
	for _, action := range o.Actions {
		var err error
		if action.Remove {
			err = applyRemoveAction(root, action)
		} else {
			err = applyUpdateAction(root, action)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func applyRemoveAction(root *yaml.Node, action Action) error {
	if action.Target == "" {
		return nil
	}

	idx := newParentIndex(root)

	p, err := jsonpath.NewPath(action.Target, config.WithPropertyNameExtension())
	if err != nil {
		return err
	}

	nodes := p.Query(root)
	if err != nil {
		return err
	}

	for _, node := range nodes {
		removeNode(idx, node)
	}

	return nil
}

func removeNode(idx parentIndex, node *yaml.Node) {
	parent := idx.getParent(node)
	if parent == nil {
		return
	}

	for i, child := range parent.Content {
		if child == node {
			switch parent.Kind {
			case yaml.MappingNode:
				if i%2 == 1 {
					// if we select a value, we should delete the key too
					parent.Content = append(parent.Content[:i-1], parent.Content[i+1:]...)
				} else {
					// if we select a key, we should delete the value
					parent.Content = append(parent.Content[:i], parent.Content[i+2:]...)
				}
				return
			case yaml.SequenceNode:
				parent.Content = append(parent.Content[:i], parent.Content[i+1:]...)
				return
			}
		}
	}
}

func applyUpdateAction(root *yaml.Node, action Action) error {
	if action.Target == "" {
		return nil
	}

	if action.Update.IsZero() {
		return nil
	}

	p, err := jsonpath.NewPath(action.Target, config.WithPropertyNameExtension())
	if err != nil {
		return err
	}

	nodes := p.Query(root)

	for _, node := range nodes {
		if err := updateNode(node, &action.Update); err != nil {
			return err
		}
	}

	return nil
}

func updateNode(node *yaml.Node, updateNode *yaml.Node) error {
	mergeNode(node, updateNode)
	return nil
}

func mergeNode(node *yaml.Node, merge *yaml.Node) {
	if node.Kind != merge.Kind {
		*node = *clone(merge)
		return
	}
	switch node.Kind {
	default:
		node.Value = merge.Value
	case yaml.MappingNode:
		mergeMappingNode(node, merge)
	case yaml.SequenceNode:
		mergeSequenceNode(node, merge)
	}
}

// mergeMappingNode will perform a shallow merge of the merge node into the main
// node.
func mergeMappingNode(node *yaml.Node, merge *yaml.Node) {
NextKey:
	for i := 0; i < len(merge.Content); i += 2 {
		mergeKey := merge.Content[i].Value
		mergeValue := merge.Content[i+1]

		for j := 0; j < len(node.Content); j += 2 {
			nodeKey := node.Content[j].Value
			if nodeKey == mergeKey {
				mergeNode(node.Content[j+1], mergeValue)
				continue NextKey
			}
		}

		node.Content = append(node.Content, merge.Content[i], clone(mergeValue))
	}
}

// mergeSequenceNode will append the merge node's content to the original node.
func mergeSequenceNode(node *yaml.Node, merge *yaml.Node) {
	node.Content = append(node.Content, clone(merge).Content...)
}

func clone(node *yaml.Node) *yaml.Node {
	newNode := &yaml.Node{
		Kind:        node.Kind,
		Style:       node.Style,
		Tag:         node.Tag,
		Value:       node.Value,
		Anchor:      node.Anchor,
		HeadComment: node.HeadComment,
		LineComment: node.LineComment,
		FootComment: node.FootComment,
	}
	if node.Alias != nil {
		newNode.Alias = clone(node.Alias)
	}
	if node.Content != nil {
		newNode.Content = make([]*yaml.Node, len(node.Content))
		for i, child := range node.Content {
			newNode.Content[i] = clone(child)
		}
	}
	return newNode
}
