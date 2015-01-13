package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ParseExample struct {
	raw           string
	expectedTerms []string
	expectedId    string
	expectedError error
}

var parseExamples = []ParseExample{
	{
		"{}",
		make([]string, 0),
		"",
		nil,
	},
	{
		`{"terms":["Government digital service","GDS"],"id":"1"}`,
		[]string{"Government digital service", "GDS"},
		"1",
		nil,
	},
}

func TestParseEntityFromJson(t *testing.T) {
	for i, example := range parseExamples {
		actual, err := EntityFromJSON(example.raw)

		assert.Equal(t, example.expectedError, err, fmt.Sprint("unexpected error in example ", i))
		assert.Equal(t, example.expectedTerms, actual.Terms, fmt.Sprint("terms differ in example ", i))
		assert.Equal(t, example.expectedId, actual.Id, fmt.Sprint("ids differ in example ", i))
	}
}
