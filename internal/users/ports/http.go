package ports

import (
	"encoding/json"
	"github.com/Antipascal/itmo-internship-task/users/domain/auth"
	"github.com/Antipascal/itmo-internship-task/users/domain/users"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type HTTPServer struct {
	AuthManager  *auth.Manager
	UsersManager *users.Manager
}

func NewHTTPServer(authManager *auth.Manager, usersManager *users.Manager) *HTTPServer {
	return &HTTPServer{AuthManager: authManager, UsersManager: usersManager}
}

const (
	tokenLength = 64
)

func (hs *HTTPServer) SetupRoutes(r *mux.Router) {
	r.Use(CorsMiddleware)
	r.HandleFunc("/ping", hs.pongHandler).Methods(http.MethodGet)

	// Service
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/user", hs.GetUserHandler).Methods(http.MethodGet)
	s.HandleFunc("/user", hs.UpdateUserHandler).Methods(http.MethodPut)
	s.HandleFunc("/users/search", hs.SearchUserHandler).Methods(http.MethodGet)
	s.Use(hs.AuthMiddleware)

	// Admin
	a := s.PathPrefix("/admin").Subrouter()
	a.HandleFunc("/users/{id}", hs.UpdateUserHandler).Methods(http.MethodPut)
	a.HandleFunc("/users/{id}", hs.GetUserHandler).Methods(http.MethodGet)
	a.HandleFunc("/users", hs.GetUsersHandler).Methods(http.MethodGet)
	a.Use(hs.AdminMiddleware)

	// Auth
	r.HandleFunc("/", hs.AuthHandler).Methods(http.MethodGet)
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	j := json.NewEncoder(w)
	err = j.Encode(map[string]string{"error": err.Error()})
	if err != nil {
		log.Println(err)
	}
}
