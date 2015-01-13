package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"path"
	"runtime"
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

func fixturePath(fixtureFile string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "data", fixtureFile)
}

func TestLoadEntities(t *testing.T) {
	config := Config{
		extractAddress: ":9999",
		entitiesPath:   fixturePath("entities.jsonl"),
		logPath:        "STDERR",
	}

	extractor := NewExtractor(&config)
	extractor.LoadEntities()

	expectedTerms := [...]string{
		"Government digital service",
		"GDS",
		"Ministry of Justice",
		"MoJ",
	}
	for i, expectedTerm := range expectedTerms {
		assert.Equal(t, expectedTerm, extractor.entities.terms[i])
	}
	expectedIds := [...]string{
		"1",
		"1",
		"2",
		"2",
	}
	for i, expectedId := range expectedIds {
		assert.Equal(t, expectedId, extractor.entities.termOffsetToEntity[i].id)
	}
}
