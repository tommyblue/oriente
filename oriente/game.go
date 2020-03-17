package oriente

func (g *Game) addPrize() {
	g.Prize = append(g.Prize, g.pickCard())
}

func (g *Game) pickCard() Card {
	c := g.Deck[len(g.Deck)-1]
	g.Deck = g.Deck[:len(g.Deck)-1]
	return c
}
