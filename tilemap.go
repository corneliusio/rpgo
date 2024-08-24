package main

import (
	"encoding/json"
	"os"
)

type TilemapLayerJSON struct {
	Data   []int `json:"data"`
	Height int   `json:"height"`
	Width  int   `json:"width"`
}

type TilemapJSON struct {
	Layers []TilemapLayerJSON `json:"layers"`
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
