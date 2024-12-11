package middlewares

import (
	"bytes"
	"net/http"
	"scm.x5.ru/x5m/go-backend/template/internal/logger"
	"strconv"
)

// ResponseRecorder это обёртка над http.ResponseWriter, которая сохраняет ответ
type ResponseRecorder struct {
	http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode int
}

// NewResponseRecorder создаёт новый ResponseRecorder
func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		Body:           new(bytes.Buffer),
		StatusCode:     http.StatusOK,
	}
}

// WriteHeader перехватывает запись статуса
func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

// Write перехватывает запись тела ответа
func (rec *ResponseRecorder) Write(b []byte) (int, error) {
	rec.Body.Write(b)                  // записываем в наш буфер
	return rec.ResponseWriter.Write(b) // записываем в оригинальный ResponseWriter
}

// ResponseLoggerMiddleware логирует все ответы
func ResponseLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Создаем буфер для перехвата ответа
		rec := NewResponseRecorder(w)
		defer func() {
			// Логируем ответ
			originalLogger := logger.LogWithUserData(r)
			responseLogger := originalLogger.With().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", rec.StatusCode).
				Logger()
			responseLogger.Info().Msg(r.Method + " " + r.URL.String() + " " + strconv.Itoa(rec.StatusCode))

		}()

		next.ServeHTTP(rec, r)
	})
}
