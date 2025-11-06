package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/theclifmeister/sample-shifter/internal/scanner"
)

var scanCmd = &cobra.Command{
	Use:   "scan [directory]",
	Short: "Scan a directory for audio sample files",
	Long:  `Scan a directory recursively to find all audio sample files.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourceDir := args[0]

		// Verify source directory exists
		if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
			fmt.Printf("Error: Directory '%s' does not exist\n", sourceDir)
			os.Exit(1)
		}

		fmt.Printf("Scanning directory: %s\n\n", sourceDir)

		samples, err := scanner.ScanDirectory(sourceDir)
		if err != nil {
			fmt.Printf("Error scanning directory: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Found %d audio sample file(s):\n\n", len(samples))

		for _, sample := range samples {
			fmt.Printf("  - %s\n", sample.OriginalPath)
		}
	},
}
