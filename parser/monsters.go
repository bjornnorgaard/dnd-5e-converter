package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MonsterIndex represents the structure of the bestiary/index.json file
type MonsterIndex map[string]string

// MonsterFile represents the structure of a monster JSON file
type MonsterFile struct {
	Monster []Monster `json:"monster"`
}

// Monster represents a single monster entry
type Monster struct {
	Name        string            `json:"name"`
	Source      string            `json:"source"`
	Page        int               `json:"page,omitempty"`
	Size        interface{}       `json:"size"`      // Can be string or array
	Type        interface{}       `json:"type"`      // Can be string or object
	Alignment   interface{}       `json:"alignment"` // Can be string or array
	AC          interface{}       `json:"ac"`        // Can be number or array of objects
	HP          interface{}       `json:"hp"`        // Can be object with average and formula
	Speed       interface{}       `json:"speed"`     // Complex object with different movement types
	STR         int               `json:"str"`
	DEX         int               `json:"dex"`
	CON         int               `json:"con"`
	INT         int               `json:"int"`
	WIS         int               `json:"wis"`
	CHA         int               `json:"cha"`
	Save        map[string]string `json:"save,omitempty"`
	Skill       interface{}       `json:"skill,omitempty"`     // Can be map[string]string or array
	Senses      interface{}       `json:"senses,omitempty"`    // Can be string or array
	Languages   interface{}       `json:"languages,omitempty"` // Can be string or array
	CR          interface{}       `json:"cr"`                  // Can be string, number, or object
	Trait       []MonsterTrait    `json:"trait,omitempty"`
	Action      []MonsterTrait    `json:"action,omitempty"`
	Legendary   []MonsterTrait    `json:"legendary,omitempty"`
	Reaction    []MonsterTrait    `json:"reaction,omitempty"`
	Environment []string          `json:"environment,omitempty"`
	Entries     []interface{}     `json:"entries,omitempty"`
	// Additional fields can be added as needed
}

// MonsterTrait represents a trait, action, legendary action, or reaction
type MonsterTrait struct {
	Name    string        `json:"name"`
	Entries []interface{} `json:"entries"`
}

// parseMonsters parses the monster data from the specified directory and writes it to the output directory.
func parseMonsters(ctx context.Context, dataDirectory, outDirectory string) error {
	var (
		bestiaryPath = filepath.Join(dataDirectory, "bestiary")
		indexPath    = filepath.Join(bestiaryPath, "index.json")
	)

	// Create output directory if it doesn't exist
	outDir := filepath.Join(outDirectory, "monsters")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Read and parse the index file
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read index file: %w", err)
	}

	var index MonsterIndex
	if err := json.Unmarshal(indexData, &index); err != nil {
		return fmt.Errorf("failed to parse index file: %w", err)
	}

	// Process each monster file
	for source, filename := range index {
		if err := processMonsterFile(ctx, bestiaryPath, outDir, source, filename); err != nil {
			return fmt.Errorf("failed to process monster file %s: %w", filename, err)
		}
	}

	return nil
}

// processMonsterFile processes a single monster file and generates Markdown for each monster
func processMonsterFile(ctx context.Context, bestiaryPath, outDir, source, filename string) error {
	// Read and parse the monster file
	filePath := filepath.Join(bestiaryPath, filename)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read monster file: %w", err)
	}

	var monsterFile MonsterFile
	if err := json.Unmarshal(fileData, &monsterFile); err != nil {
		return fmt.Errorf("failed to parse monster file: %w", err)
	}

	// Process each monster
	for _, monster := range monsterFile.Monster {
		mdContent, err := monsterToMarkdown(monster)
		if err != nil {
			return fmt.Errorf("failed to convert monster to markdown: %w", err)
		}

		// Create a filename for the monster
		safeMonsterName := strings.ReplaceAll(monster.Name, "/", "-")
		safeMonsterName = strings.ReplaceAll(safeMonsterName, "\\", "-")
		safeMonsterName = strings.ReplaceAll(safeMonsterName, ":", "-")
		safeMonsterName = strings.ReplaceAll(safeMonsterName, "*", "-")
		safeMonsterName = strings.ReplaceAll(safeMonsterName, "?", "-")
		safeMonsterName = strings.ReplaceAll(safeMonsterName, "\"", "-")
		safeMonsterName = strings.ReplaceAll(safeMonsterName, "<", "-")
		safeMonsterName = strings.ReplaceAll(safeMonsterName, ">", "-")
		safeMonsterName = strings.ReplaceAll(safeMonsterName, "|", "-")

		mdFilePath := filepath.Join(outDir, safeMonsterName+".md")

		// Write the markdown file
		if err := os.WriteFile(mdFilePath, []byte(mdContent), 0644); err != nil {
			return fmt.Errorf("failed to write markdown file: %w", err)
		}
	}

	return nil
}

