package oriente

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (g *Game) MakeAction(playerID string, action *Action) error {
	if action.Action != AttackAction && action.Action != UseAbilityAction && action.Action != PassAction {
		return fmt.Errorf("invalid action %s. Possible actions: 'attack', 'use_ability', 'pass'", action.Action)
	}
	p := g.GetPlayer(playerID)
	if p == nil {
		return fmt.Errorf("Can't find the player")
	}
	if p.ID != g.NextPlayerID {
		return fmt.Errorf("not the turn of the player")
	}

	if action.Action != PassAction {
		if err := g.canPerformAction(p); err != nil {
			return err
		}
		p.DidAction = true
		p.VisibleCard = true
		// If the player is "stealing" the action to another player, that player Action must be reset
		if g.CalledAction != nil {
			p := g.GetPlayer(g.CalledAction.SourcePlayerID)
			if p == nil {
				log.Panicf("Unknown player ID %s that already did called the action", g.CalledAction.SourcePlayerID)
			}
			p.DidAction = false
		}
		g.CalledAction = action
		// The player gets the token
		g.TokenOwnerID = p.ID
	}

	g.nextPlayerTurn()
	g.checkAndFulfillDestiny()
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
	// TODO: Check if someone picked the Geisha or has 3 Ninjas

	check := false
	// All players passed
	if g.NextPlayerID == g.TokenOwnerID {
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
		// Set NextPlayer to the TokenOwner
		g.NextPlayerID = g.TokenOwnerID
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
2 - battle
3 - at the end, add the temporary prize to the players points
4 - The next turn will be of the player that owns the token
5 - player(s) that lost the battle will get a new card. Their old cards are the prize for the winners
*/
func (g *Game) checkAndFulfillDestiny() bool {
	if g.NextPlayerID != g.TokenOwnerID || g.CalledAction == nil || g.CalledAction.SourcePlayerID != g.NextPlayerID {
		return false
	}

	// 1
	g.TempPrize = g.Prize
	g.Prize = nil

	// 2
	if g.CalledAction.Action == AttackAction {
		attacker := g.GetPlayer(g.CalledAction.SourcePlayerID)
		if attacker == nil {
			log.Fatalf("Cannot find attacker player %s", g.CalledAction.SourcePlayerID)
		}
		defender := g.GetPlayer(g.CalledAction.TargetPlayerID)
		if defender == nil {
			log.Fatalf("Cannot find defender player %s", g.CalledAction.TargetPlayerID)
		}
		// The defender shows the card
		defender.VisibleCard = true
		if defender.CurrentCard.Value >= attacker.CurrentCard.Value {
			// defender wins
			defender.Points = append(defender.Points, attacker.CurrentCard)
			attacker.CurrentCard = g.pickCard()
			attacker.VisibleCard = false
		} else {
			// attacker wins
			attacker.Points = append(attacker.Points, attacker.CurrentCard)
			defender.CurrentCard = g.pickCard()
			defender.VisibleCard = false
		}
	} else {
		// special power
	}

	// 3
	g.NextPlayer().Points = append(g.NextPlayer().Points, g.TempPrize...)
	g.TempPrize = nil

	// 4
	g.NextPlayerID = g.TokenOwnerID

	return true
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
	if p.CurrentCard.Value < g.TokenOwner().CurrentCard.Value {
		return fmt.Errorf("less power")
	}
	// If the power of the players is the same, the player must be poorer to stop the action
	if p.CurrentCard.Value == g.TokenOwner().CurrentCard.Value && len(p.Points) >= len(g.TokenOwner().Points) {
		return fmt.Errorf("too rich")
	}
	return nil
}

func (g *Game) nextPlayerTurn() {
	for i, p := range g.Players {
		if p.ID == g.NextPlayerID {
			nextP := g.Players[(i+1)%len(g.Players)].ID
			// The player target of the previous action can't play this turn, so it goes to the
			// following player
			if g.CalledAction != nil && nextP == g.CalledAction.TargetPlayerID {
				nextP = g.Players[(i+2)%len(g.Players)].ID
			}
			g.NextPlayerID = nextP
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
func (g *Game) GetPlayer(playerID string) *Player {
	for _, p := range g.Players {
		if p.ID == playerID {
			return p
		}
	}
	return nil
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

// NextPlayer returns the next player to play
func (g *Game) NextPlayer() *Player {
	p := g.GetPlayer(g.NextPlayerID)
	if p == nil {
		log.Fatal("Can't find next player")
	}
	return p
}

// TokenOwner is the player with the token
func (g *Game) TokenOwner() *Player {
	p := g.GetPlayer(g.TokenOwnerID)
	if p == nil {
		log.Fatal("Can't find token owner")
	}
	return p
}
