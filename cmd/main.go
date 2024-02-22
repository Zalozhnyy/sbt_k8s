package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zalozhnyy/sbt_k8s/internal/http-server/handlers/health"
	"github.com/Zalozhnyy/sbt_k8s/internal/http-server/handlers/objects"
	mylogger "github.com/Zalozhnyy/sbt_k8s/internal/http-server/middleware/myLogger"
	mapstorage "github.com/Zalozhnyy/sbt_k8s/internal/storage/map_storage"
	"github.com/Zalozhnyy/sbt_k8s/lib/sl"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	log := setupLogger()

	storage := mapstorage.New()
	go mapstorage.MapCleaner(ctx, log, storage)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mylogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/objects", func(r chi.Router) {
		r.Put("/{id}", objects.NewSaver(log, storage))
		r.Get("/{id}", objects.NewGetter(log, storage))
	})

	router.Route("/probes", func(r chi.Router) {
		r.Get("/liveness", health.New(log))
		r.Get("/readiness", health.New(log))
	})

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      router,
		ReadTimeout:  time.Second * 4,
		WriteTimeout: time.Second * 4,
		IdleTimeout:  time.Second * 30,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server", sl.Err(err))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("server started")

	<-done
	log.Info("stopping server")
	cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))
		return
	}

	log.Info("server stopped")

}

func setupLogger() *slog.Logger {
	var log *slog.Logger

	log = slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	return log
}
