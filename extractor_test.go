package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func fixturePath(fixtureFile string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "data", fixtureFile)
}

func exampleConfig() *Config {
	return &Config{
		extractAddress:     ":3096",
		dbConnectionString: "host=/var/run/postgresql dbname=entity-extractor_test sslmode=disable",
		logPath:            "STDERR",
	}
}

func TestLoadEntities(t *testing.T) {
	config := exampleConfig()
	extractor := NewExtractor(config)
	err := extractor.LoadEntities()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	expectedEntities := [...]Entity{

		Entity{"1",
			[]Representation{
				Representation{
					"Government digital service",
					Tokens{[]string{"government", "digital", "service"}},
				},
				Representation{
					"GDS",
					Tokens{[]string{"gds"}},
				},
			},
		},

		Entity{"2",
			[]Representation{
				Representation{
					"Ministry of Justice",
					Tokens{[]string{"ministry", "of", "justice"}},
				},
				Representation{
					"MoJ",
					Tokens{[]string{"moj"}},
				},
			},
		},
	}

	require.Equal(t, len(expectedEntities), len(extractor.entities.entities))
	for i, expectedEntity := range expectedEntities {
		e := extractor.entities.entities[i]
		assert.Equal(t, expectedEntity, e)
	}
}

type ExtractionExample struct {
	comment         string
	document        string
	expectedTermIds []string
}

var extractionExamples = []ExtractionExample{
	{
		"a document matching a single word term",
		"This document mentions GDS but it doesn't mention the Ministry of J...",
		[]string{"1"},
	},
	{
		"a document matching a multi-word term",
		"Government digital service",
		[]string{"1"},
	},
	{
		"terms are matched case sensitively",
		"gds",
		[]string{},
	},
}

func TestExtract(t *testing.T) {
	config := exampleConfig()
	extractor := NewExtractor(config)
	err := extractor.LoadEntities()
	if err != nil {
		fmt.Println("ERROR: %v", err)
		t.FailNow()
	}

	for _, example := range extractionExamples {
		matchedTermIds := extractor.Extract(example.document)
		assert.Equal(t, example.expectedTermIds, matchedTermIds, example.comment)
	}
}
