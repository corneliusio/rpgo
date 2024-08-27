package main

import (
	"math"

	"github.com/corneliusio/rpgo/entities"
)

type Camera struct {
	X, Y, tileSize float64
}

func (c *Camera) FollowSprite(sprite *entities.Sprite, screenWidth, screenHeight float64) {
	c.X = -sprite.X + (screenWidth / 2) - (c.tileSize / 2)
	c.Y = -sprite.Y + (screenHeight / 2) - (c.tileSize / 2)
}

func (c *Camera) ConstrainToLayer(layer *TilemapLayerJSON, screenWidth, screenHeight float64) {
	tilemapWidth := (float64(layer.Width) * c.tileSize)
	tilemapHeight := (float64(layer.Height-2) * c.tileSize)

	c.X = math.Min(c.X, 0)
	c.Y = math.Min(c.Y, 0)
	c.X = math.Max(c.X, -(tilemapWidth - screenWidth))
	c.Y = math.Max(c.Y, -(tilemapHeight - screenHeight))
}
