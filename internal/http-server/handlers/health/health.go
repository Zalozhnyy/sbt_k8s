package health

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.objects.NewSaver"

		log := log.With(
			slog.String("op", op),
		)

		log.Info("Health check")
		render.Status(r, http.StatusOK)
		render.PlainText(w, r, "ok")
	}
}
