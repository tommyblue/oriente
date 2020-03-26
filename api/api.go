package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/tommyblue/oriente/oriente"
	"github.com/tommyblue/oriente/utils"
)

func Run() {
	r := mux.NewRouter()

	// Generate a new match
	r.HandleFunc("/game/new/{players:[0-9]+}", newGameHandler) //.Methods("POST")
	/* A player can call the action if:
	- it's its turn
	- has the token
	When it calls:
	- the card is uncovered
	- get the token
	- tell the action ("attack" or "use_ability")
	- turn is to the next player
	*/
	r.HandleFunc("/game/{id}/{player}/call_action/{action}", actionHandler)
	// Game status for the player
	r.HandleFunc("/game/{id}/{player}", gameHandler)
	// Add a new player to a game
	r.HandleFunc("/game/{id}", newPlayerHandler)

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

func enableCors(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return r.Method == http.MethodOptions
}

func actionHandler(w http.ResponseWriter, r *http.Request) {
	if ok := enableCors(w, r); ok {
		return
	}
	vars := mux.Vars(r)
	g, ok := oriente.RunningGames[vars["id"]]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "the game doesn't exist",
		})
		return
	}

	p, ok := g.Player(vars["player"])
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "the player doesn't exist in this game",
		})
		return
	}

	if p.DidAction {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "the player already played during this era",
		})
		return
	}

	if ok := validateAction(vars["action"]); !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid action. Possible actions: 'attack', 'use_ability'",
		})
		return
	}

	g.MakeAction(p, vars["action"])

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gameStatusResponse(g, vars))
}

// Return status of the game
func gameHandler(w http.ResponseWriter, r *http.Request) {
	if ok := enableCors(w, r); ok {
		return
	}
	vars := mux.Vars(r)
	g, ok := oriente.RunningGames[vars["id"]]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if ok := g.ValidatePlayer(vars["player"]); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gameStatusResponse(g, vars))
}

// Add a new player to the game
func newPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if ok := enableCors(w, r); ok {
		return
	}
	vars := mux.Vars(r)
	g, ok := oriente.RunningGames[vars["id"]]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	playerID, ok := g.GetFreePlayer()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"game": vars["id"], "player": playerID})
}

// Create new game
func newGameHandler(w http.ResponseWriter, r *http.Request) {
	if ok := enableCors(w, r); ok {
		return
	}
	vars := mux.Vars(r)
	token := utils.TokenGenerator()
	p, err := strconv.Atoi(vars["players"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	g := oriente.NewGame(p)
	g.ID = token
	oriente.RunningGames[token] = g
	playerID, ok := g.GetFreePlayer()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"game": token, "player": playerID})
}
