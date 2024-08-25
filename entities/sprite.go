package entities

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Image  *ebiten.Image
	X, Y   float64
	Dx, Dy float64
	Speed  float64
}

func (s *Sprite) Rect(renderedTileSize float64) image.Rectangle {
	return image.Rect(
		int(s.X),
		int(s.Y),
		int(s.X+renderedTileSize),
		int(s.Y+renderedTileSize),
	)
}

func (s *Sprite) NormalizeVector() {
	if s.Dx != 0 && s.Dy != 0 {
		s.Dx *= 0.75
		s.Dy *= 0.75
	}
}
