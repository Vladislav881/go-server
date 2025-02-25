package save

import (
	resp "awesomeProject/internal/lib/api/response"
	"awesomeProject/internal/lib/random"
	"awesomeProject/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
)

const shortUrlLength = 10

type Request struct {
	URL      string `json:"url" validate:"required,url"`
	ShortUrl string `json:"short_url,omitempty"`
}

type Response struct {
	resp.Response
	ShortUrl string `json:"short_url,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, shortUrl string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.save.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body", slog.Any("err", err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("invalid request", slog.Any("error", err))

			render.JSON(w, r, map[string]interface{}{
				"error":   "validation failed",
				"details": err.Error(),
			})

			return
		}

		shortUrl := req.ShortUrl
		if shortUrl == "" {
			var err error
			for attempts := 0; attempts < 3; attempts++ {
				shortUrl, err = random.RandomString(shortUrlLength)
				if err == nil {
					break
				}
				log.Error("Error generating RandomString, retrying:", slog.Any("err", err))
			}

			if err != nil {
				log.Error("Failed to generate RandomString after multiple attempts:", slog.Any("err", err))
				return
			}
		}

		id, err := urlSaver.SaveURL(req.URL, shortUrl)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", slog.Any("err", err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		responseOK(w, r, shortUrl)

	}
}
func responseOK(w http.ResponseWriter, r *http.Request, shortUrl string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
		ShortUrl: shortUrl,
	})
}
