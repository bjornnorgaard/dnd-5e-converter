package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func main() {
	path := filepath.Join("..", "5etools-src", "data", "spells")
	files := []string{}

	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		_, err = extractSpells(file)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func extractSpells(path string) ([]Spell, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var list Spells
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json for file %s with error: %w", path, err)
	}

	return list.Spell, nil
}

type Spells struct {
	Spell []Spell `json:"spell"`
}

type Spell struct {
	Name        string   `json:"name"`
	Source      string   `json:"source"`
	Page        int      `json:"page"`
	Srd         any      `json:"srd"`
	BasicRules  bool     `json:"basicRules"`
	ReprintedAs []string `json:"reprintedAs"`
	Level       int      `json:"level"`
	School      string   `json:"school"`
	Time        []struct {
		Number int    `json:"number"`
		Unit   string `json:"unit"`
	} `json:"time"`
	Range struct {
		Type     string `json:"type"`
		Distance struct {
			Type   string `json:"type"`
			Amount int    `json:"amount"`
		} `json:"distance"`
	} `json:"range"`
	Components struct {
		V bool `json:"v"`
		S bool `json:"s"`
	} `json:"components"`
	Duration []struct {
		Type string `json:"type"`
	} `json:"duration"`
	Entries          []string `json:"entries"`
	ScalingLevelDice struct {
		Label   string `json:"label"`
		Scaling struct {
			Field1 string `json:"1"`
			Field2 string `json:"5"`
			Field3 string `json:"11"`
			Field4 string `json:"17"`
		} `json:"scaling"`
	} `json:"scalingLevelDice"`
	DamageInflict []string `json:"damageInflict"`
	SavingThrow   []string `json:"savingThrow"`
	MiscTags      []string `json:"miscTags"`
	AreaTags      []string `json:"areaTags"`
}
