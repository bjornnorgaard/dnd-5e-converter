package parser

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetSpellLevel(t *testing.T) {
	tests := []struct {
		level    int
		expected string
	}{
		{0, "Cantrip"},
		{1, "1st-level"},
		{2, "2nd-level"},
		{3, "3rd-level"},
		{4, "4th-level"},
		{9, "9th-level"},
	}

	for _, test := range tests {
		result := getSpellLevel(test.level)
		if result != test.expected {
			t.Errorf("getSpellLevel(%d) = %s; want %s", test.level, result, test.expected)
		}
	}
}

func TestGetSchoolName(t *testing.T) {
	tests := []struct {
		school   string
		expected string
	}{
		{"A", "Abjuration"},
		{"C", "Conjuration"},
		{"D", "Divination"},
		{"E", "Enchantment"},
		{"V", "Evocation"},
		{"I", "Illusion"},
		{"N", "Necromancy"},
		{"T", "Transmutation"},
		{"X", "X"}, // Unknown school should return itself
	}

	for _, test := range tests {
		result := getSchoolName(test.school)
		if result != test.expected {
			t.Errorf("getSchoolName(%s) = %s; want %s", test.school, result, test.expected)
		}
	}
}

func TestProcessSpecialFormatting(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"You hurl a bubble of acid. Choose one creature you can see within range, or choose two creatures you can see within range that are within 5 feet of each other. A target must succeed on a Dexterity saving throw or take {@damage 1d6} acid damage.",
			"You hurl a bubble of acid. Choose one creature you can see within range, or choose two creatures you can see within range that are within 5 feet of each other. A target must succeed on a Dexterity saving throw or take 1d6 acid damage.",
		},
		{
			"This spell's damage increases by {@dice 1d6} when you reach 5th level ({@damage 2d6}), 11th level ({@damage 3d6}), and 17th level ({@damage 4d6}).",
			"This spell's damage increases by 1d6 when you reach 5th level (2d6), 11th level (3d6), and 17th level (4d6).",
		},
	}

	for _, test := range tests {
		result := processSpecialFormatting(test.input)
		if result != test.expected {
			t.Errorf("processSpecialFormatting() failed\nInput: %s\nGot: %s\nWant: %s", test.input, result, test.expected)
		}
	}
}

func TestSpellToMarkdown(t *testing.T) {
	// Create a sample spell for testing
	spell := Spell{
		Name:   "Test Spell",
		Source: "PHB",
		Page:   123,
		Level:  1,
		School: "V",
		Time: []SpellTime{
			{Number: 1, Unit: "action"},
		},
		Range: SpellRange{
			Type: "point",
			Distance: SpellRangeDetail{
				Type:   "feet",
				Amount: 60,
			},
		},
		Components: SpellComponents{
			V: true,
			S: true,
			M: "a pinch of dust",
		},
		Duration: []SpellDuration{
			{Type: "instant"},
		},
		Entries: []interface{}{
			"This is a test spell description.",
			"It has multiple paragraphs.",
			map[string]interface{}{
				"type": "list",
				"items": []interface{}{
					"First effect",
					"Second effect",
				},
			},
		},
		EntriesHigher: []interface{}{
			"When cast using a spell slot of 2nd level or higher, the damage increases by {@dice 1d6} for each slot level above 1st.",
		},
		Classes: SpellClasses{
			FromClassList: []SpellClass{
				{Name: "Wizard", Source: "PHB"},
				{Name: "Sorcerer", Source: "PHB"},
			},
		},
	}

	md, err := spellToMarkdown(spell)
	if err != nil {
		t.Fatalf("spellToMarkdown() error = %v", err)
	}

	// Check that the markdown contains expected elements
	expectedElements := []string{
		"# Test Spell",
		"*1st-level Evocation*",
		"**Casting Time:** 1 action",
		"**Range:** 60 feet",
		"**Components:** V, S, M (a pinch of dust)",
		"**Duration:** Instantaneous",
		"This is a test spell description.",
		"It has multiple paragraphs.",
		"- First effect",
		"- Second effect",
		"**At Higher Levels:** When cast using a spell slot of 2nd level or higher, the damage increases by 1d6 for each slot level above 1st.",
		"**Classes:** Wizard, Sorcerer",
		"**Source:** PHB, page 123",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(md, expected) {
			t.Errorf("spellToMarkdown() output missing expected element: %s", expected)
		}
	}
}

