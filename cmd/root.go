package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sample-shifter",
	Short: "A CLI tool to organize audio sample files",
	Long: `Sample Shifter is a CLI tool that scans audio sample files,
categorizes them based on their names, and organizes them into
appropriate folders. The process is non-destructive and allows
previewing changes before applying them.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(previewCmd)
	rootCmd.AddCommand(applyCmd)
}
