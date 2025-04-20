package parser

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestItemToMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		item     Item
		expected string
	}{
		{
			name: "Basic item",
			item: Item{
				Name:   "Longsword",
				Type:   "Weapon",
				Rarity: "Common",
				Weight: 3.0,
				Value:  float64(15),
				Source: "PHB",
				Page:   149,
				Entries: []interface{}{
					"A versatile weapon that can be used with one or two hands.",
				},
			},
			expected: "# Longsword\n\n*Weapon, Common*\n\n**Weight:** 3.0 lb.\n\n**Value:** 15 gp\n\nA versatile weapon that can be used with one or two hands.\n\n**Source:** PHB, page 149\n",
		},
		{
			name: "Magic item with attunement",
			item: Item{
				Name:               "Ring of Protection",
				Type:               "Ring",
				Rarity:             "Rare",
				Weight:             0.1,
				RequiresAttunement: interface{}(true),
				Source:             "DMG",
				Page:               191,
				Entries: []interface{}{
					"You gain a +1 bonus to AC and saving throws while wearing this ring.",
				},
			},
			expected: "# Ring of Protection\n\n*Ring, Rare*\n\n*Requires attunement*\n\n**Weight:** 0.1 lb.\n\nYou gain a +1 bonus to AC and saving throws while wearing this ring.\n\n**Source:** DMG, page 191\n",
		},
		{
			name: "Item with complex value",
			item: Item{
				Name:   "Gold Piece",
				Type:   "Currency",
				Weight: 0.02,
				Value: map[string]interface{}{
					"quantity": float64(1),
					"unit":     "gp",
				},
				Source: "PHB",
				Page:   143,
			},
			expected: "# Gold Piece\n\n*Currency*\n\n**Weight:** 0.0 lb.\n\n**Value:** 1 gp\n\n**Source:** PHB, page 143\n",
		},
		{
			name: "Item with list entries",
			item: Item{
				Name:   "Bag of Holding",
				Type:   "Wondrous Item",
				Rarity: "Uncommon",
				Weight: 15.0,
				Source: "DMG",
				Page:   153,
				Entries: []interface{}{
					"This bag has an interior space considerably larger than its outside dimensions.",
					map[string]interface{}{
						"type": "list",
						"items": []interface{}{
							"The bag can hold up to 500 pounds.",
							"The bag weighs 15 pounds, regardless of its contents.",
							"Retrieving an item from the bag requires an action.",
						},
					},
				},
			},
			expected: "# Bag of Holding\n\n*Wondrous Item, Uncommon*\n\n**Weight:** 15.0 lb.\n\nThis bag has an interior space considerably larger than its outside dimensions.\n\n- The bag can hold up to 500 pounds.\n- The bag weighs 15 pounds, regardless of its contents.\n- Retrieving an item from the bag requires an action.\n\n**Source:** DMG, page 153\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := itemToMarkdown(tt.item)
			if err != nil {
				t.Fatalf("itemToMarkdown() error = %v", err)
			}
			if result != tt.expected {
				t.Errorf("itemToMarkdown() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseItems_WithMockData(t *testing.T) {
	// Create temporary directories for test
	tempDir, err := os.MkdirTemp("", "items-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dataDir := filepath.Join(tempDir, "data")
	outDir := filepath.Join(tempDir, "out")

	if err := os.MkdirAll(dataDir, 0755); err != nil {
		t.Fatalf("Failed to create data dir: %v", err)
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create out dir: %v", err)
	}

	// Create mock items.json file
	itemsFile := ItemFile{
		Item: []Item{
			{
				Name:   "Longsword",
				Type:   "Weapon",
				Rarity: "Common",
				Weight: 3.0,
				Value:  float64(15),
				Source: "PHB",
				Page:   149,
				Entries: []interface{}{
					"A versatile weapon that can be used with one or two hands.",
				},
			},
			{
				Name:               "Ring of Protection",
				Type:               "Ring",
				Rarity:             "Rare",
				Weight:             0.1,
				RequiresAttunement: interface{}(true),
				Source:             "DMG",
				Page:               191,
				Entries: []interface{}{
					"You gain a +1 bonus to AC and saving throws while wearing this ring.",
				},
			},
		},
	}

	itemsData, err := json.Marshal(itemsFile)
	if err != nil {
		t.Fatalf("Failed to marshal items data: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dataDir, "items.json"), itemsData, 0644); err != nil {
		t.Fatalf("Failed to write items.json: %v", err)
	}

	// Create mock items-base.json file
	itemsBaseFile := ItemFile{
		Item: []Item{
			{
				Name:   "Gold Piece",
				Type:   "Currency",
				Weight: 0.02,
				Value: map[string]interface{}{
					"quantity": float64(1),
					"unit":     "gp",
				},
				Source: "PHB",
				Page:   143,
			},
		},
	}

	itemsBaseData, err := json.Marshal(itemsBaseFile)
	if err != nil {
		t.Fatalf("Failed to marshal items-base data: %v", err)
	}

	if err := os.WriteFile(filepath.Join(dataDir, "items-base.json"), itemsBaseData, 0644); err != nil {
		t.Fatalf("Failed to write items-base.json: %v", err)
	}

	// Run the parser
	ctx := context.Background()
	if err := parseItems(ctx, dataDir, outDir); err != nil {
		t.Fatalf("parseItems() error = %v", err)
	}

	// Check that the output files were created
	expectedFiles := []string{
		"Longsword.md",
		"Ring of Protection.md",
		"Gold Piece.md",
	}

	itemsOutDir := filepath.Join(outDir, "items")
	files, err := os.ReadDir(itemsOutDir)
	if err != nil {
		t.Fatalf("Failed to read output directory: %v", err)
	}

	if len(files) != len(expectedFiles) {
		t.Errorf("Expected %d files, got %d", len(expectedFiles), len(files))
	}

	for _, expectedFile := range expectedFiles {
		filePath := filepath.Join(itemsOutDir, expectedFile)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("Expected file %s does not exist", expectedFile)
		}
	}

	// Check the content of one of the files
	longswordContent, err := os.ReadFile(filepath.Join(itemsOutDir, "Longsword.md"))
	if err != nil {
		t.Fatalf("Failed to read Longsword.md: %v", err)
	}

	expectedContent := "# Longsword\n\n*Weapon, Common*\n\n**Weight:** 3.0 lb.\n\n**Value:** 15 gp\n\nA versatile weapon that can be used with one or two hands.\n\n**Source:** PHB, page 149\n"
	if string(longswordContent) != expectedContent {
		t.Errorf("Longsword.md content = %v, want %v", string(longswordContent), expectedContent)
	}
}
