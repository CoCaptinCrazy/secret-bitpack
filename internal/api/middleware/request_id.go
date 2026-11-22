package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			raw := make([]byte, 12)
			if _, err := rand.Read(raw); err != nil {
				http.Error(w, "request ID generation failed", http.StatusInternalServerError)
				return
			}
			id = hex.EncodeToString(raw)
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
