package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// SpellIndex represents the structure of the spells/index.json file
type SpellIndex map[string]string

// SpellFile represents the structure of a spell JSON file
type SpellFile struct {
	Spell []Spell `json:"spell"`
}

// Spell represents a single spell entry
type Spell struct {
	Name             string                 `json:"name"`
	Source           string                 `json:"source"`
	Page             int                    `json:"page,omitempty"`
	Level            int                    `json:"level"`
	School           string                 `json:"school"`
	Time             []SpellTime            `json:"time"`
	Range            SpellRange             `json:"range"`
	Components       SpellComponents        `json:"components"`
	Duration         []SpellDuration        `json:"duration"`
	Entries          []interface{}          `json:"entries"`
	EntriesHigher    []interface{}          `json:"entriesHigher,omitempty"`
	DamageInflict    []string               `json:"damageInflict,omitempty"`
	SavingThrow      []string               `json:"savingThrow,omitempty"`
	MiscTags         []string               `json:"miscTags,omitempty"`
	AreaTags         []string               `json:"areaTags,omitempty"`
	Classes          SpellClasses           `json:"classes,omitempty"`
	Meta             map[string]interface{} `json:"meta,omitempty"`
	ScalingLevelDice interface{}            `json:"scalingLevelDice,omitempty"`
	// Fields for handling variations in spell data
	SRD         interface{} `json:"srd,omitempty"`
	BasicRules  interface{} `json:"basicRules,omitempty"`
	ReprintedAs []string    `json:"reprintedAs,omitempty"`
	// Additional fields can be added as needed
}

// SpellTime represents the casting time of a spell
type SpellTime struct {
	Number int    `json:"number"`
	Unit   string `json:"unit"`
}

// SpellRange represents the range of a spell
type SpellRange struct {
	Type     string           `json:"type"`
	Distance SpellRangeDetail `json:"distance,omitempty"`
}

// SpellRangeDetail represents the distance details of a spell range
type SpellRangeDetail struct {
	Type   string `json:"type"`
	Amount int    `json:"amount,omitempty"`
}

// SpellComponents represents the components required to cast a spell
type SpellComponents struct {
	V bool        `json:"v,omitempty"`
	S bool        `json:"s,omitempty"`
	M interface{} `json:"m,omitempty"` // Can be string or object
	R bool        `json:"r,omitempty"`
}

// SpellDuration represents the duration of a spell
type SpellDuration struct {
	Type      string                 `json:"type"`
	Duration  map[string]interface{} `json:"duration,omitempty"`
	Condition string                 `json:"condition,omitempty"`
}

// SpellClasses represents the classes that can use a spell
type SpellClasses struct {
	FromClassList []SpellClass `json:"fromClassList,omitempty"`
	// Additional class-related fields can be added as needed
}

// SpellClass represents a class that can use a spell
type SpellClass struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

// ScalingLevelDice represents the scaling damage dice for cantrips and other scaling spells
type ScalingLevelDice struct {
	Label   string            `json:"label,omitempty"`
	Scaling map[string]string `json:"scaling"`
}

// parseSpells parses the spell data from the specified directory and writes it to the output directory.
func parseSpells(ctx context.Context, dataDirectory, outDirectory string) error {
	var (
		spellsPath = filepath.Join(dataDirectory, "spells")
		indexPath  = filepath.Join(spellsPath, "index.json")
	)

	// Create output directory if it doesn't exist
	outDir := filepath.Join(outDirectory, "spells")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Read and parse the index file
	indexData, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read index file: %w", err)
	}

	var index SpellIndex
	if err := json.Unmarshal(indexData, &index); err != nil {
		return fmt.Errorf("failed to parse index file: %w", err)
	}

	// Process each spell file
	for source, filename := range index {
		if err := processSpellFile(ctx, spellsPath, outDir, source, filename); err != nil {
			return fmt.Errorf("failed to process spell file %s: %w", filename, err)
		}
	}

	return nil
}

