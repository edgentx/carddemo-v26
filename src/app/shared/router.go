package shared

import "github.com/go-chi/chi/v5"

// NewRouter creates a new chi router instance.
func NewRouter() *chi.Mux {
	return chi.NewRouter()
}
