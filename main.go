package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var (
		ctx   = context.Background()
		path  = filepath.Join("..", "5etools-src", "data", "spells")
		start = time.Now()
		files []string
	)

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

	var spells []json.RawMessage
	for _, file := range files {
		logger := slog.With(slog.String("file", file))

		if strings.Contains(file, "fluff") {
			logger.InfoContext(ctx, "skipping fluff")
			continue
		}

		if !strings.Contains(file, "spells-") {
			logger.InfoContext(ctx, "skipping non spells file")
			continue
		}

		res, err := extractRawSpellList(file)
		if err != nil {
			log.Fatal(err)
		}

		spells = append(spells, res...)
	}

	_, err = spellMarkdown(ctx, spells[0])

	slog.InfoContext(ctx, "finished parsing files",
		slog.Duration("elapsed", time.Since(start)),
		slog.Group("spells",
			slog.Int("count", len(spells)),
		))
}

func spellMarkdown(ctx context.Context, m json.RawMessage) (*string, error) {
	data, err := m.MarshalJSON()
	if err != nil {
		return nil, err
	}
	spell := map[string]any{}
	err = json.Unmarshal(data, &spell)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal spell")
	}

	for field, value := range spell {
		slog.InfoContext(ctx, "spell",
			slog.String("field", field),
			slog.Any("value", value),
		)
	}

	return nil, nil
}

func extractRawSpellList(path string) ([]json.RawMessage, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var list struct {
		Spell []json.RawMessage `json:"spell"`
	}
	err = json.Unmarshal(bytes, &list)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json for file %s with error: %w", path, err)
	}

	if len(list.Spell) == 0 {
		return nil, fmt.Errorf("no spells found in %s", path)
	}

	return list.Spell, nil
}
