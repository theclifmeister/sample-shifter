package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/theclifmeister/sample-shifter/internal/categorizer"
	"github.com/theclifmeister/sample-shifter/internal/scanner"
)

var (
	applyTargetDir            string
	previewFile               string
	dryRun                    bool
	applyNormalizeFilenames   bool
)

var applyCmd = &cobra.Command{
	Use:   "apply [source-directory]",
	Short: "Apply categorization and copy files to target directory",
	Long: `Copy audio files to their categorized folders in the target directory.
You can either specify a source directory to scan and categorize on-the-fly,
or use a previously generated preview file.`,
	Run: func(cmd *cobra.Command, args []string) {
		var categorized []categorizer.CategorizedFile

		// Load from preview file if provided
		if previewFile != "" {
			data, err := os.ReadFile(previewFile)
			if err != nil {
				fmt.Printf("Error reading preview file: %v\n", err)
				os.Exit(1)
			}

			if err := json.Unmarshal(data, &categorized); err != nil {
				fmt.Printf("Error parsing preview file: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Loaded preview from: %s\n", previewFile)
		} else {
			// Scan and categorize on-the-fly
			if len(args) != 1 {
				fmt.Println("Error: source directory required when not using --preview-file")
				os.Exit(1)
			}

			if applyTargetDir == "" {
				fmt.Println("Error: --target flag is required")
				os.Exit(1)
			}

			sourceDir := args[0]

			// Verify source directory exists
			if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
				fmt.Printf("Error: Directory '%s' does not exist\n", sourceDir)
				os.Exit(1)
			}

			fmt.Printf("Scanning: %s\n", sourceDir)

			samples, err := scanner.ScanDirectory(sourceDir)
			if err != nil {
				fmt.Printf("Error scanning directory: %v\n", err)
				os.Exit(1)
			}

			categorized = categorizer.CategorizeBatch(samples, applyTargetDir, applyNormalizeFilenames)
		}

		if len(categorized) == 0 {
			fmt.Println("No files to process.")
			return
		}

		if dryRun {
			fmt.Println("\n=== DRY RUN MODE - No files will be copied ===")
		}

		fmt.Printf("\nProcessing %d file(s)...\n\n", len(categorized))

		// Copy files
		successCount := 0
		errorCount := 0

		for _, cat := range categorized {
			fmt.Printf("Copying: %s\n  -> %s\n", cat.Sample.OriginalPath, cat.TargetPath)

			if !dryRun {
				if err := copyFile(cat.Sample.OriginalPath, cat.TargetPath); err != nil {
					fmt.Printf("  ERROR: %v\n", err)
					errorCount++
				} else {
					fmt.Println("  ✓ Success")
					successCount++
				}
			} else {
				fmt.Println("  (skipped - dry run)")
				successCount++
			}
		}

		fmt.Printf("\n=== Summary ===\n")
		fmt.Printf("Total files: %d\n", len(categorized))
		fmt.Printf("Successful: %d\n", successCount)
		if errorCount > 0 {
			fmt.Printf("Errors: %d\n", errorCount)
		}
		if dryRun {
			fmt.Println("\nThis was a dry run. Use without --dry-run to actually copy files.")
		}
	},
}

func copyFile(src, dst string) error {
	// Create target directory if it doesn't exist
	targetDir := filepath.Dir(dst)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy contents
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func init() {
	applyCmd.Flags().StringVarP(&applyTargetDir, "target", "t", "", "Target directory for organized samples")
	applyCmd.Flags().StringVarP(&previewFile, "preview-file", "p", "", "Use a previously saved preview file")
	applyCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview what would be done without actually copying files")
	applyCmd.Flags().BoolVar(&applyNormalizeFilenames, "normalize", false, "Normalize filenames (lowercase, spaces and underscores to dashes)")
}
