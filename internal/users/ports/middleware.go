package ports

import (
	"context"
	"net/http"
	"strings"
)

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (hs *HTTPServer) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is authenticated
		token, found := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !found {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ISU, isAuthorized := hs.AuthManager.GetISU(token)
		if !isAuthorized {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "ISU", ISU))
		r = r.WithContext(context.WithValue(r.Context(), "IsAdmin", hs.AuthManager.IsAdmin(ISU)))
		next.ServeHTTP(w, r)
	})
}

func (hs *HTTPServer) AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is admin
		isAdmin, ok := r.Context().Value("IsAdmin").(bool)
		if !ok || !isAdmin {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