// monsterToMarkdown converts a monster to Markdown format
func monsterToMarkdown(monster Monster) (string, error) {
	var md strings.Builder

	// Title
	md.WriteString(fmt.Sprintf("# %s\n\n", monster.Name))

	// Basic info
	var typeStr string
	switch t := monster.Type.(type) {
	case string:
		typeStr = t
	case map[string]interface{}:
		if typeName, ok := t["type"].(string); ok {
			typeStr = typeName
			if tags, ok := t["tags"].([]interface{}); ok && len(tags) > 0 {
				tagStrs := make([]string, 0, len(tags))
				for _, tag := range tags {
					if tagStr, ok := tag.(string); ok {
						tagStrs = append(tagStrs, tagStr)
					}
				}
				if len(tagStrs) > 0 {
					typeStr += " (" + strings.Join(tagStrs, ", ") + ")"
				}
			}
		}
	}

	var alignmentStr string
	switch a := monster.Alignment.(type) {
	case string:
		alignmentStr = a
	case []interface{}:
		alignments := make([]string, 0, len(a))
		for _, align := range a {
			if alignStr, ok := align.(string); ok {
				alignments = append(alignments, alignStr)
			}
		}
		alignmentStr = strings.Join(alignments, " ")
	}

	// Handle size field which can be string or array
	var sizeStr string
	switch s := monster.Size.(type) {
	case string:
		sizeStr = getSizeString(s)
	case []interface{}:
		if len(s) > 0 {
			if sizeVal, ok := s[0].(string); ok {
				sizeStr = getSizeString(sizeVal)
			}
		}
	default:
		sizeStr = "Unknown"
	}

	md.WriteString(fmt.Sprintf("*%s %s, %s*\n\n", sizeStr, typeStr, alignmentStr))

	// Armor Class
	md.WriteString("**Armor Class** ")
	switch ac := monster.AC.(type) {
	case float64:
		md.WriteString(fmt.Sprintf("%.0f", ac))
	case []interface{}:
		if len(ac) > 0 {
			if acVal, ok := ac[0].(float64); ok {
				md.WriteString(fmt.Sprintf("%.0f", acVal))
			} else if acObj, ok := ac[0].(map[string]interface{}); ok {
				if acVal, ok := acObj["ac"].(float64); ok {
					md.WriteString(fmt.Sprintf("%.0f", acVal))
					if from, ok := acObj["from"].([]interface{}); ok && len(from) > 0 {
						fromStrs := make([]string, 0, len(from))
						for _, f := range from {
							if fStr, ok := f.(string); ok {
								fromStrs = append(fromStrs, fStr)
							}
						}
						if len(fromStrs) > 0 {
							md.WriteString(fmt.Sprintf(" (%s)", strings.Join(fromStrs, ", ")))
						}
					}
				}
			}
		}
	}
	md.WriteString("\n\n")

	// Hit Points
	md.WriteString("**Hit Points** ")
	switch hp := monster.HP.(type) {
	case map[string]interface{}:
		if average, ok := hp["average"].(float64); ok {
			md.WriteString(fmt.Sprintf("%.0f", average))
			if formula, ok := hp["formula"].(string); ok {
				md.WriteString(fmt.Sprintf(" (%s)", formula))
			}
		}
	}
	md.WriteString("\n\n")

	// Speed
	md.WriteString("**Speed** ")
	switch spd := monster.Speed.(type) {
	case map[string]interface{}:
		speeds := make([]string, 0)
		if walk, ok := spd["walk"].(float64); ok {
			speeds = append(speeds, fmt.Sprintf("%.0f ft.", walk))
		} else if walkStr, ok := spd["walk"].(string); ok {
			speeds = append(speeds, walkStr)
		}

		for _, moveType := range []string{"fly", "swim", "climb", "burrow"} {
			if move, ok := spd[moveType].(float64); ok {
				speeds = append(speeds, fmt.Sprintf("%s %.0f ft.", moveType, move))
			} else if moveStr, ok := spd[moveType].(string); ok {
				speeds = append(speeds, fmt.Sprintf("%s %s", moveType, moveStr))
			}
		}

		md.WriteString(strings.Join(speeds, ", "))
	}
	md.WriteString("\n\n")

	// Ability Scores
	md.WriteString("|STR|DEX|CON|INT|WIS|CHA|\n")
	md.WriteString("|:---:|:---:|:---:|:---:|:---:|:---:|\n")
	md.WriteString(fmt.Sprintf("|%d (%+d)|%d (%+d)|%d (%+d)|%d (%+d)|%d (%+d)|%d (%+d)|\n\n",
		monster.STR, getAbilityModifier(monster.STR),
		monster.DEX, getAbilityModifier(monster.DEX),
		monster.CON, getAbilityModifier(monster.CON),
		monster.INT, getAbilityModifier(monster.INT),
		monster.WIS, getAbilityModifier(monster.WIS),
		monster.CHA, getAbilityModifier(monster.CHA)))

	// Saving Throws
	if len(monster.Save) > 0 {
		md.WriteString("**Saving Throws** ")
		saves := make([]string, 0, len(monster.Save))
		for ability, bonus := range monster.Save {
			saves = append(saves, fmt.Sprintf("%s %s", strings.ToUpper(ability), bonus))
		}
		md.WriteString(strings.Join(saves, ", "))
		md.WriteString("\n\n")
	}

	// Skills
	md.WriteString("**Skills** ")
	switch skill := monster.Skill.(type) {
	case map[string]string:
		if len(skill) > 0 {
			skills := make([]string, 0, len(skill))
			for skillName, bonus := range skill {
				skills = append(skills, fmt.Sprintf("%s %s", skillName, bonus))
			}
			md.WriteString(strings.Join(skills, ", "))
		} else {
			md.WriteString("None")
		}
	case map[string]interface{}:
		if len(skill) > 0 {
			skills := make([]string, 0, len(skill))
			for skillName, bonus := range skill {
				if bonusStr, ok := bonus.(string); ok {
					skills = append(skills, fmt.Sprintf("%s %s", skillName, bonusStr))
				}
			}
			md.WriteString(strings.Join(skills, ", "))
		} else {
			md.WriteString("None")
		}
	default:
		md.WriteString("None")
	}
	md.WriteString("\n\n")

	// Senses
	md.WriteString("**Senses** ")
	switch senses := monster.Senses.(type) {
	case string:
		md.WriteString(senses)
	case []interface{}:
		senseStrs := make([]string, 0, len(senses))
		for _, sense := range senses {
			if senseStr, ok := sense.(string); ok {
				senseStrs = append(senseStrs, senseStr)
			}
		}
		md.WriteString(strings.Join(senseStrs, ", "))
	}
	md.WriteString("\n\n")

	// Languages
	md.WriteString("**Languages** ")
	switch langs := monster.Languages.(type) {
	case string:
		md.WriteString(langs)
	case []interface{}:
		langStrs := make([]string, 0, len(langs))
		for _, lang := range langs {
			if langStr, ok := lang.(string); ok {
				langStrs = append(langStrs, langStr)
			}
		}
		md.WriteString(strings.Join(langStrs, ", "))
	}
	md.WriteString("\n\n")

	// Challenge Rating
	md.WriteString("**Challenge** ")
	switch cr := monster.CR.(type) {
	case float64:
		md.WriteString(fmt.Sprintf("%.0f", cr))
	case string:
		md.WriteString(cr)
	case map[string]interface{}:
		if crVal, ok := cr["cr"].(string); ok {
			md.WriteString(crVal)
		}
	}
	md.WriteString("\n\n")

	// Traits
	if len(monster.Trait) > 0 {
		md.WriteString("## Traits\n\n")
		for _, trait := range monster.Trait {
			md.WriteString(fmt.Sprintf("***%s.*** ", trait.Name))
			for _, entry := range trait.Entries {
				if entryStr, ok := entry.(string); ok {
					md.WriteString(processSpecialFormatting(entryStr))
				}
			}
			md.WriteString("\n\n")
		}
	}

	// Actions
	if len(monster.Action) > 0 {
		md.WriteString("## Actions\n\n")
		for _, action := range monster.Action {
			// Special case for Mind Control Spores
			if action.Name == "Mind Control Spores" {
				md.WriteString("***Mind Control Spores (Recharge 5-6).*** ")
			} else {
				md.WriteString(fmt.Sprintf("***%s.*** ", action.Name))
			}
			for _, entry := range action.Entries {
				if entryStr, ok := entry.(string); ok {
					// Handle recharge tags directly
					entryStr = processRechargeTagsDirectly(entryStr)
					// Special case for Mind Control Spores
					if action.Name == "Mind Control Spores" {
						entryStr = strings.ReplaceAll(entryStr, "{@recharge 5}", "")
					}
					md.WriteString(processSpecialFormatting(entryStr))
				}
			}
			md.WriteString("\n\n")
		}
	}

	// Reactions
	if len(monster.Reaction) > 0 {
		md.WriteString("## Reactions\n\n")
		for _, reaction := range monster.Reaction {
			md.WriteString(fmt.Sprintf("***%s.*** ", reaction.Name))
			for _, entry := range reaction.Entries {
				if entryStr, ok := entry.(string); ok {
					md.WriteString(processSpecialFormatting(entryStr))
				}
			}
			md.WriteString("\n\n")
		}
	}

	// Legendary Actions
	if len(monster.Legendary) > 0 {
		md.WriteString("## Legendary Actions\n\n")
		for _, legendary := range monster.Legendary {
			md.WriteString(fmt.Sprintf("***%s.*** ", legendary.Name))
			for _, entry := range legendary.Entries {
				if entryStr, ok := entry.(string); ok {
					md.WriteString(processSpecialFormatting(entryStr))
				}
			}
			md.WriteString("\n\n")
		}
	}

	return md.String(), nil
}

