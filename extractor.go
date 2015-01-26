package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Entities struct {
	entities []Entity
}

type Extractor struct {
	cfg      *Config
	entities Entities
}

func NewExtractor(cfg *Config) *Extractor {
	return &Extractor{cfg: cfg}
}

func (entities *Entities) addEntity(id string, termsJson string) error {
	entity, err := NewEntity(id, termsJson)
	if err != nil {
		return err
	}

	entities.entities = append(entities.entities, *entity)

	return nil
}

func (extr *Extractor) LoadEntities() error {
	return extr.loadEntitiesFromDatabase(extr.cfg.dbConnectionString)
}

func (extr *Extractor) loadEntitiesFromDatabase(dbConnectionString string) error {
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT id, terms FROM entities")
	if err != nil {
		return err
	}
	defer rows.Close()

	var entities Entities

	for rows.Next() {
		var id, terms string
		if err := rows.Scan(&id, &terms); err != nil {
			return err
		}

		if err := entities.addEntity(id, terms); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	extr.entities = entities

	return nil
}

func (extr *Extractor) Extract(text string) []string {
	matchingEntityIds := make([]string, 0, 1000)

	return matchingEntityIds
}
