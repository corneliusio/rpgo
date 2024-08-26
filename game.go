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
	baseScale        float64
	baseVector       float64
	realTileSize     float64
	renderedTileSize float64
	camera           *Camera
	player           *entities.Player
	enemies          []*entities.Enemy
	items            []*entities.Item
	tilemapImg       *ebiten.Image
	tilemapJSON      *TilemapJSON
	tilesets         []Tileset
	staticColliders  []entities.Collider
	dynamicColliders []entities.Collider
	drawOpts         ebiten.DrawImageOptions
}

func (g *Game) DrawSprite(screen *ebiten.Image, sprite *entities.Sprite) {
	g.drawOpts.GeoM.Scale(g.baseScale, g.baseScale)
	g.drawOpts.GeoM.Translate(sprite.X, sprite.Y)
	g.drawOpts.GeoM.Translate(g.camera.X, g.camera.Y)

	screen.DrawImage(
		sprite.Image.SubImage(
			image.Rect(0, 0, int(g.realTileSize), int(g.realTileSize)),
		).(*ebiten.Image),
		&g.drawOpts,
	)

	g.drawOpts.GeoM.Reset()
}

func (g *Game) CheckCollisionHorizontal(sprite *entities.Sprite) {
	sprite.X += sprite.Dx

	rect := sprite.Rect(g.renderedTileSize)

	for _, collider := range g.staticColliders {
		if collider.Self == sprite {
			continue
		}

		if collider.Rect.Overlaps(rect) {
			if sprite.Dx > 0 {
				sprite.X = float64(collider.Rect.Min.X) - g.renderedTileSize
			} else if sprite.Dx < 0 {
				sprite.X = float64(collider.Rect.Max.X)
			}
		}
	}

	for _, collider := range g.dynamicColliders {
		if collider.Self == sprite {
			continue
		}

		if collider.Rect.Overlaps(rect) {
			if sprite.Dx > 0 {
				sprite.X = float64(collider.Rect.Min.X) - g.renderedTileSize
			} else if sprite.Dx < 0 {
				sprite.X = float64(collider.Rect.Max.X)
			}
		}
	}
}

func (g *Game) CheckCollisionVertical(sprite *entities.Sprite) {
	sprite.Y += sprite.Dy

	rect := sprite.Rect(g.renderedTileSize)

	for _, collider := range g.staticColliders {
		if collider.Self == sprite {
			continue
		}

		if collider.Rect.Overlaps(rect) {
			if sprite.Dy > 0 {
				sprite.Y = float64(collider.Rect.Min.Y) - g.renderedTileSize
			} else if sprite.Dy < 0 {
				sprite.Y = float64(collider.Rect.Max.Y)
			}
		}
	}

	for _, collider := range g.dynamicColliders {
		if collider.Self == sprite {
			continue
		}

		if collider.Rect.Overlaps(rect) {
			if sprite.Dy > 0 {
				sprite.Y = float64(collider.Rect.Min.Y) - g.renderedTileSize
			} else if sprite.Dy < 0 {
				sprite.Y = float64(collider.Rect.Max.Y)
			}
		}
	}
}

func (g *Game) UpdatePlayerVectors() {
	vector := g.baseVector * g.player.Speed
	g.player.Dx = 0
	g.player.Dy = 0

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		vector = g.baseVector * 1.75
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

	g.player.NormalizeVector()
}

func (g *Game) UpdateAggroEnemyVectors(enemy *entities.Enemy) {
	vector := g.baseVector * enemy.Speed
	enemy.Dx = 0
	enemy.Dy = 0

	if enemy.X+g.renderedTileSize <= g.player.X {
		enemy.Dx = +vector
	} else if enemy.X >= g.player.X+g.renderedTileSize {
		enemy.Dx = -vector
	}

	if enemy.Y+g.renderedTileSize <= g.player.Y {
		enemy.Dy = +vector
	} else if enemy.Y >= g.player.Y+g.renderedTileSize {
		enemy.Dy = -vector
	}

	enemy.NormalizeVector()
}

