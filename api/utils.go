package api

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/tommyblue/oriente/oriente"
)

func validateAction(action string) bool {
	return action == "attack" || action == "use_ability" || action == "pass"
}

// Build the JSON response for the game status
func gameStatusResponse(g *oriente.Game, playerID string) map[string]interface{} {
	var players []map[string]interface{}
	for _, p := range g.Players {
		var points []map[string]interface{}
		for _, point := range p.Points {
			points = append(points, map[string]interface{}{"name": point.Name, "value": point.Value})
		}
		player := map[string]interface{}{
			"id":        p.ID,
			"name":      p.Name,
			"has_token": !p.DidAction,
			"points":    points,
		}
		if p.VisibleCard {
			player["card"] = p.CurrentCard.Name
			player["card_value"] = p.CurrentCard.Value
		} else {
			player["card"] = "hidden"
			player["card_value"] = "hidden"
		}
		players = append(players, player)
	}
	var tokenOwner string
	if g.TokenOwner != nil {
		tokenOwner = g.TokenOwner.ID
	}
	action := map[string]interface{}{}
	if g.CalledAction != nil {
		action["action"] = g.CalledAction.Action
		if g.CalledAction.Action != "pass" {
			action["source_player_id"] = g.CalledAction.SourcePlayerID
			action["target_player_id"] = g.CalledAction.TargetPlayerID
		}
	}
	return map[string]interface{}{
		"round":          g.Round,
		"players":        players,
		"active_players": g.ActivePlayers(),
		"game_started":   g.GameStarted(),
		"token_owner":    tokenOwner,
		"next_player":    g.NextPlayer.ID,
		"your_turn":      g.NextPlayer.ID == playerID,
		"called_action":  action,
		"prize_cards":    len(g.Prize), // Number of cards that the player playing first in the era will win
	}
}

func enableCors(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	return r.Method == http.MethodOptions
}

func (s *server) handleAndSync(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
		if err := s.game.SyncStore(); err != nil {
			log.Error(err)
		}
	}
}
