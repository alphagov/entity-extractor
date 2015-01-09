package main

import (
	"bufio"
	"encoding/json"
	"github.com/cloudflare/ahocorasick"
	"os"
)

// An entity (which has an ID and several representative strings)
type Entity struct {
	id          string
	termOffsets []uint
}

type Entities struct {
	terms              []string
	termOffsetToEntity []*Entity
}

type Extractor struct {
	cfg      *Config
	entities Entities
	matcher  *ahocorasick.Matcher
}

func NewExtractor(cfg *Config) *Extractor {
	return &Extractor{cfg: cfg}
}

func (extr *Extractor) LoadEntities() error {
	return extr.loadEntitiesFromFile(extr.cfg.entitiesPath)
}

type ParsedEntity struct {
	terms []string
	id    string
}

func EntityFromJSON(raw string) (*ParsedEntity, error) {
	var parsed ParsedEntity
	err := json.Unmarshal([]byte(raw), &parsed)
	return &parsed, err
}

func (extr *Extractor) loadEntitiesFromFile(path string) error {
	logInfo("Loading entities from", path)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var entities Entities
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		_, err := EntityFromJSON(line)
		if err != nil {
			return err
		}

		// FIXME - need to parse the JSON entity here,
		// extract each representation as a term and add
		// to the terms array.  Also need to build a map
		// from offset in the terms array to a list of
		// entity IDs.  Also need to handle multiple
		// instances of each entity.
		entities.terms = append(entities.terms, line)
		//entities.termOffsetToEntity = append(entities.termOffsetToEntity, entity)
	}

	//extr.terms = terms
	//extr.matcher = ahocorasick.NewStringMatcher(terms)
	//logInfo("Loaded", len(extr.terms), "entities")

	extr.entities = entities
	return nil
}
