package oriente

import (
	"fmt"
	"math/rand"

	"github.com/tommyblue/oriente/utils"
)

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

func (g *Game) pickCard() Card {
	c := g.Deck[len(g.Deck)-1]
	g.Deck = g.Deck[:len(g.Deck)-1]
	return c
}

func (g *Game) generatePlayers(nPlayers int) {
	// This is the deck of cards with money
	coinsDeck := []int{2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4}

	for i := 0; i < nPlayers; i++ {
		mIdx := rand.Intn(len(coinsDeck))
		coin := coinsDeck[mIdx]
		p := &Player{Name: fmt.Sprintf("pl%d", i), Money: coin, ID: utils.TokenGenerator()}
		g.Players = append(g.Players, p)
		coinsDeck = append(coinsDeck[:mIdx], coinsDeck[mIdx+1:]...)
	}

	g.Token = g.Players[0]
}

func (g *Game) generateDeck() {

	tmpDeck := []Card{}
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

	// Get 4 cards fot the base
	base := make([]Card, 4)
	copy(base, tmpDeck[0:4])
	tmpDeck = append(tmpDeck[:0], tmpDeck[4:]...)

	// Get 4 cards and insert Geisha, shuffle
	wGeisha := make([]Card, 4)
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

func generateCard(c Character, n int) []Card {
	v := []Card{}
	for i := 0; i < n; i++ {
		v = append(v, Card{
			Name:  c.String(),
			Value: c,
		})
	}
	return v
}
