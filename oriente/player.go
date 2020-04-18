package oriente

// return the points in the pile of the player (including the current card)
func (p *Player) totalPoints() int {
	var total int
	for _, p := range p.Points {
		total += int(p.Value)
	}
	total += int(p.CurrentCard.Value)
	return total
}

func (p *Player) highestCard() int {
	h := int(p.CurrentCard.Value)
	for _, c := range p.Points {
		if int(c.Value) > h {
			h = int(c.Value)
		}
	}
	return h
}
