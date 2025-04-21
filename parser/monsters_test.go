package parser

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetSizeString(t *testing.T) {
	tests := []struct {
		size     string
		expected string
	}{
		{"T", "Tiny"},
		{"S", "Small"},
		{"M", "Medium"},
		{"L", "Large"},
		{"H", "Huge"},
		{"G", "Gargantuan"},
		{"X", "X"}, // Unknown size should return itself
	}

	for _, test := range tests {
		result := getSizeString(test.size)
		if result != test.expected {
			t.Errorf("getSizeString(%s) = %s; want %s", test.size, result, test.expected)
		}
	}
}

func TestGetAlignmentString(t *testing.T) {
	tests := []struct {
		alignment string
		expected  string
	}{
		{"L", "lawful"},
		{"N", "neutral"},
		{"C", "chaotic"},
		{"G", "good"},
		{"E", "evil"},
		{"U", "unaligned"},
		{"A", "any alignment"},
		{"neutral good", "neutral good"}, // Already expanded alignment should return itself
	}

	for _, test := range tests {
		result := getAlignmentString(test.alignment)
		if result != test.expected {
			t.Errorf("getAlignmentString(%s) = %s; want %s", test.alignment, result, test.expected)
		}
	}
}

func TestGetAbilityModifier(t *testing.T) {
	tests := []struct {
		score    int
		expected int
	}{
		{1, -5},
		{2, -4},
		{3, -4},
		{8, -1},
		{10, 0},
		{12, 1},
		{15, 2},
		{20, 5},
		{30, 10},
	}

	for _, test := range tests {
		result := getAbilityModifier(test.score)
		if result != test.expected {
			t.Errorf("getAbilityModifier(%d) = %d; want %d", test.score, result, test.expected)
		}
	}
}

func TestMonsterToMarkdown(t *testing.T) {
	// Create a sample monster for testing
	monster := Monster{
		Name:      "Test Monster",
		Source:    "MM",
		Page:      123,
		Size:      "M",
		Type:      "humanoid",
		Alignment: "neutral good",
		AC:        float64(15),
		HP: map[string]interface{}{
			"average": float64(45),
			"formula": "10d8+5",
		},
		Speed: map[string]interface{}{
			"walk": float64(30),
			"fly":  float64(60),
		},
		STR: 16,
		DEX: 14,
		CON: 12,
		INT: 10,
		WIS: 8,
		CHA: 6,
		Save: map[string]string{
			"str": "+5",
			"dex": "+4",
		},
		Skill: map[string]string{
			"Perception": "+2",
			"Stealth":    "+4",
		},
		Senses:    "darkvision 60 ft., passive Perception 12",
		Languages: "Common, Elvish",
		CR:        "2",
		Trait: []MonsterTrait{
			{
				Name:    "Keen Senses",
				Entries: []interface{}{"The monster has advantage on Wisdom (Perception) checks that rely on sight."},
			},
		},
		Action: []MonsterTrait{
			{
				Name:    "Multiattack",
				Entries: []interface{}{"The monster makes two attacks: one with its bite and one with its claws."},
			},
			{
				Name:    "Bite",
				Entries: []interface{}{"Melee Weapon Attack: +5 to hit, reach 5 ft., one target. Hit: 7 (1d8 + 3) piercing damage."},
			},
		},
	}

	md, err := monsterToMarkdown(monster)
	if err != nil {
		t.Fatalf("monsterToMarkdown() error = %v", err)
	}

	// Check that the markdown contains expected elements
	expectedElements := []string{
		"# Test Monster",
		"*Medium humanoid, neutral good*",
		"**Armor Class** 15",
		"**Hit Points** 45 (10d8+5)",
		"**Speed** 30 ft., fly 60 ft.",
		"|STR|DEX|CON|INT|WIS|CHA|",
		"|16 (+3)|14 (+2)|12 (+1)|10 (+0)|8 (-1)|6 (-2)|",
		"**Saving Throws** STR +5, DEX +4",
		"**Skills** Perception +2, Stealth +4",
		"**Senses** darkvision 60 ft., passive Perception 12",
		"**Languages** Common, Elvish",
		"**Challenge** 2",
		"## Traits",
		"***Keen Senses.*** The monster has advantage on Wisdom (Perception) checks that rely on sight.",
		"## Actions",
		"***Multiattack.*** The monster makes two attacks: one with its bite and one with its claws.",
		"***Bite.*** Melee Weapon Attack: +5 to hit, reach 5 ft., one target. Hit: 7 (1d8 + 3) piercing damage.",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(md, expected) {
			t.Errorf("monsterToMarkdown() output missing expected element: %s", expected)
		}
	}
}

