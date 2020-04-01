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
	s.setupGameRoutes()
	s.setupSpaRoutes()
	s.setupCors()

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

func (s *server) setupGameRoutes() {
	g := s.router.PathPrefix("/game").Subrouter()
	// Generate a new match
	g.HandleFunc("/new/{players:[0-9]+}", s.handleAndSync(s.newGameHandler))
	// Call the action
	g.HandleFunc("/{game_id}/call_action/", s.handleAndSync(s.actionHandler)).Methods("POST").HeadersRegexp("Content-Type", "application/json")
	// Game status for the player
	g.HandleFunc("/{game_id}/{player}", s.handleAndSync(s.gameStatusHandler))
	// Add a new player to a game
	g.HandleFunc("/{game_id}", s.handleAndSync(s.newPlayerHandler))
}

func (s *server) setupSpaRoutes() {
	spa := spaHandler{staticPath: "web/build", indexPath: "index.html"}
	s.router.PathPrefix("/").Handler(spa)
}

func (s *server) setupCors() {
	s.router.Use(mux.CORSMethodMiddleware(s.router))
}
