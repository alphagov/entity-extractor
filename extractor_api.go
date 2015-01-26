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

		fmt.Printf("Reading post body\n")
		fmt.Printf("%v\n", r.Header)

		postBody := make([]byte, 100000)
		n, err := r.Body.Read(postBody)

		if err != nil {
			fmt.Printf("Error reading post body: %v\n", err)
			return
		}
		fmt.Printf("Read %v bytes into post body\n", n)

		matchedTermIds := extractor.Extract(string(postBody))
		if matchedTermIds == nil {
			fmt.Printf("Matched term ids was nil\n")
			errlog.Log(map[string]interface{}{
				"message":  "matchedTermIds was nil",
				"document": string(postBody),
				"status":   500,
			})

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var marshalled []byte
		marshalled, err = json.Marshal(matchedTermIds)
		if err != nil {
			fmt.Printf("Faild to marshal due to: %v", err)
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