// TestParseMonsters_WithMockData tests the ParseMonsters function with mock data
func TestParseMonsters_WithMockData(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "monster-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a data directory structure
	dataDir := filepath.Join(tempDir, "data")
	bestiaryDir := filepath.Join(dataDir, "bestiary")
	outDir := filepath.Join(tempDir, "out")

	if err := os.MkdirAll(bestiaryDir, 0755); err != nil {
		t.Fatalf("Failed to create bestiary dir: %v", err)
	}

	// Create a mock index.json
	indexData := map[string]string{
		"TEST": "bestiary-test.json",
	}
	indexBytes, err := json.Marshal(indexData)
	if err != nil {
		t.Fatalf("Failed to marshal index data: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bestiaryDir, "index.json"), indexBytes, 0644); err != nil {
		t.Fatalf("Failed to write index file: %v", err)
	}

	// Create a mock monster file
	monsterFile := MonsterFile{
		Monster: []Monster{
			{
				Name:      "Test Monster",
				Source:    "TEST",
				Size:      "M",
				Type:      "humanoid",
				Alignment: "neutral",
				AC:        float64(12),
				HP: map[string]interface{}{
					"average": float64(22),
					"formula": "4d8+4",
				},
				Speed: map[string]interface{}{
					"walk": float64(30),
				},
				STR:       10,
				DEX:       10,
				CON:       10,
				INT:       10,
				WIS:       10,
				CHA:       10,
				Senses:    "passive Perception 10",
				Languages: "Common",
				CR:        "1/4",
				Action: []MonsterTrait{
					{
						Name:    "Shortsword",
						Entries: []interface{}{"Melee Weapon Attack: +2 to hit, reach 5 ft., one target. Hit: 4 (1d6 + 1) piercing damage."},
					},
				},
			},
		},
	}
	monsterBytes, err := json.Marshal(monsterFile)
	if err != nil {
		t.Fatalf("Failed to marshal monster data: %v", err)
	}
	if err := os.WriteFile(filepath.Join(bestiaryDir, "bestiary-test.json"), monsterBytes, 0644); err != nil {
		t.Fatalf("Failed to write monster file: %v", err)
	}

	// Parse the monsters
	if err := parseMonsters(context.Background(), dataDir, outDir); err != nil {
		t.Fatalf("parseMonsters() error = %v", err)
	}

	// Check that the output file was created
	outFile := filepath.Join(outDir, "monsters", "Test Monster.md")
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
		"# Test Monster",
		"*Medium humanoid, neutral*",
		"**Armor Class** 12",
		"**Hit Points** 22 (4d8+4)",
		"**Speed** 30 ft.",
		"|STR|DEX|CON|INT|WIS|CHA|",
		"|10 (+0)|10 (+0)|10 (+0)|10 (+0)|10 (+0)|10 (+0)|",
		"**Senses** passive Perception 10",
		"**Languages** Common",
		"**Challenge** 1/4",
		"## Actions",
		"***Shortsword.*** Melee Weapon Attack: +2 to hit, reach 5 ft., one target. Hit: 4 (1d6 + 1) piercing damage.",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(string(mdContent), expected) {
			t.Errorf("Output file missing expected element: %s", expected)
		}
	}
}

