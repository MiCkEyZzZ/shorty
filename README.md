# Коротышка

коротышка — простой и быстрый инструмент для сокращения URL-адресов.

## 🚀 Запуск локально

```zsh
git clone git@github.com:MiCkEyZzZ/shorty.git
cd shorty
docker-compose up -d
go run cmd/main.go
```

## 🌍 Демо версия

[Перейти к приложению]()

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
- [x] Перенправление на оригинальный URL-адрес
- [x] Обновление ссылки
- [x] Удаление ссылки
- [ ] Получить статистику
