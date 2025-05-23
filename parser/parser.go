package parser

import (
	"context"
)

type Config struct {
	DataDirectory string
	OutDirectory  string
}

type Parser struct {
	Config
}

func New(config Config) *Parser {
	return &Parser{
		Config: config,
	}
}

// ParseSpells parses the spell data from the specified directory and writes it to the output directory.
func (p Parser) ParseSpells(ctx context.Context) error {
	return parseSpells(ctx, p.DataDirectory, p.OutDirectory)
}

func (p Parser) ParseMonsters(ctx context.Context) error {
	return parseMonsters(ctx, p.DataDirectory, p.OutDirectory)
}

// ParseItems parses the item data from the specified directory and writes it to the output directory.
func (p Parser) ParseItems(ctx context.Context) error {
	return parseItems(ctx, p.DataDirectory, p.OutDirectory)
}
