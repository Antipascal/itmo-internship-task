package ports

import (
	"fmt"
	"net/http"
)

func (hs *HTTPServer) pongHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "pong")
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetUserHandler returns user data by id
// If user is admin, he can get any user by adding id to path
func (hs *HTTPServer) GetUserHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

// UpdateUserHandler updates users data
// If user is admin, he can update any user by adding id to path
func (hs *HTTPServer) UpdateUserHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

// SearchUserHandler returns users by search query. Query can be:
//   - Search by phone number: "/users/search?phone=123456789".
func (hs *HTTPServer) SearchUserHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

// GetUsersHandler returns all users in database
// Admin permission required
func (hs *HTTPServer) GetUsersHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}
