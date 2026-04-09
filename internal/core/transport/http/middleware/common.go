package core_http_middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gibanator/go-server/internal/core/logger"
	core_http_response "github.com/gibanator/go-server/internal/core/transport/http/response"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	requestIDHeader = "X-Request-ID"
)

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "log", l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Recoverer() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			responseHandler := core_http_response.NewHttpResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "panic when HTTP handled")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func Tracer() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.FromContext(ctx)
			rw_with_status := core_http_response.NewResponseWriter(w)

			before := time.Now()

			log.Debug(
				"===> Incoming HTTP Request",
			)

			next.ServeHTTP(rw_with_status, r)

			log.Debug(
				"<=== HTTP Request Handled",
				zap.Int("status code:", rw_with_status.GetStatusCodeOrPanic()),
				zap.Duration("time elapsed:", time.Since(before)),
			)
		})
	}
}
