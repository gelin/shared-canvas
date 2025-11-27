package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"service": "shared-canvas",
		"time":    time.Now().Format(time.RFC3339Nano),
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func ImageHandler(w http.ResponseWriter, _ *http.Request) {
	imgHolder.WriteImagePNG(w)
}
