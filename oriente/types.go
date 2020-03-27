package oriente

type Player struct {
	ID           string // Identifier of the player
	Name         string
	CurrentCard  *Card         // Card in the hand of the player
	VisibleCard  bool          // true if the card is visible to other players
	DidAction    bool          // true if the player played its action in this era
	CalledAction *PlayerAction // the action the player wants to perform ("attack" or "use_ability")
	Money        int           // Amount of money obtained
	Managed      bool          // true if this player has been assigned to a client
}

type PlayerAction struct {
	Action   string `json:"action"`    // Can be attack or use_ability. attack requires to fill in the player id
	PlayerID string `json:"player_id"` // The player to attack
}
type Game struct {
	ID         string // ID of the game
	Players    []*Player
	Laborers   int
	Deck       []*Card
	Token      *Player // The first player playing the turn
	NextPlayer *Player // The next player to play the turn
	Prize      []*Card // The cards that will be won by the first playing player
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
