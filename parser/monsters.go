package parser

import (
	"context"
	"path/filepath"
)

// parseMonsters parses the monster data from the specified directory and writes it to the output directory.
func parseMonsters(ctx context.Context, dataDirectory, outDirectory string) error {
	var (
		spellsPath = filepath.Join(dataDirectory, "bestiary")
		indexPath  = filepath.Join(spellsPath, "index.json")
	)

	// TODO: Implement the actual parsing logic here.

	return nil
}
