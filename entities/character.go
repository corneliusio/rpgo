package entities

import "math"

type Character struct {
	*Sprite
	MaxHealth float64
	Health    float64
	Damage    float64
}

func (c *Character) EffectHealth(amount float64) {
	c.Health = math.Max(0, math.Min(c.MaxHealth, c.Health-amount))
}
