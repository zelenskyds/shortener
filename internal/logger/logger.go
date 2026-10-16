package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *responseWriter) WriteHeader(status int) {
	if w.status != 0 {
		return
	}

	w.ResponseWriter.WriteHeader(status)
	w.status = status
}

func (w *responseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}

	addedSize, err := w.ResponseWriter.Write(data)
	w.size += addedSize
	return addedSize, err
}

var Log = zap.NewNop()

func Initialize(logLevel string) error {
	level, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		return err
	}

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level = level

	log, err := loggerConfig.Build()
	if err != nil {
		return err
	}

	Log = log
	return nil
}

func HTTPHandlerLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := &responseWriter{ResponseWriter: w}

		start := time.Now()
		handler.ServeHTTP(ww, r)

		duration := time.Since(start)

		if ww.status == 0 {
			ww.status = http.StatusOK
		}

		Log.Info("HTTP request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Duration("duration", duration),
			zap.Int("status", ww.status),
			zap.Int("size", ww.size),
		)
	})
}
