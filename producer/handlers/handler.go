package handlers

import (
	"bytes"
	"consumer/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log/slog"
	"net/http"
)

type Handler struct {
	service *service.Service
	log     *slog.Logger
}

func NewHandler(service *service.Service, log *slog.Logger) *Handler {
	return &Handler{
		service: service,
		log:     log,
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("POST /messages/", h.LogMiddleware(MetricsMiddleware(http.HandlerFunc(h.AddMessageHandler))))
	mux.Handle("/metrics", h.LogMiddleware(promhttp.Handler()))

	return mux
}

func (h *Handler) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			h.log.Error("Failed to read request body", err)
		}
		r.Body = io.NopCloser(bytes.NewBuffer(body))
		h.log.Info("Received request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("body", string(body)),
		)

		rec := &statusRecorder{ResponseWriter: w, body: new(bytes.Buffer)}
		next.ServeHTTP(rec, r)

		h.log.Info("Handled request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.statusCode),
			slog.String("response", rec.body.String()),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *statusRecorder) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}
	rec.body.Write(b)
	return rec.ResponseWriter.Write(b)
}
