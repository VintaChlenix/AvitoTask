package handlers

import (
	"encoding/json"
	"net/http"
)

func RenderJSON(w http.ResponseWriter, payload interface{}) {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonPayload)
}

func RenderOK(w http.ResponseWriter) {
	http.Error(w, "200 OK", http.StatusOK)
}

func RenderBadRequest(w http.ResponseWriter, error error) {
	http.Error(w, error.Error(), http.StatusBadRequest)
}

func RenderInternalError(w http.ResponseWriter, error error) {
	http.Error(w, error.Error(), http.StatusInternalServerError)
}
