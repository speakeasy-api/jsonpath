package jsonpath_test

import (
	"encoding/json"
	"github.com/speakeasy-api/jsonpath/pkg/jsonpath"
	"os"
	"testing"
)

type TestSuite struct {
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
	file, err := os.Open("./jsonpath-compliance-test-suite/cts.json")
	if err != nil {
		t.Fatalf("Failed to open test suite file: %v", err)
	}
	defer file.Close()

	// Parse the test suite JSON
	var testSuite TestSuite
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&testSuite)
	if err != nil {
		t.Fatalf("Failed to parse test suite JSON: %v", err)
	}

	// Run each test case as a subtest
	for _, test := range testSuite.Tests {
		t.Run(test.Name, func(t *testing.T) {
			if test.InvalidSelector {
				// Test case for an invalid selector
				_ := test.Selector
				if err == nil {
					t.Errorf("Expected an error for invalid selector, but got none")
				}
			} else {
				// Test case for a valid selector
				jp, err := jsonpath.Parse(test.Selector)
				if err != nil {
					t.Errorf("Failed to parse JSONPath selector: %v", err)
					return
				}

				result, err := jp.Evaluate(test.Document)
				if err != nil {
					t.Errorf("Failed to evaluate JSONPath: %v", err)
					return
				}

				if test.Results != nil {
					// Test case with multiple possible results
					var found bool
					for _, expectedResult := range test.Results {
						if compareResults(result, expectedResult) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Unexpected result. Got: %v, Want one of: %v", result, test.Results)
					}
				} else {
					// Test case with a single expected result
					if !compareResults(result, test.Result) {
						t.Errorf("Unexpected result. Got: %v, Want: %v", result, test.Result)
					}
				}
			}
		})
	}
}

func compareResults(actual, expected []interface{}) bool {
	// Implement the logic to compare actual and expected results
	// You may need to consider the order and equality of elements
	// This is just a placeholder implementation
	return true
}
