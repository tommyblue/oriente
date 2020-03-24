package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tommyblue/oriente/oriente"
)

func Run() {
	r := mux.NewRouter()

	// Generate a new match
	r.HandleFunc("/game/new/{players:[0-9]+}", newGameHandler) //.Methods("POST")
	// Game status for the player
	r.HandleFunc("/game/{id}/{player}", gameHandler)
	// Add a new player to a game
	r.HandleFunc("/game/{id}", newPlayerHandler)

	spa := spaHandler{staticPath: "web/build", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// Return status of the game
func gameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"games": vars["id"]})
}

// Add a new player to the game
func newPlayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"game": vars["id"]})
}

// Create new game
func newGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := tokenGenerator()
	p, err := strconv.Atoi(vars["players"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	g := oriente.NewGame(p)
	oriente.RunningGames[token] = g
	json.NewEncoder(w).Encode(map[string]string{"game": token})
}
