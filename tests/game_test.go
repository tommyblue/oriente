package tests

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/tommyblue/oriente/oriente"
)

func TestGameGetPlayers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping functional test")
	}
	g, err := oriente.NewGame(4)
	if err != nil {
		t.Error(err)
	}
	if gs := g.GameStarted(); gs {
		t.Errorf("Game started, want: false, got: %t", gs)
	}

	p1_tk, ok := g.GetFreePlayer()
	if !ok {
		t.Error("Can't get player 1")
	}
	p1 := g.GetPlayer(p1_tk)
	if p1 == nil {
		t.Errorf("Can't get player 1 %s", p1_tk)
	}

	p2_tk, ok := g.GetFreePlayer()
	if !ok {
		t.Error("Can't get player 2")
	}
	p2 := g.GetPlayer(p2_tk)
	if p2 == nil {
		t.Errorf("Can't get player 2 %s", p2_tk)
	}

	p3_tk, ok := g.GetFreePlayer()
	if !ok {
		t.Error("Can't get player 3")
	}
	p3 := g.GetPlayer(p3_tk)
	if p3 == nil {
		t.Errorf("Can't get player 3 %s", p3_tk)
	}

	p4_tk, ok := g.GetFreePlayer()
	if !ok {
		t.Error("Can't get player 4")
	}
	p4 := g.GetPlayer(p4_tk)
	if p4 == nil {
		t.Errorf("Can't get player 4 %s", p4_tk)
	}

	if ap := g.ActivePlayers(); ap != 4 {
		t.Errorf("Active players, want: 4, got: %d", ap)
	}
	if gs := g.GameStarted(); !gs {
		t.Errorf("Game started, want: true, got: %t", gs)
	}
}

