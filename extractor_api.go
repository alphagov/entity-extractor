package main

import (
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

	// FIXME - add endpoint to perform extraction

	return mux
}
