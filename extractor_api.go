package main

import (
	"encoding/json"
	"net/http"
)

func NewExtractorAPI(extractor *Extractor) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.Header().Set("Allow", "GET")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/extract", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.Header().Set("Allow", "POST")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		postBody := make([]byte, 10000)
		r.Body.Read(postBody)
		matchedTermIds := extractor.Extract(string(postBody))
		marshalled, _ := json.Marshal(matchedTermIds)

		w.Write(marshalled)
	})

	return mux
}
