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

// TestFireballSpell tests the spellToMarkdown function specifically with the Fireball spell
func TestFireballSpell(t *testing.T) {
	// Create a Fireball spell using the provided JSON
	fireballJSON := `{
		"name": "Fireball",
		"source": "PHB",
		"page": 241,
		"srd": true,
		"basicRules": true,
		"otherSources": [
			{
				"source": "RMR",
				"page": 53
			}
		],
		"reprintedAs": [
			"Fireball|XPHB"
		],
		"level": 3,
		"school": "V",
		"time": [
			{
				"number": 1,
				"unit": "action"
			}
		],
		"range": {
			"type": "point",
			"distance": {
				"type": "feet",
				"amount": 150
			}
		},
		"components": {
			"v": true,
			"s": true,
			"m": "a tiny ball of bat guano and sulfur"
		},
		"duration": [
			{
				"type": "instant"
			}
		],
		"entries": [
			"A bright streak flashes from your pointing finger to a point you choose within range and then blossoms with a low roar into an explosion of flame. Each creature in a 20-foot-radius sphere centered on that point must make a Dexterity saving throw. A target takes {@damage 8d6} fire damage on a failed save, or half as much damage on a successful one.",
			"The fire spreads around corners. It ignites flammable objects in the area that aren't being worn or carried."
		],
		"entriesHigherLevel": [
			{
				"type": "entries",
				"name": "At Higher Levels",
				"entries": [
					"When you cast this spell using a spell slot of 4th level or higher, the damage increases by {@scaledamage 8d6|3-9|1d6} for each slot level above 3rd."
				]
			}
		],
		"damageInflict": [
			"fire"
		],
		"savingThrow": [
			"dexterity"
		],
		"miscTags": [
			"OBJ"
		],
		"areaTags": [
			"S"
		],
		"hasFluffImages": true
	}`

	var fireball Spell
	if err := json.Unmarshal([]byte(fireballJSON), &fireball); err != nil {
		t.Fatalf("Failed to unmarshal Fireball JSON: %v", err)
	}

	md, err := spellToMarkdown(fireball)
	if err != nil {
		t.Fatalf("spellToMarkdown() error = %v", err)
	}

	// Check that the markdown contains expected elements
	expectedElements := []string{
		"# Fireball",
		"*3rd-level Evocation*",
		"**Casting Time:** 1 action",
		"**Range:** 150 feet",
		"**Components:** V, S, M (a tiny ball of bat guano and sulfur)",
		"**Duration:** Instantaneous",
		"A bright streak flashes from your pointing finger to a point you choose within range and then blossoms with a low roar into an explosion of flame.",
		"Each creature in a 20-foot-radius sphere centered on that point must make a Dexterity saving throw.",
		"A target takes 8d6 fire damage on a failed save, or half as much damage on a successful one.",
		"The fire spreads around corners. It ignites flammable objects in the area that aren't being worn or carried.",
		"**At Higher Levels:** When you cast this spell using a spell slot of 4th level or higher, the damage increases by 1d6 for each slot level above 3rd.",
		"**Damage Type:** fire",
		"**Saving Throw:** dexterity",
		"**Source:** PHB, page 241 (SRD) (Basic Rules)",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(md, expected) {
			t.Errorf("spellToMarkdown() output missing expected element: %s", expected)
		}
	}
}

