package oriente

import "github.com/tommyblue/oriente/store"

type Oriente struct {
	RunningGames map[string]*Game
	store        *store.Store
}

type Player struct {
	ID          string // Identifier of the player
	Name        string
	CurrentCard *Card   // Card in the hand of the player
	VisibleCard bool    // true if the card is visible to other players
	DidAction   bool    // true if the player played its action in this era
	Points      []*Card // The points of the player (sum of values of the cards)
	Managed     bool    // true if this player has been assigned to a client
}

type Action struct {
	Action         string `json:"action"`           // Can be "pass", "attack" or "use_ability". attack requires to fill in the player id
	SourcePlayerID string `json:"source_player_id"` // The player performing the action
	TargetPlayerID string `json:"target_player_id"` // The player target of the action
}
type Game struct {
	ID           string `json:"id"` // ID of the game
	Round        int    // Increase after any player action
	Players      []*Player
	Laborers     int
	Deck         []*Card
	TokenOwner   *Player // The last player that called the action or the first to play in an era
	NextPlayer   *Player // The next player to play the turn
	Prize        []*Card // The cards that will be won by the first playing player
	TempPrize    []*Card // The Prize is moved here when the player starts fulfilling his destiny
	CalledAction *Action // the action the player wants to perform ("pass", "attack" or "use_ability")
}

type Card struct {
	Name  string
	Value Character
}

type Character int

func (c Character) String() string {
	return [...]string{"Ninja", "Nofu", "Akindo", "Samurai", "Daimyo", "Maho-Tsukai", "Soryo", "Shogun", "Geisha"}[c]
}

const (
	Ninja Character = iota
	Nofu
	Akindo
	Samurai
	Daimyo
	MahoTsukai
	Soryo
	Shogun
	Geisha
)
