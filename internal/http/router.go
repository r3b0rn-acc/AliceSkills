package httpserver

import (
	"AliceSkills/internal/http/middleware"
	"AliceSkills/internal/skills"
	"net/http"
)

func NewRouter(reg skills.Registry) http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("/alice/", aliceWebhookHandler(reg))

	return middleware.LoggerMiddleware(
		middleware.RecoverMiddleware(
			middleware.RequestIDMiddleware(router),
		),
	)
}