// processSpellFile processes a single spell file and generates Markdown for each spell
func processSpellFile(ctx context.Context, spellsPath, outDir, source, filename string) error {
	// Read and parse the spell file
	filePath := filepath.Join(spellsPath, filename)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read spell file: %w", err)
	}

	var spellFile SpellFile
	if err := json.Unmarshal(fileData, &spellFile); err != nil {
		return fmt.Errorf("failed to parse spell file: %w", err)
	}

	// Process each spell
	for _, spell := range spellFile.Spell {
		mdContent, err := spellToMarkdown(spell)
		if err != nil {
			return fmt.Errorf("failed to convert spell to markdown: %w", err)
		}

		// Create a filename for the spell
		safeSpellName := strings.ReplaceAll(spell.Name, "/", "-")
		safeSpellName = strings.ReplaceAll(safeSpellName, "\\", "-")
		safeSpellName = strings.ReplaceAll(safeSpellName, ":", "-")
		safeSpellName = strings.ReplaceAll(safeSpellName, "*", "-")
		safeSpellName = strings.ReplaceAll(safeSpellName, "?", "-")
		safeSpellName = strings.ReplaceAll(safeSpellName, "\"", "-")
		safeSpellName = strings.ReplaceAll(safeSpellName, "<", "-")
		safeSpellName = strings.ReplaceAll(safeSpellName, ">", "-")
		safeSpellName = strings.ReplaceAll(safeSpellName, "|", "-")

		mdFilePath := filepath.Join(outDir, safeSpellName+".md")

		// Write the markdown file
		if err := os.WriteFile(mdFilePath, []byte(mdContent), 0644); err != nil {
			return fmt.Errorf("failed to write markdown file: %w", err)
		}
	}

	return nil
}

