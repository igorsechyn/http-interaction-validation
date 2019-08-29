package e2e_test

import (
	"net/http"

	"github.com/go-chi/chi"
)

func startServer(handler http.Handler) *http.Server {
	r := chi.NewRouter()
	r.Method("POST", "/", handler)
	server := &http.Server{Addr: ":8080", Handler: r}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	return server
}
