package oriente

import "testing"

func Test_generateCard(t *testing.T) {
	c := generateCard(Nofu, 3)
	if len(c) != 3 {
		t.Errorf("want: %d, got %d", 3, len(c))
	}
	for _, p := range c {
		if p.Value != Nofu {
			t.Errorf("want %d, got %d", Nofu, p.Value)
		}
	}
}

func Test_generateDeck(t *testing.T) {
	g := &Game{}
	g.generateDeck()

	// first 4 cards must not contain Geisha, Ninja or Shogun
	for i := 0; i < 4; i++ {
		c := g.Deck[i]
		if c.Value == Geisha {
			t.Errorf("Geisha in first 4 cards")
		}
		if c.Value == Ninja {
			t.Errorf("Ninja in first 4 cards")
		}
		if c.Value == Shogun {
			t.Errorf("Shogun in first 4 cards")
		}
	}

	// then 5 cards containing the Geisha but not Ninja or Shogun
	geisha := false
	for i := 4; i < 9; i++ {
		c := g.Deck[i]
		if c.Value == Geisha {
			geisha = true
		}
		if c.Value == Ninja {
			t.Errorf("Ninja in first cards 5-10")
		}
		if c.Value == Shogun {
			t.Errorf("Shogun in cards 5-10")
		}
	}
	if !geisha {
		t.Errorf("Can't find the Geisha")
	}

	var ninja int
	shogun := false
	for i := 9; i < len(g.Deck); i++ {
		c := g.Deck[i]
		if c.Value == Geisha {
			t.Errorf("Geisha in last cards")
		}
		if c.Value == Ninja {
			ninja++
		}
		if c.Value == Shogun {
			shogun = true
		}
	}
	if !shogun {
		t.Errorf("Can't find the Shogun")
	}
	if ninja != 3 {
		t.Errorf("Can't find 3 Ninjas")
	}
}

func Test_generatePlayers(t *testing.T) {
	g := &Game{}
	g.generateDeck()
	nPlayers := 7
	if err := g.generatePlayers(nPlayers); err != nil {
		t.Error(err)
	}

	if len(g.Players) != nPlayers {
		t.Errorf("players: want: %d, got: %d", nPlayers, len(g.Players))
	}
	if g.TokenOwnerID != g.Players[0].ID {
		t.Errorf("Wrong TokenOwner")
	}
	if g.NextPlayerID != g.Players[0].ID {
		t.Errorf("NextPlayer want %s, got %s", g.Players[0].ID, g.NextPlayerID)
	}
	for i, p := range g.Players {
		if p.CurrentCard == nil {
			t.Errorf("Player %d doesn't have a card", i)
		}
		if p.ID == "" {
			t.Errorf("Player %d doesn't have an ID", i)
		}
		if len(p.Points) != 1 {
			t.Errorf("Player %d points, want: %d, got: %d", i, 1, len(p.Points))
		}
		if p.Points[0].Value < 2 || p.Points[0].Value > 4 {
			t.Errorf("Player %d: want card value 2 <= v <= 4, got %d", i, p.Points[0].Value)
		}
	}
}

func Test_checkGeisha(t *testing.T) {
	g, _ := NewGame(4)
	haveWinner := g.checkGeisha()
	if haveWinner {
		t.Fatalf("newGame, winner want: %t, got %t", false, haveWinner)
	}

	// Add a Geisha
	g.Players[0].CurrentCard = generateCard(Geisha, 1)[0]
	haveWinner = g.checkGeisha()
	if !haveWinner {
		t.Fatalf("newGame, winner want: %t, got %t", true, haveWinner)
	}
}

func Test_checkWinnerNinja(t *testing.T) {
	g, _ := NewGame(4)
	winner := g.checkNinjas()
	if winner != nil {
		t.Fatalf("newGame, winner want: %v, got %v", nil, winner.ID)
	}
	p1 := g.Players[0]
	ninjas := generateCard(Ninja, 3)
	// Add Ninjas to p1 points
	p1.Points = append(p1.Points, ninjas[0])
	winner = g.checkNinjas()
	if winner != nil {
		t.Fatalf("1 ninja, winner want: %v, got %v", nil, winner.ID)
	}

	p1.Points = append(p1.Points, ninjas[1])
	winner = g.checkNinjas()
	if winner != nil {
		t.Fatalf("2 ninja, winner want: %v, got %v", nil, winner.ID)
	}

	p1.Points = append(p1.Points, ninjas[2])
	winner = g.checkNinjas()
	if winner == nil {
		t.Fatalf("3 ninja, winner want: %v, got %v", p1.ID, nil)
	}
}
