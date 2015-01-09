package main

import (
	"strings"
	"testing"
)

type ParseExample struct {
	raw   string
	terms []string
	id    string
	err   error
}

var parseExamples = []ParseExample{
	{
		"{}",
		[]string{},
		"",
		nil,
	},
	{
		`{"terms":["Government digital service","GDS"],"id":"1"}`,
		[]string{"Government digital service", "GDS"},
		"1",
		nil,
	},
	{
		`{"`,
		[]string{},
		"",
		nil,
	},
}

func TestParseEntityFromJson(t *testing.T) {
	for i, example := range parseExamples {
		entity, err := EntityFromJSON(example.raw)

		if err != example.err {
			t.Error("unexpected error in example", i, "error:", err, "expected:", example.err)
		}

		joinedTerms := strings.Join(entity.terms, ",")
		joinedExpectedTerms := strings.Join(example.terms, ",")
		if joinedTerms != joinedExpectedTerms {
			t.Error("in example", i, "entity.terms != example.terms", joinedTerms, joinedExpectedTerms)
		}
		if entity.id != example.id {
			t.Error("in example", i, "entity.id != example.id")
		}
	}
}
