package ports

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

func (hs *HTTPServer) AuthHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	// If no code, redirect to itmo openid auth
	if code == "" {
		state, err := randStringBytes(tokenLength)
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
