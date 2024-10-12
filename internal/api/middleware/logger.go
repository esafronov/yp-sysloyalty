package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// RequestLogger middleware for logging incomming requests
func RequestLogger(Log *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			timeStamp := time.Now()
			rd := &ResponseData{
				size:   0,
				status: 0,
			}
			wr := &LoggerResponseWriter{
				data: rd,
				rw:   w,
			}
			h.ServeHTTP(wr, r)
			duration := time.Since(timeStamp).Milliseconds()
			Log.Info("request",
				zap.String("method", r.Method),
				zap.String("URI", r.RequestURI),
				zap.Int64("duration", duration),
				zap.Int("status", rd.status),
				zap.Int("size", rd.size))
		})
	}
}

type ResponseData struct {
	size   int
	status int
}

type LoggerResponseWriter struct {
	rw   http.ResponseWriter
	data *ResponseData
}

func (w *LoggerResponseWriter) Write(b []byte) (int, error) {
	len, err := w.rw.Write(b)
	w.data.size += len
	return len, err
}

func (w *LoggerResponseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *LoggerResponseWriter) WriteHeader(statusCode int) {
	w.rw.WriteHeader(statusCode)
	w.data.status = statusCode
}
