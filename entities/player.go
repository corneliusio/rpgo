package entities

type Player struct {
	*Sprite
	*Character
	Damage float64
}
