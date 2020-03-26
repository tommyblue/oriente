package oriente

var RunningGames map[string]*Game

func Initialize() {
	RunningGames = make(map[string]*Game)
}

func NewGame(nPlayers int) *Game {
	g := &Game{}
	g.generateDeck()
	g.addPrize()
	g.generatePlayers(nPlayers)

	return g
}
