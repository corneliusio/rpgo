package main

import (
	"log"

	"github.com/corneliusio/rpgo/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	TPS := 120

	ebiten.SetTPS(TPS)
	ebiten.SetWindowTitle("RPGo")
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	playerImg, _, err := ebitenutil.NewImageFromFile("assets/images/character.png")
	if err != nil {
		log.Fatal(err)
	}

	skeletonImg, _, err := ebitenutil.NewImageFromFile("assets/images/skeleton.png")
	if err != nil {
		log.Fatal(err)
	}

	potionImg, _, err := ebitenutil.NewImageFromFile("assets/images/potion.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/images/floor.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemap, err := NewTileMapJSON("assets/tilemap.json")
	if err != nil {
		log.Fatal(err)
	}

	tilesets, err := tilemap.GenerateTilesets()
	if err != nil {
		log.Fatal(err)
	}

	realTileSize := float64(16)
	baseScale := ebiten.Monitor().DeviceScaleFactor() * 1.5
	renderedTileSize := realTileSize * baseScale
	game := Game{
		baseVector:       float64(210 / TPS),
		baseScale:        baseScale,
		realTileSize:     realTileSize,
		renderedTileSize: renderedTileSize,
		tilemapJSON:      tilemap,
		tilemapImg:       tilemapImg,
		tilesets:         tilesets,
		staticColliders:  []entities.Collider{},
		dynamicColliders: []entities.Collider{},
		camera:           &Camera{tileSize: renderedTileSize},
		drawOpts:         ebiten.DrawImageOptions{},
		player: &entities.Player{
			Character: &entities.Character{
				Sprite:    &entities.Sprite{Image: playerImg, X: 800.0, Y: 800.0, Speed: 1},
				MaxHealth: 12,
				Health:    8,
				Damage:    1,
			},
		},
		enemies: []*entities.Enemy{
			{
				Character: &entities.Character{
					Sprite:    &entities.Sprite{Image: skeletonImg, X: 200.0, Y: 150.0, Speed: 0.75},
					MaxHealth: 8,
					Health:    8,
					Damage:    2,
				},
				Aggro: false,
			},
			{
				Character: &entities.Character{
					Sprite:    &entities.Sprite{Image: skeletonImg, X: 400.0, Y: 300.0, Speed: 0.75},
					MaxHealth: 8,
					Health:    8,
					Damage:    2,
				},
				Aggro: true,
			},
			{
				Character: &entities.Character{
					Sprite:    &entities.Sprite{Image: skeletonImg, X: 600.0, Y: 450.0, Speed: 0.75},
					MaxHealth: 8,
					Health:    8,
					Damage:    2,
				},
				Aggro: true,
			},
		},
		items: []*entities.Item{
			{
				Sprite: &entities.Sprite{Image: potionImg, X: 400.0, Y: 100.0},
				Damage: -4,
			},
			{
				Sprite: &entities.Sprite{Image: potionImg, X: 100.0, Y: 200.0},
				Damage: -4,
			},
			{
				Sprite: &entities.Sprite{Image: potionImg, X: 450.0, Y: 450.0},
				Damage: -4,
			},
			{
				Sprite: &entities.Sprite{Image: potionImg, X: 550.0, Y: 250.0},
				Damage: -4,
			},
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
