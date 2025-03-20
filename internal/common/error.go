package common

import "errors"

var (
	// Общие ошибки
	ErrInvalidID        = errors.New("некорректный идентификатор")
	ErrInvalidLimit     = errors.New("Неверный лимит ввода")
	ErrInvalidOffset    = errors.New("Неверное смещение ввода")
	ErrURLNotFound      = errors.New("URL-адрес не найден")
	ErrInvalidParam     = errors.New("неверный параметр")
	ErrRequestBodyParse = errors.New("не удалось обработать тело запроса")
	ErrMethodNotAllowed = errors.New("метод не поддерживается")
	ErrInvalidRequest   = errors.New("некорректный запрос")
	ErrNotFound         = errors.New("ошибка поиска")

	// Ошибки авторизации
	ErrBadRequest             = errors.New("неверный формат URL")
	ErrUserRegistrationFailed = errors.New("ошибка регистрации пользователя")
	ErrAuthFailed             = errors.New("ошибка при авторизации")
	ErrUnauthorized           = errors.New("ощибка нет авторизации")
	ErrInvalidToken           = errors.New("ощибка не валидный токен")
	ErrForbidden              = errors.New("ощибка доступ запрещён")
	UserContextKey            = errors.New("ощибка используй ключ")

	// Ошибки пользователя.
	ErrorGetUsers       = errors.New("не удалось получить список пользователей")
	ErrUserNotFound     = errors.New("пользователь не найден")
	ErrUserUpdateFailed = errors.New("не удалось обновить пользователя")
	ErrUserDeleteFailed = errors.New("не удалось удалить пользователя")

	// Ошибки ссылок.
	ErrLinkCreateUR         = errors.New("не удалось создать сокращённый URL")
	ErrLinkNotFound         = errors.New("ссылка с таким идентификатором не найдена")
	ErrLinkHashNotProvided  = errors.New("hash не указан")
	ErrLinkUpdateLinkFailed = errors.New("ошибка при обновлении ссылки")
	ErrLinkDeleteFailed     = errors.New("ошибка при удаления ссылки")
	ErrLinkBlockFailed      = errors.New("ошибка при попытке заблокировать ссылку")
	ErrUnBlockFailed        = errors.New("ошибка при попытке разблокировать ссылку")

	ErrClickWriteFailed = errors.New("ошибка записи при клике")
)
