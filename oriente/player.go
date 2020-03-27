package oriente

func (p *Player) totalPoints() int {
	var total int
	for _, p := range p.Points {
		total += p
	}
	return total
}
