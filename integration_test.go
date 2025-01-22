package jsonpath_test

import (
	"encoding/json"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath"
	"github.com/speakeasy-api/jsonpath/pkg/yaml"
	"github.com/stretchr/testify/require"
	"os"
	"slices"
	"strings"
	"testing"
)

type FullTestSuite struct {
	Description string `json:"description"`
	Tests       []Test `json:"tests"`
}
type Test struct {
	Name            string          `json:"name"`
	Selector        string          `json:"selector"`
	Document        interface{}     `json:"document"`
	Result          []interface{}   `json:"result"`
	Results         [][]interface{} `json:"results"`
	InvalidSelector bool            `json:"invalid_selector"`
	Tags            []string        `json:"tags"`
}

func TestJSONPathComplianceTestSuite(t *testing.T) {
	// Read the test suite JSON file
	file, err := os.ReadFile("./jsonpath-compliance-test-suite/cts.json")
	require.NoError(t, err, "Failed to read test suite file")
	// alter the file to delete any unicode tests: these break the yaml library we use..
	var testSuite FullTestSuite
	json.Unmarshal(file, &testSuite)
	for i := 0; i < len(testSuite.Tests); i++ {
		// if Tags contains "unicode", delete it
		// (they break the yaml parser)
		shouldDelete := slices.Contains(testSuite.Tests[i].Tags, "unicode")
		// delete some other new line / unicode tests -- these also break the yaml parser
		shouldDelete = shouldDelete || strings.Contains(testSuite.Tests[i].Name, "line feed")
		shouldDelete = shouldDelete || strings.Contains(testSuite.Tests[i].Name, "carriage return")
		shouldDelete = shouldDelete || strings.Contains(testSuite.Tests[i].Name, "u2028")
		shouldDelete = shouldDelete || strings.Contains(testSuite.Tests[i].Name, "u2029")
		if shouldDelete {
			testSuite.Tests = append(testSuite.Tests[:i], testSuite.Tests[i+1:]...)
			i--
		}
	}

	// Run each test case as a subtest
	for _, test := range testSuite.Tests {
		t.Run(test.Name, func(t *testing.T) {
			// Test case for a valid selector
			jp, err := jsonpath.NewPath(test.Selector)
			if test.InvalidSelector {
				require.Error(t, err, "Expected an error for invalid selector, but got none for path", jp.String())
				return
			} else {
				require.NoError(t, err, "Failed to parse JSONPath selector", jp.String())
			}

			// expect stability of ToString()
			stringified := jp.String()
			recursive, err := jsonpath.NewPath(stringified)
			require.NoError(t, err, "Failed to parse recursive JSONPath selector. expected=%s got=%s", test.Selector, jp.String())
			require.Equal(t, stringified, recursive.String(), "JSONPath selector does not match test case")
			// interface{} to yaml.Node
			toYAML := func(i interface{}) *yaml.Node {
				o, err := yaml.Marshal(i)
				require.NoError(t, err, "Failed to marshal interface to yaml")
				n := new(yaml.Node)
				err = yaml.Unmarshal(o, n)
				require.NoError(t, err, "Failed to unmarshal yaml to yaml.Node")
				// unwrap the document node
				if n.Kind == yaml.DocumentNode && len(n.Content) == 1 {
					n = n.Content[0]
				}
				return n
			}

			result := jp.Query(toYAML(test.Document))

			if test.Results != nil {
				expectedResults := make([][]*yaml.Node, 0)
				for _, expectedResult := range test.Results {
					expected := make([]*yaml.Node, 0)
					for _, expectedResult := range expectedResult {
						expected = append(expected, toYAML(expectedResult))
					}
					expectedResults = append(expectedResults, expected)
				}

				// Test case with multiple possible results
				var found bool
				for i, _ := range test.Results {
					if match, msg := compareResults(result, expectedResults[i]); match {
						found = true
						break
					} else {
						t.Log(msg)
					}
				}
				if !found {
					t.Errorf("Unexpected result. Got: %v, Want one of: %v", result, test.Results)
				}
			} else {
				expectedResult := make([]*yaml.Node, 0)
				for _, res := range test.Result {
					expectedResult = append(expectedResult, toYAML(res))
				}
				// Test case with a single expected result
				if match, msg := compareResults(result, expectedResult); !match {
					t.Error(msg)
				}
			}
		})
	}
}

func compareResults(actual, expected []*yaml.Node) (bool, string) {
	actualStr, err := yaml.Marshal(actual)
	if err != nil {
		return false, "Failed to serialize actual result: " + err.Error()
	}

	expectedStr, err := yaml.Marshal(expected)
	if err != nil {
		return false, "Failed to serialize expected result: " + err.Error()
	}

	if string(actualStr) == string(expectedStr) {
		return true, ""
	}

	// Generate a nice diff string
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(expectedStr)),
		B:        difflib.SplitLines(string(actualStr)),
		FromFile: "Expected",
		ToFile:   "Actual",
		Context:  3,
	}
	diffStr, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return false, "Failed to generate diff: " + err.Error()
	}

	return false, diffStr
}