// TestControlWaterSpell tests the spellToMarkdown function specifically with the Control Water spell
func TestControlWaterSpell(t *testing.T) {
	// Create a Control Water spell using the provided JSON
	controlWaterJSON := `{
		"name": "Control Water",
		"source": "XPHB",
		"page": 256,
		"freeRules2024": true,
		"level": 4,
		"school": "T",
		"time": [
			{
				"number": 1,
				"unit": "action"
			}
		],
		"range": {
			"type": "point",
			"distance": {
				"type": "feet",
				"amount": 300
			}
		},
		"components": {
			"v": true,
			"s": true,
			"m": "a mixture of water and dust"
		},
		"duration": [
			{
				"type": "timed",
				"duration": {
					"type": "minute",
					"amount": 10
				},
				"concentration": true
			}
		],
		"entries": [
			"Until the spell ends, you control any water inside an area you choose that is a Cube up to 100 feet on a side, using one of the following effects. As a {@action Magic|XPHB} action on your later turns, you can repeat the same effect or choose a different one.",
			{
				"type": "entries",
				"name": "Flood",
				"entries": [
					"You cause the water level of all standing water in the area to rise by as much as 20 feet. If you choose an area in a large body of water, you instead create a 20-foot tall wave that travels from one side of the area to the other and then crashes. Any Huge or smaller vehicles in the wave's path are carried with it to the other side. Any Huge or smaller vehicles struck by the wave have a {@chance 25|||Capsizes!|No effect} chance of capsizing.",
					"The water level remains elevated until the spell ends or you choose a different effect. If this effect produced a wave, the wave repeats on the start of your next turn while the flood effect lasts."
				]
			},
			{
				"type": "entries",
				"name": "Part Water",
				"entries": [
					"You part water in the area and create a trench. The trench extends across the spell's area, and the separated water forms a wall to either side. The trench remains until the spell ends or you choose a different effect. The water then slowly fills in the trench over the course of the next round until the normal water level is restored."
				]
			},
			{
				"type": "entries",
				"name": "Redirect Flow",
				"entries": [
					"You cause flowing water in the area to move in a direction you choose, even if the water has to flow over obstacles, up walls, or in other unlikely directions. The water in the area moves as you direct it, but once it moves beyond the spell's area, it resumes its flow based on the terrain. The water continues to move in the direction you chose until the spell ends or you choose a different effect."
				]
			},
			{
				"type": "entries",
				"name": "Whirlpool",
				"entries": [
					"You cause a whirlpool to form in the center of the area, which must be at least 50 feet square and 25 feet deep. The whirlpool lasts until you choose a different effect or the spell ends. The whirlpool is 5 feet wide at the base, up to 50 feet wide at the top, and 25 feet tall. Any creature in the water and within 25 feet of the whirlpool is pulled 10 feet toward it. When a creature enters the whirlpool for the first time on a turn or ends its turn there, it makes a Strength saving throw. On a failed save, the creature takes {@damage 2d8} Bludgeoning damage. On a successful save, the creature takes half as much damage. A creature can swim away from the whirlpool only if it first takes an action to pull away and succeeds on a Strength ({@skill Athletics}) check against your spell save DC."
				]
			}
		],
		"damageInflict": [
			"bludgeoning"
		],
		"savingThrow": [
			"strength"
		],
		"abilityCheck": [
			"strength"
		],
		"miscTags": [
			"FMV",
			"OBJ"
		],
		"areaTags": [
			"C"
		]
	}`

	var controlWater Spell
	if err := json.Unmarshal([]byte(controlWaterJSON), &controlWater); err != nil {
		t.Fatalf("Failed to unmarshal Control Water JSON: %v", err)
	}

	md, err := spellToMarkdown(controlWater)
	if err != nil {
		t.Fatalf("spellToMarkdown() error = %v", err)
	}

	// Check that the markdown contains expected elements
	expectedElements := []string{
		"# Control Water",
		"*4th-level Transmutation*",
		"**Casting Time:** 1 action",
		"**Range:** 300 feet",
		"**Components:** V, S, M (a mixture of water and dust)",
		"**Duration:** Concentration, up to 10 minute",
		"Until the spell ends, you control any water inside an area you choose that is a Cube up to 100 feet on a side, using one of the following effects. As a Magic action on your later turns, you can repeat the same effect or choose a different one.",
		"**Flood**",
		"You cause the water level of all standing water in the area to rise by as much as 20 feet.",
		"**Part Water**",
		"You part water in the area and create a trench.",
		"**Redirect Flow**",
		"You cause flowing water in the area to move in a direction you choose",
		"**Whirlpool**",
		"You cause a whirlpool to form in the center of the area",
		"On a failed save, the creature takes 2d8 Bludgeoning damage.",
		"**Damage Type:** bludgeoning",
		"**Saving Throw:** strength",
		"**Source:** XPHB, page 256",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(md, expected) {
			t.Errorf("spellToMarkdown() output missing expected element: %s", expected)
		}
	}
}
