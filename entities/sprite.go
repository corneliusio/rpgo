package entities

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Image  *ebiten.Image
	X, Y   float64
	Dx, Dy float64
	Speed  float64
}
