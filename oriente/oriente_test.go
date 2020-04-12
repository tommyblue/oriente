package oriente

import (
	"testing"
)

func TestNewGamePlayers(t *testing.T) {
	tests := []struct {
		name     string
		nPlayers int
		want     bool
	}{
		{"less than 4", 3, false},
		{"more than 12", 13, false},
		{"low bound", 4, true},
		{"high bound", 12, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGame(tt.nPlayers)
			if (err != nil) == tt.want {
				t.Errorf("got %v, want %v", err, tt.want)
			}
		})
	}
}

func TestNewGame(t *testing.T) {
	nPlayers := 5
	g, err := NewGame(nPlayers)
	if err != nil {
		t.Error(err)
	}
	if len(g.Players) != nPlayers {
		t.Errorf("players, want: %d, got: %d", nPlayers, len(g.Players))
	}
	if g.Deck == nil {
		t.Errorf("Deck not generated")
	}
	if len(g.Players) == 0 {
		t.Errorf("Players not generated")
	}
	if len(g.Prize) == 0 {
		t.Errorf("Prize not generated")
	}
}
