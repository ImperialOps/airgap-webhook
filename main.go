package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/validate", validate)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func writeJson(w http.ResponseWriter, code int, v any) error {
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func validate(w http.ResponseWriter, r *http.Request) {
	if err := writeJson(w, http.StatusOK, map[string]string{
		"message": "hello there",
	}); err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{
			"error": "could not write response",
		})
	}
}
