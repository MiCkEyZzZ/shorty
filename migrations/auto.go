package main

import (
	"os"
	"os/user"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shorty/internal/link"
	"shorty/internal/stat"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// db.Migrator().DropTable(&models.User{})
	// db.Migrator().DropTable(&models.Link{})
	db.AutoMigrate(&link.Link{}, &user.User{}, &stat.Stat{})
}
