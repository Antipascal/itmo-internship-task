package ports

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Antipascal/itmo-internship-task/users/adapters"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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
func (hs *HTTPServer) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var user adapters.UserDTO
	var err error
	if r.Context().Value("IsAdmin") == true && mux.Vars(r)["id"] != "" {
		var ISU int
		ISU, err = strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			writeError(w, errors.New("ID need to be valid ISU number"))
			return
		}
		user, err = hs.UsersManager.GetUser(ISU)
	} else {
		user, err = hs.UsersManager.GetUser(r.Context().Value("ISU").(int))
	}

	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j := json.NewEncoder(w)
	err = j.Encode(user)
	if err != nil {
		writeError(w, err)
		return
	}
	return
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