func TestGame(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping functional test")
	}
	f := filepath.Join("testdata", "NewGame.json")
	raw_g, err := ioutil.ReadFile(f)
	if err != nil {
		t.Fatalf("opening %q: %v", f, err)
	}
	g, err := oriente.LoadGame(raw_g)
	if err != nil {
		t.Fatalf("unmarshalling %q: %v", raw_g, err)
	}
	if !g.GameStarted() {
		t.Fatalf("Game not started")
	}

	// 1st era

	if len(g.Prize) != 1 {
		t.Fatalf("Prize want: %d, got %d", 1, len(g.Prize))
	}
	round := 0
	if g.Round != round {
		t.Fatalf("Round want %d, got %d", round, g.Round)
	}
	if err := g.MakeAction(g.NextPlayerID, &oriente.Action{Action: "pass"}); err != nil {
		t.Fatalf("e1 a1 %v", err)
	}
	round++
	if g.NextPlayerID != g.Players[1].ID {
		t.Fatalf("Player: want %s, got: %s", g.Players[1].ID, g.NextPlayerID)
	}
	if g.Round != round {
		t.Fatalf("Round want %d, got %d", round, g.Round)
	}
	// make a full loop of "pass"
	if err := g.MakeAction(g.NextPlayerID, &oriente.Action{Action: "pass"}); err != nil {
		t.Fatalf("e1 a2 %v", err)
	}
	round++
	if g.Round != round {
		t.Fatalf("Round want %d, got %d", round, g.Round)
	}
	if err := g.MakeAction(g.NextPlayerID, &oriente.Action{Action: "pass"}); err != nil {
		t.Fatalf("e1 a3 %v", err)
	}
	round++
	if g.Round != round {
		t.Fatalf("Round want %d, got %d", round, g.Round)
	}
	if err := g.MakeAction(g.NextPlayerID, &oriente.Action{Action: "pass"}); err != nil {
		t.Fatalf("e1 a4 %v", err)
	}
	round++
	if g.Round != round {
		t.Fatalf("Round want %d, got %d", round, g.Round)
	}
	if g.NextPlayerID != g.Players[0].ID {
		t.Fatalf("Player: want %s, got: %s", g.Players[0].ID, g.NextPlayerID)
	}
	if len(g.Prize) != 2 {
		t.Fatalf("Prize want: %d, got %d", 2, len(g.Prize))
	}

	// 2nd era
	if g.CalledAction != nil {
		t.Fatalf("called action: want: %v, got: %v", nil, g.CalledAction)
	}
	// player_0 (mahotsukai) attacks player_1 (shogun)
	a := &oriente.Action{
		Action:         "attack",
		SourcePlayerID: g.NextPlayerID,
		TargetPlayerID: g.Players[1].ID,
	}
	if err := g.MakeAction(g.NextPlayerID, a); err != nil {
		t.Fatalf("e2 a1 %v", err)
	}
	round++
	if g.Round != round {
		t.Fatalf("Round want %d, got %d", round, g.Round)
	}
	if g.CalledAction == nil {
		t.Fatalf("called action: want: %v, got: %v", a, g.CalledAction)
	}
	if !g.Players[0].VisibleCard {
		t.Fatalf("p0 visible card, want: %t, got: %t", true, g.Players[0].VisibleCard)
	}
	if !g.Players[0].DidAction {
		t.Fatalf("p0 did action, want: %t, got: %t", true, g.Players[0].DidAction)
	}
	if g.CalledAction.TargetPlayerID == g.NextPlayerID {
		t.Fatalf("target player %s must not be the next (%s)", g.CalledAction.TargetPlayerID, g.NextPlayerID)
	}
	if g.NextPlayerID != g.Players[2].ID {
		t.Fatalf("e2 a1 next player, want %s got %s", g.Players[2].ID, g.NextPlayerID)
	}
	// player 1 is under attack, can't make actions
	a = &oriente.Action{
		Action:         oriente.AttackAction,
		SourcePlayerID: g.Players[1].ID,
		TargetPlayerID: g.Players[2].ID,
	}
	p0Card := g.Players[0].CurrentCard
	p1Card := g.Players[1].CurrentCard
	if err := g.MakeAction(g.Players[1].ID, a); err.Error() != "not the turn of the player" {
		t.Fatalf("player %s shouldn't do the action, err: %v", a.SourcePlayerID, err)
	}

	// player 2 and player 3 pass
	if err := g.MakeAction(g.NextPlayerID, &oriente.Action{Action: "pass"}); err != nil {
		t.Fatalf("e2 a2 %v", err)
	}
	round++
	if g.Round != round {
		t.Fatalf("Round want %d, got %d", round, g.Round)
	}
	if err := g.MakeAction(g.NextPlayerID, &oriente.Action{Action: "pass"}); err != nil {
		t.Fatalf("e2 a3 %v", err)
	}
	round++
	if g.Round != round {
		t.Fatalf("Round want %d, got %d", round, g.Round)
	}

	// Player 1 fulfills his destiny

	// p1 card: Maho-Tsukai (5)
	// p2 card: Shogun (7)
	// p2 wins

	// Points for p1, are now 3, the initial 1 + the 2 prizes
	if p := len(g.Players[0].Points); p != 3 {
		t.Fatalf("p1 points, want %d, got %d", 3, p)
	}
	// p1 card is new, not visible
	if g.Players[0].VisibleCard {
		t.Fatalf("p1 card visibility, want: %t, got: %t", false, g.Players[0].VisibleCard)
	}
	// p2 card is visible
	if !g.Players[1].VisibleCard {
		t.Fatalf("p2 card visibility, want: %t, got: %t", true, g.Players[1].VisibleCard)
	}
	if len(g.Players[1].Points) != 2 {
		t.Fatalf("p2 prize, want: %d, got: %d", 2, len(g.Players[1].Points))
	}

	if p1Card != g.Players[1].CurrentCard {
		t.Fatalf("p1 card, want: %s, got: %s", p1Card.Name, g.Players[1].CurrentCard.Name)
	}

	if p0Card == g.Players[0].CurrentCard {
		t.Fatalf("p0 card, want: different (%s), got: same (%s)", p0Card.Name, g.Players[0].CurrentCard.Name)
	}

	// p1 is the token owner, as such it's the next to play
	if g.NextPlayerID != g.TokenOwnerID || g.NextPlayerID != g.Players[0].ID {
		t.Fatalf("Wrong next player. Token owner: %s, nextPlayer: %s, want: %s", g.TokenOwnerID, g.NextPlayerID, g.Players[0].ID)
	}

	// New era
	// the prize is 1 again
	if len(g.Prize) != 1 {
		t.Fatalf("e2, prize, want: %d, got: %d", 1, len(g.Prize))
	}
}
