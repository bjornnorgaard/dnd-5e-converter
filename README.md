# D&D 5e Converter

A tool for converting D&D 5e JSON data to Markdown files for easy reference during gameplay.

## Overview

This project converts JSON files containing D&D 5e content (spells, creatures, and items) into well-formatted Markdown files. These Markdown files can then be used with tools like Obsidian for quick reference during gameplay.

## Features

### Spell Parser

The spell parser converts JSON spell data into Markdown files with the following features:

- Converts all spell details (name, level, school, casting time, range, components, duration, etc.)
- Handles special formatting in spell descriptions (damage, dice, spell references, etc.)
- Supports various spell structures (cantrips with scaling, concentration spells, etc.)
- Generates well-structured Markdown with appropriate headers and sections
- Preserves metadata like SRD and Basic Rules flags

#### Example Output

```markdown
# Acid Splash

*Cantrip Conjuration*

**Casting Time:** 1 action

**Range:** 60 feet

**Components:** V, S

**Duration:** Instantaneous

You hurl a bubble of acid. Choose one creature you can see within range, or choose two creatures you can see within range that are within 5 feet of each other. A target must succeed on a Dexterity saving throw or take 1d6 acid damage.

This spell's damage increases by 1d6 when you reach 5th level (2d6), 11th level (3d6), and 17th level (4d6).

**Scaling:**
*acid damage*
- 1st level: 1d6
- 5th level: 2d6
- 11th level: 3d6
- 17th level: 4d6

**Damage Type:** acid

**Saving Throw:** dexterity

**Classes:** Artificer, Sorcerer, Wizard

**Source:** PHB, page 211 (SRD) (Basic Rules)
```

## Usage

To use the converter, run the following command:

```bash
go run main.go
```

This will read the JSON files from the specified data directory and write the Markdown files to the specified output directory.

## Configuration

The converter can be configured by modifying the `Config` struct in `main.go`:

```go
converter := parser.New(parser.Config{
    DataDirectory: filepath.Join("..", "5etools-src", "data"),
    OutDirectory:  filepath.Join(".", "out"),
})
```

- `DataDirectory`: The directory containing the JSON data files
- `OutDirectory`: The directory where the Markdown files will be written
