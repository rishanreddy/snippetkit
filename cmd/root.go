package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var start time.Time

// Root command for snippetkit CLI
var rootCmd = &cobra.Command{
	Use:   "snippetkit",
	Short: "SnippetKit is a CLI tool for managing code snippets",
	Long: `SnippetKit is a command-line tool that helps developers easily 
add, manage, and retrieve code snippets from the snippetkit registry.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚀 Welcome to SnippetKit!\nUse 'snippetkit --help' to see available commands.")
	},
	// Start Timer BEFORE every command runs
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		start = time.Now()
	},

	// Show execution time AFTER every command runs
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		elapsed := time.Since(start)
		fmt.Printf("\nProcessed in %d ms\n", elapsed.Milliseconds())
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌ Error:", err)
		os.Exit(1)
	}
}

func init() {
	// Here you can add global flags if needed
	rootCmd.CompletionOptions.DisableDefaultCmd = false
	// Example: rootCmd.PersistentFlags().StringVar(&flagVar, "flag", "", "A global flag description")
}
