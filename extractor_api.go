package main

import (
	"encoding/json"
	"fmt"
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
		if matchedTermIds == nil {
			errlog.Log(map[string]interface{}{
				"message":  "matchedTermIds was nil",
				"document": string(postBody),
				"status":   500,
			})

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		marshalled, err := json.Marshal(matchedTermIds)
		if err != nil {
			errlog.Log(map[string]interface{}{
				"message": "Failed to marshal matched terms to Json",
				"error":   fmt.Sprintf("%v", err),
				"status":  500,
			})

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(marshalled)
	})

	return mux
}
