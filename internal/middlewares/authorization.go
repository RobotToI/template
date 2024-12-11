package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"scm.x5.ru/x5m/go-backend/packages/zlogger"
	"scm.x5.ru/x5m/go-backend/template/internal/api"
	"scm.x5.ru/x5m/go-backend/template/internal/dto"
	"strings"
)

// UserAuthMiddleware извлекает токен из заголовка Authorization и помещает данные пользователя в context
func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию авторизация только на /public, остальные игнорируем
		if !strings.HasPrefix(r.RequestURI, "/public") {
			next.ServeHTTP(w, r.WithContext(r.Context()))
			return
		}

		logger := zlogger.LoadOrCreateFromCtx(r.Context())
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logger.Error().Msg("Missing Authorization header")
			err := fmt.Errorf("missing Authorization header")
			api.SendErrorJSON(w, r, http.StatusBadRequest, err, "", api.ErrUnauthorization)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			logger.Error().Msg("Invalid token format")
			err := fmt.Errorf("invalid token format")
			api.SendErrorJSON(w, r, http.StatusBadRequest, err, "", api.ErrUnauthorization)
			return
		}

		// Разбираем токен без проверки подписи
		claims := &jwt.MapClaims{}
		parsedToken, _, err := new(jwt.Parser).ParseUnverified(token, claims)
		_ = parsedToken
		if err != nil {
			logger.Error().Msg("Invalid token format")
			err := fmt.Errorf("invalid token format")
			api.SendErrorJSON(w, r, http.StatusBadRequest, err, "", api.ErrUnauthorization)
			return
		}

		// Преобразуем claims в JSON, а затем в UserData
		claimsJSON, err := json.Marshal(claims)
		if err != nil {
			logger.Error().Msg("Failed to decode token claims")
			err := fmt.Errorf("failed to decode token claims")
			api.SendErrorJSON(w, r, http.StatusBadRequest, err, "", api.ErrUnauthorization)
			return
		}

		userData := &dto.AuthUserDTO{}
		err = json.Unmarshal(claimsJSON, userData)
		if err != nil {
			logger.Error().Msg("Invalid token data")
			err := fmt.Errorf("invalid token data")
			api.SendErrorJSON(w, r, http.StatusBadRequest, err, "", api.ErrUnauthorization)
			return
		}

		// Добавляем данные пользователя в контекст
		ctx := context.WithValue(r.Context(), dto.UserContextKey, userData)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
