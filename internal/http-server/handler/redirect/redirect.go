package redirect

import (
	resp "awesomeProject/internal/lib/api/response"
	"awesomeProject/internal/storage"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type URLGetter interface {
	GetURL(ShortUrl string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		shortUrl := chi.URLParam(r, "shortUrl")
		if shortUrl == "" {
			log.Info("shortUrl is empty")

			render.JSON(w, r, resp.Error("invalid request"))

			return
		}

		resURL, err := urlGetter.GetURL(shortUrl)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "shortUrl", shortUrl)

			render.JSON(w, r, resp.Error("short-url not found"))

			return
		}
		if err != nil {
			log.Error("failed to get url", slog.Any("error", err))

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("got url", slog.String("url", resURL))

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
