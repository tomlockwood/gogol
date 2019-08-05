package gol

import (
	"encoding/json"
	"io/ioutil"
)

// SaveContent used with the save function to write to a file
type SaveContent struct {
	Rules []Rule    `json:"rules"`
	Grid  [][]uint8 `json:"grid"`
}

// Save game of life to a file
func Save(G SaveContent, Filename string) {
	json, err := json.Marshal(G)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(Filename, json, 0644)
}

// Load a game from file
func Load(Filename string) Options {
	data, err := ioutil.ReadFile(Filename)
	if err != nil {
		panic(err)
	}
	gs := SaveContent{}
	json.Unmarshal(data, &gs)
	return Options{
		X:          0,
		Y:          0,
		Grid:       gs.Grid,
		RuleNumber: 0,
		Rules:      Rules{gs.Rules}}
}
