package entities

type Enemy struct {
	*Sprite
	MaxHealth float64
	Health    float64
	Damage    float64
	Aggro     bool
}
