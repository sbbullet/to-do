package util

import (
	"encoding/json"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(response)
}

func RespondWithOk(w http.ResponseWriter, data interface{}) {
	RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
	})
}

func RespondWithBadRequest(w http.ResponseWriter, errorMsg string) {
	RespondWithJSON(w, http.StatusBadRequest, map[string]interface{}{
		"success": false,
		"error":   errorMsg,
	})
}

func RespondWithValidationErrors(w http.ResponseWriter, validationErrors interface{}) {
	RespondWithJSON(w, http.StatusUnprocessableEntity, map[string]interface{}{
		"success": false,
		"errors":  validationErrors,
	})
}

func RespondWithInternalServerError(w http.ResponseWriter) {
	RespondWithJSON(w, http.StatusInternalServerError, map[string]interface{}{
		"success": false,
		"error":   "Something went wrong on the server. Please, try after a while",
	})
}
