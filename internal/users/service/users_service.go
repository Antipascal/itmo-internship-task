package service

import (
	"github.com/Antipascal/itmo-internship-task/users/adapters"
	"github.com/Antipascal/itmo-internship-task/users/domain/auth"
	"github.com/Antipascal/itmo-internship-task/users/ports"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func NewService() error {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("empty POSTGRES_DSN")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	authRepo, err := adapters.NewAuthPostgres(db)
	if err != nil {
		log.Fatal(err)
	}

	usersRepo, err := adapters.NewUsersPostgres(db)
	if err != nil {
		log.Fatal(err)
	}

	am, err := auth.NewManager(authRepo, usersRepo)
	if err != nil {
		log.Fatal(err)
	}

	s := ports.NewHTTPServer(*am)
	r := mux.NewRouter()
	s.SetupRoutes(r)
	return http.ListenAndServe(":8080", r)
}
