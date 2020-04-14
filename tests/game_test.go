package tests

import (
	"encoding/json"
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
	_, ok = g.GetPlayer(p1_tk)
	if !ok {
		t.Errorf("Can't get player 1 %s", p1_tk)
	}

	p2_tk, ok := g.GetFreePlayer()
	if !ok {
		t.Error("Can't get player 2")
	}
	_, ok = g.GetPlayer(p2_tk)
	if !ok {
		t.Errorf("Can't get player 2 %s", p2_tk)
	}

	p3_tk, ok := g.GetFreePlayer()
	if !ok {
		t.Error("Can't get player 3")
	}
	_, ok = g.GetPlayer(p3_tk)
	if !ok {
		t.Errorf("Can't get player 3 %s", p3_tk)
	}

	p4_tk, ok := g.GetFreePlayer()
	if !ok {
		t.Error("Can't get player 4")
	}
	_, ok = g.GetPlayer(p4_tk)
	if !ok {
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
	var g oriente.Game
	if err := json.Unmarshal(raw_g, &g); err != nil {
		t.Fatalf("unmarshalling %q: %v", raw_g, err)
	}
	if !g.GameStarted() {
		t.Fatalf("Game not started")
	}
	if len(g.Prize) != 1 {
		t.Fatalf("Prize want: %d, got %d", 1, len(g.Prize))
	}
	round := 0
	if g.Round != round {
		t.Fatalf("Round want %d, got %d", round, g.Round)
	}
	if err := g.MakeAction(g.NextPlayer, &oriente.Action{Action: "pass"}); err != nil {
		t.Fatalf("Initial action %v", err)
	}
	round++
	if g.NextPlayer != g.Players[1] {
		t.Fatalf("Player: want %s, got: %s", g.Players[1].ID, g.NextPlayer.ID)
	}
	if g.Round != round {
		t.Fatalf("Round want %d, got %d", round, g.Round)
	}
	// make a full loop of "pass"
	if err := g.MakeAction(g.NextPlayer, &oriente.Action{Action: "pass"}); err != nil {
		t.Fatalf("Action 2 %v", err)
	}
	round++
	if err := g.MakeAction(g.NextPlayer, &oriente.Action{Action: "pass"}); err != nil {
		t.Fatalf("Action 2 %v", err)
	}
	round++
	if err := g.MakeAction(g.NextPlayer, &oriente.Action{Action: "pass"}); err != nil {
		t.Fatalf("Action 2 %v", err)
	}
	round++
	if g.NextPlayer.ID != g.Players[0].ID {
		t.Fatalf("Player: want %s, got: %s", g.Players[0].ID, g.NextPlayer.ID)
	}
	if len(g.Prize) != 2 {
		t.Fatalf("Prize want: %d, got %d", 2, len(g.Prize))
	}
}
