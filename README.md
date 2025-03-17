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

- Регистрация нового пользователя
- Авторизация существующего пользователя
- Обновление данных у пользователя
- Удаление пользователя
- Получения списка пользователей
- получить пользователя по идентификатору
- Создание коротких ссылок
- Перенправление на оригинальный URL-адрес
- Обновление ссылки
- Удаление ссылки
- Получить статистику
