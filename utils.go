package main

import (
	"encoding/json"
	"net/http"
)

func encodeAsJson(toEncode any, w http.ResponseWriter, statusCode int) *httpError {
	// ResponseWriter requires that you set all headers before calling WriteHeader
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(toEncode); err != nil {
		return &httpError{err, http.StatusInternalServerError}
	}
	return nil
}