// spellToMarkdown converts a spell to Markdown format
func spellToMarkdown(spell Spell) (string, error) {
	var md strings.Builder

	// Title
	md.WriteString(fmt.Sprintf("# %s\n\n", spell.Name))

	// Basic info
	md.WriteString(fmt.Sprintf("*%s %s*\n\n", getSpellLevel(spell.Level), getSchoolName(spell.School)))

	// Casting Time
	md.WriteString("**Casting Time:** ")
	for i, time := range spell.Time {
		if i > 0 {
			md.WriteString(", ")
		}
		md.WriteString(fmt.Sprintf("%d %s", time.Number, time.Unit))
	}
	md.WriteString("\n\n")

	// Range
	md.WriteString("**Range:** ")
	if spell.Range.Type == "point" {
		if spell.Range.Distance.Type == "self" {
			md.WriteString("Self")
		} else if spell.Range.Distance.Type == "touch" {
			md.WriteString("Touch")
		} else if spell.Range.Distance.Type == "sight" {
			md.WriteString("Sight")
		} else {
			md.WriteString(fmt.Sprintf("%d %s", spell.Range.Distance.Amount, spell.Range.Distance.Type))
		}
	} else if spell.Range.Type == "radius" {
		md.WriteString(fmt.Sprintf("%d-%s radius", spell.Range.Distance.Amount, spell.Range.Distance.Type))
	} else if spell.Range.Type == "cone" {
		md.WriteString(fmt.Sprintf("%d-%s cone", spell.Range.Distance.Amount, spell.Range.Distance.Type))
	} else if spell.Range.Type == "line" {
		md.WriteString(fmt.Sprintf("%d-%s line", spell.Range.Distance.Amount, spell.Range.Distance.Type))
	} else if spell.Range.Type == "cube" {
		md.WriteString(fmt.Sprintf("%d-%s cube", spell.Range.Distance.Amount, spell.Range.Distance.Type))
	} else {
		md.WriteString(spell.Range.Type)
	}
	md.WriteString("\n\n")

	// Components
	md.WriteString("**Components:** ")
	var components []string
	if spell.Components.V {
		components = append(components, "V")
	}
	if spell.Components.S {
		components = append(components, "S")
	}
	if spell.Components.M != nil {
		mStr := "M"
		switch m := spell.Components.M.(type) {
		case string:
			mStr = fmt.Sprintf("M (%s)", m)
		case map[string]interface{}:
			if text, ok := m["text"].(string); ok {
				mStr = fmt.Sprintf("M (%s)", text)
			}
		}
		components = append(components, mStr)
	}
	if spell.Components.R {
		components = append(components, "R")
	}
	md.WriteString(strings.Join(components, ", "))
	md.WriteString("\n\n")

	// Duration
	md.WriteString("**Duration:** ")
	for i, duration := range spell.Duration {
		if i > 0 {
			md.WriteString(", ")
		}

		if duration.Type == "instant" {
			md.WriteString("Instantaneous")
		} else if duration.Type == "timed" {
			if duration.Duration != nil {
				if amount, ok := duration.Duration["amount"].(float64); ok {
					unit, _ := duration.Duration["type"].(string)
					md.WriteString(fmt.Sprintf("%d %s", int(amount), unit))
				}
			}
		} else if duration.Type == "permanent" {
			md.WriteString("Until dispelled")
			if duration.Condition != "" {
				md.WriteString(fmt.Sprintf(" or %s", duration.Condition))
			}
		} else if duration.Type == "concentration" {
			md.WriteString("Concentration")
			if duration.Duration != nil {
				if amount, ok := duration.Duration["amount"].(float64); ok {
					unit, _ := duration.Duration["type"].(string)
					md.WriteString(fmt.Sprintf(", up to %d %s", int(amount), unit))
				}
			}
		} else {
			md.WriteString(duration.Type)
		}
	}
	md.WriteString("\n\n")

	// Description
	for _, entry := range spell.Entries {
		switch e := entry.(type) {
		case string:
			// Process special formatting in the text
			processedText := processSpecialFormatting(e)
			md.WriteString(processedText)
			md.WriteString("\n\n")
		case map[string]interface{}:
			// Handle other entry types like lists, tables, etc.
			if entryType, ok := e["type"].(string); ok {
				if entryType == "list" {
					if items, ok := e["items"].([]interface{}); ok {
						for _, item := range items {
							if itemStr, ok := item.(string); ok {
								md.WriteString(fmt.Sprintf("- %s\n", processSpecialFormatting(itemStr)))
							} else if itemMap, ok := item.(map[string]interface{}); ok {
								if itemName, ok := itemMap["name"].(string); ok {
									md.WriteString(fmt.Sprintf("- **%s**", itemName))
									if itemEntry, ok := itemMap["entry"].(string); ok {
										md.WriteString(fmt.Sprintf(": %s", processSpecialFormatting(itemEntry)))
									}
									md.WriteString("\n")
								}
							}
						}
						md.WriteString("\n")
					}
				} else if entryType == "table" {
					// Handle tables
					md.WriteString("| ")
					if colLabels, ok := e["colLabels"].([]interface{}); ok {
						for i, col := range colLabels {
							if colStr, ok := col.(string); ok {
								md.WriteString(colStr)
								if i < len(colLabels)-1 {
									md.WriteString(" | ")
								}
							}
						}
					}
					md.WriteString(" |\n")

					md.WriteString("| ")
					if colLabels, ok := e["colLabels"].([]interface{}); ok {
						for i := range colLabels {
							md.WriteString("---")
							if i < len(colLabels)-1 {
								md.WriteString(" | ")
							}
						}
					}
					md.WriteString(" |\n")

					if rows, ok := e["rows"].([]interface{}); ok {
						for _, row := range rows {
							if rowArr, ok := row.([]interface{}); ok {
								md.WriteString("| ")
								for i, cell := range rowArr {
									if cellStr, ok := cell.(string); ok {
										md.WriteString(processSpecialFormatting(cellStr))
									} else if cellMap, ok := cell.(map[string]interface{}); ok {
										if cellText, ok := cellMap["text"].(string); ok {
											md.WriteString(processSpecialFormatting(cellText))
										}
									}
									if i < len(rowArr)-1 {
										md.WriteString(" | ")
									}
								}
								md.WriteString(" |\n")
							}
						}
					}
					md.WriteString("\n")
				}
				// Add handling for other entry types as needed
			}
		}
	}

	// Scaling Level Dice (for cantrips)
	if spell.ScalingLevelDice != nil {
		md.WriteString("**Scaling:**\n")

		// Handle different types of scalingLevelDice
		switch sld := spell.ScalingLevelDice.(type) {
		case map[string]interface{}:
			// Handle object format
			if label, ok := sld["label"].(string); ok && label != "" {
				md.WriteString(fmt.Sprintf("*%s*\n", label))
			}

			if scaling, ok := sld["scaling"].(map[string]interface{}); ok {
				var levels []string
				for level := range scaling {
					levels = append(levels, level)
				}

				// Sort levels numerically
				sort.Slice(levels, func(i, j int) bool {
					a, _ := strconv.Atoi(levels[i])
					b, _ := strconv.Atoi(levels[j])
					return a < b
				})

				for _, level := range levels {
					dice := scaling[level]
					levelText := formatLevel(level)
					md.WriteString(fmt.Sprintf("- %s: %v\n", levelText, dice))
				}
			}
		case []interface{}:
			// Handle array format
			for _, item := range sld {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if level, ok := itemMap["level"].(float64); ok {
						levelText := formatLevel(fmt.Sprintf("%d", int(level)))
						if dice, ok := itemMap["dice"].(map[string]interface{}); ok {
							if diceCount, ok := dice["count"].(float64); ok {
								if diceSize, ok := dice["faces"].(float64); ok {
									md.WriteString(fmt.Sprintf("- %s: %dd%d\n", levelText, int(diceCount), int(diceSize)))
								}
							}
						} else if diceStr, ok := itemMap["dice"].(string); ok {
							md.WriteString(fmt.Sprintf("- %s: %s\n", levelText, diceStr))
						}
					}
				}
			}
		}

		md.WriteString("\n")
	}

	// Higher Levels
	if len(spell.EntriesHigher) > 0 {
		md.WriteString("**At Higher Levels:** ")
		for _, entry := range spell.EntriesHigher {
			if entryStr, ok := entry.(string); ok {
				md.WriteString(processSpecialFormatting(entryStr))
				md.WriteString("\n\n")
			} else if entryMap, ok := entry.(map[string]interface{}); ok {
				if entryType, ok := entryMap["type"].(string); ok && entryType == "entries" {
					if entries, ok := entryMap["entries"].([]interface{}); ok {
						for _, subEntry := range entries {
							if subEntryStr, ok := subEntry.(string); ok {
								md.WriteString(processSpecialFormatting(subEntryStr))
								md.WriteString("\n\n")
							}
						}
					}
				}
			}
		}
	}

	// Damage Type
	if len(spell.DamageInflict) > 0 {
		md.WriteString("**Damage Type:** ")
		md.WriteString(strings.Join(spell.DamageInflict, ", "))
		md.WriteString("\n\n")
	}

	// Saving Throw
	if len(spell.SavingThrow) > 0 {
		md.WriteString("**Saving Throw:** ")
		md.WriteString(strings.Join(spell.SavingThrow, ", "))
		md.WriteString("\n\n")
	}

	// Classes
	if spell.Classes.FromClassList != nil && len(spell.Classes.FromClassList) > 0 {
		md.WriteString("**Classes:** ")
		var classNames []string
		for _, class := range spell.Classes.FromClassList {
			classNames = append(classNames, class.Name)
		}
		md.WriteString(strings.Join(classNames, ", "))
		md.WriteString("\n\n")
	}

	// Source
	md.WriteString(fmt.Sprintf("**Source:** %s", spell.Source))
	if spell.Page > 0 {
		md.WriteString(fmt.Sprintf(", page %d", spell.Page))
	}

	// Handle SRD field which can be bool or string
	switch srd := spell.SRD.(type) {
	case bool:
		if srd {
			md.WriteString(" (SRD)")
		}
	case string:
		if srd == "true" || srd == "1" {
			md.WriteString(" (SRD)")
		}
	}

	// Handle BasicRules field which can be bool or string
	switch br := spell.BasicRules.(type) {
	case bool:
		if br {
			md.WriteString(" (Basic Rules)")
		}
	case string:
		if br == "true" || br == "1" {
			md.WriteString(" (Basic Rules)")
		}
	}

	md.WriteString("\n")

	return md.String(), nil
}

