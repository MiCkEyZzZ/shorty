package common

import "errors"

var (
	ErrMethodNotAllowed = errors.New("метод не поддерживается")
	ErrBadRequest       = errors.New("неверный формат URL")
	ErrInvalidRequest   = errors.New("некорректный запрос")

	ErrUserRegistrationFailed = errors.New("ошибка регистрации пользователя")
	ErrAuthFailed             = errors.New("ошибка при авторизации")

	ErrInvalidID    = errors.New("ошибка парсинга ID")
	ErrNotFound     = errors.New("ошибка поиска")
	ErrUpdateFailed = errors.New("ошибка при обновлении")

	ErrInvalidParam     = errors.New("неверный параметр")
	ErrRequestBodyParse = errors.New("не удалось обработать тело запроса")
)
