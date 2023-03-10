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
func (hs *HTTPServer) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var ISU int
	if r.Context().Value("IsAdmin") == true && mux.Vars(r)["id"] != "" {
		var err error
		ISU, err = strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			writeError(w, errors.New("ID need to be valid ISU number"))
			return
		}
	} else {
		ISU = r.Context().Value("ISU").(int)
	}

	var user adapters.UserDTO
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		writeError(w, err)
		return
	}
	user.ISU = ISU
	err = hs.UsersManager.UpdateUser(user)

	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// SearchUserHandler returns users by search query. Query can be:
//   - Search by phone number: "/users/search?phone=123456789".
func (hs *HTTPServer) SearchUserHandler(w http.ResponseWriter, r *http.Request) {
	phone := r.URL.Query().Get("phone")
	if phone != "" {
		user, err := hs.UsersManager.GetUserByPhone(phone)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			j := json.NewEncoder(w)
			err = j.Encode(user)
			if err != nil {
				writeError(w, err)
			}
		} else {
			writeError(w, err)
		}
	}
}

// GetUsersHandler returns users in database with offset and limit
// given as a query parameters. Default values for limit is 20 and offset is 0.
func (hs *HTTPServer) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 20
	}

	users, err := hs.UsersManager.GetUsers(offset, limit)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j := json.NewEncoder(w)
	err = j.Encode(users)
	if err != nil {
		writeError(w, err)
		return
	}
	return
}
