package service

import (
	"github.com/Antipascal/itmo-internship-task/users/adapters"
	"github.com/Antipascal/itmo-internship-task/users/domain/auth"
	"github.com/Antipascal/itmo-internship-task/users/domain/users"
	"github.com/Antipascal/itmo-internship-task/users/ports"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func NewService() error {
	// Setup database
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

	// Setup managers
	am, err := auth.NewManager(authRepo, usersRepo)
	if err != nil {
		log.Fatal(err)
	}

	um := users.NewManager(authRepo, usersRepo)

	// Setup HTTP server
	s := ports.NewHTTPServer(am, um)
	r := mux.NewRouter()
	s.SetupRoutes(r)
	return http.ListenAndServe(":8080", r)
}
