package main

import "math"

type Camera struct {
	X, Y float64
}

func NewCamera(x, y float64) *Camera {
	return &Camera{
		X: x,
		Y: y,
	}
}

func (c *Camera) FollowTarget(X, Y, screenWidth, screenHeight float64) {
	c.X = -X + screenWidth/2
	c.Y = -Y + screenHeight/2
}

func (c *Camera) Constrain(tilemapWidth, tilemapHeight, screenWidth, screenHeight float64) {
	c.X = math.Min(c.X, 0)
	c.Y = math.Min(c.Y, 0)
	c.X = math.Max(c.X, tilemapWidth-screenWidth)
	c.Y = math.Max(c.Y, tilemapHeight-screenHeight)
}
