package oriente

import (
	"fmt"

	"github.com/tommyblue/oriente/store"
)

func Initialize(s *store.Store) *Oriente {
	return &Oriente{
		store:        s,
		runningGames: make(map[string]*Game),
	}
}

func NewGame(nPlayers int) *Game {
	g := &Game{}
	g.generateDeck()
	g.addPrize()
	g.generatePlayers(nPlayers)

	return g
}

func (o *Oriente) AddGame(token string, g *Game) {
	o.runningGames[token] = g
}

func (o *Oriente) GetGame(token string) (*Game, error) {
	g, ok := o.runningGames[token]
	if !ok {
		return nil, fmt.Errorf("game not found")
	}
	return g, nil
}
