package oriente

import (
	"math/rand"
)

func NewGame(players int) *Game {
	moneyDeck := []int{2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4}
	i := rand.Intn(len(moneyDeck))
	p1 := Player{Name: "p1", Money: moneyDeck[i]}
	removeIndex(moneyDeck, i)
	i = rand.Intn(len(moneyDeck))
	p2 := Player{Name: "p2", Money: moneyDeck[i]}

	g := &Game{
		Players: []Player{p1, p2},
		Deck:    newDeck(),
		Token:   &p1,
	}

	g.addPrize()

	return g
}

func removeIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

func newDeck() []Card {
	deck := []Card{}

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
	deck = append(deck, base...)
	deck = append(deck, wGeisha...)
	deck = append(deck, tmpDeck...)
	return deck
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
