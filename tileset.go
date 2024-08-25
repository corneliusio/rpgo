package main

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Tileset interface {
	Image(id int, tileWidth, tileHeight float64) *ebiten.Image
}

type TilesetJSON struct {
	Gid    int    `json:"firstgid"`
	Source string `json:"source"`
}

type UniformTilesetJSON struct {
	Path string `json:"image"`
}

type UniformTileset struct {
	img *ebiten.Image
	gid int
}

func (u *UniformTileset) Image(id int, tileWidth, tileHeight float64) *ebiten.Image {
	id -= u.gid

	w := int(tileWidth)
	h := int(tileHeight)

	srcX := id % 22
	srcY := id / 22
	srcX *= w
	srcY *= h

	return u.img.SubImage(
		image.Rect(
			srcX,
			srcY,
			srcX+w,
			srcY+h,
		),
	).(*ebiten.Image)
}

type TileJSON struct {
	Id     int    `json:"id"`
	Path   string `json:"image"`
	Width  int    `json:"imagewidth"`
	Height int    `json:"imageheight"`
}

type DynamicTilesetJSON struct {
	Tiles []*TileJSON `json:"tiles"`
}

type DynamicTileset struct {
	imgs []*ebiten.Image
	gid  int
}

func (d *DynamicTileset) Image(id int, tileWidth, tileHeight float64) *ebiten.Image {
	fmt.Println(id, d.gid)
	id -= d.gid
	return d.imgs[id]
}

func NewTileset(path string, gid int) (Tileset, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if strings.Contains(path, "buildings") {
		var dynamicTilesetJSON DynamicTilesetJSON
		err = json.Unmarshal(content, &dynamicTilesetJSON)
		if err != nil {
			return nil, err
		}

		dynamicTileset := DynamicTileset{}

		dynamicTileset.gid = gid
		dynamicTileset.imgs = make([]*ebiten.Image, len(dynamicTilesetJSON.Tiles))

		for i, tileJSON := range dynamicTilesetJSON.Tiles {
			dynamicPath := strings.Replace(tileJSON.Path, "../", "assets/", 1)
			img, _, err := ebitenutil.NewImageFromFile(dynamicPath)
			if err != nil {
				return nil, err
			}

			dynamicTileset.imgs[i] = img
		}

		return &dynamicTileset, nil
	}

	var uniformTilesetJSON UniformTilesetJSON
	err = json.Unmarshal(content, &uniformTilesetJSON)
	if err != nil {
		return nil, err
	}

	uniformTileset := UniformTileset{}
	uniformPath := strings.Replace(uniformTilesetJSON.Path, "../", "assets/", 1)
	img, _, err := ebitenutil.NewImageFromFile(uniformPath)
	if err != nil {
		return nil, err
	}

	uniformTileset.img = img
	uniformTileset.gid = gid

	return &uniformTileset, nil
}
