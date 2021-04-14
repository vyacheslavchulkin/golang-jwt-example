package controllers

import (
	"encoding/json"
	"net/http"
)

func notFound(w http.ResponseWriter) {
	response := responseJson{
		Status:   "error",
		Message:  "Not Found",
		Response: nil,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(response)
}

func invalidPassword(w http.ResponseWriter) {
	response := responseJson{
		Status:   "error",
		Message:  "Invalid uuid or password",
		Response: nil,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(response)
}

func serverError(w http.ResponseWriter) {
	response := responseJson{
		Status:   "error",
		Message:  "Server error",
		Response: nil,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
}

func badRequest(w http.ResponseWriter) {
	response := responseJson{
		Status:   "error",
		Message:  "Bad Request",
		Response: nil,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

func unauthorized(w http.ResponseWriter) {
	response := responseJson{
		Status:   "error",
		Message:  "Unauthorized",
		Response: nil,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(response)
}

func unprocessableEntity(w http.ResponseWriter) {
	response := responseJson{
		Status:   "error",
		Message:  "Unprocessable Entity",
		Response: nil,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(response)
}