// getSizeString returns the full name of a size from its abbreviation
func getSizeString(size string) string {
	switch size {
	case "T":
		return "Tiny"
	case "S":
		return "Small"
	case "M":
		return "Medium"
	case "L":
		return "Large"
	case "H":
		return "Huge"
	case "G":
		return "Gargantuan"
	default:
		return size
	}
}

// getAbilityModifier calculates the ability modifier for a given ability score
func getAbilityModifier(score int) int {
	// The formula for ability modifiers in D&D 5e is floor((score - 10) / 2)
	// For negative numbers, we need to handle the rounding differently
	if score < 10 {
		return (score - 11) / 2
	}
	return (score - 10) / 2
}

// processRechargeTagsDirectly handles the {@recharge X} format in monster descriptions
// Example: {@recharge 5} -> (Recharge 5-6)
func processRechargeTagsDirectly(text string) string {
	// Simple regex-like replacement for {@recharge X}
	for {
		start := strings.Index(text, "{@recharge ")
		if start == -1 {
			break
		}

		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		rechargeText := text[start+10 : end]
		var displayText string
		if rechargeText == "0" {
			displayText = "(Recharge after a Short or Long Rest)"
		} else {
			displayText = "(Recharge " + rechargeText + "-6)"
		}

		text = text[:start] + displayText + text[end+1:]
	}

	return text
}
