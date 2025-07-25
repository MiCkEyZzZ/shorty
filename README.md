# Shorty

Shorty is a simple and fast URL shortening tool.

## 🚀  Run locally

```zsh
git clone git@github.com:MiCkEyZzZ/shorty.git
cd shorty
docker-compose up -d
go run cmd/main.go
```

## 🛠 Technologies Used

- [go](https://go.dev/)
- [docker-compose](https://docs.docker.com/compose/)
- [jwt-go](https://github.com/golang-jwt/jwt)
- [bcrypt](https://github.com/golang/crypto)
- [gorm](https://github.com/go-gorm/gorm)
- [postgresql](https://www.postgresql.org/)
- [validator](https://github.com/go-playground/validator)
- [godotenv](https://github.com/joho/godotenv)
- [minikube](https://minikube.sigs.k8s.io/docs/)
- [docker](https://www.docker.com/)
- [gihub action](https://github.com/features/actions)
- [keubernetes](https://kubernetes.io/)

## 📌 Features

- [x] Register a new user
- [x] Authenticate an existing user
- [x] Update user information
- [x] Delete a user
- [x] Create shortened links
- [x] Retrieve a list of shortened links with pagination (admin only)
- [x] Redirect to the original URL
- [x] Update a link
- [x] Delete a link
- [x] Retrieve link statistics
- [x] Retrieve statistics for all links (admin only)
- [x] Click count per day/month (admin only)
- [x] Total number of created links (admin only)
- [ ] Number of active/inactive links (admin only)
- [ ] Top 10 most popular links by clicks (admin only)
- [x] Number of blocked links (admin only)
- [x] Number of deleted links (admin only)
- [x] Block unwanted links (admin only)
- [x] Delete unwanted links (admin only)
- [x] Retrieve all users (admin only)
- [ ] Total number of users (admin only)
- [x] Retrieve a user by ID (admin only)
- [ ] Number of active users in the last 24h/week/month (admin only)
- [x] Update a user by ID (admin only)
- [x] Delete a user by ID (admin only)
- [x] Block a user by ID (admin only)
- [x] Number of blocked users (admin only)
- [ ] Number of users who created at least one link (admin only)
- [ ] Top 10 users by number of created links (admin only)
- [ ] Number of new users per day/week/month (admin only)
- [ ] Number of new links per day/week/month (admin only)
- [ ] Average number of clicks per user (admin only)
- [ ] Average lifespan of a link (time between creation and last click) (admin only)

## Примеры запросов

- **Авторизация пользователя:**

```
{
  "name": "твоё имя",
  "email": "твой адрес эл.почты",
  "password": "test",
  "role": "user" или "admin"
}
```

- **Аутентификация пользователя:**

```
{
  "email": "твой адрес эл.почты",
  "password": "test",
  "role": "user" или "admin"
}
```

## License

This project is licensed under the MIT License. The full license text is
available in the [License](./LICENSE).
