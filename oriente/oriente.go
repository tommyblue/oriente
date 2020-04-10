package oriente

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/tommyblue/oriente/store"
)

func Initialize(s *store.Store) *Oriente {
	o := &Oriente{
		store:        s,
		RunningGames: make(map[string]*Game),
	}
	if err := o.loadGames(); err != nil {
		log.Fatal(err)
	}
	return o
}

func NewGame(nPlayers int) *Game {
	g := &Game{}
	g.generateDeck()
	g.addPrize()
	g.generatePlayers(nPlayers)

	return g
}

func (o *Oriente) AddGame(token string, g *Game) {
	o.RunningGames[token] = g
}

func (o *Oriente) GetGame(token string) (*Game, error) {
	g, ok := o.RunningGames[token]
	if !ok {
		return nil, fmt.Errorf("game not found")
	}
	return g, nil
}

func (o *Oriente) SyncStore() error {
	log.Info("Syncing")
	for token, game := range o.RunningGames {
		g, err := json.Marshal(game)
		if err != nil {
			return err
		}
		if err := o.store.SyncGame(token, g); err != nil {
			return err
		}
	}
	return nil
}

func (o *Oriente) loadGames() error {
	games, err := o.store.LoadGames()
	if err != nil {
		return err
	}
	for token, g := range games {
		var game Game
		if err := json.Unmarshal(g, &game); err != nil {
			return err
		}
		log.Infof("Found game %s", token)
		o.RunningGames[token] = &game
	}
	return nil
}