// TestParseSpells_WithMockData tests the ParseSpells function with mock data
func TestParseSpells_WithMockData(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "spell-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a data directory structure
	dataDir := filepath.Join(tempDir, "data")
	spellsDir := filepath.Join(dataDir, "spells")
	outDir := filepath.Join(tempDir, "out")

	if err := os.MkdirAll(spellsDir, 0755); err != nil {
		t.Fatalf("Failed to create spells dir: %v", err)
	}

	// Create a mock index.json
	indexData := map[string]string{
		"TEST": "spells-test.json",
	}
	indexBytes, err := json.Marshal(indexData)
	if err != nil {
		t.Fatalf("Failed to marshal index data: %v", err)
	}
	if err := os.WriteFile(filepath.Join(spellsDir, "index.json"), indexBytes, 0644); err != nil {
		t.Fatalf("Failed to write index file: %v", err)
	}

	// Create a mock spell file
	spellFile := SpellFile{
		Spell: []Spell{
			{
				Name:   "Test Spell",
				Source: "TEST",
				Level:  0,
				School: "V",
				Time: []SpellTime{
					{Number: 1, Unit: "action"},
				},
				Range: SpellRange{
					Type: "point",
					Distance: SpellRangeDetail{
						Type:   "feet",
						Amount: 30,
					},
				},
				Components: SpellComponents{
					V: true,
				},
				Duration: []SpellDuration{
					{Type: "instant"},
				},
				Entries: []interface{}{
					"This is a test cantrip.",
					"It does {@damage 1d4} damage.",
				},
			},
		},
	}
	spellBytes, err := json.Marshal(spellFile)
	if err != nil {
		t.Fatalf("Failed to marshal spell data: %v", err)
	}
	if err := os.WriteFile(filepath.Join(spellsDir, "spells-test.json"), spellBytes, 0644); err != nil {
		t.Fatalf("Failed to write spell file: %v", err)
	}

	// Parse the spells
	if err := parseSpells(context.Background(), dataDir, outDir); err != nil {
		t.Fatalf("parseSpells() error = %v", err)
	}

	// Check that the output file was created
	outFile := filepath.Join(outDir, "spells", "Test Spell.md")
	if _, err := os.Stat(outFile); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %s", outFile)
	}

	// Read the output file and check its contents
	mdContent, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	// Check that the markdown contains expected elements
	expectedElements := []string{
		"# Test Spell",
		"*Cantrip Evocation*",
		"**Casting Time:** 1 action",
		"**Range:** 30 feet",
		"**Components:** V",
		"**Duration:** Instantaneous",
		"This is a test cantrip.",
		"It does 1d4 damage.",
		"**Source:** TEST",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(string(mdContent), expected) {
			t.Errorf("Output file missing expected element: %s", expected)
		}
	}
}

// TestSpellToMarkdown_EdgeCases tests the spellToMarkdown function with various edge cases
func TestSpellToMarkdown_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		spell    Spell
		expected []string
	}{
		{
			name: "Spell with scaling level dice",
			spell: Spell{
				Name:   "Edge Case Cantrip",
				Source: "TEST",
				Level:  0,
				School: "V",
				Time: []SpellTime{
					{Number: 1, Unit: "action"},
				},
				Range: SpellRange{
					Type: "point",
					Distance: SpellRangeDetail{
						Type:   "feet",
						Amount: 30,
					},
				},
				Components: SpellComponents{
					V: true,
				},
				Duration: []SpellDuration{
					{Type: "instant"},
				},
				Entries: []interface{}{
					"This is a test cantrip with scaling.",
				},
				ScalingLevelDice: map[string]interface{}{
					"label": "acid damage",
					"scaling": map[string]interface{}{
						"1":  "1d6",
						"5":  "2d6",
						"11": "3d6",
						"17": "4d6",
					},
				},
			},
			expected: []string{
				"**Scaling:**",
				"*acid damage*",
				"- 1st level: 1d6",
				"- 5th level: 2d6",
				"- 11th level: 3d6",
				"- 17th level: 4d6",
			},
		},
		{
			name: "Spell with concentration duration",
			spell: Spell{
				Name:   "Concentration Spell",
				Source: "TEST",
				Level:  1,
				School: "C",
				Time: []SpellTime{
					{Number: 1, Unit: "action"},
				},
				Range: SpellRange{
					Type: "point",
					Distance: SpellRangeDetail{
						Type: "self",
					},
				},
				Components: SpellComponents{
					V: true,
					S: true,
				},
				Duration: []SpellDuration{
					{
						Type: "concentration",
						Duration: map[string]interface{}{
							"amount": float64(10),
							"type":   "minute",
						},
					},
				},
				Entries: []interface{}{
					"This is a concentration spell.",
				},
			},
			expected: []string{
				"**Duration:** Concentration, up to 10 minute",
			},
		},
		{
			name: "Spell with complex components",
			spell: Spell{
				Name:   "Complex Components",
				Source: "TEST",
				Level:  3,
				School: "N",
				Time: []SpellTime{
					{Number: 1, Unit: "minute"},
				},
				Range: SpellRange{
					Type: "point",
					Distance: SpellRangeDetail{
						Type: "touch",
					},
				},
				Components: SpellComponents{
					V: true,
					S: true,
					M: map[string]interface{}{
						"text":    "a diamond worth at least 300 gp, which the spell consumes",
						"cost":    float64(300),
						"consume": true,
					},
					R: true,
				},
				Duration: []SpellDuration{
					{Type: "instant"},
				},
				Entries: []interface{}{
					"This spell has complex components.",
				},
			},
			expected: []string{
				"**Components:** V, S, M (a diamond worth at least 300 gp, which the spell consumes), R",
			},
		},
		{
			name: "Spell with SRD and Basic Rules flags",
			spell: Spell{
				Name:       "SRD Spell",
				Source:     "PHB",
				Level:      1,
				School:     "A",
				SRD:        true,
				BasicRules: true,
				Time: []SpellTime{
					{Number: 1, Unit: "action"},
				},
				Range: SpellRange{
					Type: "point",
					Distance: SpellRangeDetail{
						Type:   "feet",
						Amount: 30,
					},
				},
				Components: SpellComponents{
					V: true,
				},
				Duration: []SpellDuration{
					{Type: "instant"},
				},
				Entries: []interface{}{
					"This is an SRD spell.",
				},
			},
			expected: []string{
				"**Source:** PHB (SRD) (Basic Rules)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md, err := spellToMarkdown(tt.spell)
			if err != nil {
				t.Fatalf("spellToMarkdown() error = %v", err)
			}

			for _, expected := range tt.expected {
				if !strings.Contains(md, expected) {
					t.Errorf("spellToMarkdown() output missing expected element: %s", expected)
				}
			}
		})
	}
}
