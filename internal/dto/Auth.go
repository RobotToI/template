package dto

import (
	"scm.x5.ru/x5m/go-backend/template/pkg"
)

// UserContextKey Ключ для хранения данных пользователя в контексте
const UserContextKey pkg.ContextKey = "AuthUserDTOkey"

// AuthUserDTO Структура для хранения данных авторизированного пользователя
type AuthUserDTO struct {
	Azp               string `json:"azp"`
	SessionState      string `json:"session_state"`
	PreferredUsername string `json:"preferred_username"`
	Subscription      string `json:"subscription"`
	CipID             string `json:"cip_id"`
	X5ID              string `json:"x5id"`
}
