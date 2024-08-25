package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"cornelius.dev/ebiten/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	baseScale   float64
	baseVector  float64
	tileSize    float64
	camera      *Camera
	player      *entities.Player
	enemies     []*entities.Enemy
	potions     []*entities.Potion
	tilemapImg  *ebiten.Image
	tilemapJSON *TilemapJSON
	tilesets    []Tileset
}

func (g *Game) Update() error {
	playerVector := g.baseVector

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		playerVector = g.baseVector * 1.5
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.X += playerVector
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.X -= playerVector
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Y += playerVector
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Y -= playerVector
	}

	for _, enemy := range g.enemies {
		if enemy.Aggro == true {
			if enemy.X+g.tileSize*g.baseScale <= g.player.X {
				enemy.X += g.baseVector / 2
			} else if enemy.X >= g.player.X+g.tileSize*g.baseScale {
				enemy.X -= g.baseVector / 2
			}

			if enemy.Y+g.tileSize*g.baseScale <= g.player.Y {
				enemy.Y += g.baseVector / 2
			} else if enemy.Y >= g.player.Y+g.tileSize*g.baseScale {
				enemy.Y -= g.baseVector / 2
			}
		}

		if g.player.X+g.tileSize*g.baseScale >= enemy.X && g.player.X <= enemy.X+g.tileSize && g.player.Y+g.tileSize*g.baseScale >= enemy.Y && g.player.Y <= enemy.Y+g.tileSize {
			enemy.EffectHealth(g.player.Damage)
			g.player.EffectHealth(enemy.Damage)

			if enemy.Health == 0 {
				g.RemoveEnemy(enemy)
			}

			if g.player.Health == 0 {
				log.Fatal("Game Over")
			}
		}
	}

	for _, potion := range g.potions {
		if g.player.X+g.tileSize*g.baseScale >= potion.X && g.player.X <= potion.X+g.tileSize && g.player.Y+g.tileSize*g.baseScale >= potion.Y && g.player.Y <= potion.Y+g.tileSize {
			g.player.EffectHealth(potion.Damage)
			g.RemovePotion(potion)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})
	opts := ebiten.DrawImageOptions{}

	g.camera.FollowTarget(
		g.player.X-float64(g.tileSize/2),
		g.player.Y-float64(g.tileSize/2),
		1280,
		960,
	)

	for n, layer := range g.tilemapJSON.Layers {
		if n == 0 {
			g.camera.Constrain(
				float64(layer.Width)*g.tileSize*g.baseScale,
				float64(layer.Height)*g.tileSize*g.baseScale,
				1280,
				960,
			)
		}

		for i, id := range layer.Data {
			if id == 0 {
				continue
			}

			x := float64(i % layer.Width)
			y := float64(i / layer.Height)
			x *= g.tileSize * g.baseScale
			y *= g.tileSize * g.baseScale

			image := g.tilesets[n].Image(id, g.tileSize, g.tileSize)
			offset := (float64(image.Bounds().Dy()) + g.tileSize) * g.baseScale

			opts.GeoM.Scale(g.baseScale, g.baseScale)
			opts.GeoM.Translate(x, y-offset)
			opts.GeoM.Translate(g.camera.X, g.camera.Y)

			screen.DrawImage(image, &opts)

			opts.GeoM.Reset()
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%d", int(g.player.Health)))

	opts.GeoM.Scale(g.baseScale, g.baseScale)
	opts.GeoM.Translate(g.player.X, g.player.Y)
	opts.GeoM.Translate(g.camera.X, g.camera.Y)

	screen.DrawImage(
		g.player.Image.SubImage(
			image.Rect(0, 0, int(g.tileSize), int(g.tileSize)),
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
				image.Rect(0, 0, int(g.tileSize), int(g.tileSize)),
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
				image.Rect(0, 0, int(g.tileSize), int(g.tileSize)),
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
