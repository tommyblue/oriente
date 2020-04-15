package api

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/tommyblue/oriente/oriente"
)

type server struct {
	game   *oriente.Oriente
	router *mux.Router
	srv    *http.Server
}

func Run(game *oriente.Oriente) {
	s := &server{
		game:   game,
		router: mux.NewRouter(),
	}
	s.start()
}

func (s *server) start() {
	s.router.Use(loggingMiddleware)
	s.routes()

	s.srv = &http.Server{
		Handler:      s.router,
		Addr:         ":8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(s.srv.ListenAndServe())
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
