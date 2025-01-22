package overlay

import (
	"bytes"
	"fmt"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath"
	"github.com/speakeasy-api/jsonpath/pkg/yaml"
	"sort"
	"strings"
)

func getIndents(node *yaml.Node) []int {
	indentValues := []int{}
	if node.Kind == yaml.MappingNode {
		var nodeColumnDelta int
		// find the child key column
		// if it's valid (might not be if we've dynamically constructed yaml)
		if len(node.Content) > 0 && node.Line != node.Content[0].Line && node.Content[0].Column > node.Column {
			nodeColumnDelta = node.Content[0].Column - node.Column
			indentValues = append(indentValues, nodeColumnDelta)
		}
		// Now do this recursively
		for i, child := range node.Content {
			// values only
			if i%2 == 1 {
				indentValues = append(indentValues, getIndents(child)...)
			}
		}
		return indentValues
	} else if node.Kind == yaml.SequenceNode {
		var nodeColumnDelta int
		// find the child key column
		// if it's valid (might not be if we've dynamically constructed yaml)
		if len(node.Content) > 0 && node.Line != node.Content[0].Line && node.Content[0].Column > node.Column {
			nodeColumnDelta = node.Content[0].Column - node.Column
			indentValues = append(indentValues, nodeColumnDelta)
		}
		// Now do this recursively
		for _, child := range node.Content {
			indentValues = append(indentValues, getIndents(child)...)
		}
	}
	return indentValues
}

func CalculateIndent(node *yaml.Node) int {
	indents := getIndents(node)
	if len(indents) == 0 {
		return 2
	}
	// find the median indent
	sort.Ints(indents)
	median := indents[len(indents)/2]
	return median
}

func IsAllFlowStyle(y *yaml.Node) bool {
	for _, node := range y.Content {
		if !IsAllFlowStyle(node) {
			return false
		}
	}
	if len(y.Content) > 0 && y.Style != yaml.FlowStyle {
		return false
	}
	return true
}

type Encoder interface {
	Encode(value interface{}) error
}

func ApplyOverlay(originalYAML, overlayYAML string) (string, error) {
	var orig yaml.Node
	err := yaml.Unmarshal([]byte(originalYAML), &orig)
	if err != nil {
		return "", fmt.Errorf("failed to parse original schema: %w", err)
	}
	// Unwrap the document node if it exists and has only one content node
	if orig.Kind == yaml.DocumentNode && len(orig.Content) == 1 {
		orig = *orig.Content[0]
	}

	indent := CalculateIndent(&orig)

	var overlay Overlay
	err = yaml.Unmarshal([]byte(overlayYAML), &overlay)
	if err != nil {
		return "", fmt.Errorf("failed to parse overlay schema: %w", err)
	}

	err = overlay.ApplyTo(&orig)
	if err != nil {
		return "", fmt.Errorf("failed to apply overlay: %w", err)
	}

	buf := bytes.NewBuffer([]byte{})
	yamlEncoder := yaml.NewEncoder(buf)
	yamlEncoder.SetIndent(indent)
	err = yamlEncoder.Encode(&orig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return buf.String(), nil
}

func toJSONIndent(indent int) string {
	// " " * the length of the indent
	return strings.Repeat(" ", indent)
}

func CalculateOverlay(originalYAML, targetYAML, existingOverlay string) (string, error) {
	var orig yaml.Node
	err := yaml.Unmarshal([]byte(originalYAML), &orig)
	if err != nil {
		return "", fmt.Errorf("failed to parse source schema: %w", err)
	}
	var target yaml.Node
	err = yaml.Unmarshal([]byte(targetYAML), &target)
	if err != nil {
		return "", fmt.Errorf("failed to parse target schema: %w", err)
	}

	// we go from the original to a new version, then look at the extra overlays on top
	// of that, then add that to the existing overlay
	var existingOverlayDocument Overlay
	err = yaml.Unmarshal([]byte(existingOverlay), &existingOverlayDocument)
	if err != nil {
		return "", fmt.Errorf("failed to parse overlay schema: %w", err)
	}
	// now modify the original using the existing overlay
	err = existingOverlayDocument.ApplyTo(&orig)
	if err != nil {
		return "", fmt.Errorf("failed to apply existing overlay: %w", err)
	}

	newOverlay, err := Compare("example overlay", &orig, target)
	if err != nil {
		return "", fmt.Errorf("failed to compare schemas: %w", err)
	}
	// special case, is there only one action and it targets the same as the last overlayDocument.Actions item entry, we'll just replace it.
	if len(newOverlay.Actions) == 1 && len(existingOverlayDocument.Actions) > 0 && newOverlay.Actions[0].Target == existingOverlayDocument.Actions[len(existingOverlayDocument.Actions)-1].Target {
		existingOverlayDocument.Actions[len(existingOverlayDocument.Actions)-1] = newOverlay.Actions[0]
	} else {
		// Otherwise, we'll just append the new overlay to the existing overlay
		existingOverlayDocument.Actions = append(existingOverlayDocument.Actions, newOverlay.Actions...)
	}

	out, err := yaml.Marshal(existingOverlayDocument)
	if err != nil {
		return "", fmt.Errorf("failed to marshal schema: %w", err)
	}

	return string(out), nil
}

func GetInfo(originalYAML string) (string, error) {
	var orig yaml.Node
	err := yaml.Unmarshal([]byte(originalYAML), &orig)
	if err != nil {
		return "", fmt.Errorf("failed to parse source schema: %w", err)
	}

	titlePath, err := jsonpath.NewPath("$.info.title")
	if err != nil {
		return "", err
	}
	versionPath, err := jsonpath.NewPath("$.info.version")
	if err != nil {
		return "", err
	}
	descriptionPath, err := jsonpath.NewPath("$.info.version")
	if err != nil {
		return "", err
	}
	toString := func(node []*yaml.Node) string {
		if len(node) == 0 {
			return ""
		}
		return node[0].Value
	}

	return `{
    "title": "` + toString(titlePath.Query(&orig)) + `",
    "version": "` + toString(versionPath.Query(&orig)) + `",
    "description": "` + toString(descriptionPath.Query(&orig)) + `"
}`, nil
}
