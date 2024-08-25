package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type TilemapLayerJSON struct {
	Name   string `json:"name"`
	Data   []int  `json:"data"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type TilemapJSON struct {
	Layers   []*TilemapLayerJSON `json:"layers"`
	Tilesets []*TilesetJSON      `json:"tilesets"`
}

func (t *TilemapJSON) GenerateTilesets() ([]Tileset, error) {
	tilesets := make([]Tileset, len(t.Tilesets))

	for i, tilesetJSON := range t.Tilesets {
		path := fmt.Sprintf("assets/%s", tilesetJSON.Source)
		tileset, err := NewTileset(path, tilesetJSON.Gid)
		if err != nil {
			return nil, err
		}

		tilesets[i] = tileset
	}

	return tilesets, nil
}

func NewTileMapJSON(filepath string) (*TilemapJSON, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	tilemap := TilemapJSON{}
	err = json.Unmarshal(content, &tilemap)
	if err != nil {
		return nil, err
	}

	return &tilemap, nil
}
