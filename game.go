package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"cornelius.dev/ebiten/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
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
	colliders   []image.Rectangle
}

func (g *Game) CalcTileSize() float64 {
	return g.tileSize * g.baseScale
}

func (g *Game) CheckCollisionHorizontal(sprite *entities.Sprite, colliders []image.Rectangle) {
	rect := image.Rect(
		int(sprite.X),
		int(sprite.Y),
		int(sprite.X+g.CalcTileSize()),
		int(sprite.Y+g.CalcTileSize()),
	)

	for _, collider := range g.colliders {
		if collider.Overlaps(rect) {
			if sprite.Dx > 0 {
				sprite.X = float64(collider.Min.X) - g.CalcTileSize()
			} else if sprite.Dx < 0 {
				sprite.X = float64(collider.Max.X)
			}
		}
	}
}

func (g *Game) CheckCollisionVertical(sprite *entities.Sprite, colliders []image.Rectangle) {
	rect := image.Rect(
		int(sprite.X),
		int(sprite.Y),
		int(sprite.X+g.CalcTileSize()),
		int(sprite.Y+g.CalcTileSize()),
	)

	for _, collider := range g.colliders {
		if collider.Overlaps(rect) {
			if sprite.Dy > 0 {
				sprite.Y = float64(collider.Min.Y) - g.CalcTileSize()
			} else if sprite.Dy < 0 {
				sprite.Y = float64(collider.Max.Y)
			}
		}
	}
}

func (g *Game) UpdatePlayerVectors() {
	vector := g.baseVector * g.player.Speed
	g.player.Dx = 0
	g.player.Dy = 0

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		vector = g.baseVector * 1.5
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.Dx = +vector
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.Dx = -vector
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.Dy = +vector
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.Dy = -vector
	}
}

func (g *Game) Update() error {
	g.UpdatePlayerVectors()

	g.player.X += g.player.Dx
	g.CheckCollisionHorizontal(g.player.Sprite, g.colliders)

	g.player.Y += g.player.Dy
	g.CheckCollisionVertical(g.player.Sprite, g.colliders)

	for _, enemy := range g.enemies {
		if enemy.Aggro == true {
			// enemy.Dx = 0
			// enemy.Dy = 0

			if enemy.X+g.CalcTileSize() <= g.player.X {
				enemy.Dx = +g.baseVector * enemy.Speed
			} else if enemy.X >= g.player.X+g.CalcTileSize() {
				enemy.Dx = -g.baseVector * enemy.Speed
			}

			enemy.X += enemy.Dx
			g.CheckCollisionHorizontal(enemy.Sprite, g.colliders)

			if enemy.Y+g.CalcTileSize() <= g.player.Y {
				enemy.Dy = +g.baseVector * enemy.Speed
			} else if enemy.Y >= g.player.Y+g.CalcTileSize() {
				enemy.Dy = -g.baseVector * enemy.Speed
			}

			enemy.Y += enemy.Dy
			g.CheckCollisionVertical(enemy.Sprite, g.colliders)
		}

		if g.player.X+g.CalcTileSize() >= enemy.X && g.player.X <= enemy.X+g.tileSize && g.player.Y+g.CalcTileSize() >= enemy.Y && g.player.Y <= enemy.Y+g.tileSize {
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
		if g.player.X+g.CalcTileSize() >= potion.X && g.player.X <= potion.X+g.tileSize && g.player.Y+g.CalcTileSize() >= potion.Y && g.player.Y <= potion.Y+g.tileSize {
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
				float64(layer.Width)*g.CalcTileSize(),
				float64(layer.Height)*g.CalcTileSize(),
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
			x *= g.CalcTileSize()
			y *= g.CalcTileSize()

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

	for _, collider := range g.colliders {
		vector.StrokeRect(
			screen,
			float32(collider.Min.X)+float32(g.camera.X),
			float32(collider.Min.Y)+float32(g.camera.Y),
			float32(collider.Dx()),
			float32(collider.Dy()),
			1.0,
			color.RGBA{255, 0, 0, 255},
			true,
		)
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
