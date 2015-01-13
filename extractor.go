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

	//extr.terms = terms
	//extr.matcher = ahocorasick.NewStringMatcher(terms)
	//logInfo("Loaded", len(extr.terms), "entities")

	extr.entities = entities
	return nil
}
