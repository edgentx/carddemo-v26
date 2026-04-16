package main

import (
	"net/http"

	"github.com/carddemo/project/src/app/shared"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Mount("/", shared.Routers())

	http.ListenAndServe(":8080", r)
}
