package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/theclifmeister/sample-shifter/internal/categorizer"
	"github.com/theclifmeister/sample-shifter/internal/scanner"
	"github.com/theclifmeister/sample-shifter/internal/stats"
)

var (
	targetDir          string
	outputFile         string
	normalizeFilenames bool
	configFile         string
)

var previewCmd = &cobra.Command{
	Use:   "preview [source-directory]",
	Short: "Preview how files will be categorized and organized",
	Long: `Preview the categorization of audio files without making any changes.
This shows where each file will be copied to when you run the apply command.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourceDir := args[0]

		// Verify source directory exists
		if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
			fmt.Printf("Error: Directory '%s' does not exist\n", sourceDir)
			os.Exit(1)
		}

		if targetDir == "" {
			fmt.Println("Error: --target flag is required")
			os.Exit(1)
		}

		fmt.Printf("Scanning: %s\n", sourceDir)
		fmt.Printf("Target: %s\n\n", targetDir)

		// Scan for sample files
		samples, err := scanner.ScanDirectory(sourceDir)
		if err != nil {
			fmt.Printf("Error scanning directory: %v\n", err)
			os.Exit(1)
		}

		if len(samples) == 0 {
			fmt.Println("No audio sample files found.")
			return
		}

		// Create categorizer with config
		cat, err := categorizer.NewCategorizerFromFile(configFile)
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			os.Exit(1)
		}

		// Categorize files
		categorized := cat.CategorizeBatch(samples, targetDir, normalizeFilenames)

		// Display initial summary
		fmt.Printf("Preview: Found %d file(s) to categorize\n\n", len(categorized))

		// Display detailed file list first
		stats.DisplayDetailedFileList(categorized)

		// Display statistics summary at the end
		stats.DisplayStats(categorized)

		// Save preview to file if requested
		if outputFile != "" {
			savePreview(categorized, outputFile)
		}
	},
}

func savePreview(categorized []categorizer.CategorizedFile, filename string) {
	data, err := json.MarshalIndent(categorized, "", "  ")
	if err != nil {
		fmt.Printf("Error creating preview file: %v\n", err)
		return
	}

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Error creating directory for preview file: %v\n", err)
			return
		}
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		fmt.Printf("Error saving preview file: %v\n", err)
		return
	}

	fmt.Printf("Preview saved to: %s\n", filename)
	fmt.Println("Use this file with the 'apply' command to execute the categorization.")
}

func init() {
	previewCmd.Flags().StringVarP(&targetDir, "target", "t", "", "Target directory for organized samples (required)")
	previewCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Save preview to JSON file for later use with apply command")
	previewCmd.Flags().BoolVar(&normalizeFilenames, "normalize", false, "Normalize filenames (lowercase, spaces and underscores to dashes)")
	previewCmd.Flags().StringVarP(&configFile, "config", "c", "", "Path to category configuration JSON file (optional, uses default if not provided)")
}
