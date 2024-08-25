package entities

type Enemy struct {
	*Sprite
	*Character
	Damage float64
	Aggro  bool
}