// getSpellLevel returns a string representation of the spell level
func getSpellLevel(level int) string {
	if level == 0 {
		return "Cantrip"
	}

	switch level {
	case 1:
		return "1st-level"
	case 2:
		return "2nd-level"
	case 3:
		return "3rd-level"
	default:
		return fmt.Sprintf("%dth-level", level)
	}
}

// getSchoolName returns the full name of a spell school
func getSchoolName(school string) string {
	schools := map[string]string{
		"A": "Abjuration",
		"C": "Conjuration",
		"D": "Divination",
		"E": "Enchantment",
		"V": "Evocation",
		"I": "Illusion",
		"N": "Necromancy",
		"T": "Transmutation",
	}

	if fullName, ok := schools[school]; ok {
		return fullName
	}
	return school
}

// formatLevel formats a level number as a string (e.g., "1" -> "1st level")
func formatLevel(level string) string {
	switch level {
	case "1":
		return "1st level"
	case "2":
		return "2nd level"
	case "3":
		return "3rd level"
	case "5":
		return "5th level"
	case "11":
		return "11th level"
	case "17":
		return "17th level"
	default:
		return level + "th level"
	}
}

// processSpecialFormatting handles the special formatting in spell descriptions
// The JSON data contains special formatting tags like {@damage X} and {@dice X}
// that need to be converted to plain text for the Markdown output.
func processSpecialFormatting(text string) string {
	// Handle {@damage X} format - converts damage tags to plain text
	text = processDamageTag(text)

	// Handle {@dice X} format - converts dice tags to plain text
	text = processDiceTag(text)

	// Handle {@spell X} format - converts spell references to plain text
	text = processSpellTag(text)

	// Handle {@item X} format - converts item references to plain text
	text = processItemTag(text)

	// Handle {@creature X} format - converts creature references to plain text
	text = processCreatureTag(text)

	// Handle {@condition X} format - converts condition references to plain text
	text = processConditionTag(text)

	// Handle {@hazard X} format - converts hazard references to plain text
	text = processHazardTag(text)

	// Handle {@atk X} format - converts attack tags to plain text
	text = processAttackTag(text)

	// Handle {@hit X} format - converts hit bonus tags to plain text
	text = processHitTag(text)

	// Handle {@h} format - converts hit tags to plain text
	text = processHTag(text)

	// Handle {@dc X} format - converts DC tags to plain text
	text = processDCTag(text)

	// Handle {@recharge X} format - converts recharge tags to plain text
	text = processRechargeTag(text)

	return text
}

