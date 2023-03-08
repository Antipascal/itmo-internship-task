package service

import (
	"github.com/Antipascal/itmo-internship-task/users/adapters"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func NewService() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("empty POSTGRES_DSN")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	_, err = adapters.NewAuthPostgres(db)
	if err != nil {
		log.Fatal(err)
	}

	_, err = adapters.NewUsersPostgres(db)
	if err != nil {
		log.Fatal(err)
	}
}
