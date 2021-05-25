package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondWithUriError(w http.ResponseWriter, r *http.Request, err error) {
	respondWithBadRequestError(w, fmt.Sprintf("Invalid request: %s -> %s", r.RequestURI, err.Error()))
}

func respondWithBadURI(w http.ResponseWriter, r *http.Request) {
	respondWithBadRequestError(w, fmt.Sprintf("Invalid request: %s", r.RequestURI))
}

func respondWithBadRequestError(w http.ResponseWriter, message string) {
	respondWithError(w, http.StatusBadRequest, message)
}

func respondWithError(w http.ResponseWriter, httpStatusCode int, message string) {
	msg := map[string]interface{}{
		"error": message,
	}
	respondWithErrorMap(w, httpStatusCode, msg)
}

func respondWithErrorMap(w http.ResponseWriter, httpStatusCode int, msg map[string]interface{}) {
	jsonBytes, _ := json.Marshal(msg)
	respondWithJSON(w, httpStatusCode, jsonBytes)
}

func respondWithJSON(w http.ResponseWriter, httpStatusCode int, jsonBytes []byte) {
	// https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	_, err := w.Write(jsonBytes)
	if err != nil {
		print(err)
	}
}