// processDamageTag handles the {@damage X} format in spell descriptions
// Example: {@damage 1d6} -> 1d6
func processDamageTag(text string) string {
	// Simple regex-like replacement for {@damage X}
	for {
		start := strings.Index(text, "{@damage ")
		if start == -1 {
			break
		}

		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		damageText := text[start+9 : end]
		text = text[:start] + damageText + text[end+1:]
	}

	return text
}

// processDiceTag handles the {@dice X} format in spell descriptions
// Example: {@dice 1d6} -> 1d6
func processDiceTag(text string) string {
	// Simple regex-like replacement for {@dice X}
	for {
		start := strings.Index(text, "{@dice ")
		if start == -1 {
			break
		}

		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		diceText := text[start+7 : end]
		text = text[:start] + diceText + text[end+1:]
	}

	return text
}

// processSpellTag handles the {@spell X} format in spell descriptions
// Example: {@spell fireball} -> fireball
func processSpellTag(text string) string {
	// Simple regex-like replacement for {@spell X}
	for {
		start := strings.Index(text, "{@spell ")
		if start == -1 {
			break
		}

		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		// Extract the spell name, handling the case where there might be a pipe character
		// Format can be {@spell name} or {@spell name|display text}
		spellText := text[start+8 : end]
		parts := strings.Split(spellText, "|")
		displayText := parts[0]
		if len(parts) > 1 {
			displayText = parts[1]
		}

		text = text[:start] + displayText + text[end+1:]
	}

	return text
}

