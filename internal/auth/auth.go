package auth

import (
	"context"
	"scm.x5.ru/x5m/go-backend/template/internal/dto"
)

// GetAuthDataFromCtx Получение данных пользователя из контекста
func GetAuthDataFromCtx(ctx context.Context) (*dto.AuthUserDTO, bool) {
	userData, ok := ctx.Value(dto.UserContextKey).(*dto.AuthUserDTO)
	return userData, ok
}
