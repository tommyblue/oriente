package oriente

import (
	"fmt"
	"math/rand"

	log "github.com/sirupsen/logrus"

	"github.com/tommyblue/oriente/utils"
)

func validateAction(action string) bool {
	return action == "attack" || action == "use_ability" || action == "pass"
}

func (g *Game) MakeAction(p *Player, action *Action) error {
	if !validateAction(action.Action) {
		return fmt.Errorf("invalid action %s. Possible actions: 'attack', 'use_ability', 'pass'", action.Action)
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
	g.nextPlayerTurn()
	g.Round++
	// The loop has ended, the destiny can be fulfilled
	if g.NextPlayer == p { // TODO: is this check correct?
		g.fulfillDestiny()
	}
	// g.endEra() // TODO
	return nil
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
	found := false
	for _, p := range g.Players {
		if found {
			g.NextPlayer = p
			return
		}

		if p.ID == g.NextPlayer.ID {
			found = true
		}
	}

	// TODO: is this always true?
	g.NextPlayer = g.Players[0]
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

func (g *Game) generatePlayers(nPlayers int) error {
	if g.Deck == nil {
		return fmt.Errorf("Must generate the deck before players")
	}
	// This is the deck of cards with money
	coinsDeck := []*Card{
		&Card{Name: "2 Coins", Value: 2},
		&Card{Name: "2 Coins", Value: 2},
		&Card{Name: "2 Coins", Value: 2},
		&Card{Name: "2 Coins", Value: 2},
		&Card{Name: "3 Coins", Value: 3},
		&Card{Name: "3 Coins", Value: 3},
		&Card{Name: "3 Coins", Value: 3},
		&Card{Name: "3 Coins", Value: 3},
		&Card{Name: "4 Coins", Value: 4},
		&Card{Name: "4 Coins", Value: 4},
		&Card{Name: "4 Coins", Value: 4},
		&Card{Name: "4 Coins", Value: 4},
	}

	for i := 0; i < nPlayers; i++ {
		mIdx := rand.Intn(len(coinsDeck))
		coin := coinsDeck[mIdx]
		p := &Player{
			ID:          utils.IDGenerator(),
			Name:        fmt.Sprintf("player_%d", i),
			CurrentCard: g.pickCard(),
		}
		p.Points = append(p.Points, coin)
		g.Players = append(g.Players, p)
		coinsDeck = append(coinsDeck[:mIdx], coinsDeck[mIdx+1:]...)
	}

	g.TokenOwner = g.Players[0]
	g.NextPlayer = g.Players[0]
	return nil
}

func (g *Game) generateDeck() {

	tmpDeck := []*Card{}
	// Add 12 Nofu
	tmpDeck = append(tmpDeck, generateCard(Nofu, 12)...)
	// Add 7 Akindo
	tmpDeck = append(tmpDeck, generateCard(Akindo, 7)...)
	// Add 6 Samurai
	tmpDeck = append(tmpDeck, generateCard(Samurai, 6)...)
	// Add 5 Daimyo
	tmpDeck = append(tmpDeck, generateCard(Daimyo, 5)...)
	// Add 4 Moho-Tsukai
	tmpDeck = append(tmpDeck, generateCard(MahoTsukai, 4)...)
	// Add 2 Soryo
	tmpDeck = append(tmpDeck, generateCard(Soryo, 2)...)

	// Shuffle
	rand.Shuffle(len(tmpDeck), func(i, j int) { tmpDeck[i], tmpDeck[j] = tmpDeck[j], tmpDeck[i] })

	// Get 4 cards for the base
	base := make([]*Card, 4)
	copy(base, tmpDeck[0:4])
	tmpDeck = append(tmpDeck[:0], tmpDeck[4:]...)

	// Get 4 cards and insert Geisha, shuffle
	wGeisha := make([]*Card, 4)
	copy(wGeisha, tmpDeck[0:4])
	tmpDeck = append(tmpDeck[:0], tmpDeck[4:]...)
	wGeisha = append(wGeisha, generateCard(Geisha, 1)...)
	rand.Shuffle(len(wGeisha), func(i, j int) { wGeisha[i], wGeisha[j] = wGeisha[j], wGeisha[i] })

	// Get remaining cards and add 1 Shogun and 3 Ninja, shuffle
	tmpDeck = append(tmpDeck, generateCard(Shogun, 1)...)
	tmpDeck = append(tmpDeck, generateCard(Ninja, 3)...)
	rand.Shuffle(len(tmpDeck), func(i, j int) { tmpDeck[i], tmpDeck[j] = tmpDeck[j], tmpDeck[i] })

	// Compile the deck: 4 cards + (4 cards + 1 Geisha) + all other cards
	g.Deck = append(g.Deck, base...)
	g.Deck = append(g.Deck, wGeisha...)
	g.Deck = append(g.Deck, tmpDeck...)
}

func generateCard(c Character, n int) []*Card {
	v := []*Card{}
	for i := 0; i < n; i++ {
		v = append(v, &Card{
			Name:  c.String(),
			Value: c,
		})
	}
	return v
}
