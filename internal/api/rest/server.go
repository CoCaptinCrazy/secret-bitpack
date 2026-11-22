package rest

import (
	"encoding/json"
	"net/http"
)

type Server struct {
	Mux *http.ServeMux
}

func New() *Server {
	server := &Server{Mux: http.NewServeMux()}
	server.Mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})
	return server
}

func (s *Server) ListenAndServe(address string) error {
	return http.ListenAndServe(address, s.Mux)
}
