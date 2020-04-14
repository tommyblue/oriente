package oriente

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (g *Game) MakeAction(p *Player, action *Action) error {
	if action.Action != "attack" && action.Action != "use_ability" && action.Action != "pass" {
		return fmt.Errorf("invalid action %s. Possible actions: 'attack', 'use_ability', 'pass'", action.Action)
	}

	if p != g.NextPlayer {
		return fmt.Errorf("not the turn of the player")
	}

	if action.Action != "pass" {
		if err := g.canPerformAction(p); err != nil {
			return err
		}
		p.DidAction = true
		p.VisibleCard = true
		// If the player is "stealing" the action to another player, that player Action must be reset
		if g.CalledAction != nil {
			p, ok := g.GetPlayer(g.CalledAction.SourcePlayerID)
			if !ok {
				log.Panicf("Unknown player ID %s that already did called the action", g.CalledAction.SourcePlayerID)
			}
			p.DidAction = false
		}
		g.CalledAction = action
		// The player gets the token
		g.TokenOwner = p
	}

	// The loop has ended, the destiny can be fulfilled
	// if g.NextPlayer == p { // TODO: is this check correct?
	// 	g.fulfillDestiny()
	// }
	g.nextPlayerTurn()
	g.checkEndEra()
	g.Round++
	return nil
}

/* The era ends if:
- all player passed
or
- nobody owns the action token
or
- all players with the action token passed
At the end of the era:
- a new card is put on the prize deck
- all players get a new action token
- the TokenOwner will play the next turn
*/
func (g *Game) checkEndEra() {
	check := false
	// All players passed
	if g.NextPlayer.ID == g.TokenOwner.ID {
		check = true
	}
	// Check if at least one player still has the action token
	allDidAction := true
	for _, p := range g.Players {
		if !p.DidAction {
			allDidAction = false
			break
		}
	}
	if check || allDidAction {
		// Set NextPlayer to nil so that the next turn will be played by the TokenOwner
		g.NextPlayer = g.TokenOwner
		// Add a card to the prize
		g.Prize = append(g.Prize, g.pickCard()) // TODO: check if there's at least one card left in the deck
		// Reset the action tokens
		for _, p := range g.Players {
			p.DidAction = false
		}
	}
}

/* When the player fulfill his destiny, these things happen:
1 - if there's a prize, move the prize in the turn's temporary prize and empty the era prize
2 - at the end, add the temporary prize to the players points
3 - The next turn will be of the player that owns the token
*/
func (g *Game) fulfillDestiny() {
	// 1
	g.TempPrize = g.Prize
	g.Prize = nil
	// 2
	g.NextPlayer.Points = append(g.NextPlayer.Points, g.TempPrize...)
	g.TempPrize = nil
}

func (g *Game) canPerformAction(p *Player) error {
	// First action
	if g.CalledAction == nil {
		return nil
	}
	// A player can't perform an action if already done during this era
	if p.DidAction {
		return fmt.Errorf("Already did action")
	}
	// The player target of the called action can't play
	if g.CalledAction.TargetPlayerID == p.ID {
		return fmt.Errorf("called player")
	}
	// Player with less power can't stop the action
	if p.CurrentCard.Value < g.TokenOwner.CurrentCard.Value {
		return fmt.Errorf("less power")
	}
	// If the power of the players is the same, the player must be poorer to stop the action
	if p.CurrentCard.Value == g.TokenOwner.CurrentCard.Value && len(p.Points) >= len(g.TokenOwner.Points) {
		return fmt.Errorf("too rich")
	}
	return nil
}

func (g *Game) nextPlayerTurn() {
	for i, p := range g.Players {
		if p.ID == g.NextPlayer.ID {
			g.NextPlayer = g.Players[(i+1)%len(g.Players)]
			return
		}
	}
	log.Fatal("Can't find next player")
}

func (g *Game) GameStarted() bool {
	return len(g.Players) == g.ActivePlayers()
}

// ActivePlayers return the number of currently active players (ready to play)
func (g *Game) ActivePlayers() int {
	var players int
	for _, p := range g.Players {
		if p.Managed {
			players++
		}
	}
	return players
}

// Player return the player
func (g *Game) GetPlayer(playerID string) (*Player, bool) {
	for _, p := range g.Players {
		if p.ID == playerID {
			return p, true
		}
	}
	return nil, false
}

// GetFreePlayer return the ID of the first available spot in the game
func (g *Game) GetFreePlayer() (string, bool) {
	for _, p := range g.Players {
		if !p.Managed {
			p.Managed = true
			return p.ID, true
		}
	}
	return "", false
}

func (g *Game) addPrize() {
	g.Prize = append(g.Prize, g.pickCard())
}

func (g *Game) pickCard() *Card {
	c := g.Deck[len(g.Deck)-1]
	g.Deck = g.Deck[:len(g.Deck)-1]
	return c
}
