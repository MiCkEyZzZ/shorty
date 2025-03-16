# Shorty

shorty is a simple and fast tool for shortening URLs.

## 🚀 Run locally

```zsh
git clone git@github.com:MiCkEyZzZ/shorty.git
cd shorty
docker-compose up -d
go run cmd/main.go
```

## 🌍 Demo version

[Go to app]()

## 🛠 Technologies used

- [go](https://go.dev/)
- [docker-compose](https://docs.docker.com/compose/)
- [gorm](https://github.com/go-gorm/gorm)
- [postgresql](https://www.postgresql.org/)
- [validator](https://github.com/go-playground/validator)
- [godotenv](https://github.com/joho/godotenv)

## 📌 App functionality

- Registration
- Authorization
- Creating short links
- Redirect to the original URL
- Updating short links
- Removing short links
- Get statistics on short links
