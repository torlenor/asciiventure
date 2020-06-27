package game

func (g *Game) regenerationSystem() {
	for _, e := range g.entities {
		if e.Health != nil && e.IsDead == nil && e.Health.Regeneration != 0 && e.Health.CurrentHP < e.Health.HP {
			e.Health.CurrentHP += e.Health.Regeneration
			if e.Health.CurrentHP > e.Health.HP {
				e.Health.CurrentHP = e.Health.HP
			} else if e.Health.CurrentHP < 0 {
				g.killEntity(e)
			}
		}
	}
}
