package api

import "github.com/gorilla/mux"

func (s *server) routes() {

	/* Routes */

	g := s.router.PathPrefix("/game").Subrouter()
	// Generate a new match
	g.HandleFunc("/new/{players:[0-9]+}", s.handleAndSync(s.handleGameNew))
	// Call the action
	g.HandleFunc("/{game_id}/call_action/", s.handleAndSync(s.handleGameAction)).Methods("POST").HeadersRegexp("Content-Type", "application/json")
	// Game status for the player
	g.HandleFunc("/{game_id}/{player}", s.handleAndSync(s.handleGameStatus))
	// Add a new player to a game
	g.HandleFunc("/{game_id}", s.handleAndSync(s.handleGamePlayerNew))

	// SPA
	spa := spaHandler{staticPath: "web/build", indexPath: "index.html"}
	s.router.PathPrefix("/").Handler(spa)

	// CORS
	s.router.Use(mux.CORSMethodMiddleware(s.router))
}
