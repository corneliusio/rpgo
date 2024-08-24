package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"

	"cornelius.dev/ebiten/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	baseScale   float64
	baseVector  float64
	tileSize    int
	camera      *Camera
	player      *entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	tilemapImg  *ebiten.Image
	tilemapJSON *TilemapJSON
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += g.baseVector
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= g.baseVector
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += g.baseVector
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= g.baseVector
	}

	for _, enemy := range g.enemies {
		if enemy.Aggro == true {
			if enemy.X+float64(g.tileSize)*g.baseScale <= g.player.X {
				enemy.X += g.baseVector / 2
			} else if enemy.X >= g.player.X+float64(g.tileSize)*g.baseScale {
				enemy.X -= g.baseVector / 2
			}

			if enemy.Y+float64(g.tileSize)*g.baseScale <= g.player.Y {
				enemy.Y += g.baseVector / 2
			} else if enemy.Y >= g.player.Y+float64(g.tileSize)*g.baseScale {
				enemy.Y -= g.baseVector / 2
			}
		}

		if g.player.X+float64(g.tileSize)*g.baseScale >= enemy.X && g.player.X <= enemy.X+float64(g.tileSize) && g.player.Y+float64(g.tileSize)*g.baseScale >= enemy.Y && g.player.Y <= enemy.Y+float64(g.tileSize) {
			enemy.Health = math.Max(0, enemy.Health-10)
			g.player.Health = math.Max(0, g.player.Health-5)

			if enemy.Health == 0 {
				g.RemoveEnemy(enemy)
			}

			if g.player.Health == 0 {
				log.Fatal("Game Over")
			}
		}
	}

	for _, potion := range g.potions {
		if g.player.X+float64(g.tileSize)*g.baseScale >= potion.X && g.player.X <= potion.X+float64(g.tileSize) && g.player.Y+float64(g.tileSize)*g.baseScale >= potion.Y && g.player.Y <= potion.Y+float64(g.tileSize) {
			g.player.Health = math.Max(0, math.Min(g.player.MaxHealth, g.player.Health-potion.Damage))
			g.RemovePotion(potion)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	opts := ebiten.DrawImageOptions{}

	g.camera.FollowTarget(g.player.X-float64(g.tileSize/2), g.player.Y-float64(g.tileSize/2), 1280, 960)

	for n, layer := range g.tilemapJSON.Layers {
		if n == 0 {
			g.camera.Constrain(
				float64(layer.Width),
				float64(layer.Height),
				1280,
				960,
			)
		}

		for i, id := range layer.Data {
			x := i % layer.Width
			y := i / layer.Height
			x *= g.tileSize * int(g.baseScale)
			y *= g.tileSize * int(g.baseScale)

			srcX := (id - 1) % 22
			srcY := (id - 1) / 22
			srcX *= g.tileSize
			srcY *= g.tileSize

			opts.GeoM.Scale(g.baseScale, g.baseScale)
			opts.GeoM.Translate(float64(x), float64(y))
			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			screen.DrawImage(
				g.tilemapImg.SubImage(image.Rect(srcX, srcY, srcX+g.tileSize, srcY+g.tileSize)).(*ebiten.Image),
				&opts,
			)

			opts.GeoM.Reset()
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%d", int(g.player.Health)))

	opts.GeoM.Scale(g.baseScale, g.baseScale)
	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.camera.X, g.camera.Y)

	screen.DrawImage(
		g.player.Image.SubImage(
			image.Rect(0, 0, g.tileSize, g.tileSize),
		).(*ebiten.Image),
		&opts,
	)

	opts.GeoM.Reset()

	for _, sprite := range g.enemies {
		opts.GeoM.Scale(g.baseScale, g.baseScale)
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(
			sprite.Image.SubImage(
				image.Rect(0, 0, g.tileSize, g.tileSize),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}

	for _, sprite := range g.potions {
		opts.GeoM.Scale(g.baseScale, g.baseScale)
		opts.GeoM.Translate(sprite.X, sprite.Y)
		opts.GeoM.Translate(g.camera.X, g.camera.Y)
		screen.DrawImage(
			sprite.Image.SubImage(
				image.Rect(0, 0, g.tileSize, g.tileSize),
			).(*ebiten.Image),
			&opts,
		)

		opts.GeoM.Reset()
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func (g *Game) PlaceEnemy(enemy *entities.Enemy) {
	g.enemies = append(g.enemies, enemy)
}

func (g *Game) RemoveEnemy(enemy *entities.Enemy) {
	for i, e := range g.enemies {
		if e == enemy {
			g.enemies = append(g.enemies[:i], g.enemies[i+1:]...)
			break
		}
	}
}

func (g *Game) PlacePotion(potion *entities.Potion) {
	g.potions = append(g.potions, potion)
}

func (g *Game) RemovePotion(potion *entities.Potion) {
	for i, p := range g.potions {
		if p == potion {
			g.potions = append(g.potions[:i], g.potions[i+1:]...)
			break
		}
	}
}
