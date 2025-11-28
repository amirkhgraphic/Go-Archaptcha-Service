package initializers

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	if dbURL := os.Getenv("DB_URL"); dbURL != "" {
		DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
		if err != nil {
			panic("failed to connect database: " + err.Error())
		}
		return
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/data.db"
	}

	if dir := filepath.Dir(dbPath); dir != "" && dir != "." {
		_ = os.MkdirAll(dir, 0o755)
	}

	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
}
