package api

import (
	"encoding/json"
	"fmt"
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
func (s *server) handleGameAction(w http.ResponseWriter, r *http.Request) {
	if ok := enableCors(w, r); ok {
		return
	}
	vars := mux.Vars(r)
	g, err := s.game.GetGame(vars["game_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "the game doesn't exist"
		log.Error(msg)
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

	if err := g.MakeAction(a.SourcePlayerID, &a); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Error(err)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("%v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gameStatusResponse(g, a.SourcePlayerID))
}

// Return status of the game
func (s *server) handleGameStatus(w http.ResponseWriter, r *http.Request) {
	if ok := enableCors(w, r); ok {
		return
	}
	vars := mux.Vars(r)
	g, err := s.game.GetGame(vars["game_id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p := g.GetPlayer(vars["player"])
	if p == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(gameStatusResponse(g, p.ID))
}

// Add a new player to the game
func (s *server) handleGamePlayerNew(w http.ResponseWriter, r *http.Request) {
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
	log.Infof("Adding player %s", playerID)
	json.NewEncoder(w).Encode(map[string]string{"game": vars["game_id"], "player": playerID})
}

// Create new game
func (s *server) handleGameNew(w http.ResponseWriter, r *http.Request) {
	if ok := enableCors(w, r); ok {
		return
	}

	vars := mux.Vars(r)
	p, err := strconv.Atoi(vars["players"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	g, err := oriente.NewGame(p)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.game.AddGame(g)

	playerID, ok := g.GetFreePlayer()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"game": g.ID, "player": playerID})
}
