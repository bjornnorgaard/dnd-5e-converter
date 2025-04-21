package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ItemFile struct {
	Item []Item `json:"item"`
}

type Item struct {
	Name               string        `json:"name"`
	Type               string        `json:"type"`
	Rarity             string        `json:"rarity,omitempty"`
	Weight             float64       `json:"weight,omitempty"`
	Value              interface{}   `json:"value,omitempty"`
	Source             string        `json:"source"`
	Page               int           `json:"page,omitempty"`
	Entries            []interface{} `json:"entries,omitempty"`
	Attunement         interface{}   `json:"attunement,omitempty"`
	Tier               string        `json:"tier,omitempty"`
	RequiresAttunement interface{}   `json:"reqAttune,omitempty"`
	// Additional fields can be added as needed
}

// parseItems parses the item data from the specified directory and writes it to the output directory.
func parseItems(ctx context.Context, dataDirectory, outDirectory string) error {
	// Create output directory if it doesn't exist
	outDir := filepath.Join(outDirectory, "items")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Process the main items.json file
	if err := processItemFile(ctx, dataDirectory, outDir, "items.json"); err != nil {
		return fmt.Errorf("failed to process items.json: %w", err)
	}

	// Process the items-base.json file
	if err := processItemFile(ctx, dataDirectory, outDir, "items-base.json"); err != nil {
		return fmt.Errorf("failed to process items-base.json: %w", err)
	}

	return nil
}

// processItemFile processes a single item file and generates Markdown for each item
func processItemFile(ctx context.Context, dataDirectory, outDir, filename string) error {
	// Read and parse the item file
	filePath := filepath.Join(dataDirectory, filename)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read item file: %w", err)
	}

	var itemFile ItemFile
	if err := json.Unmarshal(fileData, &itemFile); err != nil {
		return fmt.Errorf("failed to parse item file: %w", err)
	}

	// Process each item
	for _, item := range itemFile.Item {
		mdContent, err := itemToMarkdown(item)
		if err != nil {
			return fmt.Errorf("failed to convert item to markdown: %w", err)
		}

		// Create a filename for the item
		safeItemName := strings.ReplaceAll(item.Name, "/", "-")
		safeItemName = strings.ReplaceAll(safeItemName, "\\", "-")
		safeItemName = strings.ReplaceAll(safeItemName, ":", "-")
		safeItemName = strings.ReplaceAll(safeItemName, "*", "-")
		safeItemName = strings.ReplaceAll(safeItemName, "?", "-")
		safeItemName = strings.ReplaceAll(safeItemName, "\"", "-")
		safeItemName = strings.ReplaceAll(safeItemName, "<", "-")
		safeItemName = strings.ReplaceAll(safeItemName, ">", "-")
		safeItemName = strings.ReplaceAll(safeItemName, "|", "-")

		mdFilePath := filepath.Join(outDir, safeItemName+".md")

		// Write the markdown file
		if err := os.WriteFile(mdFilePath, []byte(mdContent), 0644); err != nil {
			return fmt.Errorf("failed to write markdown file: %w", err)
		}
	}

	return nil
}

// itemToMarkdown converts an item to Markdown format
func itemToMarkdown(item Item) (string, error) {
	var md strings.Builder

	// Title
	md.WriteString(fmt.Sprintf("# %s\n\n", item.Name))

	// Basic info
	var typeRarity string
	if item.Rarity != "" {
		if item.Type != "" {
			typeRarity = fmt.Sprintf("*%s, %s*", item.Type, item.Rarity)
		} else {
			typeRarity = fmt.Sprintf("*%s*", item.Rarity)
		}
	} else {
		typeRarity = fmt.Sprintf("*%s*", item.Type)
	}
	md.WriteString(typeRarity + "\n\n")

	// Attunement
	if item.RequiresAttunement != nil {
		switch reqAttune := item.RequiresAttunement.(type) {
		case bool:
			if reqAttune {
				md.WriteString("*Requires attunement*\n\n")
			}
		case string:
			if reqAttune != "" {
				md.WriteString(fmt.Sprintf("*Requires attunement %s*\n\n", reqAttune))
			} else {
				md.WriteString("*Requires attunement*\n\n")
			}
		}
	} else if item.Attunement != nil {
		switch a := item.Attunement.(type) {
		case string:
			md.WriteString(fmt.Sprintf("*Requires attunement %s*\n\n", a))
		case bool:
			if a {
				md.WriteString("*Requires attunement*\n\n")
			}
		}
	}

	// Weight and value
	if item.Weight > 0 {
		md.WriteString(fmt.Sprintf("**Weight:** %.1f lb.\n\n", item.Weight))
	}

	if item.Value != nil {
		switch v := item.Value.(type) {
		case float64:
			md.WriteString(fmt.Sprintf("**Value:** %.0f gp\n\n", v))
		case map[string]interface{}:
			if quantity, ok := v["quantity"].(float64); ok {
				if unit, ok := v["unit"].(string); ok {
					md.WriteString(fmt.Sprintf("**Value:** %.0f %s\n\n", quantity, unit))
				}
			}
		}
	}

	// Description
	if item.Entries != nil && len(item.Entries) > 0 {
		for _, entry := range item.Entries {
			switch e := entry.(type) {
			case string:
				md.WriteString(processSpecialFormatting(e) + "\n\n")
			case map[string]interface{}:
				if entryType, ok := e["type"].(string); ok {
					if entryType == "list" && e["items"] != nil {
						if items, ok := e["items"].([]interface{}); ok {
							for _, item := range items {
								if itemStr, ok := item.(string); ok {
									md.WriteString("- " + processSpecialFormatting(itemStr) + "\n")
								}
							}
							md.WriteString("\n")
						}
					} else if entryType == "table" {
						// Handle tables if needed
					}
				}
			}
		}
	}

	// Source
	md.WriteString(fmt.Sprintf("**Source:** %s", item.Source))
	if item.Page > 0 {
		md.WriteString(fmt.Sprintf(", page %d", item.Page))
	}
	md.WriteString("\n")

	return md.String(), nil
}

// Use the existing processSpecialFormatting function from spells.go
