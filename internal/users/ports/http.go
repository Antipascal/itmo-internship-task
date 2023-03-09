package ports

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Antipascal/itmo-internship-task/users/domain/auth"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

type HTTPServer struct {
	AuthManager auth.Manager
}

func NewHTTPServer(authManager auth.Manager) *HTTPServer {
	return &HTTPServer{AuthManager: authManager}
}

func (hs *HTTPServer) SetupRoutes(r *mux.Router) {
	r.Use(CorsMiddleware)

	// Service
	s := r.PathPrefix("/api").Subrouter()
	s.HandleFunc("/ping", hs.pongHandler).Methods(http.MethodGet)
	s.HandleFunc("/user", hs.UserInfoHandler).Methods(http.MethodGet)
	s.HandleFunc("/user", hs.pongHandler).Methods(http.MethodPut)
	s.HandleFunc("/users/search", hs.pongHandler).Methods(http.MethodGet)
	s.HandleFunc("/users/{id}", hs.pongHandler).Methods(http.MethodPut)
	s.HandleFunc("/users", hs.pongHandler).Methods(http.MethodGet)
	s.Use(hs.AuthMiddleware)

	// Auth
	r.HandleFunc("/", hs.AuthHandler).Methods(http.MethodGet)
}

func (hs *HTTPServer) pongHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "pong")
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (hs *HTTPServer) AuthHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	// If no code, redirect to itmo openid auth
	if code == "" {
		state, err := randStringBytes(64)
		if err != nil {
			writeError(w, err)
			return
		}
		SetCookie(w, "state", state)
		http.Redirect(w, r, hs.AuthManager.GetAuthURL(state), http.StatusFound)
		return
	}

	// Check that state in cookie and in url are equal
	s, err := r.Cookie("state")
	state := s.Value
	if err != nil || state != r.URL.Query().Get("state") {
		writeError(w, errors.New("invalid state"))
		return
	}

	token, err := hs.AuthManager.GetUserAccessToken(code, state)
	if err != nil {
		writeError(w, err)
		return
	}

	j := json.NewEncoder(w)
	err = j.Encode(map[string]string{"token": token})
	if err != nil {
		writeError(w, err)
		return
	}
}

func SetCookie(w http.ResponseWriter, name string, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		HttpOnly: true,
	})
}

func randStringBytes(n int) (string, error) {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (hs *HTTPServer) UserInfoHandler(w http.ResponseWriter, _ *http.Request) {
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
