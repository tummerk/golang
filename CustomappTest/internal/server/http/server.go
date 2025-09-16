package httpServer

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
)

type GeneratorService interface {
	Generate() float64
}

type Server struct {
	GeneratorService GeneratorService
}

func NewServer(s GeneratorService) *Server {
	return &Server{
		GeneratorService: s,
	}
}
func (s *Server) Generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	response := struct {
		Result float64 `json:"result"`
	}{
		Result: s.GeneratorService.Generate(),
	}
	slog.Info("response", response)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (s *Server) RegisterRoutes(r chi.Router) {
	r.Get("/get", s.Generate)
}
