package api

import "github.com/tommyblue/oriente/oriente"

func validateAction(action string) bool {
	return action == "attack" || action == "use_ability"
}

func gameStatusResponse(g *oriente.Game, vars map[string]string) map[string]interface{} {
	var players []map[string]interface{}
	for _, p := range g.Players {
		player := map[string]interface{}{
			"id":        p.ID,
			"name":      p.Name,
			"has_token": !p.DidAction,
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
	var token string
	if g.Token != nil {
		token = g.Token.ID
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
		"players":        players,
		"active_players": g.ActivePlayers(),
		"game_started":   g.GameStarted(),
		"token":          token,
		"next_player":    g.NextPlayer.ID,
		"your_turn":      g.NextPlayer.ID == vars["player"],
		"called_action":  action,
	}
}