func (g *Game) DrawLayer(screen *ebiten.Image, layer *TilemapLayerJSON, tsi int) {
	for i, id := range layer.Data {
		if id == 0 {
			continue
		}

		x := float64(i % layer.Width)
		y := float64(i / layer.Height)
		x *= g.renderedTileSize
		y *= g.renderedTileSize

		img := g.tilesets[tsi].Image(id, g.realTileSize, g.realTileSize)
		offset := (float64(img.Bounds().Dy()) + g.realTileSize) * g.baseScale

		g.drawOpts.GeoM.Scale(g.baseScale, g.baseScale)
		g.drawOpts.GeoM.Translate(x, y-offset)
		g.drawOpts.GeoM.Translate(g.camera.X, g.camera.Y)

		screen.DrawImage(img, &g.drawOpts)

		g.drawOpts.GeoM.Reset()

		if layer.Name == "objects" {
			bottom := y - offset + (float64(img.Bounds().Dy()) * g.baseScale)
			top := bottom - g.renderedTileSize*2 - g.renderedTileSize/4
			width := float64(img.Bounds().Dx()) * g.baseScale

			g.staticColliders = append(g.staticColliders, entities.Collider{
				Self: nil,
				Rect: image.Rect(
					int(x),
					int(top),
					int(x+width),
					int(bottom),
				),
			})
		}
	}
}

func (g *Game) Update() error {
	g.UpdatePlayerVectors()
	g.CheckCollisionHorizontal(g.player.Sprite)
	g.CheckCollisionVertical(g.player.Sprite)

	playerRect := g.player.Rect(g.renderedTileSize)

	for _, enemy := range g.enemies {
		if enemy.Aggro == true {
			g.UpdateAggroEnemyVectors(enemy)
			g.CheckCollisionHorizontal(enemy.Sprite)
			g.CheckCollisionVertical(enemy.Sprite)
		}

		if playerRect.Overlaps(enemy.Rect(g.renderedTileSize)) {
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

	for _, item := range g.items {
		if playerRect.Overlaps(item.Rect(g.renderedTileSize)) {
			if item.Damage != 0 {
				g.player.EffectHealth(item.Damage)
			}

			g.RemoveItem(item)
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{120, 180, 255, 255})

	screenWidth, screenHeight := ebiten.WindowSize()
	floor, layers := g.tilemapJSON.Layers[0], g.tilemapJSON.Layers[1:]

	g.camera.FollowTarget(g.player.Sprite, float64(screenWidth), float64(screenHeight), g)
	g.camera.Constrain(floor, float64(screenWidth), float64(screenHeight), g)

	g.DrawLayer(screen, floor, 0)

	g.dynamicColliders = []entities.Collider{}

	g.DrawSprite(screen, g.player.Sprite)
	g.dynamicColliders = append(g.dynamicColliders, entities.Collider{
		Self: g.player.Sprite,
		Rect: g.player.Rect(g.renderedTileSize),
	})

	for _, sprite := range g.items {
		g.DrawSprite(screen, sprite.Sprite)
	}

	for _, sprite := range g.enemies {
		g.DrawSprite(screen, sprite.Sprite)
		g.dynamicColliders = append(g.dynamicColliders, entities.Collider{
			Self: sprite.Sprite,
			Rect: sprite.Rect(g.renderedTileSize),
		})
	}

	for n, layer := range layers {
		g.DrawLayer(screen, layer, n+1)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%d", int(g.player.Health)))
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("%d %d", int(ebiten.ActualTPS()), int(ebiten.ActualFPS())))

	// for _, collider := range g.staticColliders {
	// 	vector.StrokeRect(
	// 		screen,
	// 		float32(collider.Rect.Min.X)+float32(g.camera.X),
	// 		float32(collider.Rect.Min.Y)+float32(g.camera.Y),
	// 		float32(collider.Rect.Dx()),
	// 		float32(collider.Rect.Dy()),
	// 		1.0,
	// 		color.RGBA{255, 0, 0, 255},
	// 		true,
	// 	)
	// }

	// for _, collider := range g.dynamicColliders {
	// 	vector.StrokeRect(
	// 		screen,
	// 		float32(collider.Rect.Min.X)+float32(g.camera.X),
	// 		float32(collider.Rect.Min.Y)+float32(g.camera.Y),
	// 		float32(collider.Rect.Dx()),
	// 		float32(collider.Rect.Dy()),
	// 		1.0,
	// 		color.RGBA{255, 0, 0, 255},
	// 		true,
	// 	)
	// }
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

func (g *Game) PlaceItem(item *entities.Item) {
	g.items = append(g.items, item)
}

func (g *Game) RemoveItem(item *entities.Item) {
	for i, p := range g.items {
		if p == item {
			g.items = append(g.items[:i], g.items[i+1:]...)
			break
		}
	}
}
