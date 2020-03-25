package oriente

type Player struct {
	Name        string
	CurrentCard *Card
	Action      bool
	Money       int
	ID          string
	Managed     bool // true if this player has been assigned to a client
}

type Game struct {
	Players  []*Player
	Laborers int
	Deck     []Card
	Token    *Player // The first player playing the turn
	Prize    []Card  // The cards that will be won by the first playing player
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