// TestMonsterToMarkdown_EdgeCases tests the monsterToMarkdown function with various edge cases
func TestMonsterToMarkdown_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		monster     Monster
		expected    []string
		notExpected []string
	}{
		{
			name: "Monster with complex type",
			monster: Monster{
				Name:   "Complex Type Monster",
				Source: "TEST",
				Size:   "L",
				Type: map[string]interface{}{
					"type": "humanoid",
					"tags": []interface{}{"elf", "shapechanger"},
				},
				Alignment: "chaotic neutral",
				AC:        float64(14),
				HP: map[string]interface{}{
					"average": float64(65),
					"formula": "10d10+10",
				},
				Speed: map[string]interface{}{
					"walk": float64(30),
				},
				STR:       16,
				DEX:       14,
				CON:       12,
				INT:       10,
				WIS:       8,
				CHA:       6,
				Senses:    "darkvision 60 ft.",
				Languages: "Common, Elvish",
				CR:        "3",
			},
			expected: []string{
				"*Large humanoid (elf, shapechanger), chaotic neutral*",
			},
		},
		{
			name: "Monster with array alignment",
			monster: Monster{
				Name:      "Alignment Array Monster",
				Source:    "TEST",
				Size:      "M",
				Type:      "humanoid",
				Alignment: []interface{}{"chaotic", "evil"},
				AC:        float64(13),
				HP: map[string]interface{}{
					"average": float64(45),
					"formula": "10d8+5",
				},
				Speed: map[string]interface{}{
					"walk": float64(30),
				},
				STR:       10,
				DEX:       10,
				CON:       10,
				INT:       10,
				WIS:       10,
				CHA:       10,
				Senses:    "passive Perception 10",
				Languages: "Common",
				CR:        "1",
			},
			expected: []string{
				"*Medium humanoid, chaotic evil*",
			},
		},
		{
			name: "Monster with complex AC",
			monster: Monster{
				Name:      "Complex AC Monster",
				Source:    "TEST",
				Size:      "M",
				Type:      "humanoid",
				Alignment: "neutral",
				AC: []interface{}{
					map[string]interface{}{
						"ac":   float64(16),
						"from": []interface{}{"natural armor", "shield"},
					},
				},
				HP: map[string]interface{}{
					"average": float64(45),
					"formula": "10d8+5",
				},
				Speed: map[string]interface{}{
					"walk": float64(30),
				},
				STR:       10,
				DEX:       10,
				CON:       10,
				INT:       10,
				WIS:       10,
				CHA:       10,
				Senses:    "passive Perception 10",
				Languages: "Common",
				CR:        "1",
			},
			expected: []string{
				"**Armor Class** 16 (natural armor, shield)",
			},
		},
		{
			name: "Monster with array senses",
			monster: Monster{
				Name:      "Array Senses Monster",
				Source:    "TEST",
				Size:      "M",
				Type:      "humanoid",
				Alignment: "neutral",
				AC:        float64(13),
				HP: map[string]interface{}{
					"average": float64(45),
					"formula": "10d8+5",
				},
				Speed: map[string]interface{}{
					"walk": float64(30),
				},
				STR:       10,
				DEX:       10,
				CON:       10,
				INT:       10,
				WIS:       10,
				CHA:       10,
				Senses:    []interface{}{"darkvision 60 ft.", "tremorsense 30 ft.", "passive Perception 10"},
				Languages: "Common",
				CR:        "1",
			},
			expected: []string{
				"**Senses** darkvision 60 ft., tremorsense 30 ft., passive Perception 10",
			},
		},
		{
			name: "Monster with empty languages",
			monster: Monster{
				Name:      "Basilisk",
				Source:    "TEST",
				Size:      "M",
				Type:      "monstrosity",
				Alignment: "U",
				AC:        float64(15),
				HP: map[string]interface{}{
					"average": float64(52),
					"formula": "8d8+16",
				},
				Speed: map[string]interface{}{
					"walk": float64(20),
				},
				STR:       16,
				DEX:       8,
				CON:       15,
				INT:       2,
				WIS:       8,
				CHA:       7,
				Senses:    "darkvision 60 ft.",
				Languages: "",
				CR:        "3",
			},
			expected: []string{
				"*Medium monstrosity, unaligned*",
				"**Senses** darkvision 60 ft.",
				"**Challenge** 3",
			},
			notExpected: []string{
				"**Languages**",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			md, err := monsterToMarkdown(tt.monster)
			if err != nil {
				t.Fatalf("monsterToMarkdown() error = %v", err)
			}

			for _, expected := range tt.expected {
				if !strings.Contains(md, expected) {
					t.Errorf("monsterToMarkdown() output missing expected element: %s", expected)
				}
			}

			for _, notExpected := range tt.notExpected {
				if strings.Contains(md, notExpected) {
					t.Errorf("monsterToMarkdown() output contains unexpected element: %s", notExpected)
				}
			}
		})
	}
}
