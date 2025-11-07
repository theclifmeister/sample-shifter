# GitHub Copilot Instructions for Sample Shifter

## Project Overview

Sample Shifter is a CLI tool written in Go that organizes audio sample files by automatically categorizing them based on filename keywords. The tool is designed to be:

- **Non-destructive**: Original files remain untouched; copies are made to target directory
- **Preview-first**: Users can preview categorization before applying changes
- **Format-agnostic**: Supports multiple audio formats (WAV, MP3, FLAC, AIF, AIFF, OGG, M4A, WMA, AAC)

## Project Structure

```
sample-shifter/
├── main.go                          # Entry point
├── cmd/                             # CLI commands using cobra
│   ├── root.go                      # Root command and CLI setup
│   ├── scan.go                      # Scan command
│   ├── preview.go                   # Preview command
│   └── apply.go                     # Apply command
└── internal/                        # Internal packages
    ├── scanner/                     # File scanning logic
    │   ├── scanner.go              # Recursive directory scanning for audio files
    │   └── scanner_test.go         # Tests for scanner
    └── categorizer/                 # Categorization logic
        ├── categorizer.go          # Keyword-based categorization
        └── categorizer_test.go     # Tests for categorizer
```

## Key Technologies

- **Language**: Go 1.24.9
- **CLI Framework**: Cobra (github.com/spf13/cobra v1.10.1)
- **Testing**: Standard Go testing package

## Coding Conventions

### Go Style

- Follow idiomatic Go conventions and style guides
- Use `gofmt` for code formatting
- Keep functions small and focused
- Use descriptive variable names (no single-letter variables except for loop counters)
- Add comments for exported functions and types
- Use constants for fixed values (e.g., category names, file extensions)

### Project-Specific Patterns

1. **Package Organization**:
   - `cmd/` package handles CLI interface and user interaction
   - `internal/` packages contain core business logic
   - Keep packages focused on single responsibility

2. **Error Handling**:
   - Return errors from functions; handle them at appropriate level
   - Provide context in error messages
   - Use `fmt.Errorf` for wrapping errors with context

3. **File Operations**:
   - All file operations should be safe and non-destructive by default
   - Validate paths before operations
   - Use `filepath` package for cross-platform path handling

4. **Categorization**:
   - Categories are processed in priority order (defined in `categoryPriority`)
   - First matching category wins (prevents files from being double-categorized)
   - Keywords should be lowercase for case-insensitive matching

## Testing

- Use standard Go testing package (`testing`)
- Test files should be named `*_test.go` alongside source files
- Write table-driven tests where appropriate
- Run tests with: `go test ./...`
- Aim for good coverage of core logic (scanner and categorizer packages)

### Running Tests

```bash
go test ./...                    # Run all tests
go test -v ./...                # Run with verbose output
go test ./internal/categorizer  # Run specific package tests
```

## Building and Running

### Build

```bash
go build -o sample-shifter .
```

### Run

```bash
# Scan directory
./sample-shifter scan /path/to/samples

# Preview categorization
./sample-shifter preview /path/to/samples --target /path/to/organized

# Apply changes
./sample-shifter apply /path/to/samples --target /path/to/organized
```

## Development Workflow

1. **Before making changes**: Run tests to ensure current state
2. **After making changes**: 
   - Run tests: `go test ./...`
   - Build: `go build -o sample-shifter .`
   - Test manually with sample files if needed
3. **For new features**:
   - Add tests first (TDD approach preferred)
   - Implement feature
   - Update README.md if user-facing
   - Keep changes focused and minimal

## Common Tasks

### Adding a New Category

1. Add constant in `internal/categorizer/categorizer.go` (e.g., `CategoryNewType`)
2. Add to `categoryPriority` slice in desired priority position
3. Add keywords to `keywords` map
4. Add tests in `internal/categorizer/categorizer_test.go`
5. Update README.md categories section

### Adding a New Audio Format

1. Add extension to `AudioExtensions` in `internal/scanner/scanner.go`
2. Ensure extension is lowercase with leading dot (e.g., `.mp4`)
3. Add test case in `internal/scanner/scanner_test.go`
4. Update README.md supported formats list

### Adding a New CLI Command

1. Create new file in `cmd/` directory (e.g., `cmd/newcommand.go`)
2. Define command using Cobra pattern (see existing commands)
3. Register command in `cmd/root.go` init function
4. Update README.md with command documentation

## Code Quality Guidelines

- **Simplicity**: Prefer simple, readable code over clever solutions
- **Consistency**: Follow existing patterns in the codebase
- **Documentation**: Update README.md for user-facing changes
- **Testing**: Add tests for new functionality
- **Error messages**: Make them clear and actionable for users
- **Dependencies**: Minimize external dependencies; only add when necessary

## Common Pitfalls to Avoid

- Don't modify original sample files; always copy to target directory
- Don't assume case-sensitive matching; always use `strings.ToLower()`
- Don't skip path validation; always check if directories exist
- Don't add categories without considering priority order
- Don't forget to update both code and documentation together

## When Adding New Features

Consider these questions:
- Does this align with the tool's non-destructive philosophy?
- Should this have a preview mode?
- What error cases need handling?
- How will this work on different operating systems?
- Does the README need updating?
- Are tests needed?
