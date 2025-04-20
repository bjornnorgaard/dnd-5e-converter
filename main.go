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

	err := converter.ParseSpells(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
