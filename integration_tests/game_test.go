package integration_tests

import (
	"fmt"
	"testing"

	"github.com/tommyblue/oriente/oriente"
)

func TestGame(t *testing.T) {
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
	p1, ok := g.GetPlayer(p1_tk)
	if !ok {
		t.Errorf("Can't get player 1 %s", p1_tk)
	}

	p2_tk, ok := g.GetFreePlayer()
	if !ok {
		t.Error("Can't get player 2")
	}
	p2, ok := g.GetPlayer(p2_tk)
	if !ok {
		t.Errorf("Can't get player 2 %s", p2_tk)
	}

	p3_tk, ok := g.GetFreePlayer()
	if !ok {
		t.Error("Can't get player 3")
	}
	p3, ok := g.GetPlayer(p3_tk)
	if !ok {
		t.Errorf("Can't get player 3 %s", p3_tk)
	}

	p4_tk, ok := g.GetFreePlayer()
	if !ok {
		t.Error("Can't get player 4")
	}
	p4, ok := g.GetPlayer(p4_tk)
	if !ok {
		t.Errorf("Can't get player 4 %s", p4_tk)
	}

	if ap := g.ActivePlayers(); ap != 4 {
		t.Errorf("Active players, want: 4, got: %d", ap)
	}
	if gs := g.GameStarted(); !gs {
		t.Errorf("Game started, want: true, got: %t", gs)
	}
	fmt.Println(p1, p2, p3, p4)
}
