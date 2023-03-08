package ports

import (
	"encoding/json"
	"fmt"
	"github.com/Antipascal/itmo-internship-task/users/domain/auth"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type HTTPServer struct {
	AuthManager auth.Manager
}

func NewHTTPServer(authManager auth.Manager) *HTTPServer {
	return &HTTPServer{AuthManager: authManager}
}

func (hs *HTTPServer) SetupRoutes(r *mux.Router) {
	r.Use(CorsMiddleware)

	// service
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/ping", hs.pongHandler).Methods(http.MethodGet)
	s.HandleFunc("/user", hs.UserInfoHandler).Methods(http.MethodGet)
	s.HandleFunc("/user", hs.pongHandler).Methods(http.MethodPut)
	s.HandleFunc("/users/search", hs.pongHandler).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", hs.pongHandler).Methods(http.MethodPut)
	s.HandleFunc("/users", hs.pongHandler).Methods(http.MethodGet)
	s.Use(AuthMiddleware)

	// auth
	r.HandleFunc("/", hs.AuthHandler).Methods(http.MethodGet)
}

func (hs *HTTPServer) pongHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "pong")
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (hs *HTTPServer) AuthHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	authType := r.URL.Query().Get("state")
	if code == "" || authType == "" {
		w.Header().Set("Content-Type", "application/json")
		j := json.NewEncoder(w)
		err := j.Encode(map[string]string{"itmoid_link": auth.ItmoAuthLink})
		if err != nil {
			writeError(w, err)
			return
		}
		return
	}

	var token string
	var err error
	switch authType {
	case "itmoid":
		token, err = hs.AuthManager.AuthItmoCode(code)
		if err != nil {
			writeError(w, err)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	j := json.NewEncoder(w)
	err = j.Encode(map[string]string{"token": token})
	if err != nil {
		writeError(w, err)
		return
	}
}

func (hs *HTTPServer) UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
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
