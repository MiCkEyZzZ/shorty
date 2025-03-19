package common

import "errors"

var (
	ErrMethodNotAllowed = errors.New("метод не поддерживается")
	ErrBadRequest       = errors.New("некорректный запрос")

	ErrUserRegistrationFailed = errors.New("ошибка регистрации пользователя")
	ErrAuthFailed             = errors.New("ошибка при авторизации")

	ErrInvalidParam     = errors.New("неверный параметр")
	ErrRequestBodyParse = errors.New("не удалось обработать тело запроса")
)
