package main

import (
	"context"
	"github.com/bjornnorgaard/dnd-5e-converter/parser"
	"log"
	"path/filepath"
)

func main() {
	converter := parser.New(parser.Config{
		DataDirectory: filepath.Join("..", "5etools-src", "data"),
		OutDirectory:  filepath.Join(".", "out"),
	})

	ctx := context.Background()

	err := converter.ParseSpells(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = converter.ParseMonsters(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = converter.ParseItems(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
