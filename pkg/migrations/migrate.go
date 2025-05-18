package migrations

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// RunUp applies the database schema by running the necessary migrations.
// Currently, it performs AutoMigrate on the Vacancy model.
func RunUp(m any) {
	loadEnv()
	db := openDB()
	err := db.AutoMigrate(m)
	if err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}
	log.Println("✅ Migration UP completed")
}

// RunDown rolls back the database schema changes by dropping the tables
// related to the Vacancy model.
func RunDown(m any) {
	loadEnv()
	db := openDB()
	err := db.Migrator().DropTable(m)
	if err != nil {
		log.Fatalf("❌ DropTable failed: %v", err)
	}
	log.Println("✅ Migration DOWN completed")
}

// openDB opens a connection to the PostgreSQL database using the DSN
// from environment variables.
func openDB() *gorm.DB {
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	return db
}

// loadEnv loads environment variables from a .env file.
// If the file is missing, it logs a warning and uses existing defaults.
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ .env not found, using defaults")
	}
}
