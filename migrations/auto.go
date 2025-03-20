package main

import (
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shorty/internal/models"
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
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Link{})
	db.Migrator().DropTable(&models.Stat{})
	db.AutoMigrate(&models.Link{}, &models.User{}, &models.Stat{})
}
