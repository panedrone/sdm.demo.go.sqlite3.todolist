package main

import (
	"encoding/json"
	"net/http"
)

func respondWithBabRequestError(w http.ResponseWriter, message string) {
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
