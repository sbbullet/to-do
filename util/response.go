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