// processItemTag handles the {@item X} format in spell descriptions
// Example: {@item potion of healing} -> potion of healing
func processItemTag(text string) string {
	// Simple regex-like replacement for {@item X}
	for {
		start := strings.Index(text, "{@item ")
		if start == -1 {
			break
		}

		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		// Extract the item name, handling the case where there might be a pipe character
		// Format can be {@item name} or {@item name|display text}
		itemText := text[start+7 : end]
		parts := strings.Split(itemText, "|")
		displayText := parts[0]
		if len(parts) > 1 {
			displayText = parts[1]
		}

		text = text[:start] + displayText + text[end+1:]
	}

	return text
}

// processCreatureTag handles the {@creature X} format in spell descriptions
// Example: {@creature goblin} -> goblin
func processCreatureTag(text string) string {
	// Simple regex-like replacement for {@creature X}
	for {
		start := strings.Index(text, "{@creature ")
		if start == -1 {
			break
		}

		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		// Extract the creature name, handling the case where there might be a pipe character
		// Format can be {@creature name} or {@creature name|display text}
		creatureText := text[start+11 : end]
		parts := strings.Split(creatureText, "|")
		displayText := parts[0]
		if len(parts) > 1 {
			displayText = parts[1]
		}

		text = text[:start] + displayText + text[end+1:]
	}

	return text
}

// processConditionTag handles the {@condition X} format in spell descriptions
// Example: {@condition poisoned} -> poisoned
func processConditionTag(text string) string {
	// Simple regex-like replacement for {@condition X}
	for {
		start := strings.Index(text, "{@condition ")
		if start == -1 {
			break
		}

		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		conditionText := text[start+12 : end]
		text = text[:start] + conditionText + text[end+1:]
	}

	return text
}

// processHazardTag handles the {@hazard X} format in spell descriptions
// Example: {@hazard burning|XPHB} -> burning
func processHazardTag(text string) string {
	// Simple regex-like replacement for {@hazard X}
	for {
		start := strings.Index(text, "{@hazard ")
		if start == -1 {
			break
		}

		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		// Extract the hazard name, handling the case where there might be a pipe character
		// Format can be {@hazard name} or {@hazard name|source}
		hazardText := text[start+9 : end]
		parts := strings.Split(hazardText, "|")
		displayText := parts[0]

		text = text[:start] + displayText + text[end+1:]
	}

	return text
}

// processAttackTag handles the {@atk X} format in monster descriptions
// Example: {@atk mw} -> (melee weapon)
func processAttackTag(text string) string {
	// Simple regex-like replacement for {@atk X}
	for {
		start := strings.Index(text, "{@atk ")
		if start == -1 {
			break
		}

		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		atkText := text[start+6 : end]
		var displayText string
		switch atkText {
		case "mw":
			displayText = "(melee weapon)"
		case "rw":
			displayText = "(ranged weapon)"
		case "ms":
			displayText = "(melee spell)"
		case "rs":
			displayText = "(ranged spell)"
		default:
			displayText = ""
		}

		text = text[:start] + displayText + text[end+1:]
	}

	return text
}

// processHitTag handles the {@hit X} format in monster descriptions
// Example: {@hit 13} -> +13
func processHitTag(text string) string {
	// Simple regex-like replacement for {@hit X}
	for {
		start := strings.Index(text, "{@hit ")
		if start == -1 {
			break
		}

		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		hitText := text[start+6 : end]
		displayText := "+" + hitText

		text = text[:start] + displayText + text[end+1:]
	}

	return text
}

// processHTag handles the {@h} format in monster descriptions
// Example: {@h} -> Hit:
func processHTag(text string) string {
	// Simple replacement for {@h}
	return strings.ReplaceAll(text, "{@h}", "Hit:")
}

// processDCTag handles the {@dc X} format in monster descriptions
// Example: {@dc 19} -> DC 19
func processDCTag(text string) string {
	// Simple regex-like replacement for {@dc X}
	for {
		start := strings.Index(text, "{@dc ")
		if start == -1 {
			break
		}

		end := strings.Index(text[start:], "}")
		if end == -1 {
			break
		}
		end += start

		dcText := text[start+5 : end]
		displayText := "DC " + dcText

		text = text[:start] + displayText + text[end+1:]
	}

	return text
}

// processRechargeTag handles the {@recharge X} format in monster descriptions
// Example: {@recharge 5} -> (Recharge 5-6)
func processRechargeTag(text string) string {
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
