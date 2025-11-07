package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/theclifmeister/sample-shifter/internal/categorizer"
	"github.com/theclifmeister/sample-shifter/internal/scanner"
)

var (
	targetDir          string
	outputFile         string
	normalizeFilenames bool
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

		// Categorize files
		categorized := categorizer.CategorizeBatch(samples, targetDir, normalizeFilenames)

		// Group by category and subcategory
		categoryGroups := make(map[categorizer.Category][]categorizer.CategorizedFile)
		subcategoryGroups := make(map[categorizer.Category]map[string][]categorizer.CategorizedFile)

		for _, cat := range categorized {
			categoryGroups[cat.Category] = append(categoryGroups[cat.Category], cat)

			if subcategoryGroups[cat.Category] == nil {
				subcategoryGroups[cat.Category] = make(map[string][]categorizer.CategorizedFile)
			}
			subcatKey := cat.Subcategory
			if subcatKey == "" {
				subcatKey = "(no subcategory)"
			}
			subcategoryGroups[cat.Category][subcatKey] = append(subcategoryGroups[cat.Category][subcatKey], cat)
		}

		totalFiles := len(categorized)

		// Display initial summary
		fmt.Printf("Preview: Found %d file(s) to categorize\n\n", totalFiles)

		// Sort categories by count (descending)
		type categoryCount struct {
			category categorizer.Category
			count    int
		}
		var categoryCounts []categoryCount
		for cat, files := range categoryGroups {
			categoryCounts = append(categoryCounts, categoryCount{cat, len(files)})
		}
		// Simple bubble sort by count
		for i := 0; i < len(categoryCounts)-1; i++ {
			for j := 0; j < len(categoryCounts)-i-1; j++ {
				if categoryCounts[j].count < categoryCounts[j+1].count {
					categoryCounts[j], categoryCounts[j+1] = categoryCounts[j+1], categoryCounts[j]
				}
			}
		}

		// Display detailed file list first
		fmt.Println("=== DETAILED FILE LIST ===")
		fmt.Println()

		for _, cc := range categoryCounts {
			category := cc.category
			files := categoryGroups[category]
			fmt.Printf("Category: %s (%d files)\n", category, len(files))
			for _, file := range files {
				fmt.Printf("  %s\n    -> %s\n", file.Sample.OriginalPath, file.TargetPath)
			}
			fmt.Println()
		}

		// Display statistics summary at the end
		fmt.Println("=== CATEGORIZATION STATISTICS ===")
		fmt.Println()

		// Display category statistics
		fmt.Printf("%-20s %10s %10s\n", "Category", "Count", "Percentage")
		fmt.Println("--------------------------------------------------")

		for _, cc := range categoryCounts {
			percentage := float64(cc.count) * 100.0 / float64(totalFiles)
			fmt.Printf("%-20s %10d %9.1f%%\n", cc.category, cc.count, percentage)
		}
		fmt.Println()

		// Display subcategory statistics for each category
		fmt.Println("=== SUBCATEGORY BREAKDOWN ===")
		fmt.Println()

		for _, cc := range categoryCounts {
			category := cc.category
			categoryTotal := cc.count
			subcats := subcategoryGroups[category]

			fmt.Printf("%s (%d files)\n", category, categoryTotal)
			fmt.Printf("  %-30s %10s %10s\n", "Subcategory", "Count", "% of Cat")
			fmt.Println("  --------------------------------------------------")

			// Sort subcategories by count
			type subcatCount struct {
				name  string
				count int
			}
			var subcatCounts []subcatCount
			for subcat, files := range subcats {
				subcatCounts = append(subcatCounts, subcatCount{subcat, len(files)})
			}
			// Simple bubble sort
			for i := 0; i < len(subcatCounts)-1; i++ {
				for j := 0; j < len(subcatCounts)-i-1; j++ {
					if subcatCounts[j].count < subcatCounts[j+1].count {
						subcatCounts[j], subcatCounts[j+1] = subcatCounts[j+1], subcatCounts[j]
					}
				}
			}

			for _, sc := range subcatCounts {
				percentage := float64(sc.count) * 100.0 / float64(categoryTotal)
				fmt.Printf("  %-30s %10d %9.1f%%\n", sc.name, sc.count, percentage)
			}
			fmt.Println()
		}

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
}
