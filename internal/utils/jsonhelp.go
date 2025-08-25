package utils

import (
	"encoding/json"
	"net/http"
)

func JsonHelper(w http.ResponseWriter, status int, payload any) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte("encode error"))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, err = w.Write(payloadBytes)
}
