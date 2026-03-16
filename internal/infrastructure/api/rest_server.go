package api

import (
	"encoding/json"
	"net/http"

	"github.com/JoseGusnay/ursus/internal/application/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type RESTServer struct {
	service *service.MemoryService
	router  *chi.Mux
}

func NewRESTServer(svc *service.MemoryService) *RESTServer {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	s := &RESTServer{
		service: svc,
		router:  r,
	}

	s.registerRoutes()
	return s
}

func (s *RESTServer) registerRoutes() {
	s.router.Route("/v1", func(r chi.Router) {
		r.Get("/memories", s.handleList)
		r.Post("/memories", s.handleStore)
		r.Get("/search", s.handleSearch)
		r.Delete("/memories/{id}", s.handleDelete)
	})
}

func (s *RESTServer) handleList(w http.ResponseWriter, r *http.Request) {
	results, err := s.service.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)
}

func (s *RESTServer) handleStore(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Content  string `json:"content"`
		Metadata string `json:"metadata"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := s.service.Store(r.Context(), body.Content, body.Metadata, "", "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (s *RESTServer) handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	results, err := s.service.Search(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)
}

func (s *RESTServer) handleDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.service.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *RESTServer) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.router)
}
