package main

import (
	"encoding/json"
	"net/http"
)

func response(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
