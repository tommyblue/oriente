package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/tommyblue/oriente/oriente"
)

/* A player can call the action if:
- it's its turn
- has the token
When it calls:
- the card is uncovered
- get the token
- tell the action ("pass", "attack" or "use_ability")
- turn is to the next player
*/
// TODO: Move some game logic into oriente package
func (s *server) actionHandler(w http.ResponseWriter, r *http.Request) {
	if ok := enableCors(w, r); ok {
		return
	}
	vars := mux.Vars(r)
	g, err := s.game.GetGame(vars["game_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "the game doesn't exist"
		log.Info(msg)
		json.NewEncoder(w).Encode(map[string]string{
			"error": msg,
		})
		return
	}

	var a oriente.Action
	err = json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p, ok := g.GetPlayer(a.SourcePlayerID)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		msg := "the player doesn't exist in this game"
		log.Info(msg)
		json.NewEncoder(w).Encode(map[string]string{
			"error": msg,
		})
		return
	}

	if p.DidAction {
		w.WriteHeader(http.StatusBadRequest)
		msg := "the player already played during this era"
		log.Info(msg)
		json.NewEncoder(w).Encode(map[string]string{
			"error": msg,
		})
		return
	}
	if ok := validateAction(a.Action); !ok {
		w.WriteHeader(http.StatusBadRequest)
		msg := "invalid action. Possible actions: 'attack', 'use_ability', 'pass'"
		log.Info(msg)
		json.NewEncoder(w).Encode(map[string]string{
			"error": msg,
		})
		return
	}

	if ok := g.MakeAction(p, &a); !ok {
		w.WriteHeader(http.StatusBadRequest)
		msg := "Cannot perform the action"
		log.Error(msg)
		json.NewEncoder(w).Encode(map[string]string{
			"error": msg,
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gameStatusResponse(g, p.ID))
}

// Return status of the game
func (s *server) gameStatusHandler(w http.ResponseWriter, r *http.Request) {
	if ok := enableCors(w, r); ok {
		return
	}
	vars := mux.Vars(r)
	g, err := s.game.GetGame(vars["game_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, ok := g.GetPlayer(vars["player"])
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gameStatusResponse(g, p.ID))
}

// Add a new player to the game
func (s *server) newPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if ok := enableCors(w, r); ok {
		return
	}
	vars := mux.Vars(r)
	g, err := s.game.GetGame(vars["game_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	playerID, ok := g.GetFreePlayer()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Debugf("Adding player %s", playerID)
	json.NewEncoder(w).Encode(map[string]string{"game": vars["game_id"], "player": playerID})
}

// Create new game
func (s *server) newGameHandler(w http.ResponseWriter, r *http.Request) {
	if ok := enableCors(w, r); ok {
		return
	}

	vars := mux.Vars(r)
	p, err := strconv.Atoi(vars["players"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	g := oriente.NewGame(p)
	s.game.AddGame(g)

	playerID, ok := g.GetFreePlayer()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"game": g.ID, "player": playerID})
}
