package main

import (
	"math"

	"github.com/corneliusio/rpgo/entities"
)

type Camera struct {
	game *Game
	X, Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

func (c *Camera) FollowTarget(sprite *entities.Sprite, screenWidth, screenHeight float64, game *Game) {
	c.X = -sprite.X + float64(game.renderedTileSize/2) + screenWidth/2
	c.Y = -sprite.Y + float64(game.renderedTileSize/2) + screenHeight/2
}

func (c *Camera) Constrain(layer *TilemapLayerJSON, screenWidth, screenHeight float64, game *Game) {
	tilemapWidth := (float64(layer.Width) * game.renderedTileSize)
	tilemapHeight := (float64(layer.Height) * game.renderedTileSize) - (game.renderedTileSize * 2)

	c.X = math.Min(c.X, 0)
	c.Y = math.Min(c.Y, 0)
	c.X = math.Max(c.X, -(tilemapWidth - screenWidth))
	c.Y = math.Max(c.Y, -(tilemapHeight - screenHeight))
}
