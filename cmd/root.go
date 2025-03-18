package cmd

import (
	"fmt"
	"os"
	"snippetkit/internal"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	titleStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ea580c")) // Orange
	labelStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#fbbf24"))            // Light Orange
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#10b981"))            // Green
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444"))            // Red
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#f59e0b"))            // Yellow
	divider      = lipgloss.NewStyle().Foreground(lipgloss.Color("#444")).Render(strings.Repeat("─", 40))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6b7280")) // Gray
)

var rootCmd = &cobra.Command{
	Use:   "snippetkit",
	Short: "SnippetKit - Easily manage reusable code snippets",
	Long:  `SnippetKit CLI allows you to search, add, and manage code snippets.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		internal.LoadConfig() // Load config before executing commands
		internal.InitLogger()
	},
	Run: func(cmd *cobra.Command, args []string) {
		version := "internal.GetVersion() " // Pull the version dynamically
		internal.Info(fmt.Sprintf("SnippetKit CLI v%s", version), nil)
		cmd.Help() // Display the help command
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Global Persistent Flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "Specify config file (default is $HOME/.snippetkit/config.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	// Add logging flag
	rootCmd.PersistentFlags().Bool("logging", true, "Enable or disable logging")
	viper.BindPFlag("logging_enabled", rootCmd.PersistentFlags().Lookup("logging"))
}
