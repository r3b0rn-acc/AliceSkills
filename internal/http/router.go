package httpserver

import (
	"AliceSkills/internal/http/middleware"
	"AliceSkills/internal/utils"
	"net/http"
)

func NewRouter() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		utils.JsonHelper(w, http.StatusOK, map[string]any{"status": "ok"})
	})
	router.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		utils.JsonHelper(w, http.StatusOK, map[string]any{"status": "ready"})
	})

	return middleware.LoggerMiddleware(
		middleware.RecoverMiddleware(
			middleware.RequestIDMiddleware(router),
		),
	)
}
