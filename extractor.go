package main

import (
	"bufio"
	"encoding/json"
	"github.com/cloudflare/ahocorasick"
	"os"
)

// An entity (which has an ID and several representative strings)
type Entity struct {
	id string
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

type ParsedEntity struct {
	Terms []string
	Id    string
}

func NewExtractor(cfg *Config) *Extractor {
	return &Extractor{cfg: cfg}
}

func (self *Entities) addEntity(parsedEntity ParsedEntity) {
	entity := Entity{parsedEntity.Id}
	for _, term := range parsedEntity.Terms {
		self.terms = append(self.terms, term)
		self.termOffsetToEntity = append(self.termOffsetToEntity, &entity)
	}
}

func (extr *Extractor) LoadEntities() error {
	return extr.loadEntitiesFromFile(extr.cfg.entitiesPath)
}

func EntityFromJSON(raw string) (ParsedEntity, error) {
	parsed := ParsedEntity{make([]string, 0), ""}
	err := json.Unmarshal([]byte(raw), &parsed)
	return parsed, err
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

		entity, err := EntityFromJSON(line)
		if err != nil {
			return err
		}
		entities.addEntity(entity)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	extr.matcher = ahocorasick.NewStringMatcher(entities.terms)
	logInfo("Loaded", len(entities.terms), "terms")

	extr.entities = entities
	return nil
}

func (extr *Extractor) Extract(text string) []string {
	matchIndexes := extr.matcher.Match([]byte(text))
	matchingTermIds := make([]string, 0, 1000)
	for _, termIndex := range matchIndexes {
		matchingTermIds = append(matchingTermIds, extr.entities.termOffsetToEntity[termIndex].id)
	}

	return matchingTermIds
}
