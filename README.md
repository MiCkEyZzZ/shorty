# Коротышка

коротышка — простой и быстрый инструмент для сокращения URL-адресов.

## Структура проекта (Layered Architecture)

shorty
├── cmd
│   └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── models
│   │   ├── link.go
│   │   ├── user.go
│   │   └── stat.go
│   ├── payload
│   │   ├── auth_payload.go
│   │   ├── link_payload.go
│   │   ├── stat_payload.go
│   │   └── user_payload.go
│   ├── repository
│   │   ├── link_repository.go
│   │   ├── user_repository.go
│   │   └── stat_repository.go
│   ├── service
│   │   ├── auth_service.go
│   │   ├── link_service.go
│   │   ├── user_service.go
│   │   └── stat_service.go
│   └── handler
│       ├── auth_handler.go
│       ├── link_handler.go
│       ├── user_handler.go
│       └── stat_handler.go
├── migrations
│   └── auto.go
├── pkg
│   ├── db
│   │   └── db.go
│   ├── di
│   │   └── interfaces.go
│   ├── event
│   │   └── event_bus.go
│   ├── jwt
│   │   └── jwt.go
│   ├── middleware
│   │   ├── auth_middleware.go
│   │   ├── chain.go
│   │   ├── cors.go
│   │   └── log.go
│   ├── req
│   │   ├── decode.go
│   │   ├── handle.go
│   │   └── validate.go
│   └── res
│       └── res.go
├── .env
├── .gitignore
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md

## 🚀 Запуск локально

```zsh
git clone git@github.com:MiCkEyZzZ/shorty.git
cd shorty
docker-compose up -d
go run cmd/main.go
```

## 🛠 Используемые технологии

- [go](https://go.dev/)
- [docker-compose](https://docs.docker.com/compose/)
- [jwt-go](https://github.com/golang-jwt/jwt)
- [bcrypt](https://github.com/golang/crypto)
- [gorm](https://github.com/go-gorm/gorm)
- [postgresql](https://www.postgresql.org/)
- [validator](https://github.com/go-playground/validator)
- [godotenv](https://github.com/joho/godotenv)

## 📌 Функционал приложения

- [x] Регистрация нового пользователя
- [x] Авторизация существующего пользователя
- [x] Обновление данных у пользователя
- [x] Удаление пользователя
- [x] Получения списка пользователей
- [x] получить пользователя по идентификатору
- [x] Создание коротких ссылок
- [x] Получить список коротких ссылок + лимит и смещение
- [x] Перенправление на оригинальный URL-адрес
- [x] Обновление ссылки
- [x] Удаление ссылки
- [x] Получить статистику
