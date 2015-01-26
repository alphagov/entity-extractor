package main

import (
	"bytes"
	"github.com/blevesearch/segment"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

type Tokens struct {
	words []string
}

func (tokens *Tokens) add(word string) {
	tokens.words = append(tokens.words, word)
}

// Tokenise a representation of an entity into words
func TokeniseEntityRepresentation(text string) (*Tokens, error) {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}

	normaliser := transform.Chain(
		norm.NFD,
		transform.RemoveFunc(isMn),
		norm.NFKC,
	)

	reader := transform.NewReader(
		bytes.NewReader([]byte(strings.ToLower(text))),
		normaliser,
	)

	segmenter := segment.NewWordSegmenter(reader)
	result := Tokens{}
	for segmenter.Segment() {
		if segmenter.Type() != 0 {
			result.add(segmenter.Text())
		}
	}
	if err := segmenter.Err(); err != nil {
		return nil, err
	}

	return &result, nil
}
