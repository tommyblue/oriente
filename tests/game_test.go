package tests

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/tommyblue/oriente/oriente"
)

func TestGameGetPlayers(t *testing.T) {
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
	a := oriente.Action{
		Action:         "attack",
		SourcePlayerID: g.NextPlayerID,
		TargetPlayerID: g.Players[1].ID,
	}
	if err := g.MakeAction(g.NextPlayerID, &a); err != nil {
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
}
