package main

import (
	"log"

	"cornelius.dev/ebiten/entities"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	ebiten.SetWindowSize(1280, 960)
	ebiten.SetWindowTitle("Ebiten")
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

	tilemapImg, _, err := ebitenutil.NewImageFromFile("assets/maps/tileset-floor.png")
	if err != nil {
		log.Fatal(err)
	}

	tilemapJSON, err := NewTileMapJSON("assets/maps/tilesets/tilemap.json")
	if err != nil {
		log.Fatal(err)
	}

	game := Game{
		baseScale:   ebiten.Monitor().DeviceScaleFactor() * 1.5,
		baseVector:  4,
		tileSize:    16,
		tilemapJSON: tilemapJSON,
		tilemapImg:  tilemapImg,
		camera:      NewCamera(50.0, 50.0),
		player: &entities.Player{
			Sprite:    &entities.Sprite{Image: playerImg, X: 275.0, Y: 275.0},
			MaxHealth: 100,
			Health:    80,
		},
		enemies: []*entities.Enemy{
			{
				Sprite:    &entities.Sprite{Image: skeletonImg, X: 200.0, Y: 150.0},
				MaxHealth: 50,
				Health:    50,
				Aggro:     false,
			},
			{
				Sprite:    &entities.Sprite{Image: skeletonImg, X: 400.0, Y: 300.0},
				MaxHealth: 50,
				Health:    50,
				Aggro:     true,
			},
			{
				Sprite:    &entities.Sprite{Image: skeletonImg, X: 600.0, Y: 450.0},
				MaxHealth: 50,
				Health:    50,
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
