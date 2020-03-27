package api

import "github.com/tommyblue/oriente/oriente"

func validateAction(action string) bool {
	return action == "attack" || action == "use_ability"
}

func gameStatusResponse(g *oriente.Game, vars map[string]string) map[string]interface{} {
	var players []map[string]interface{}
	for _, p := range g.Players {
		action := map[string]interface{}{}
		if p.CalledAction != nil {
			action["action"] = p.CalledAction.Action
			if p.CalledAction.Action == "attack" {
				action["player_id"] = p.CalledAction.PlayerID
			}
		}
		player := map[string]interface{}{
			"id":        p.ID,
			"name":      p.Name,
			"has_token": !p.DidAction,
			"action":    action,
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
	return map[string]interface{}{
		"players":        players,
		"active_players": g.ActivePlayers(),
		"game_started":   g.GameStarted(),
		"next_player":    g.NextPlayer.ID,
		"your_turn":      g.NextPlayer.ID == vars["player"],
	}
}
