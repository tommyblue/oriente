package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Run() {
	r := mux.NewRouter()

	gameRoutes(r)

	spa := spaHandler{staticPath: "web/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	r.Use(mux.CORSMethodMiddleware(r))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func gameRoutes(r *mux.Router) {
	g := r.PathPrefix("/game").Subrouter()
	// Generate a new match
	g.HandleFunc("/new/{players:[0-9]+}", newGameHandler) //.Methods("POST")
	// Call the action
	g.HandleFunc("/{id}/{player}/call_action/", actionHandler).Methods("POST").HeadersRegexp("Content-Type", "application/json")
	// Game status for the player
	g.HandleFunc("/{id}/{player}", gameStatusHandler)
	// Add a new player to a game
	g.HandleFunc("/{id}", newPlayerHandler)
}

func enableCors(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return r.Method == http.MethodOptions
}
