package main

import (
	"image"
	"log"

	"cornelius.dev/ebiten/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
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

	baseScale := ebiten.Monitor().DeviceScaleFactor() * 1.5
	game := Game{
		baseScale:   baseScale,
		baseVector:  3,
		tileSize:    16,
		tilemapJSON: tilemap,
		tilemapImg:  tilemapImg,
		tilesets:    tilesets,
		colliders:   []image.Rectangle{},
		camera:      NewCamera(0.0, 0.0),
		player: &entities.Player{
			Sprite:    &entities.Sprite{Image: playerImg, X: 275.0, Y: 275.0, Speed: 1},
			Character: &entities.Character{MaxHealth: 100, Health: 80},
			Damage:    10,
		},
		enemies: []*entities.Enemy{
			{
				Sprite:    &entities.Sprite{Image: skeletonImg, X: 200.0, Y: 150.0, Speed: 0.5},
				Character: &entities.Character{MaxHealth: 50, Health: 50},
				Damage:    5,
				Aggro:     false,
			},
			{
				Sprite:    &entities.Sprite{Image: skeletonImg, X: 400.0, Y: 300.0, Speed: 0.5},
				Character: &entities.Character{MaxHealth: 50, Health: 50},
				Damage:    5,
				Aggro:     true,
			},
			{
				Sprite:    &entities.Sprite{Image: skeletonImg, X: 600.0, Y: 450.0, Speed: 0.5},
				Character: &entities.Character{MaxHealth: 50, Health: 50},
				Damage:    5,
				Aggro:     true,
			},
		},
		potions: []*entities.Potion{
			{
				Sprite: &entities.Sprite{Image: potionImg, X: 400.0, Y: 100.0},
				Damage: -20,
			},
			{
				Sprite: &entities.Sprite{Image: potionImg, X: 100.0, Y: 200.0},
				Damage: -20,
			},
			{
				Sprite: &entities.Sprite{Image: potionImg, X: 450.0, Y: 450.0},
				Damage: -20,
			},
			{
				Sprite: &entities.Sprite{Image: potionImg, X: 550.0, Y: 250.0},
				Damage: -20,
			},
		},
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
