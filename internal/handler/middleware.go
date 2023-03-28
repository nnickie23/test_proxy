package handler

import (
	"context"
	"net/http"

	"github.com/rs/cors"
)


func (a *handler) middleware(h http.Handler) http.Handler {
	h = cors.AllowAll().Handler(h)
	h = a.mwRecovery(h)
	return h
}

func (a *handler) mwRecovery(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cancelCtx, cancel := context.WithCancel(r.Context())
		r = r.WithContext(cancelCtx)
		defer func() {
			if err := recover(); err != nil {
				cancel()
				w.WriteHeader(http.StatusInternalServerError)
				a.logger.Errorw(
					"Panic in http handler",
					err,
					"method", r.Method,
					"path", r.URL,
				)
			}
		}()
		h.ServeHTTP(w, r)
	})
}
