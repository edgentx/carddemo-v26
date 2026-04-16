package shared

import "net/http"

// Routers returns a list of sub-routers to be mounted in main.
// Currently returns an empty router to satisfy compilation.
func Routers() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
}
