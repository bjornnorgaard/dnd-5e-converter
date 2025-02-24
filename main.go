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
	)

	files, err := findFiles(path)
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

	md, err := spellMarkdown(ctx, spells[0])
	if err != nil {
		log.Fatal(err)
	}

	spellsPath := filepath.Join("out", "spells")
	err = os.MkdirAll(spellsPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(filepath.Join(spellsPath, md.FileName), []byte(md.Markdown), 0644)
	if err != nil {
		log.Fatal(err)
	}

	slog.InfoContext(ctx, "finished parsing files",
		slog.Duration("elapsed", time.Since(start)),
		slog.Group("spells",
			slog.Int("count", len(spells)),
		))
}

func findFiles(path string) ([]string, error) {
	var files []string
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
		return nil, fmt.Errorf("failed to walk directory '%s': %w", path, err)
	}
	return files, nil
}

type RenderResult struct {
	Markdown string
	FileName string
}

var schoolMap = map[string]string{
	"A":  "Abjuration",
	"C":  "Conjuration",
	"D":  "Divination",
	"E":  "Enchantment",
	"EV": "Evocation",
	"I":  "Illusion",
	"N":  "Necromancy",
}

func spellMarkdown(ctx context.Context, m json.RawMessage) (*RenderResult, error) {
	data, err := m.MarshalJSON()
	if err != nil {
		return nil, err
	}
	spell := map[string]any{}
	err = json.Unmarshal(data, &spell)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal spell")
	}

	res := RenderResult{}

	var structured struct {
		name   string
		level  int
		school string
		time   struct {
			number int
			unit   string
		}
		todo []string
	}

	for field, value := range spell {
		switch field {
		case "name":
			structured.name = value.(string)
			res.FileName = fmt.Sprintf("%s.md", value.(string))
		case "level":
			structured.level = int(value.(float64))
		case "school":
			school, ok := schoolMap[value.(string)]
			if !ok {
				slog.WarnContext(ctx, "unknown school", slog.String("school", value.(string)))
			}
			structured.school = school
		default:
			slog.WarnContext(ctx, "unhandled field", slog.String("field", field), slog.Any("value", value))
			structured.todo = append(structured.todo, fmt.Sprintf("#todo %s: %v", field, value))
		}
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("# %s", structured.name))
	lines = append(lines, fmt.Sprintf("*Level %d %s*", structured.level, structured.school))
	lines = append(lines, structured.todo...)
	res.Markdown = strings.Join(lines, "\n\n")

	return &res, nil
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
