package entities

type Player struct {
	*Sprite
	MaxHealth float64
	Health    float64
	Damage    float64
}
