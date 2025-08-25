package app

import (
	httpserver "AliceSkills/internal/http"
	"AliceSkills/internal/skills"
	"AliceSkills/internal/skills/gpt"
	"AliceSkills/pkg/config"
	"context"
	"errors"
	"fmt"
	"net/http"
)

type App struct {
	cfg config.Config
	srv *http.Server
	reg skills.Registry
}

func New(cfg config.Config) (*App, error) {
	reg := skills.NewRegistry()
	_ = reg.Register(gptquestion.New())

	router := httpserver.NewRouter(reg)

	srv := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           router,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.WriteTimeout,
	}

	return &App{
		cfg: cfg,
		srv: srv,
		reg: reg,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	stoppedCh := make(chan struct{})

	go func() {
		defer close(stoppedCh)
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	fmt.Printf("listening on %s (env=%s)\n", a.cfg.Addr(), a.cfg.Env)

	select {
	case <-ctx.Done():
	case err := <-errCh:
		return err
	}

	fmt.Printf("shutting down... timeout=%s\n", a.cfg.ShutdownTimeout)

	shCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
	defer cancel()

	if err := a.srv.Shutdown(shCtx); err != nil {
		return fmt.Errorf("server shutdown: %w", err)
	}
	<-stoppedCh
	fmt.Println("server stopped")
	return nil
}
