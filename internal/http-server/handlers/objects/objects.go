package objects

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	resp "github.com/Zalozhnyy/sbt_k8s/internal/api/response"
	"github.com/Zalozhnyy/sbt_k8s/internal/storage"
	"github.com/Zalozhnyy/sbt_k8s/lib/sl"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type Request struct {
	Data string `json:"data"`
}

type JsonDto struct {
	RrawJson    string
	ExpiredTime time.Time
}

type JsonSaver interface {
	Save(id string, obj JsonDto) error
}

type JsonGetter interface {
	Get(id string) (JsonDto, error)
}

func NewGetter(log *slog.Logger, getter JsonGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.objects.NewGetter"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		log.Info("GET BY ID", slog.String("id", id))

		js, err := getter.Get(id)
		if errors.Is(err, storage.ErrDoNotExists) || errors.Is(err, storage.ErrExpired) {
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		render.JSON(w, r, resp.OkWithData(js.RrawJson))

	}
}

func NewSaver(log *slog.Logger, saver JsonSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "handlers.objects.NewSaver"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		dto := JsonDto{}

		t_raw := r.Header.Get("Expires")
		if t_raw != "" {
			parsedTime, err := time.Parse(time.RFC3339, t_raw)
			if err != nil {
				dto.ExpiredTime = parsedTime
			}
		}

		var req Request
		defer r.Body.Close()
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// TODO is json validation

		dto.RrawJson = req.Data
		id := chi.URLParam(r, "id")

		if err = saver.Save(id, dto); err != nil {
			render.JSON(w, r, resp.Error("failed to save"))
		}

		render.JSON(w, r, resp.OK())
	}
}
