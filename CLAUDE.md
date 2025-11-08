# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Sample Shifter is a Go CLI tool that organizes audio sample files by automatically categorizing them based on filename keywords. The tool is non-destructive (original files remain untouched), preview-first (users can preview categorization before applying changes), and format-agnostic (supports WAV, MP3, FLAC, AIF, AIFF, OGG, M4A, WMA, AAC).

## Build and Development Commands

### Build
```bash
go build -o sample-shifter .
# Or using the shorthand:
go build -o ss
```

### Run Tests
```bash
go test ./...                    # Run all tests
go test -v ./...                # Run with verbose output
go test ./internal/categorizer  # Run specific package tests
go test ./internal/scanner      # Run scanner tests
```

### Running the CLI
```bash
# Scan directory
./ss scan /path/to/samples

# Preview categorization
./ss preview /path/to/samples --target /path/to/organized

# Preview with output file
./ss preview /path/to/samples --target /path/to/organized --output preview.json

# Apply changes directly
./ss apply /path/to/samples --target /path/to/organized

# Apply from preview file
./ss apply --preview-file preview.json

# Dry run to see what would happen
./ss apply /path/to/samples --target /path/to/organized --dry-run

# Normalize filenames (lowercase, spaces/underscores to dashes)
./ss preview /path/to/samples --target /path/to/organized --normalize
```

## Architecture

### Project Structure
```
sample-shifter/
├── main.go                          # Entry point - calls cmd.Execute()
├── cmd/                             # CLI commands using Cobra framework
│   ├── root.go                      # Root command definition
│   ├── scan.go                      # Scan command - lists audio files
│   ├── preview.go                   # Preview command - shows categorization plan
│   └── apply.go                     # Apply command - executes file copying
└── internal/                        # Core business logic
    ├── scanner/                     # File discovery
    │   ├── scanner.go              # Recursive directory scanning for audio files
    │   └── scanner_test.go
    └── categorizer/                 # File categorization
        ├── categorizer.go          # Keyword-based categorization logic
        └── categorizer_test.go
```

### Key Architectural Patterns

**Separation of Concerns:**
- `cmd/` package handles CLI interface and user interaction
- `internal/scanner` handles file discovery (what files exist)
- `internal/categorizer` handles classification (where files should go)
- Each package has a single, focused responsibility

**Data Flow:**
1. Scanner finds audio files and returns `[]SampleFile`
2. Categorizer processes samples and returns `[]CategorizedFile`
3. Apply command copies files to their target locations

**Priority-Based Categorization:**
The categorizer uses a priority system defined in `categoryPriority` slice. Categories are checked in order, and the first matching category wins. This prevents files from being double-categorized. For example, OneShot has highest priority, so "bass_shot.wav" becomes oneshots/bass rather than bass/.

**Subcategory System:**
Many categories support subcategories for more granular organization. The `subcategoryKeywords` map defines keyword-to-subfolder mappings. The categorizer prefers the longest matching keyword (most specific match). If a category has subcategory support but no keyword matches, files go to an "uncategorized" subfolder within that category.

**Non-Destructive Philosophy:**
All file operations copy rather than move. Original files are never modified or deleted. This safety-first approach runs throughout the codebase.

## Code Conventions

### Category Management
- Category constants are defined in `internal/categorizer/categorizer.go`
- `categoryPriority` slice determines check order (higher priority = earlier in list)
- Keywords must be lowercase for case-insensitive matching
- When adding a new category, update: category constant, categoryPriority, keywords map, tests, and README

### Subcategory Management
- Defined in `subcategoryKeywords` map with Category -> keyword -> subfolder structure
- Longest matching keyword wins (enables specific matches like "bass shot" over "shot")
- If category has subcategory support but no match found, file goes to "uncategorized" subfolder

### File Extensions
- Defined in `scanner.AudioExtensions` slice
- Must be lowercase with leading dot (e.g., `.wav`, `.mp3`)
- When adding new format, update: AudioExtensions, tests, and README

### Testing
- Use table-driven tests where appropriate
- Test files live alongside source files (*_test.go)
- Focus on core logic (scanner and categorizer packages)
- Run tests before and after changes

### Error Handling
- Return errors from functions; handle at appropriate level
- Provide context in error messages using `fmt.Errorf`
- Validate paths before operations
- Use `filepath` package for cross-platform compatibility

## Common Development Patterns

### Adding a New Category
1. Add constant in `internal/categorizer/categorizer.go` (e.g., `CategoryNewType`)
2. Insert into `categoryPriority` slice at desired priority position
3. Add keywords to `keywords` map
4. Optionally add subcategory support in `subcategoryKeywords` map
5. Add test cases in `internal/categorizer/categorizer_test.go`
6. Update README.md categories section

### Adding a New Audio Format
1. Add extension to `AudioExtensions` in `internal/scanner/scanner.go`
2. Ensure lowercase with leading dot (e.g., `.opus`)
3. Add test case in `internal/scanner/scanner_test.go`
4. Update README.md supported formats list

### Adding a New CLI Command
1. Create new file in `cmd/` directory (e.g., `cmd/newcommand.go`)
2. Define command using Cobra pattern (see existing commands)
3. Register command in `cmd/root.go` init function
4. Update README.md with command documentation

## Important Constraints

- Don't modify original sample files; always copy to target directory
- Don't assume case-sensitive matching; always use `strings.ToLower()`
- Don't skip path validation; always check if directories exist
- Don't add categories without considering priority order
- Category priority matters: earlier in `categoryPriority` = higher priority
- Subcategories use longest keyword match for specificity
- All file paths should use `filepath` package for cross-platform support
