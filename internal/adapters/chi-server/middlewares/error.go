package middlewares

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"

	"github.com/diogorodriguesc/boilerplate-go/internal/adapters/chi-server/handlers"
	applicationerrors "github.com/diogorodriguesc/boilerplate-go/internal/application/errors"
)

type statusCodeRecorder struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (rec *statusCodeRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *statusCodeRecorder) Write(b []byte) (int, error) {
	if rec.statusCode == 0 {
		rec.WriteHeader(http.StatusOK)
	}
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}

type handlerErrCtxKey struct{}

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := fn(w, r)
	if err == nil {
		return
	}
	if errPtr, ok := r.Context().Value(handlerErrCtxKey{}).(*error); ok {
		*errPtr = err
	}
}

func statusCodeForError(err error) int {
	switch {
	case errors.Is(err, applicationerrors.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, applicationerrors.ErrDuplicateEntry):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func SetErrorMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var handlerErr error
			ctx := context.WithValue(r.Context(), handlerErrCtxKey{}, &handlerErr)

			recorder := &statusCodeRecorder{ResponseWriter: w}
			next.ServeHTTP(recorder, r.WithContext(ctx))

			if handlerErr != nil {
				recorder.WriteHeader(statusCodeForError(handlerErr))
				render.JSON(recorder, r, handlers.MapErrorIntoErrorResponse(handlerErr))
			}

			if recorder.statusCode < http.StatusBadRequest || recorder.statusCode == http.StatusNotFound{
				return
			}

			log := log.Error().
				Int("status_code", recorder.statusCode).
				Str("method", r.Method).
				Str("path", r.URL.Path)

			var errResp handlers.ErrorResponse
			if err := json.Unmarshal(recorder.body.Bytes(), &errResp); err == nil && errResp.Error.Message != "" {
				log.Str("error", errResp.Error.Message)
			}

			log.Msg("request failed")
		})
	}
}
