package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	var reqData struct {
		Message string `json:"message"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		h.log.Error("Failed to decode request body", err)
		return
	}
	err = h.service.SendMessage(r.Context(), reqData.Message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		h.log.Error("Failed to send message", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
