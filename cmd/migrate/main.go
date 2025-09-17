package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"lab2/internal/app/ds"
	"lab2/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(
		&ds.Period{},
		&ds.Company{},
		&ds.PeriodsApplication{},
		&ds.SelectedPeriod{},
		&ds.User{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}