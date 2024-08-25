package entities

import "image"

type Collider struct {
	Self *Sprite
	Rect image.Rectangle
}
