package logger

import (
	"github.com/rs/zerolog"
	"net/http"
	"scm.x5.ru/x5m/go-backend/packages/zlogger"
	"scm.x5.ru/x5m/go-backend/template/internal/auth"
)

// LogWithUserData добавляет данные пользователя в логгер и возвращает обновленный логгер
func LogWithUserData(r *http.Request) *zerolog.Logger {
	originalLogger := zlogger.LoadOrCreateFromCtx(r.Context())
	logger := originalLogger.With().Logger() // Создаем новый экземпляр логгера

	userData, ok := auth.GetAuthDataFromCtx(r.Context())
	if ok {
		logger = logger.With().
			Str("handler-name", r.Method+" "+r.URL.String()).
			Str("session_state", userData.SessionState).
			Str("cip_id", userData.CipID).
			Str("x5id", userData.X5ID).
			Logger()
	}

	return &logger
}
