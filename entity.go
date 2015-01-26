package main

import (
	"encoding/json"
)

type Representation struct {
	raw    string
	tokens Tokens
}

func NewRepresentation(text string) (*Representation, error) {
	tokens, err := TokeniseEntityRepresentation(text)
	if err != nil {
		return nil, err
	}
	return &Representation{text, *tokens}, nil
}

// An entity (which has an ID and several representative strings)
type Entity struct {
	id              string
	representations []Representation
}

func (ent *Entity) addRepresentation(text string) error {
	repr, err := NewRepresentation(text)
	if err != nil {
		return err
	}
	ent.representations = append(ent.representations, *repr)
	return nil
}

func NewEntity(id string, termsJson string) (*Entity, error) {
	var terms []string
	err := json.Unmarshal([]byte(termsJson), &terms)
	if err != nil {
		return nil, err
	}

	result := &Entity{id, []Representation{}}
	for _, term := range terms {
		err = result.addRepresentation(term)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
