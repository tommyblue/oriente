package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tommyblue/oriente/oriente"
	"github.com/tommyblue/oriente/utils"
)

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

	var a oriente.PlayerAction
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if ok := validateAction(a.Action); !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid action. Possible actions: 'attack', 'use_ability'",
		})
		return
	}

	g.MakeAction(p, &a)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gameStatusResponse(g, vars))
}

// Return status of the game
func gameStatusHandler(w http.ResponseWriter, r *http.Request) {
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
