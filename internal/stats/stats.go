package stats

import (
	"fmt"
	"sort"

	"github.com/theclifmeister/sample-shifter/internal/categorizer"
)

// categoryCount is a helper struct for sorting categories by count
type categoryCount struct {
	category categorizer.Category
	count    int
}

// DisplayStats shows categorization statistics for a set of categorized files
func DisplayStats(categorized []categorizer.CategorizedFile) {
	if len(categorized) == 0 {
		return
	}

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

	// Sort categories by count (descending)
	var categoryCounts []categoryCount
	for cat, files := range categoryGroups {
		categoryCounts = append(categoryCounts, categoryCount{cat, len(files)})
	}
	// Sort by count (descending)
	sort.Slice(categoryCounts, func(i, j int) bool {
		return categoryCounts[i].count > categoryCounts[j].count
	})

	// Display statistics summary
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
		// Sort by count (descending)
		sort.Slice(subcatCounts, func(i, j int) bool {
			return subcatCounts[i].count > subcatCounts[j].count
		})

		for _, sc := range subcatCounts {
			percentage := float64(sc.count) * 100.0 / float64(categoryTotal)
			fmt.Printf("  %-30s %10d %9.1f%%\n", sc.name, sc.count, percentage)
		}
		fmt.Println()
	}
}

// DisplayDetailedFileList shows a detailed list of all categorized files
func DisplayDetailedFileList(categorized []categorizer.CategorizedFile) {
	if len(categorized) == 0 {
		return
	}

	// Group by category
	categoryGroups := make(map[categorizer.Category][]categorizer.CategorizedFile)
	for _, cat := range categorized {
		categoryGroups[cat.Category] = append(categoryGroups[cat.Category], cat)
	}

	// Sort categories by count (descending)
	var categoryCounts []categoryCount
	for cat, files := range categoryGroups {
		categoryCounts = append(categoryCounts, categoryCount{cat, len(files)})
	}
	// Sort by count (descending)
	sort.Slice(categoryCounts, func(i, j int) bool {
		return categoryCounts[i].count > categoryCounts[j].count
	})

	// Display detailed file list
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
}
