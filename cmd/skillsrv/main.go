package main

import (
	httpserver "AliceSkills/internal/http"
	"AliceSkills/pkg/config"
	"net/http"
)

func main() {
	cfg := config.MustLoad()
	router := httpserver.NewRouter()
	err := http.ListenAndServe(cfg.Addr(), router)
	if err != nil {
		return
	}
}
